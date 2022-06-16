package relation_chain_test

import (
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/relation_chain"
	"github.com/RuiChen0101/unfiyql/pkg/element"
	"github.com/stretchr/testify/assert"
)

func TestFindLCA(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
	}
	rc, _ := relation_chain.BuildRelationChain(&el)

	lca1, err := rc.FindLowestCommonParent("tableC", "tableA")
	assert.Nil(t, err)
	assert.Equal(t, "tableA", lca1)

	lca2, err := rc.FindLowestCommonParent("tableC", "tableB")
	assert.Nil(t, err)
	assert.Equal(t, "tableB", lca2)

	lca3, err := rc.FindLowestCommonParent("tableC", "tableD")
	assert.Nil(t, err)
	assert.Equal(t, "tableA", lca3)
}

func TestFindRelationPath(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
	}
	rc, _ := relation_chain.BuildRelationChain(&el)

	path1 := rc.FindRelationPath("tableA", "tableC")
	expect1 := []relation_chain.RelationChainNode{
		{
			FromField: "fieldA2",
			FromTable: "tableA",
			ToField:   "fieldB2",
			ToTable:   "tableB",
		},
		{
			FromField: "fieldB1",
			FromTable: "tableB",
			ToField:   "fieldC",
			ToTable:   "tableC",
		},
	}
	assert.EqualValues(t, expect1, *path1)

	path2 := rc.FindRelationPath("tableD", "tableA")
	expect2 := []relation_chain.RelationChainNode{
		{
			FromField: "fieldD",
			FromTable: "tableD",
			ToField:   "fieldA1",
			ToTable:   "tableA",
		},
	}
	assert.EqualValues(t, expect2, *path2)

	path3 := rc.FindRelationPath("tableB", "tableD")
	assert.Nil(t, path3)
}

func TestIsParentOfAndIsDescendantOf(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
	}
	rc, _ := relation_chain.BuildRelationChain(&el)

	assert.True(t, rc.IsParentOf("tableA", "tableC"))
	assert.False(t, rc.IsParentOf("tableC", "tableC"))
	assert.True(t, rc.IsDescendantOf("tableC", "tableB"))
	assert.False(t, rc.IsDescendantOf("tableA", "tableD"))
}
