package expression_tree_test

import (
	"testing"

	"github.com/RuiChen0101/UnifyQL_go/internal/expression_tree"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewBinaryOperatorNode(t *testing.T) {
	node := expression_tree.BinaryOperatorNode{OpType: "AND"}
	node.OutputTarget = "tableA"
	assert.Nil(t, node.GetLeftNode())
	assert.Nil(t, node.GetRightNode())
	assert.Equal(t, "AND", node.OpType)
	assert.Equal(t, "tableA", node.OutputTarget)
}

func TestSetNode(t *testing.T) {
	node := expression_tree.BinaryOperatorNode{OpType: "AND"}
	node.SetLeftNode(&expression_tree.BinaryOperatorNode{OpType: "AND"})
	node.SetRightNode(&expression_tree.BinaryOperatorNode{OpType: "OR"})

	assert.Equal(t, "AND", (*node.GetLeftNode()).(*expression_tree.BinaryOperatorNode).OpType)
	assert.Equal(t, "OR", (*node.GetRightNode()).(*expression_tree.BinaryOperatorNode).OpType)
}
