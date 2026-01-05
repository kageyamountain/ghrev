package twoapprove

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/kageyamountain/ghrev/internal/common/config"
	"github.com/kageyamountain/ghrev/internal/common/runtimeoption"
	"github.com/kageyamountain/ghrev/internal/infrastructure/gateway/mygithub"
)

type UseCase struct {
	runtimeOptions *runtimeoption.Options
	appConfig      *config.AppConfig
	githubGateway  mygithub.Gateway
}

func NewUseCase(
	runtimeOptions *runtimeoption.Options,
	appConfig *config.AppConfig,
	githubGateway mygithub.Gateway,
) *UseCase {
	return &UseCase{
		runtimeOptions: runtimeOptions,
		appConfig:      appConfig,
		githubGateway:  githubGateway,
	}
}

func (u *UseCase) Do(ctx context.Context) error {
	pullRequests, err := u.githubGateway.FindPullRequests(ctx, u.runtimeOptions.RepositoryOwner.String(), u.runtimeOptions.RepositoryName.String())
	if err != nil {
		return fmt.Errorf("failed to find pull requests: %w", err)
	}

	targetCount := 0
	for _, pullRequest := range pullRequests {
		// TODO errgroupで並列化
		// PR作成日時が対象期間外の場合はスキップ
		createdAt := *pullRequest.CreatedAt.GetTime()
		if createdAt.Before(u.runtimeOptions.CreatedAtFrom.Time()) || createdAt.After(u.runtimeOptions.CreatedAtTo.Time()) {
			continue
		}

		// 特定のラベルが付与されている場合はスキップ
		ignoreLabels := u.runtimeOptions.IgnoreLabels
		shouldSkip := false
		if len(ignoreLabels) > 0 {
			for _, label := range pullRequest.Labels {
				if slices.Contains(ignoreLabels, label.GetName()) {
					shouldSkip = true
					break
				}
			}
		}
		if shouldSkip {
			fmt.Println("skip: findy計測除外")
			continue
		}

		// PRの初回Open日時を取得
		openedAt, err2 := u.githubGateway.FindPullRequestFirstOpenTime(ctx, u.runtimeOptions.RepositoryOwner.String(), u.runtimeOptions.RepositoryName.String(), pullRequest.GetNumber())
		if err2 != nil {
			slog.ErrorContext(ctx, "failed to get pull request open time", slog.Any("error", err2), slog.Any("pullRequest", pullRequest))
			continue
		}
		// 作成時からOpenしている場合はPR作成日時と同値をセット
		if openedAt == nil {
			openedAt = &createdAt
		}

		// approveを取得
		approveTimes, err2 := u.githubGateway.FindPullRequestApproveTimes(ctx, u.runtimeOptions.RepositoryOwner.String(), u.runtimeOptions.RepositoryName.String(), pullRequest.GetNumber())
		if err2 != nil {
			slog.ErrorContext(ctx, "failed to get pull request approve times", slog.Any("error", err2), slog.Any("pullRequest", pullRequest))
			continue
		}

		// approveが2名未満の場合はスキップ
		if len(approveTimes) < 2 {
			continue
		}
		targetCount++

		// PRの初回Open日時から2人目のapprove日時までの経過時間を算出してコンソール出力
		timeToSecondApproval := approveTimes[1].Sub(*openedAt)
		fmt.Printf("%s %.2f時間\n", *pullRequest.HTMLURL, timeToSecondApproval.Hours())
	}

	if targetCount == 0 {
		fmt.Println("2名以上のApproveのあるPRが見つかりませんでした")
		return nil
	}

	return nil
}
