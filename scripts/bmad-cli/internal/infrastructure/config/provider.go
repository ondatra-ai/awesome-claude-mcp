package config

type ConfigProvider interface {
	GetString(key string) string
	GetInt(key string) int
	SetDefault(key string, value interface{})
}
