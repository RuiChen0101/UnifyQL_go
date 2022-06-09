package expression_tree

import (
	"errors"
	"regexp"
)

type ConditionNode struct {
	ConditionStr string
	OutputTarget string
}

func NewConditionNode(conditionStr string) (*ConditionNode, error) {
	node := ConditionNode{
		ConditionStr: conditionStr,
	}

	reg := regexp.MustCompile(`[\(]?\s*(\w*)\.(\w*)\s*([\&\|\+\-\*\/])?\s*([^)]*)?\s*[\)]?\s*(=|!=|<|<=|>|>=|LIKE|NOT IN|IN|IS NULL|IS NOT NULL)\s*(.*)`)
	captured := reg.FindStringSubmatch(conditionStr)
	if len(captured) != 7 {
		return nil, errors.New("ConditionNode: invalid format")
	}

	node.OutputTarget = captured[1]

	if !isValidConditionValue(captured[5], captured[6]) {
		return nil, errors.New("ConditionNode: invalid format")
	}

	return &node, nil
}

func isValidConditionValue(op string, value string) bool {
	if m, err := regexp.MatchString(`=|!=|<|<=|>|>=|LIKE`, op); m && err == nil {
		c, err2 := regexp.MatchString(`^\d+$|^"[^"]+"$`, value)
		return value != "" && c && err2 == nil
	} else if op == "IS NULL" || op == "IS NOT NULL" {
		return value == ""
	} else if op == "IN" || op == "NOT IN" {
		c, err := regexp.MatchString(`^\(([^\(\),]+,)*([^\(\),]+)\)$`, value)
		return value != "" && c && err == nil
	}
	return false
}

func (cn *ConditionNode) SetLeftNode(node ExpressionTreeNode) {
	panic("ConditionNode: cannot set left node")
}

func (cn *ConditionNode) SetRightNode(node ExpressionTreeNode) {
	panic("ConditionNode: cannot set right node")
}

func (cn *ConditionNode) GetLeftNode() *ExpressionTreeNode {
	return nil
}

func (cn *ConditionNode) GetRightNode() *ExpressionTreeNode {
	return nil
}
