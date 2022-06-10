package expression_tree

import (
	"errors"
	"fmt"
	"strings"
)

type ExpressionTreeBuilder struct {
	ExpressionTree ExpressionTreeNode
	nodeStack      []ExpressionTreeNode
}

func (eb *ExpressionTreeBuilder) BuildCondition(condition string) error {
	cond, err := NewConditionNode(condition)
	if err != nil {
		eb.nodeStack = append(eb.nodeStack, &brokenConditionNode{Condition: condition})
		return eb.tryRestoreBrokenCondition()
	}
	if eb.ExpressionTree == nil {
		eb.ExpressionTree = cond
		return nil
	}
	if _, ok := eb.ExpressionTree.(*BinaryOperatorNode); !ok {
		return errors.New("ExpressionTreeBuilder: root is not binary operator node")
	}
	if eb.ExpressionTree.GetRightNode() != nil {
		return errors.New("ExpressionTreeBuilder: right node is not empty")
	}
	eb.ExpressionTree.SetRightNode(cond)
	return nil
}

func (eb *ExpressionTreeBuilder) BuildOr() error {
	if eb.ExpressionTree == nil {
		return errors.New("ExpressionTreeBuilder: dangling OR operator")
	}
	node := BinaryOperatorNode{
		OpType: "OR",
	}
	node.SetLeftNode(eb.ExpressionTree)
	eb.nodeStack = append(eb.nodeStack, &node)
	eb.ExpressionTree = nil
	return nil
}

func (eb *ExpressionTreeBuilder) BuildAnd() error {
	if eb.ExpressionTree == nil {
		return errors.New("ExpressionTreeBuilder: dangling AND operator")
	}
	node := BinaryOperatorNode{
		OpType: "AND",
	}
	node.SetLeftNode(eb.ExpressionTree)
	eb.ExpressionTree = &node
	return nil
}

func (eb *ExpressionTreeBuilder) StartBuildParentheses() {
	if eb.ExpressionTree != nil {
		eb.nodeStack = append(eb.nodeStack, eb.ExpressionTree)
		eb.ExpressionTree = nil
	}
	eb.nodeStack = append(eb.nodeStack, &parenthesesMarkerNode{})
}

func (eb *ExpressionTreeBuilder) EndBuildParentheses() error {
	if len(eb.nodeStack) == 0 {
		return nil
	}
	top := eb.safePop()
	if _, ok := top.(*brokenConditionNode); ok {
		eb.safePop()
		eb.nodeStack = append(eb.nodeStack, &brokenConditionNode{Condition: fmt.Sprintf("(%s)", top.(*brokenConditionNode).Condition)})
		return eb.tryRestoreBrokenCondition()
	}
	if eb.ExpressionTree == nil {
		return errors.New("ExpressionTreeBuilder: empty parentheses")
	}
	for _, ok := top.(*parenthesesMarkerNode); !ok; {
		top.SetRightNode(eb.ExpressionTree)
		eb.ExpressionTree = top
		top = eb.safePop()
		_, ok = top.(*parenthesesMarkerNode)
	}
	return nil
}

func (eb *ExpressionTreeBuilder) Flush() error {
	if eb.ExpressionTree == nil {
		return errors.New("ExpressionTreeBuilder: empty tree")
	} else if bn, ok := eb.ExpressionTree.(*BinaryOperatorNode); ok && (bn.GetLeftNode() == nil || bn.GetRightNode() == nil) {
		return fmt.Errorf("ExpressionTreeBuilder: dangling %s operator", bn.OpType)
	}

	for len(eb.nodeStack) != 0 {
		top := eb.nodeStack[len(eb.nodeStack)-1]
		eb.nodeStack = eb.nodeStack[:len(eb.nodeStack)-1]
		if _, ok := top.(*brokenConditionNode); ok {
			return errors.New("ExpressionTreeBuilder: unresolved broken node")
		}
		top.SetRightNode(eb.ExpressionTree)
		eb.ExpressionTree = top
	}
	return nil
}

func (eb *ExpressionTreeBuilder) safePop() ExpressionTreeNode {
	length := len(eb.nodeStack)
	if length > 0 {
		top := eb.nodeStack[length-1]
		eb.nodeStack = eb.nodeStack[:length-1]
		return top
	}
	return nil
}

func (eb *ExpressionTreeBuilder) tryRestoreBrokenCondition() error {
	length := len(eb.nodeStack)
	if length < 2 {
		return nil
	}
	bn1, ok1 := eb.nodeStack[length-1].(*brokenConditionNode)
	bn2, ok2 := eb.nodeStack[length-2].(*brokenConditionNode)
	if ok1 && ok2 {
		eb.nodeStack = eb.nodeStack[:length-2]
		return eb.BuildCondition(strings.TrimSpace(bn2.Condition) + " " + strings.TrimSpace(bn1.Condition))
	}
	return nil
}
