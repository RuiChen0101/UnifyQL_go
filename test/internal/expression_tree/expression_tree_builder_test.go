package expression_tree_test

import (
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/expression_tree"
	"github.com/stretchr/testify/assert"
)

func TestBuildCondition(t *testing.T) {
	builder := expression_tree.ExpressionTreeBuilder{}
	err := builder.BuildCondition("(tableB.fieldB & 2) != 0")
	assert.Nil(t, err)
	builder.Flush()

	cn := builder.ExpressionTree.(*expression_tree.ConditionNode)

	assert.Equal(t, "(tableB.fieldB & 2) != 0", cn.ConditionStr)
}

func TestBuildAnd(t *testing.T) {
	builder := expression_tree.ExpressionTreeBuilder{}
	err := builder.BuildCondition("(tableB.fieldB & 2) != 0")
	assert.Nil(t, err)
	builder.BuildAnd()
	err = builder.BuildCondition("tableA.fieldA IN (\"0912\",\"0934\")")
	assert.Nil(t, err)
	builder.Flush()

	bn := builder.ExpressionTree.(*expression_tree.BinaryOperatorNode)
	lcn := (*bn.GetLeftNode()).(*expression_tree.ConditionNode)
	rcn := (*bn.GetRightNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "AND", bn.OpType)
	assert.Equal(t, "(tableB.fieldB & 2) != 0", lcn.ConditionStr)
	assert.Equal(t, "tableA.fieldA IN (\"0912\",\"0934\")", rcn.ConditionStr)
}

func TestBuildOr(t *testing.T) {
	builder := expression_tree.ExpressionTreeBuilder{}
	err := builder.BuildCondition("(tableB.fieldB & 2) != 0")
	assert.Nil(t, err)
	builder.BuildOr()
	err = builder.BuildCondition("tableA.fieldA IN (\"0912\",\"0934\")")
	assert.Nil(t, err)
	builder.Flush()

	bn := builder.ExpressionTree.(*expression_tree.BinaryOperatorNode)
	lcn := (*bn.GetLeftNode()).(*expression_tree.ConditionNode)
	rcn := (*bn.GetRightNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "OR", bn.OpType)
	assert.Equal(t, "(tableB.fieldB & 2) != 0", lcn.ConditionStr)
	assert.Equal(t, "tableA.fieldA IN (\"0912\",\"0934\")", rcn.ConditionStr)
}

func TestBuildAndOr(t *testing.T) {
	builder := expression_tree.ExpressionTreeBuilder{}
	err := builder.BuildCondition("(tableB.fieldB & 2) != 0")
	assert.Nil(t, err)
	builder.BuildOr()
	err = builder.BuildCondition("tableA.fieldA IN (\"0912\",\"0934\")")
	assert.Nil(t, err)
	builder.BuildAnd()
	err = builder.BuildCondition("tableC.fieldC LIKE \"O%\"")
	assert.Nil(t, err)
	builder.Flush()

	bn := builder.ExpressionTree.(*expression_tree.BinaryOperatorNode)
	lcn := (*bn.GetLeftNode()).(*expression_tree.ConditionNode)
	rbn := (*bn.GetRightNode()).(*expression_tree.BinaryOperatorNode)
	llcn := (*rbn.GetLeftNode()).(*expression_tree.ConditionNode)
	lrcn := (*rbn.GetRightNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "OR", bn.OpType)
	assert.Equal(t, "(tableB.fieldB & 2) != 0", lcn.ConditionStr)
	assert.Equal(t, "AND", rbn.OpType)
	assert.Equal(t, "tableA.fieldA IN (\"0912\",\"0934\")", llcn.ConditionStr)
	assert.Equal(t, "tableC.fieldC LIKE \"O%\"", lrcn.ConditionStr)
}

func TestBuildWithParentheses(t *testing.T) {
	builder := expression_tree.ExpressionTreeBuilder{}
	err := builder.BuildCondition("(tableB.fieldB & 2) != 0")
	assert.Nil(t, err)
	builder.BuildAnd()
	builder.StartBuildParentheses()
	err = builder.BuildCondition("tableA.fieldA IN (\"0912\",\"0934\")")
	assert.Nil(t, err)
	builder.BuildOr()
	err = builder.BuildCondition("tableC.fieldC LIKE \"O%\"")
	assert.Nil(t, err)
	builder.EndBuildParentheses()
	builder.Flush()

	bn := builder.ExpressionTree.(*expression_tree.BinaryOperatorNode)
	lcn := (*bn.GetLeftNode()).(*expression_tree.ConditionNode)
	rbn := (*bn.GetRightNode()).(*expression_tree.BinaryOperatorNode)
	llcn := (*rbn.GetLeftNode()).(*expression_tree.ConditionNode)
	lrcn := (*rbn.GetRightNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "AND", bn.OpType)
	assert.Equal(t, "(tableB.fieldB & 2) != 0", lcn.ConditionStr)
	assert.Equal(t, "OR", rbn.OpType)
	assert.Equal(t, "tableA.fieldA IN (\"0912\",\"0934\")", llcn.ConditionStr)
	assert.Equal(t, "tableC.fieldC LIKE \"O%\"", lrcn.ConditionStr)
}
