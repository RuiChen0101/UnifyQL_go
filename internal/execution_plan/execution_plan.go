package execution_plan

import (
	"errors"
	"fmt"

	"github.com/RuiChen0101/unfiyql/internal/expression_tree"
	"github.com/RuiChen0101/unfiyql/internal/service_lookup"
	"github.com/RuiChen0101/unfiyql/internal/utility"
	"github.com/RuiChen0101/unfiyql/pkg/element"
)

type ExecutionPlan struct {
	Operation  int
	Query      string
	With       []string
	Link       []string
	Where      string
	OrderBy    []string
	Limit      []int
	Dependency map[string]ExecutionPlan
}

func GenerateExecutionPlan(tree expression_tree.ExpressionTreeNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*ExecutionPlan, error) {
	if node, ok := tree.(*expression_tree.OutputTargetNode); ok {
		return buildOutputTarget(node, lookup, idGenerator)
	} else if node, ok := tree.(*expression_tree.BinaryOperatorNode); ok {
		return buildBinaryOperator(node, lookup, idGenerator)
	} else if node, ok := tree.(*expression_tree.ConditionNode); ok {
		return buildCondition(node, lookup, idGenerator)
	} else if node, ok := tree.(*expression_tree.RelationNode); ok {
		return buildRelation(node, lookup, idGenerator)
	}
	return nil, errors.New("ExecutionPlan: Invalid expression tree")
}

func buildOutputTarget(node *expression_tree.OutputTargetNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*ExecutionPlan, error) {
	var query string
	if node.QueryField == "" {
		query = node.OutputTarget
	} else {
		query = node.OutputTarget + "." + node.QueryField
	}
	if node.GetLeftNode() == nil {
		return &ExecutionPlan{
			Operation: node.Operation,
			Query:     query,
			OrderBy:   node.OrderBy,
			Limit:     node.Limit,
		}, nil
	}
	plan, err := GenerateExecutionPlan(*node.GetLeftNode(), lookup, idGenerator)

	if err != nil {
		return nil, err
	}

	return &ExecutionPlan{
		Operation:  node.Operation,
		Query:      query,
		With:       plan.With,
		Link:       plan.Link,
		Where:      plan.Where,
		OrderBy:    node.OrderBy,
		Limit:      node.Limit,
		Dependency: plan.Dependency,
	}, nil
}

func buildBinaryOperator(node *expression_tree.BinaryOperatorNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*ExecutionPlan, error) {
	if node.GetLeftNode() == nil || node.GetRightNode() == nil {
		return nil, errors.New("ExecutionPlan: Invalid expression tree")
	}
	leftPlan, err := GenerateExecutionPlan(*node.GetLeftNode(), lookup, idGenerator)
	if err != nil {
		return nil, err
	}

	rightPlan, err := GenerateExecutionPlan(*node.GetRightNode(), lookup, idGenerator)
	if err != nil {
		return nil, err
	}

	plan := &ExecutionPlan{
		Operation:  element.UnifyQLOperation.Query,
		Query:      node.OutputTarget,
		With:       append(leftPlan.With, rightPlan.With...),
		Link:       append(leftPlan.Link, rightPlan.Link...),
		Where:      fmt.Sprintf("(%s %s %s)", leftPlan.Where, node.OpType, rightPlan.Where),
		Dependency: leftPlan.Dependency,
	}

	for k, v := range rightPlan.Dependency {
		plan.Dependency[k] = v
	}

	return plan, nil
}

func buildCondition(node *expression_tree.ConditionNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*ExecutionPlan, error) {
	return &ExecutionPlan{
		Operation:  element.UnifyQLOperation.Query,
		Query:      node.OutputTarget,
		Where:      node.ConditionStr,
		Dependency: map[string]ExecutionPlan{},
	}, nil
}

func buildRelation(node *expression_tree.RelationNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*ExecutionPlan, error) {
	plan, err := GenerateExecutionPlan(*node.GetLeftNode(), lookup, idGenerator)
	if err != nil {
		return nil, err
	}

	if lookup.IsAllFromSameService([]string{plan.Query, node.ToTable}) {
		return &ExecutionPlan{
			Operation:  element.UnifyQLOperation.Query,
			Query:      node.ToTable,
			With:       append(plan.With, plan.Query),
			Link:       append(plan.Link, fmt.Sprintf("%s.%s=%s.%s", node.FromTable, node.FromField, node.ToTable, node.ToField)),
			Where:      plan.Where,
			Dependency: plan.Dependency,
		}, nil
	} else {
		plan.Query = node.FromTable + "." + node.FromField
		dependencyId := idGenerator.NanoId8()
		return &ExecutionPlan{
			Operation: element.UnifyQLOperation.Query,
			Query:     node.ToTable,
			Where:     fmt.Sprintf("%s.%s IN {%s}", node.ToTable, node.ToField, dependencyId),
			Dependency: map[string]ExecutionPlan{
				dependencyId: *plan,
			},
		}, nil
	}
}
