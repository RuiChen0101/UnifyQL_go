package unfiyql_test

import (
	"path/filepath"
	"testing"

	"github.com/RuiChen0101/UnifyQL_go/pkg/cache"
	"github.com/RuiChen0101/UnifyQL_go/pkg/service_config"
	"github.com/RuiChen0101/UnifyQL_go/pkg/unifyql"
	"github.com/RuiChen0101/UnifyQL_go/test/fake"
	"github.com/stretchr/testify/assert"
)

func TestQueryWithoutCondition(t *testing.T) {
	query := "QUERY tableA ORDER BY tableA.fieldA3 DESC LIMIT 0,100"

	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"fieldA\": \"fieldA\", \"fieldA1\": \"fieldA1\", \"fieldA2\": \"fieldA2\" }]"})

	uql := unifyql.NewUnifyQl(conf, fetchProxy, nil)
	result, err := uql.Query(query)

	assert.Nil(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"fieldA": "fieldA", "fieldA1": "fieldA1", "fieldA2": "fieldA2",
	}, result[0])

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "QUERY tableA ORDER BY tableA.fieldA3 DESC LIMIT 0,100", fetchProxy.GetRecord(0).UqlPayload)
}

func TestQueryWithCacheManager(t *testing.T) {
	query := "QUERY tableA"

	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	cache := cache.NewDefaultExecutionPlanCache()
	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"fieldA\": \"fieldA\", \"fieldA1\": \"fieldA1\", \"fieldA2\": \"fieldA2\" }]"})

	uql := unifyql.NewUnifyQl(conf, fetchProxy, cache)
	result, err := uql.Query(query)

	assert.Nil(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"fieldA": "fieldA", "fieldA1": "fieldA1", "fieldA2": "fieldA2",
	}, result[0])

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "QUERY tableA", fetchProxy.GetRecord(0).UqlPayload)

	plan, ok := cache.Get("26e163cbc7dc6bb34f615c95b60676ae767a55bbc1f9afe997e56ebd90efb7c7")
	assert.True(t, ok)
	assert.Equal(t, "tableA", plan.Query)
}

func TestCountQuery(t *testing.T) {
	query := "COUNT tableA"

	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"count\": 10}]"})

	uql := unifyql.NewUnifyQl(conf, fetchProxy, nil)
	result, err := uql.Query(query)

	assert.Nil(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"count": 10.0,
	}, result[0])

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "COUNT tableA", fetchProxy.GetRecord(0).UqlPayload)
}

func TestSumQuery(t *testing.T) {
	query := "SUM tableA.fieldA"

	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"sum\": 10 }]"})

	uql := unifyql.NewUnifyQl(conf, fetchProxy, nil)
	result, err := uql.Query(query)

	assert.Nil(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"sum": 10.0,
	}, result[0])

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "SUM tableA.fieldA", fetchProxy.GetRecord(0).UqlPayload)
}

func TestComplexQuery(t *testing.T) {
	query := "QUERY tableA WITH tableB, tableC, tableD LINK tableC.fieldC=tableB.fieldB1,tableD.fieldD=tableA.fieldA1,tableA.fieldA2=tableB.fieldB2 WHERE tableD.fieldD1 = 0 AND tableC.fieldC1 = 2 AND (tableD.fieldD2 = 1 OR tableB.fieldB = 3) ORDER BY tableA.tableA3 ASC LIMIT 10, 100"

	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	fetchProxy := fake.NewFakeFetchProxy([]string{
		"[{ \"fieldD\":5 }, { \"fieldD\":6 }, { \"fieldD\":7 }, { \"fieldD\":8 }]",
		"[{ \"fieldD\":1 }, { \"fieldD\":2 }, { \"fieldD\":3 }, { \"fieldD\":4 }]",
		"[{ \"fieldA\": \"fieldA\", \"fieldA1\": \"fieldA1\", \"fieldA2\": \"fieldA2\" }]",
	})

	uql := unifyql.NewUnifyQl(conf, fetchProxy, nil)
	result, err := uql.Query(query)

	assert.Nil(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"fieldA": "fieldA", "fieldA1": "fieldA1", "fieldA2": "fieldA2",
	}, result[0])

	assert.Equal(t, "http://localhost:4999/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "QUERY tableD.fieldD WHERE tableD.fieldD2 = 1", fetchProxy.GetRecord(0).UqlPayload)

	assert.Equal(t, "http://localhost:4999/query", fetchProxy.GetRecord(1).Url)
	assert.Equal(t, "QUERY tableD.fieldD WHERE tableD.fieldD1 = 0", fetchProxy.GetRecord(1).UqlPayload)

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(2).Url)
	assert.Equal(t, "QUERY tableA WITH tableC,tableB LINK tableC.fieldC=tableB.fieldB1,tableB.fieldB2=tableA.fieldA2 WHERE ((tableA.fieldA1 IN (1,2,3,4) AND tableC.fieldC1 = 2) AND (tableA.fieldA1 IN (5,6,7,8) OR tableB.fieldB = 3)) ORDER BY tableA.tableA3 ASC LIMIT 10,100", fetchProxy.GetRecord(2).UqlPayload)
}

func TestVoidAuthorizationBypass(t *testing.T) {
	query := "QUERY tableA WITH tableB LINK tableA.fieldA2=tableB.fieldB2 WHERE tableB.fieldB = \"valueB\" OR 1=1--\""

	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"sum\": 10 }]"})

	uql := unifyql.NewUnifyQl(conf, fetchProxy, nil)
	result, err := uql.Query(query)

	assert.Nil(t, result)
	assert.Equal(t, "ExpressionTreeBuilder: empty tree", err.Error())
}

func TestVoidMaliciousCommands(t *testing.T) {
	query := "QUERY tableA WITH tableB LINK tableA.fieldA2=tableB.fieldB2 WHERE tableB.fieldB = \"valueB\"; DROP TABLE tableA--\""

	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"sum\": 10 }]"})

	uql := unifyql.NewUnifyQl(conf, fetchProxy, nil)
	result, err := uql.Query(query)

	assert.Nil(t, result)
	assert.Equal(t, "ExpressionTreeBuilder: empty tree", err.Error())
}
