package story

type PerformanceRequirements struct {
	ConnectionEstablishment string `yaml:"connection_establishment" json:"connection_establishment"`
	MessageProcessing       string `yaml:"message_processing" json:"message_processing"`
	ConcurrentConnections   string `yaml:"concurrent_connections" json:"concurrent_connections"`
	MemoryUsage             string `yaml:"memory_usage" json:"memory_usage"`
}
