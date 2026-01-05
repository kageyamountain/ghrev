package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	_ "time/tzdata"

	"github.com/google/go-github/v80/github"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/kageyamountain/ghrev/internal/infrastructure/gateway/mygithub"

	"github.com/kageyamountain/ghrev/internal/common/runtimeoption"
	"github.com/kageyamountain/ghrev/internal/feature/twoapprove"

	"github.com/kageyamountain/ghrev/internal/common/config"
	"github.com/kageyamountain/ghrev/internal/common/log"
)

func main() {
	ctx := context.Background()

	// LogContextの設定
	executionID := uuid.New().String()
	logContext := &sync.Map{}
	logContext.Store("log_type", log.LogTypeApp)
	logContext.Store("execution_id", executionID)
	ctx = context.WithValue(ctx, log.LogContextKey, logContext)

	// logger設定
	customLogHandler := log.NewCustomLogHandler(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)
	slog.SetDefault(slog.New(customLogHandler))
	slog.InfoContext(ctx, "process started")

	// 実行時引数を取得
	runtimeOptions, err := runtimeoption.NewOptions()
	if err != nil {
		slog.ErrorContext(ctx, "failed to initialize runtime options", slog.Any("error", err))
		os.Exit(1)
	}
	slog.InfoContext(ctx, "runtime options", slog.Any("options", runtimeOptions))

	// 環境変数をAppConfigへマッピング
	appConfig, err := config.Load()
	if err != nil {
		slog.ErrorContext(ctx, "failed to initialize app config", slog.Any("error", err))
		os.Exit(1)
	}

	// modeに応じたユースケースを選択
	useCase, err := selectUseCase(ctx, runtimeOptions, appConfig)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get use case", slog.Any("error", err))
		os.Exit(1)
	}

	// ユースケース実行
	err = useCase.Do(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to do use case", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "completed successfully")
	os.Exit(0)
}

type useCase interface {
	Do(ctx context.Context) error
}

func selectUseCase(ctx context.Context, runtimeOptions *runtimeoption.Options, appConfig *config.AppConfig) (useCase, error) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: appConfig.GitHub.PersonalAccessToken},
	)
	httpClient := oauth2.NewClient(ctx, tokenSource)
	githubClient := github.NewClient(httpClient)
	githubGateway := mygithub.NewGateway(appConfig, githubClient)

	//exhaustive:enforce
	//nolint:gocritic
	switch runtimeOptions.Mode {
	case runtimeoption.ModeTwoApprove:
		return twoapprove.NewUseCase(runtimeOptions, appConfig, githubGateway), nil
	}

	return nil, fmt.Errorf("invalid mode. mode: %s", runtimeOptions.Mode)
}
