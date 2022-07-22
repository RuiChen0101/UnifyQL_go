package unifyql

import (
	"strings"

	"github.com/RuiChen0101/UnifyQL_go/internal/execution_plan"
	"github.com/RuiChen0101/UnifyQL_go/internal/expression_tree"
	"github.com/RuiChen0101/UnifyQL_go/internal/plan_executor"
	"github.com/RuiChen0101/UnifyQL_go/internal/relation_chain"
	"github.com/RuiChen0101/UnifyQL_go/internal/relation_linking"
	"github.com/RuiChen0101/UnifyQL_go/internal/service_lookup"
	"github.com/RuiChen0101/UnifyQL_go/internal/utility"
	"github.com/RuiChen0101/UnifyQL_go/pkg/cache"
	"github.com/RuiChen0101/UnifyQL_go/pkg/element"
	"github.com/RuiChen0101/UnifyQL_go/pkg/service_config"
)

type UnifyQl struct {
	configSource service_config.ServiceConfigSource
	fetchProxy   utility.FetchProxy
	cacheManager cache.ExecutionPlanCache
}

func NewUnifyQl(
	configSource service_config.ServiceConfigSource,
	fetchProxy utility.FetchProxy,
	cacheManager cache.ExecutionPlanCache,
) UnifyQl {
	return UnifyQl{
		configSource: configSource,
		fetchProxy:   fetchProxy,
		cacheManager: cacheManager,
	}
}

func (uql *UnifyQl) Query(query string) ([]interface{}, error) {
	serviceLookup := service_lookup.NewServiceLookup(uql.configSource)
	idGenerator := utility.DefaultIdGenerator{}
	element, err := element.ExtractElement(strings.Replace(query, "\n", " ", -1))
	if err != nil {
		return nil, err
	}

	relationChain, err := relation_chain.BuildRelationChain(element)
	if err != nil {
		return nil, err
	}

	expressionTree, err := expression_tree.ParseExpressionTree(element)
	if err != nil {
		return nil, err
	}

	linker := relation_linking.NewRelationLinker(relationChain, expressionTree)
	if _, err := linker.Link(); err != nil {
		return nil, err
	}
	expandedTree := linker.GetExpressionTree()

	executionPlan, err := execution_plan.GenerateExecutionPlan(expandedTree, &serviceLookup, &idGenerator)
	if err != nil {
		return nil, err
	}

	result, err := plan_executor.ExecutePlan("root", executionPlan, serviceLookup, uql.fetchProxy)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
