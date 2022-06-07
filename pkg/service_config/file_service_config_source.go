package service_config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type FileServiceConfigSource struct {
	serviceConfig map[string]ServiceConfig
	tableMapping  map[string]string
}

func NewFileServiceConfigSource(fileName string) (*FileServiceConfigSource, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	configs := []ServiceConfig{}
	json.Unmarshal(byteValue, &configs)

	result := FileServiceConfigSource{
		serviceConfig: map[string]ServiceConfig{},
		tableMapping:  map[string]string{},
	}
	for _, c := range configs {
		serviceName := c.ServiceName
		result.serviceConfig[serviceName] = c
		for _, t := range c.Tables {
			result.tableMapping[t] = serviceName
		}
	}
	return &result, nil
}

func (fsc *FileServiceConfigSource) GetServiceConfigs() map[string]ServiceConfig {
	return fsc.serviceConfig
}

func (fsc *FileServiceConfigSource) GetTableMapping() map[string]string {
	return fsc.tableMapping
}
