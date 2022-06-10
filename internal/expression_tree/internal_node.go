package expression_tree

type brokenConditionNode struct {
	Condition string
}

func (bn *brokenConditionNode) SetLeftNode(node ExpressionTreeNode) {
	panic("brokenConditionNode: cannot set left node")
}

func (bn *brokenConditionNode) SetRightNode(node ExpressionTreeNode) {
	panic("brokenConditionNode: cannot set right node")
}

func (bn *brokenConditionNode) GetLeftNode() *ExpressionTreeNode {
	return nil
}

func (bn *brokenConditionNode) GetRightNode() *ExpressionTreeNode {
	return nil
}

type parenthesesMarkerNode struct{}

func (pn *parenthesesMarkerNode) SetLeftNode(node ExpressionTreeNode) {
	panic("parenthesesMarkerNode: cannot set left node")
}

func (pn *parenthesesMarkerNode) SetRightNode(node ExpressionTreeNode) {
	panic("parenthesesMarkerNode: cannot set right node")
}

func (pn *parenthesesMarkerNode) GetLeftNode() *ExpressionTreeNode {
	return nil
}

func (pn *parenthesesMarkerNode) GetRightNode() *ExpressionTreeNode {
	return nil
}
