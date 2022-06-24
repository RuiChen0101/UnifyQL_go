package service_lookup

import "github.com/RuiChen0101/UnifyQL_go/pkg/service_config"

type ServiceLookup struct {
	source service_config.ServiceConfigSource
}

func NewServiceLookup(source service_config.ServiceConfigSource) ServiceLookup {
	return ServiceLookup{
		source: source,
	}
}

func (sl *ServiceLookup) GetServiceConfig(serviceName string) service_config.ServiceConfig {
	return sl.source.GetServiceConfigs()[serviceName]
}

func (sl *ServiceLookup) GetServiceNameByTable(table string) string {
	return sl.source.GetTableMapping()[table]
}

func (sl *ServiceLookup) IsAllFromSameService(tables []string) bool {
	if len(tables) == 0 {
		return false
	}
	refService := ""
	for _, t := range tables {
		if refService == "" {
			refService = sl.source.GetTableMapping()[t]
		} else if val, exist := sl.source.GetTableMapping()[t]; !exist || refService != val {
			return false
		}
	}
	return refService != ""
}
