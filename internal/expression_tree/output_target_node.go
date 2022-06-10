package expression_tree

type OutputTargetNode struct {
	Operation    int
	OutputTarget string
	QueryField   string
	OrderBy      []string
	Limit        []int
	leftNode     *ExpressionTreeNode
}

func (on *OutputTargetNode) SetLeftNode(node ExpressionTreeNode) {
	on.leftNode = &node
}

func (on *OutputTargetNode) SetRightNode(node ExpressionTreeNode) {
	panic("OutputTargetNode: cannot set right node")
}

func (on *OutputTargetNode) GetLeftNode() *ExpressionTreeNode {
	return on.leftNode
}

func (on *OutputTargetNode) GetRightNode() *ExpressionTreeNode {
	return nil
}
