package service_config

type ServiceConfig struct {
	ServiceName string   `json:"serviceName"`
	Url         string   `json:"url"`
	Tables      []string `json:"tables"`
}
