package approval

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/kageyamountain/ghrev/internal/common/progress"
	"github.com/kageyamountain/ghrev/internal/domain/aggregate/mygithub"
)

type UseCase struct {
	runtimeOptions *RuntimeOptions
	githubGateway  mygithub.Gateway
}

func NewUseCase(
	runtimeOptions *RuntimeOptions,
	githubGateway mygithub.Gateway,
) *UseCase {
	return &UseCase{
		runtimeOptions: runtimeOptions,
		githubGateway:  githubGateway,
	}
}

const maxConcurrency = 10

func (u *UseCase) Do(ctx context.Context) error {
	progressStopFunc := progress.Start("計測中")

	summaries, err := u.githubGateway.FindPullRequestSummaries(ctx, u.runtimeOptions.Owner, u.runtimeOptions.Name, u.runtimeOptions.CreatedAtFrom)
	if err != nil {
		progressStopFunc()
		return fmt.Errorf("failed to find pull request summaries: %w", err)
	}

	resultRows := make([]string, len(summaries))
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(maxConcurrency)
	for i, summary := range summaries {
		eg.Go(func() error {
			row, err2 := u.measureApprovalTime(egCtx, summary)
			if err2 != nil {
				return err2
			}
			resultRows[i] = row
			return nil
		})
	}
	err = eg.Wait()
	progressStopFunc()
	if err != nil {
		return fmt.Errorf("failed to measure approval times: %w", err)
	}

	var header bool
	for _, resultRow := range resultRows {
		if resultRow == "" {
			continue
		}
		if !header {
			fmt.Println("URL 所要時間 変更行数")
			header = true
		}
		fmt.Println(resultRow)
	}
	if !header {
		fmt.Printf("%d名以上のApproveのあるPRが見つかりませんでした\n", u.runtimeOptions.RequiredApprovals)
	}

	return nil
}

func (u *UseCase) measureApprovalTime(ctx context.Context, summary *mygithub.PullRequestSummary) (string, error) {
	if !summary.IsCreatedWithin(u.runtimeOptions.CreatedAtFrom, u.runtimeOptions.CreatedAtTo) {
		return "", nil
	}

	if summary.ContainsAnyLabel(u.runtimeOptions.IgnoreLabels) {
		return "", nil
	}

	detail, err := u.githubGateway.FindPullRequestDetail(ctx, u.runtimeOptions.Owner, u.runtimeOptions.Name, summary)
	if err != nil {
		return "", fmt.Errorf("find pull request detail (PR #%d): %w", summary.Number, err)
	}

	duration, ok := detail.TimeToNthApproval(u.runtimeOptions.RequiredApprovals)
	if !ok {
		return "", nil
	}

	return fmt.Sprintf("%s %.2f時間 +%d/-%d", detail.HTMLURL, duration.Hours(), detail.Additions, detail.Deletions), nil
}
