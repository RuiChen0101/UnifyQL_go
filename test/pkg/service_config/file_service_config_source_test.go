package service_config_test

import (
	"path/filepath"
	"testing"

	"github.com/RuiChen0101/UnifyQL_go/pkg/service_config"
	"github.com/stretchr/testify/assert"
)

func TestLoadServiceConfigFromFile(t *testing.T) {
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, err := service_config.NewFileServiceConfigSource(path)
	expectedServiceConfigs := map[string]service_config.ServiceConfig{
		"aService": {
			ServiceName: "aService",
			Url:         "http://localhost:5000/query",
			Tables: []string{
				"tableA",
				"tableB",
				"tableC",
			},
		},
		"bService": {
			ServiceName: "bService",
			Url:         "http://localhost:4999/query",
			Tables: []string{
				"tableD",
				"tableE",
				"tableF",
			},
		},
	}
	expectedTableMapping := map[string]string{
		"tableA": "aService",
		"tableB": "aService",
		"tableC": "aService",
		"tableD": "bService",
		"tableE": "bService",
		"tableF": "bService",
	}
	assert.Nil(t, err)
	assert.EqualValues(t, expectedServiceConfigs, conf.GetServiceConfigs())
	assert.EqualValues(t, expectedTableMapping, conf.GetTableMapping())
}
