package expression_tree_test

import (
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/expression_tree"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewConditionNode(t *testing.T) {
	node1, err := expression_tree.NewConditionNode("tableA.fieldA = \"1234 5678\"")
	assert.Nil(t, err)

	node2, err := expression_tree.NewConditionNode("tableA.fieldA IN (\"0912\",\"0934\")")
	assert.Nil(t, err)

	node3, err := expression_tree.NewConditionNode("(tableB.fieldB & 2) != 0")
	assert.Nil(t, err)

	node4, err := expression_tree.NewConditionNode("tableB.fieldB IS NOT NULL")
	assert.Nil(t, err)

	assert.Equal(t, "tableA", (*node1).OutputTarget)
	assert.Equal(t, "tableA.fieldA = \"1234 5678\"", (*node1).ConditionStr)

	assert.Equal(t, "tableA", (*node2).OutputTarget)
	assert.Equal(t, "tableA.fieldA IN (\"0912\",\"0934\")", (*node2).ConditionStr)

	assert.Equal(t, "tableB", (*node3).OutputTarget)
	assert.Equal(t, "(tableB.fieldB & 2) != 0", (*node3).ConditionStr)

	assert.Equal(t, "tableB", (*node4).OutputTarget)
	assert.Equal(t, "tableB.fieldB IS NOT NULL", (*node4).ConditionStr)
}

func TestInvalidFormatError(t *testing.T) {
	node1, err := expression_tree.NewConditionNode("1=1")
	assert.Nil(t, node1)
	assert.Equal(t, "ConditionNode: invalid format", err.Error())

	node2, err := expression_tree.NewConditionNode("tableA.fieldA=\"valueA\"; DROP Database tableA;--\"")
	assert.Nil(t, node2)
	assert.Equal(t, "ConditionNode: invalid format", err.Error())
}

func TestSetNodePanic(t *testing.T) {
	node, _ := expression_tree.NewConditionNode("tableA.fieldA = \"1234 5678\"")
	assert.Panics(t, func() {
		node.SetLeftNode(&expression_tree.BinaryOperatorNode{})
	})
	assert.Panics(t, func() {
		node.SetRightNode(&expression_tree.BinaryOperatorNode{})
	})
}

func TestGetNodeNil(t *testing.T) {
	node, _ := expression_tree.NewConditionNode("tableA.fieldA = \"1234 5678\"")
	assert.Nil(t, node.GetLeftNode())
	assert.Nil(t, node.GetRightNode())
}
