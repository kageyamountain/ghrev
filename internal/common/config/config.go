package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	GitHub GitHubConfig `envconfig:"GITHUB"`
}

type GitHubConfig struct {
	PersonalAccessToken string `envconfig:"GITHUB_PERSONAL_ACCESS_TOKEN" required:"true"`
}

// Load 環境変数を構造体にマッピングする
func Load() (*AppConfig, error) {
	var appConfig AppConfig
	err := envconfig.Process("", &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
