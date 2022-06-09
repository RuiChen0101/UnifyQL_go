package expression_tree

type BinaryOperatorNode struct {
	OpType       string
	OutputTarget string
	leftNode     ExpressionTreeNode
	rightNode    ExpressionTreeNode
}

func NewBinaryOperatorNode(opType string) (*BinaryOperatorNode, error) {
	return &BinaryOperatorNode{
		OpType: opType,
	}, nil
}

func (bn *BinaryOperatorNode) SetLeftNode(node ExpressionTreeNode) {
	bn.leftNode = node
}

func (bn *BinaryOperatorNode) SetRightNode(node ExpressionTreeNode) {
	bn.rightNode = node
}

func (bn *BinaryOperatorNode) GetLeftNode() *ExpressionTreeNode {
	return &bn.leftNode
}

func (bn *BinaryOperatorNode) GetRightNode() *ExpressionTreeNode {
	return &bn.rightNode
}
