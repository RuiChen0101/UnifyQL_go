package expression_tree_test

import (
	"testing"

	"github.com/RuiChen0101/UnifyQL_go/internal/expression_tree"
	"github.com/RuiChen0101/UnifyQL_go/pkg/element"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyWhere(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		QueryField:  "fieldA",
		OrderBy:     []string{"tableA.fieldA3 DESC"},
		Limit:       []int{0, 100},
	}
	tree, err := expression_tree.ParseExpressionTree(&el)
	assert.Nil(t, err)

	on := tree.(*expression_tree.OutputTargetNode)

	assert.Nil(t, on.GetLeftNode())
	assert.Nil(t, on.GetRightNode())
	assert.Equal(t, element.UnifyQLOperation.Query, on.Operation)
	assert.Equal(t, "tableA", on.OutputTarget)
	assert.Equal(t, "fieldA", on.QueryField)
	assert.EqualValues(t, []string{"tableA.fieldA3 DESC"}, on.OrderBy)
	assert.EqualValues(t, []int{0, 100}, on.Limit)
}

func TestParseWithSingleWhere(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "tableD.fieldD1 = 1",
	}
	tree, err := expression_tree.ParseExpressionTree(&el)
	assert.Nil(t, err)

	on := tree.(*expression_tree.OutputTargetNode)
	cn := (*on.GetLeftNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "tableA", on.OutputTarget)
	assert.Equal(t, "tableD.fieldD1 = 1", cn.ConditionStr)
}

func TestParseWithComplexWhere(t *testing.T) {
	el := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		Where:       "(tableB.fieldB & 2) != 0 AND (tableA.fieldA IN (\"0912\",\"0934\") OR tableC.fieldC LIKE \"O%\")",
	}
	tree, err := expression_tree.ParseExpressionTree(&el)
	assert.Nil(t, err)

	on := tree.(*expression_tree.OutputTargetNode)
	bn := (*on.GetLeftNode()).(*expression_tree.BinaryOperatorNode)
	lcn := (*bn.GetLeftNode()).(*expression_tree.ConditionNode)
	rbn := (*bn.GetRightNode()).(*expression_tree.BinaryOperatorNode)
	llcn := (*rbn.GetLeftNode()).(*expression_tree.ConditionNode)
	lrcn := (*rbn.GetRightNode()).(*expression_tree.ConditionNode)

	assert.Equal(t, "tableA", on.OutputTarget)
	assert.Equal(t, "AND", bn.OpType)
	assert.Equal(t, "(tableB.fieldB & 2) != 0", lcn.ConditionStr)
	assert.Equal(t, "OR", rbn.OpType)
	assert.Equal(t, "tableA.fieldA IN (\"0912\",\"0934\")", llcn.ConditionStr)
	assert.Equal(t, "tableC.fieldC LIKE \"O%\"", lrcn.ConditionStr)
}

func TestAuthorizationBypassError(t *testing.T) {
	tree, err := expression_tree.ParseExpressionTree(&element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		Where:       "tableA.fieldA=\"valueA\" OR 1=1--\"",
	})
	assert.Nil(t, tree)
	assert.EqualError(t, err, "ExpressionTreeBuilder: empty tree")

	tree, err = expression_tree.ParseExpressionTree(&element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		Where:       "tableA.fieldA=123 OR 1=1--",
	})
	assert.Nil(t, tree)
	assert.EqualError(t, err, "ExpressionTreeBuilder: empty tree")

	tree, err = expression_tree.ParseExpressionTree(&element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		Where:       "tableA.fieldA IN (123) OR 1=1--)",
	})
	assert.Nil(t, tree)
	assert.EqualError(t, err, "ExpressionTreeBuilder: empty tree")
}

func TestMaliciousCommandsError(t *testing.T) {
	tree, err := expression_tree.ParseExpressionTree(&element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		Where:       "tableA.fieldA=\"valueA\"; DROP TABLE tableA--\"",
	})
	assert.Nil(t, tree)
	assert.EqualError(t, err, "ExpressionTreeBuilder: empty tree")

	tree, err = expression_tree.ParseExpressionTree(&element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		Where:       "tableA.fieldA=123; DROP TABLE tableA--",
	})
	assert.Nil(t, tree)
	assert.EqualError(t, err, "ExpressionTreeBuilder: empty tree")

	tree, err = expression_tree.ParseExpressionTree(&element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		Where:       "tableA.fieldA IN (123); DROP TABLE tableA--)",
	})
	assert.Nil(t, tree)
	assert.EqualError(t, err, "ExpressionTreeBuilder: unresolved broken node")
}
