package expression_tree

import (
	"github.com/RuiChen0101/unfiyql/internal/element"
	"github.com/RuiChen0101/unfiyql/internal/utility"
)

func ParseExpressionTree(el element.UnifyQLElement) (ExpressionTreeNode, error) {
	outputNode := OutputTargetNode{
		Operation:    el.Operation,
		OutputTarget: el.QueryTarget,
		QueryField:   el.QueryField,
		OrderBy:      el.OrderBy,
		Limit:        el.Limit,
	}
	if el.Where == "" {
		return &outputNode, nil
	}

	tokens := utility.RegSplit(el.Where, `\s*(AND|OR|\(|\))\s*`)
	builder := ExpressionTreeBuilder{}
	for _, t := range tokens {
		switch t {
		case "OR":
			if err := builder.BuildOr(); err != nil {
				return nil, err
			}
		case "AND":
			if err := builder.BuildAnd(); err != nil {
				return nil, err
			}
		case "(":
			builder.StartBuildParentheses()
		case ")":
			if err := builder.EndBuildParentheses(); err != nil {
				return nil, err
			}
		case "":
		default:
			if err := builder.BuildCondition(t); err != nil {
				return nil, err
			}
		}
	}
	if err := builder.Flush(); err != nil {
		return nil, err
	}
	outputNode.SetLeftNode(builder.ExpressionTree)
	return &outputNode, nil
}
