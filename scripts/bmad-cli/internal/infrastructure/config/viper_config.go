package config

import (
	"github.com/spf13/viper"
)

type ViperConfig struct {
	viper *viper.Viper
}

func NewViperConfig() *ViperConfig {
	v := viper.New()

	v.SetConfigFile("./scripts/bmad-cli/bmad-cli.yml")
	v.SetConfigType("yaml")
	v.ReadInConfig()

	v.SetDefault("engine.type", "claude")
	v.SetDefault("templates.prompts.heuristic", "./templates/heuristic.prompt.tpl")
	v.SetDefault("templates.prompts.apply", "./templates/apply.prompt.tpl")

	return &ViperConfig{viper: v}
}

func (c *ViperConfig) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *ViperConfig) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *ViperConfig) SetDefault(key string, value interface{}) {
	c.viper.SetDefault(key, value)
}
