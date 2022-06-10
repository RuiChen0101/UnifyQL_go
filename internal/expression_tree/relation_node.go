package expression_tree

type RelationNode struct {
	ToTable   string
	ToField   string
	FromTable string
	FromField string
	leftNode  *ExpressionTreeNode
}

func (rn *RelationNode) SetLeftNode(node ExpressionTreeNode) {
	rn.leftNode = &node
}

func (rn *RelationNode) SetRightNode(node ExpressionTreeNode) {
	panic("RelationNode: cannot set right node")
}

func (rn *RelationNode) GetLeftNode() *ExpressionTreeNode {
	return rn.leftNode
}

func (rn *RelationNode) GetRightNode() *ExpressionTreeNode {
	return nil
}
