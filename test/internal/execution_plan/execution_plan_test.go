package execution_plan_test

import (
	"path/filepath"
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/execution_plan"
	"github.com/RuiChen0101/unfiyql/internal/expression_tree"
	"github.com/RuiChen0101/unfiyql/internal/relation_chain"
	"github.com/RuiChen0101/unfiyql/internal/relation_linking"
	"github.com/RuiChen0101/unfiyql/internal/service_lookup"
	"github.com/RuiChen0101/unfiyql/pkg/element"
	"github.com/RuiChen0101/unfiyql/pkg/service_config"
	"github.com/RuiChen0101/unfiyql/test/fake"
	"github.com/stretchr/testify/assert"
)

func TestGenerateWithoutCondition(t *testing.T) {
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

	plan, err := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	assert.Nil(t, err)

	assert.Equal(t, element.UnifyQLOperation.Query, plan.Operation)
	assert.Equal(t, "tableA", plan.Query)
	assert.EqualValues(t, []string{"tableA.fieldA3 DESC"}, plan.OrderBy)
	assert.EqualValues(t, []int{0, 100}, plan.Limit)
}

func TestGenerateSpecialOperation(t *testing.T) {
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

	plan, err := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)

	assert.Nil(t, err)

	assert.Equal(t, element.UnifyQLOperation.Sum, plan.Operation)
	assert.Equal(t, "tableA.fieldA", plan.Query)
}

func TestGenerateSingleConditionInSameService(t *testing.T) {
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

	plan, err := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)
	assert.Nil(t, err)

	assert.Equal(t, element.UnifyQLOperation.Query, plan.Operation)
	assert.Equal(t, "tableA", plan.Query)
	assert.Equal(t, "(tableC.fieldC1 & 2) != 0", plan.Where)
	assert.EqualValues(t, []string{"tableC.fieldC=tableB.fieldB1", "tableB.fieldB2=tableA.fieldA2"}, plan.Link)
	assert.EqualValues(t, []string{"tableC", "tableB"}, plan.With)
}

func TestGenerateSingleConditionFromDifferentService(t *testing.T) {
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

	plan, err := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)
	assert.Nil(t, err)

	depPlan := plan.Dependency["12345678"]
	assert.Equal(t, "tableA", plan.Query)
	assert.Equal(t, "tableA.fieldA1 IN {12345678}", plan.Where)
	assert.Equal(t, "tableD.fieldD", depPlan.Query)
	assert.Equal(t, "tableD.fieldD1 != 0", depPlan.Where)
}

func TestGenerateMultipleCondition(t *testing.T) {
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

	plan, err := execution_plan.GenerateExecutionPlan(linkedTree, &lookup, fakeId)
	assert.Nil(t, err)

	depPlan := plan.Dependency["12345678"]

	assert.Equal(t, "tableA", plan.Query)
	assert.Equal(t, "(tableB.fieldB = 0 AND tableA.fieldA1 IN {12345678})", plan.Where)
	assert.EqualValues(t, []string{"tableB.fieldB2=tableA.fieldA2"}, plan.Link)
	assert.EqualValues(t, []string{"tableB"}, plan.With)
	assert.Equal(t, "tableD.fieldD", depPlan.Query)
	assert.Equal(t, "tableD.fieldD1 = 1", depPlan.Where)
}
