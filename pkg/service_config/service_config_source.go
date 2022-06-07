package service_config

type ServiceConfigSource interface {
	GetServiceConfigs() map[string]ServiceConfig
	GetTableMapping() map[string]string
}
