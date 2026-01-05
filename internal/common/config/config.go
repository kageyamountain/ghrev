package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	GitHub GitHubConfig `envconfig:"GITHUB"`
}

type GitHubConfig struct {
	PersonalAccessToken string `envconfig:"PERSONAL_ACCESS_TOKEN" required:"true"`
}

// Load APPプレフィックスを持つ環境変数を構造体マッピングする
func Load() (*AppConfig, error) {
	var appConfig AppConfig
	err := envconfig.Process("APP", &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
