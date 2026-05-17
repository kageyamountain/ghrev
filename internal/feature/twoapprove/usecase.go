package twoapprove

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

	summaries, err := u.githubGateway.FindAllPullRequestSummaries(ctx, owner, name)
	if err != nil {
		return fmt.Errorf("failed to find pull request summaries: %w", err)
	}

	targetCount := 0
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

		duration, ok := detail.TimeToSecondApproval()
		if !ok {
			continue
		}
		targetCount++

		fmt.Printf("%s %.2f時間\n", detail.HTMLURL, duration.Hours())
	}

	if targetCount == 0 {
		fmt.Println("2名以上のApproveのあるPRが見つかりませんでした")
		return nil
	}

	return nil
}
