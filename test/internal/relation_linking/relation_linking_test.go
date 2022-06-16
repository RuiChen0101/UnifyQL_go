package relation_linking_test

import (
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/expression_tree"
	"github.com/RuiChen0101/unfiyql/internal/relation_chain"
	"github.com/RuiChen0101/unfiyql/internal/relation_linking"
	"github.com/RuiChen0101/unfiyql/pkg/element"
	"github.com/stretchr/testify/assert"
)

func TestLinkWithoutCondition(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		OrderBy:     []string{"tableA.fieldA3 DESC"},
		Limit:       []int{0, 100},
	}
	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)

	linker := relation_linking.NewRelationLinker(rc, tree)
	resultTable, err := linker.Link()
	linkedTree := linker.GetExpressionTree().(*expression_tree.OutputTargetNode)

	assert.Nil(t, err)
	assert.Equal(t, "tableA", resultTable)
	assert.Equal(t, "tableA", linkedTree.OutputTarget)
}

func TestLinkSingleCondition(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "(tableC.fieldC1 & 2) != 0",
	}
	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)

	linker := relation_linking.NewRelationLinker(rc, tree)
	resultTable, err := linker.Link()
	assert.Nil(t, err)

	linkedTree := linker.GetExpressionTree().(*expression_tree.OutputTargetNode)
	rl1 := (*linkedTree.GetLeftNode()).(*expression_tree.RelationNode)
	rl2 := (*rl1.GetLeftNode()).(*expression_tree.RelationNode)
	cn := (*rl2.GetLeftNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "tableA", resultTable)
	assert.Equal(t, "tableB", rl1.FromTable)
	assert.Equal(t, "tableA", rl1.ToTable)
	assert.Equal(t, "tableC", rl2.FromTable)
	assert.Equal(t, "tableB", rl2.ToTable)
	assert.Equal(t, "(tableC.fieldC1 & 2) != 0", cn.ConditionStr)
}

func TestLinkBinaryOpNode(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "tableA.fieldA = 0 AND tableD.fieldD1 = 1",
	}

	tree, _ := expression_tree.ParseExpressionTree(&el)
	rc, _ := relation_chain.BuildRelationChain(&el)

	linker := relation_linking.NewRelationLinker(rc, tree)
	resultTable, err := linker.Link()
	assert.Nil(t, err)

	linkedTree := linker.GetExpressionTree().(*expression_tree.OutputTargetNode)
	bn := (*linkedTree.GetLeftNode()).(*expression_tree.BinaryOperatorNode)
	cn1 := (*bn.GetLeftNode()).(*expression_tree.ConditionNode)
	rl := (*bn.GetRightNode()).(*expression_tree.RelationNode)
	cn2 := (*rl.GetLeftNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "tableA", resultTable)
	assert.Equal(t, "tableD", rl.FromTable)
	assert.Equal(t, "tableA", rl.ToTable)
	assert.Equal(t, "tableA", bn.OutputTarget)
	assert.Equal(t, "AND", bn.OpType)
	assert.Equal(t, "tableA.fieldA = 0", cn1.ConditionStr)
	assert.Equal(t, "tableD.fieldD1 = 1", cn2.ConditionStr)
}
