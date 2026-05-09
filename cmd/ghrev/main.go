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
	"github.com/kageyamountain/ghrev/internal/common/log"
	"github.com/kageyamountain/ghrev/internal/common/subcommand"
	"github.com/kageyamountain/ghrev/internal/feature/help"
	"github.com/kageyamountain/ghrev/internal/feature/twoapprove"
	"github.com/kageyamountain/ghrev/internal/feature/version"
	"github.com/kageyamountain/ghrev/internal/infrastructure/gateway/mygithub"
	"golang.org/x/oauth2"

	"github.com/kageyamountain/ghrev/internal/common/config"
)

// ビルド時に -ldflags で書き換える
var applicationVersion = "dev version"

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

	// LogContextの設定
	executionID := uuid.New().String()
	logContext := &sync.Map{}
	logContext.Store("log_type", log.LogTypeApp)
	logContext.Store("execution_id", executionID)
	logContext.Store("args", os.Args)
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

	// 環境変数を構造体へマッピング
	appConfig, err := config.Load()
	if err != nil {
		slog.ErrorContext(ctx, "failed to initialize app config", slog.Any("error", err))
		return 1
	}

	// 最小限のバリデーション
	if len(os.Args) < 2 {
		slog.ErrorContext(ctx, "subcommand is required")
		return 1
	}

	// サブコマンド名を取得
	subCommandName, err := subcommand.ParseName(os.Args[1])
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse subcommand", slog.Any("error", err))
		return 1
	}

	// サブコマンドに対応するユースケースをビルド
	optionArgs := os.Args[2:]
	useCase, err := buildUseCase(ctx, appConfig, subCommandName, optionArgs)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get use case", slog.Any("error", err))
		return 1
	}

	// ユースケース実行
	err = useCase.Do(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to do use case", slog.Any("error", err))
		return 1
	}

	slog.InfoContext(ctx, "completed successfully")
	return 0
}

type useCase interface {
	Do(ctx context.Context) error
}

func buildUseCase(ctx context.Context, appConfig *config.AppConfig, subCommandName subcommand.Name, optionArgs []string) (useCase, error) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: appConfig.GitHub.PersonalAccessToken},
	)
	httpClient := oauth2.NewClient(ctx, tokenSource)
	githubClient := github.NewClient(httpClient)
	githubGateway := mygithub.NewGateway(appConfig, githubClient)

	//exhaustive:enforce
	switch subCommandName {
	case subcommand.Help:
		return help.NewUseCase(), nil
	case subcommand.Version:
		return version.NewUseCase(applicationVersion), nil
	case subcommand.TwoApprove:
		runtimeOptions, err := twoapprove.NewRuntimeOptions(optionArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize runtime options. err: %w", err)
		}
		return twoapprove.NewUseCase(runtimeOptions, appConfig, githubGateway), nil
	}

	return nil, fmt.Errorf("invalid subcommand. subCommandName: %s", subCommandName)
}
