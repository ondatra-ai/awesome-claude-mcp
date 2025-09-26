package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ViperConfig struct {
	viper *viper.Viper
}

func NewViperConfig() (*ViperConfig, error) {
	v := viper.New()

	v.SetConfigFile("bmad-cli.yml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return &ViperConfig{viper: v}, nil
}

func (c *ViperConfig) GetString(key string) string {
	if !c.viper.IsSet(key) {
		return ""
	}
	return c.viper.GetString(key)
}

func (c *ViperConfig) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *ViperConfig) SetDefault(key string, value interface{}) {
	c.viper.SetDefault(key, value)
}
