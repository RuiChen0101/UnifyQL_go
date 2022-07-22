package cache

import (
	"github.com/RuiChen0101/UnifyQL_go/internal/execution_plan"
)

type DefaultExecutionPlanCache struct {
	data map[string]*execution_plan.ExecutionPlan
}

func NewDefaultExecutionPlanCache() *DefaultExecutionPlanCache {
	return &DefaultExecutionPlanCache{
		data: map[string]*execution_plan.ExecutionPlan{},
	}
}

func (cache *DefaultExecutionPlanCache) Set(key string, plan *execution_plan.ExecutionPlan) {
	cache.data[key] = plan
}

func (cache *DefaultExecutionPlanCache) Get(key string) (*execution_plan.ExecutionPlan, bool) {
	value, ok := cache.data[key]
	return value, ok
}

func (cache *DefaultExecutionPlanCache) FreeSpace() {
	cache.data = map[string]*execution_plan.ExecutionPlan{}
}
