package expression_tree

type ExpressionTreeNode interface {
	SetLeftNode(node ExpressionTreeNode)
	SetRightNode(node ExpressionTreeNode)
	GetLeftNode() *ExpressionTreeNode
	GetRightNode() *ExpressionTreeNode
}
