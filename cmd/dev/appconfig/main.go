package main

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/kageyamountain/ghrev/internal/common/config"
)

func main() {
	// dev環境変数のロード
	err := godotenv.Load(".env.dev")
	if err != nil {
		slog.Error("failed to load .env file.", slog.Any("err", err))
		return
	}

	// 環境変数をAppConfigへマッピング
	appConfig, err := config.Load()
	if err != nil {
		slog.Error("failed to AppConfig Load.", slog.Any("err", err))
		return
	}

	// 出力
	fmt.Printf("%+v", appConfig)
}
