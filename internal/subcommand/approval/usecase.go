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
	owner := u.runtimeOptions.RepositoryOwner.String()
	name := u.runtimeOptions.RepositoryName.String()
	createdAtFrom := u.runtimeOptions.CreatedAtFrom.Time()
	createdAtTo := u.runtimeOptions.CreatedAtTo.Time()
	ignoreLabels := u.runtimeOptions.IgnoreLabels.Strings()
	requiredApprovals := u.runtimeOptions.RequiredApprovals.Int()

	summaries, err := u.githubGateway.FindAllPullRequestSummaries(ctx, owner, name)
	if err != nil {
		return fmt.Errorf("failed to find pull request summaries: %w", err)
	}

	var resultRows []string
	for _, summary := range summaries {
		// TODO errgroupで並列化
		if !summary.IsCreatedWithin(createdAtFrom, createdAtTo) {
			continue
		}

		if summary.ContainsAnyLabel(ignoreLabels) {
			continue
		}

		detail, err2 := u.githubGateway.FindPullRequestDetail(ctx, owner, name, summary)
		if err2 != nil {
			slog.ErrorContext(ctx, "failed to find pull request detail", slog.Any("error", err2), slog.Any("pullRequestSummary", summary))
			continue
		}

		duration, ok := detail.TimeToNthApproval(requiredApprovals)
		if !ok {
			continue
		}

		resultRows = append(resultRows, fmt.Sprintf("%s %.2f時間 +%d/-%d", detail.HTMLURL, duration.Hours(), detail.Additions, detail.Deletions))
	}

	if len(resultRows) == 0 {
		fmt.Printf("%d名以上のApproveのあるPRが見つかりませんでした\n", requiredApprovals)
		return nil
	}

	fmt.Println("URL 所要時間 変更行数")
	for _, resultRow := range resultRows {
		fmt.Println(resultRow)
	}

	return nil
}
