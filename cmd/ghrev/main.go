package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	_ "time/tzdata"

	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/google/go-github/v80/github"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/kageyamountain/ghrev/internal/common/logger"
	"github.com/kageyamountain/ghrev/internal/infrastructure/gateway/mygithub"
	"github.com/kageyamountain/ghrev/internal/subcommand"
	"github.com/kageyamountain/ghrev/internal/subcommand/approval"
	"github.com/kageyamountain/ghrev/internal/subcommand/firstreview"
	"github.com/kageyamountain/ghrev/internal/subcommand/help"
	"github.com/kageyamountain/ghrev/internal/subcommand/version"
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
	logContext := logger.NewLogContext()
	logContext.Set("log_type", logger.LogTypeApp)
	logContext.Set("execution_id", executionID)
	logContext.Set("args", os.Args)
	ctx = logger.WithLogContext(ctx, logContext)

	// logger設定
	customLogHandler := logger.NewCustomLogHandler(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelError,
			},
		),
	)
	slog.SetDefault(slog.New(customLogHandler))
	slog.InfoContext(ctx, "process started")

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
	useCase, err := buildUseCase(ctx, subCommandName, optionArgs)
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

func buildUseCase(ctx context.Context, subCommandName subcommand.Name, optionArgs []string) (useCase, error) {
	ghAuthToken, source := auth.TokenForHost("github.com")
	if ghAuthToken == "" {
		return nil, errors.New("failed to get github token. `gh auth login` is required")
	}
	slog.InfoContext(ctx, "github token resolved", slog.Any("source", source))

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghAuthToken},
	)
	httpClient := oauth2.NewClient(ctx, tokenSource)
	githubClient := github.NewClient(httpClient)
	githubGateway := mygithub.NewGateway(githubClient)

	//exhaustive:enforce
	switch subCommandName {
	case subcommand.Help:
		return help.NewUseCase(), nil
	case subcommand.Version:
		return version.NewUseCase(applicationVersion), nil
	case subcommand.Approval:
		runtimeOptions, err := approval.NewRuntimeOptions(optionArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize runtime options. err: %w", err)
		}
		return approval.NewUseCase(runtimeOptions, githubGateway), nil
	case subcommand.FirstReview:
		runtimeOptions, err := firstreview.NewRuntimeOptions(optionArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize runtime options. err: %w", err)
		}
		return firstreview.NewUseCase(runtimeOptions, githubGateway), nil
	}

	return nil, fmt.Errorf("invalid subcommand. subCommandName: %s", subCommandName)
}
