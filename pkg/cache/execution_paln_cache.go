package cache

import "github.com/RuiChen0101/UnifyQL_go/internal/execution_plan"

type ExecutionPlanCache interface {
	Set(key string, plan *execution_plan.ExecutionPlan)
	Get(key string) (*execution_plan.ExecutionPlan, bool)
	FreeSpace()
}
