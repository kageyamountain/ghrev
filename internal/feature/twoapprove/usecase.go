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

	pullRequests, err := u.githubGateway.FindAllPullRequests(ctx, owner, name)
	if err != nil {
		return fmt.Errorf("failed to find pull requests: %w", err)
	}

	ignoreLabels := u.runtimeOptions.IgnoreLabels.Strings()
	createdAtFrom := u.runtimeOptions.CreatedAtFrom.Time()
	createdAtTo := u.runtimeOptions.CreatedAtTo.Time()

	targetCount := 0
	for _, pullRequest := range pullRequests {
		// TODO errgroupで並列化
		if !pullRequest.IsCreatedWithin(createdAtFrom, createdAtTo) {
			continue
		}

		if pullRequest.ContainsAnyLabel(ignoreLabels) {
			continue
		}

		openedAt, err2 := u.githubGateway.FindPullRequestFirstOpenedAt(ctx, owner, name, pullRequest.Number)
		if err2 != nil {
			slog.ErrorContext(ctx, "failed to get pull request open time", slog.Any("error", err2), slog.Any("pullRequest", pullRequest))
			continue
		}
		pullRequest.OpenedAt = openedAt

		reviews, err2 := u.githubGateway.FindPullRequestReviews(ctx, owner, name, pullRequest.Number)
		if err2 != nil {
			slog.ErrorContext(ctx, "failed to get pull request reviews", slog.Any("error", err2), slog.Any("pullRequest", pullRequest))
			continue
		}
		pullRequest.Reviews = reviews

		duration, ok := pullRequest.TimeToSecondApproval()
		if !ok {
			continue
		}
		targetCount++

		fmt.Printf("%s %.2f時間\n", pullRequest.HTMLURL, duration.Hours())
	}

	if targetCount == 0 {
		fmt.Println("2名以上のApproveのあるPRが見つかりませんでした")
		return nil
	}

	return nil
}
