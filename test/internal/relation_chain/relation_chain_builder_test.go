package relation_chain_test

import (
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/element"
	"github.com/RuiChen0101/unfiyql/internal/relation_chain"
	"github.com/stretchr/testify/assert"
)

func TestBuildRelationChain(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
	}
	rc, err := relation_chain.BuildRelationChain(el)
	assert.Nil(t, err)

	forwardMap := rc.GetForwardRelationMap()
	assert.EqualValues(t, relation_chain.RelationChainNode{
		FromField: "fieldA2",
		FromTable: "tableA",
		ToField:   "fieldB2",
		ToTable:   "tableB",
	}, forwardMap["tableA"]["tableB"])
	assert.EqualValues(t, relation_chain.RelationChainNode{
		FromField: "fieldA1",
		FromTable: "tableA",
		ToField:   "fieldD",
		ToTable:   "tableD",
	}, forwardMap["tableA"]["tableD"])
	assert.EqualValues(t, relation_chain.RelationChainNode{
		FromField: "fieldB1",
		FromTable: "tableB",
		ToField:   "fieldC",
		ToTable:   "tableC",
	}, forwardMap["tableB"]["tableC"])

	backwardMap := rc.GetBackwardRelationMap()
	assert.EqualValues(t, relation_chain.RelationChainNode{
		FromField: "fieldD",
		FromTable: "tableD",
		ToField:   "fieldA1",
		ToTable:   "tableA",
	}, backwardMap["tableD"]["tableA"])
	assert.EqualValues(t, relation_chain.RelationChainNode{
		FromField: "fieldC",
		FromTable: "tableC",
		ToField:   "fieldB1",
		ToTable:   "tableB",
	}, backwardMap["tableC"]["tableB"])
	assert.EqualValues(t, relation_chain.RelationChainNode{
		FromField: "fieldB2",
		FromTable: "tableB",
		ToField:   "fieldA2",
		ToTable:   "tableA",
	}, backwardMap["tableB"]["tableA"])

}

func TestBuildRelationChainWithoutWithAndLink(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
	}
	rc, err := relation_chain.BuildRelationChain(el)
	assert.Nil(t, err)

	emptyMap := map[string]map[string]relation_chain.RelationChainNode{}
	assert.EqualValues(t, emptyMap, rc.GetForwardRelationMap())
	assert.EqualValues(t, emptyMap, rc.GetBackwardRelationMap())
}

func TestInvalidFormatError(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldCtableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
	}
	rc, err := relation_chain.BuildRelationChain(el)
	assert.EqualError(t, err, "RelationChain: tableC.fieldCtableB.fieldB1 invalid format")
	assert.Nil(t, rc)
}

func TestUndefinedTableError(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
	}
	rc, err := relation_chain.BuildRelationChain(el)
	assert.EqualError(t, err, "RelationChain: tableD.fieldD=tableA.fieldA1 using undefined table")
	assert.Nil(t, rc)
}
