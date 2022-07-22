package cache

import (
	"github.com/RuiChen0101/UnifyQL_go/internal/execution_plan"
)

type DefaultExecutionPlanCache struct {
	data     map[string]*execution_plan.ExecutionPlan
	hits     map[string]int
	maxCache int
}

func NewDefaultExecutionPlanCache() *DefaultExecutionPlanCache {
	return &DefaultExecutionPlanCache{
		data:     map[string]*execution_plan.ExecutionPlan{},
		hits:     map[string]int{},
		maxCache: 100,
	}
}

func (cache *DefaultExecutionPlanCache) Set(key string, plan *execution_plan.ExecutionPlan) {
	cache.data[key] = plan
}

func (cache *DefaultExecutionPlanCache) Get(key string) (*execution_plan.ExecutionPlan, error) {

}

func (cache *DefaultExecutionPlanCache) FreeUnused() {

}
