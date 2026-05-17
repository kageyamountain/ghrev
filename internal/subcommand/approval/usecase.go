package approval

import (
	"context"
	"fmt"
	"log/slog"

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

func (u *UseCase) Do(ctx context.Context) error {
	summaries, err := u.githubGateway.FindAllPullRequestSummaries(ctx, u.runtimeOptions.Owner, u.runtimeOptions.Name)
	if err != nil {
		return fmt.Errorf("failed to find pull request summaries: %w", err)
	}

	var resultRows []string
	for _, summary := range summaries {
		// TODO errgroupで並列化
		row := u.measureApprovalTime(ctx, summary)
		if row != "" {
			resultRows = append(resultRows, row)
		}
	}

	if len(resultRows) == 0 {
		fmt.Printf("%d名以上のApproveのあるPRが見つかりませんでした\n", u.runtimeOptions.RequiredApprovals)
		return nil
	}

	fmt.Println("URL 所要時間 変更行数")
	for _, resultRow := range resultRows {
		fmt.Println(resultRow)
	}

	return nil
}

func (u *UseCase) measureApprovalTime(ctx context.Context, summary *mygithub.PullRequestSummary) string {
	if !summary.IsCreatedWithin(u.runtimeOptions.CreatedAtFrom, u.runtimeOptions.CreatedAtTo) {
		return ""
	}

	if summary.ContainsAnyLabel(u.runtimeOptions.IgnoreLabels) {
		return ""
	}

	detail, err := u.githubGateway.FindPullRequestDetail(ctx, u.runtimeOptions.Owner, u.runtimeOptions.Name, summary)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find pull request detail", slog.Any("error", err), slog.Any("pullRequestSummary", summary))
		return ""
	}

	duration, ok := detail.TimeToNthApproval(u.runtimeOptions.RequiredApprovals)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s %.2f時間 +%d/-%d", detail.HTMLURL, duration.Hours(), detail.Additions, detail.Deletions)
}
