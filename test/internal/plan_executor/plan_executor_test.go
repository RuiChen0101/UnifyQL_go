package plan_executor_test

import (
	"path/filepath"
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/execution_plan"
	"github.com/RuiChen0101/unfiyql/internal/expression_tree"
	"github.com/RuiChen0101/unfiyql/internal/plan_executor"
	"github.com/RuiChen0101/unfiyql/internal/relation_chain"
	"github.com/RuiChen0101/unfiyql/internal/relation_linking"
	"github.com/RuiChen0101/unfiyql/internal/service_lookup"
	"github.com/RuiChen0101/unfiyql/pkg/element"
	"github.com/RuiChen0101/unfiyql/pkg/service_config"
	"github.com/RuiChen0101/unfiyql/test/fake"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithoutCondition(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		OrderBy:     []string{"tableA.fieldA3 DESC"},
		Limit:       []int{0, 100},
	}
	fakeId := &fake.FakeIdGenerator{}
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)

	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)
	linker := relation_linking.NewRelationLinker(rc, tree)
	linker.Link()
	linkedTree := linker.GetExpressionTree()
	plan, _ := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"fieldA\": \"fieldA\", \"fieldA1\": \"fieldA1\", \"fieldA2\": \"fieldA2\" }]"})

	result, err := plan_executor.ExecutePlan("root", plan, lookup, fetchProxy)

	assert.Nil(t, err)
	assert.Equal(t, "root", result.Id)
	assert.EqualValues(t, map[string]interface{}{
		"fieldA": "fieldA", "fieldA1": "fieldA1", "fieldA2": "fieldA2",
	}, result.Data[0])

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "QUERY tableA ORDER BY tableA.fieldA3 DESC LIMIT 0,100", fetchProxy.GetRecord(0).UqlPayload)
}

func TestExecuteSpecialOperation(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Sum,
		QueryTarget: "tableA",
		QueryField:  "fieldA",
	}
	fakeId := &fake.FakeIdGenerator{}
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)

	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)
	linker := relation_linking.NewRelationLinker(rc, tree)
	linker.Link()
	linkedTree := linker.GetExpressionTree()
	plan, _ := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	fetchProxy := fake.NewFakeFetchProxy([]string{"[{\"sum\": 10}]"})

	result, err := plan_executor.ExecutePlan("root", plan, lookup, fetchProxy)

	assert.Nil(t, err)
	assert.Equal(t, "root", result.Id)
	assert.EqualValues(t, map[string]interface{}{
		"sum": 10.0,
	}, result.Data[0])

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "SUM tableA.fieldA", fetchProxy.GetRecord(0).UqlPayload)
}

func TestExecuteSingleConditionInSameService(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "(tableC.fieldC1 & 2) != 0",
	}
	fakeId := &fake.FakeIdGenerator{}
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)

	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)
	linker := relation_linking.NewRelationLinker(rc, tree)
	linker.Link()
	linkedTree := linker.GetExpressionTree()
	plan, _ := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	fetchProxy := fake.NewFakeFetchProxy([]string{"[{ \"fieldA\": \"fieldA\", \"fieldA1\": \"fieldA1\", \"fieldA2\": \"fieldA2\" }]"})

	result, err := plan_executor.ExecutePlan("root", plan, lookup, fetchProxy)

	assert.Nil(t, err)
	assert.Equal(t, "root", result.Id)
	assert.EqualValues(t, map[string]interface{}{
		"fieldA": "fieldA", "fieldA1": "fieldA1", "fieldA2": "fieldA2",
	}, result.Data[0])

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "QUERY tableA WITH tableC,tableB LINK tableC.fieldC=tableB.fieldB1,tableB.fieldB2=tableA.fieldA2 WHERE (tableC.fieldC1 & 2) != 0", fetchProxy.GetRecord(0).UqlPayload)
}

