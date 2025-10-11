package config

import (
	"log/slog"

	"bmad-cli/internal/pkg/errors"
	"github.com/spf13/viper"
)

type ViperConfig struct {
	viper *viper.Viper
}

func NewViperConfig() (*ViperConfig, error) {
	v := viper.New()

	v.SetConfigFile("./bmad-cli.yml")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		slog.Error("Failed to read config file", "error", err)

		return nil, errors.ErrReadConfigFailed(err)
	}

	return &ViperConfig{viper: v}, nil
}

func (c *ViperConfig) GetString(key string) string {
	if !c.viper.IsSet(key) {
		return ""
	}

	return c.viper.GetString(key)
}
