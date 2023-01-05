package execution_plan

import (
	"errors"
	"fmt"

	"github.com/RuiChen0101/UnifyQL_go/internal/expression_tree"
	"github.com/RuiChen0101/UnifyQL_go/internal/service_lookup"
	"github.com/RuiChen0101/UnifyQL_go/internal/utility"
	"github.com/RuiChen0101/UnifyQL_go/pkg/element"
	"github.com/RuiChen0101/UnifyQL_go/pkg/execution_plan"
)

func GenerateExecutionPlan(tree expression_tree.ExpressionTreeNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*execution_plan.ExecutionPlan, error) {
	if node, ok := tree.(*expression_tree.OutputTargetNode); ok {
		return buildOutputTarget(node, lookup, idGenerator)
	} else if node, ok := tree.(*expression_tree.BinaryOperatorNode); ok {
		return buildBinaryOperator(node, lookup, idGenerator)
	} else if node, ok := tree.(*expression_tree.ConditionNode); ok {
		return buildCondition(node, lookup, idGenerator)
	} else if node, ok := tree.(*expression_tree.RelationNode); ok {
		return buildRelation(node, lookup, idGenerator)
	}
	return nil, errors.New("execution_plan.ExecutionPlan: Invalid expression tree")
}

func buildOutputTarget(node *expression_tree.OutputTargetNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*execution_plan.ExecutionPlan, error) {
	var query string
	if node.QueryField == "" {
		query = node.OutputTarget
	} else {
		query = node.OutputTarget + "." + node.QueryField
	}
	if node.GetLeftNode() == nil {
		return &execution_plan.ExecutionPlan{
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

	return &execution_plan.ExecutionPlan{
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

func buildBinaryOperator(node *expression_tree.BinaryOperatorNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*execution_plan.ExecutionPlan, error) {
	if node.GetLeftNode() == nil || node.GetRightNode() == nil {
		return nil, errors.New("execution_plan.ExecutionPlan: Invalid expression tree")
	}
	leftPlan, err := GenerateExecutionPlan(*node.GetLeftNode(), lookup, idGenerator)
	if err != nil {
		return nil, err
	}

	rightPlan, err := GenerateExecutionPlan(*node.GetRightNode(), lookup, idGenerator)
	if err != nil {
		return nil, err
	}

	plan := &execution_plan.ExecutionPlan{
		Operation:  element.UnifyQLOperation.Query,
		Query:      node.OutputTarget,
		With:       unique(append(leftPlan.With, rightPlan.With...)),
		Link:       unique(append(leftPlan.Link, rightPlan.Link...)),
		Where:      fmt.Sprintf("(%s %s %s)", leftPlan.Where, node.OpType, rightPlan.Where),
		Dependency: leftPlan.Dependency,
	}

	for k, v := range rightPlan.Dependency {
		plan.Dependency[k] = v
	}

	return plan, nil
}

func buildCondition(node *expression_tree.ConditionNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*execution_plan.ExecutionPlan, error) {
	return &execution_plan.ExecutionPlan{
		Operation:  element.UnifyQLOperation.Query,
		Query:      node.OutputTarget,
		Where:      node.ConditionStr,
		Dependency: map[string]execution_plan.ExecutionPlan{},
	}, nil
}

func buildRelation(node *expression_tree.RelationNode, lookup *service_lookup.ServiceLookup, idGenerator utility.IdGenerator) (*execution_plan.ExecutionPlan, error) {
	plan, err := GenerateExecutionPlan(*node.GetLeftNode(), lookup, idGenerator)
	if err != nil {
		return nil, err
	}
	var newPlan execution_plan.ExecutionPlan
	if lookup.IsAllFromSameService([]string{plan.Query, node.ToTable}) {
		newPlan = execution_plan.ExecutionPlan{
			Operation:  element.UnifyQLOperation.Query,
			Query:      node.ToTable,
			With:       unique(append(plan.With, plan.Query)),
			Link:       unique(append(plan.Link, fmt.Sprintf("%s.%s=%s.%s", node.FromTable, node.FromField, node.ToTable, node.ToField))),
			Where:      plan.Where,
			Dependency: plan.Dependency,
		}
	} else {
		plan.Query = node.FromTable + "." + node.FromField
		dependencyId := idGenerator.NanoId8()
		newPlan = execution_plan.ExecutionPlan{
			Operation: element.UnifyQLOperation.Query,
			Query:     node.ToTable,
			Where:     fmt.Sprintf("%s.%s IN {%s}", node.ToTable, node.ToField, dependencyId),
			Dependency: map[string]execution_plan.ExecutionPlan{
				dependencyId: *plan,
			},
		}
	}
	return &newPlan, nil
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