func TestExecuteSingleConditionFromDifferentService(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "tableD.fieldD1 != 0",
	}
	fakeId := &fake.FakeIdGenerator{}
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)

	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)
	linker := relation_linking.NewRelationLinker(rc, tree)
	linker.Link()
	linkedTree := linker.GetExpressionTree()
	plan, _ := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	fetchProxy := fake.NewFakeFetchProxy([]string{
		"[{ \"fieldD\":1 }, { \"fieldD\":2 }, { \"fieldD\":\"three\" }, { \"fieldD\":\"four\" }]",
		"[{ \"fieldA\": \"fieldA\", \"fieldA1\": \"fieldA1\", \"fieldA2\": \"fieldA2\" }]",
	})

	result, err := plan_executor.ExecutePlan("root", plan, lookup, fetchProxy)

	assert.Nil(t, err)
	assert.Equal(t, "root", result.Id)
	assert.EqualValues(t, map[string]interface{}{
		"fieldA": "fieldA", "fieldA1": "fieldA1", "fieldA2": "fieldA2",
	}, result.Data[0])

	assert.Equal(t, "http://localhost:4999/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "QUERY tableD.fieldD WHERE tableD.fieldD1 != 0", fetchProxy.GetRecord(0).UqlPayload)

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(1).Url)
	assert.Equal(t, "QUERY tableA WHERE tableA.fieldA1 IN (1,2,\"three\",\"four\")", fetchProxy.GetRecord(1).UqlPayload)
}

func TestExecuteMultipleCondition(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "tableB.fieldB = 0 AND tableD.fieldD1 = 1",
	}
	fakeId := &fake.FakeIdGenerator{}
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)

	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)
	linker := relation_linking.NewRelationLinker(rc, tree)
	linker.Link()
	linkedTree := linker.GetExpressionTree()
	plan, _ := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	fetchProxy := fake.NewFakeFetchProxy([]string{
		"[{ \"fieldD\":1 }, { \"fieldD\":2 }, { \"fieldD\":\"three\" }, { \"fieldD\":\"four\" }]",
		"[{ \"fieldA\": \"fieldA\", \"fieldA1\": \"fieldA1\", \"fieldA2\": \"fieldA2\" }]",
	})

	result, err := plan_executor.ExecutePlan("root", plan, lookup, fetchProxy)

	assert.Nil(t, err)
	assert.Equal(t, "root", result.Id)
	assert.EqualValues(t, map[string]interface{}{
		"fieldA": "fieldA", "fieldA1": "fieldA1", "fieldA2": "fieldA2",
	}, result.Data[0])

	assert.Equal(t, "http://localhost:4999/query", fetchProxy.GetRecord(0).Url)
	assert.Equal(t, "QUERY tableD.fieldD WHERE tableD.fieldD1 = 1", fetchProxy.GetRecord(0).UqlPayload)

	assert.Equal(t, "http://localhost:5000/query", fetchProxy.GetRecord(1).Url)
	assert.Equal(t, "QUERY tableA WITH tableB LINK tableB.fieldB2=tableA.fieldA2 WHERE (tableB.fieldB = 0 AND tableA.fieldA1 IN (1,2,\"three\",\"four\"))", fetchProxy.GetRecord(1).UqlPayload)
}

func TestSubQueryErrorCascadeReturn(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "tableB.fieldB = 0 AND tableD.fieldD1 = 1",
	}
	fakeId := &fake.FakeIdGenerator{}
	path, _ := filepath.Abs("../../data/serviceConfig.json")
	conf, _ := service_config.NewFileServiceConfigSource(path)
	lookup := service_lookup.NewServiceLookup(conf)

	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)
	linker := relation_linking.NewRelationLinker(rc, tree)
	linker.Link()
	linkedTree := linker.GetExpressionTree()
	plan, _ := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	fetchProxy := fake.NewFakeFetchProxy([]string{
		"[{ \"fieldD\":1 }, { \"fieldD\":2 }, { \"fieldD\":\"three\" }, { \"fieldD\":\"four\" }]",
	})

	result, err := plan_executor.ExecutePlan("root", plan, lookup, fetchProxy)

	assert.Nil(t, result)
	assert.Equal(t, "Too many request", err.Error())
}
