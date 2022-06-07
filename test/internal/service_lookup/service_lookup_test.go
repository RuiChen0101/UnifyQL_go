package service_lookup_test

import (
	"path/filepath"
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/service_lookup"
	"github.com/RuiChen0101/unfiyql/pkg/service_config"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceConfig(t *testing.T) {
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)
	config := service_config.ServiceConfig{
		ServiceName: "aService",
		Url:         "http://localhost:5000/query",
		Tables: []string{
			"tableA",
			"tableB",
			"tableC",
		},
	}
	assert.EqualValues(t, config, lookup.GetServiceConfig("aService"))
}

func TestGetNameByTable(t *testing.T) {
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)
	assert.EqualValues(t, "aService", lookup.GetServiceNameByTable("tableA"))
}

func TestSameService(t *testing.T) {
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)
	assert.True(t, lookup.IsAllFromSameService([]string{"tableA", "tableB"}))
	assert.True(t, lookup.IsAllFromSameService([]string{"tableD", "tableE"}))
}

func TestDifferentService(t *testing.T) {
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)
	assert.False(t, lookup.IsAllFromSameService([]string{"tableA", "tableF"}))
	assert.False(t, lookup.IsAllFromSameService([]string{"tableD", "tableC"}))
}
