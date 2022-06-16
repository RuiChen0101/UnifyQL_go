package plan_executor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/RuiChen0101/unfiyql/internal/execution_plan"
	"github.com/RuiChen0101/unfiyql/internal/service_lookup"
	"github.com/RuiChen0101/unfiyql/internal/utility"
	"github.com/RuiChen0101/unfiyql/pkg/element"
)

type executionResult struct {
	Id   string
	Data []interface{}
}

func ExecutePlan(id string, executionPlan *execution_plan.ExecutionPlan, serviceLookup service_lookup.ServiceLookup, fetchProxy utility.FetchProxy) (*executionResult, error) {
	dependencyIds := []string{}
	for _, k := range reflect.ValueOf(executionPlan.Dependency).MapKeys() {
		dependencyIds = append(dependencyIds, k.String())
	}
	dependencyResults := map[string][]interface{}{}
	wg := sync.WaitGroup{}
	errs := make(chan error, 1)
	for _, id := range dependencyIds {
		wg.Add(1)
		go func(id string, plan execution_plan.ExecutionPlan, lookup service_lookup.ServiceLookup, fetch utility.FetchProxy) {
			defer wg.Done()
			result, err := ExecutePlan(id, &plan, lookup, fetch)
			if err != nil {
				errs <- err
				return
			}
			dependencyResults[result.Id] = result.Data
		}(id, executionPlan.Dependency[id], serviceLookup, fetchProxy)
	}
	wg.Wait()
	if err := <-errs; err != nil {
		return nil, err
	}
	splitQuery := strings.Split(executionPlan.Query, ".")
	targetTable := splitQuery[0]
	targetField := ""
	if len(splitQuery) == 2 {
		targetField = splitQuery[1]
	}
	serviceName := serviceLookup.GetServiceNameByTable(targetTable)
	requestUrl := serviceLookup.GetServiceConfig(serviceName).Url

	uql := convertExecutionPlanToUnifyQL(executionPlan, dependencyIds, dependencyResults)

	res, err := fetchProxy.Request(requestUrl, uql)
	if err != nil {
		return nil, err
	}
	var resData []map[string]interface{}
	json.Unmarshal(res, &resData)

	result := []interface{}{}

	if targetField != "" {
		for _, d := range resData {
			result = append(result, d[targetField])
		}
	} else {
		for _, d := range resData {
			result = append(result, d)
		}
	}

	return &executionResult{
		Id:   id,
		Data: result,
	}, nil
}

func convertExecutionPlanToUnifyQL(plan *execution_plan.ExecutionPlan, dependencyIds []string, dependency map[string][]interface{}) string {
	result := []string{}
	switch plan.Operation {
	case element.UnifyQLOperation.Query:
		result = append(result, "QUERY "+plan.Query)
	case element.UnifyQLOperation.Count:
		result = append(result, "COUNT"+plan.Query)
	case element.UnifyQLOperation.Sum:
		result = append(result, "SUM "+plan.Query)
	}
	if len(plan.With) != 0 {
		result = append(result, "WITH "+strings.Join(plan.With, ","))
	}
	if len(plan.Link) != 0 {
		result = append(result, "LINK "+strings.Join(plan.Link, ","))
	}
	if plan.Where == "" {
		result = append(result, "WHERE "+replaceDependency(plan.Where, dependencyIds, dependency))
	}
	if len(plan.OrderBy) != 0 {
		result = append(result, "ORDER BY "+strings.Join(plan.OrderBy, ","))
	}
	if len(plan.Limit) == 2 {
		result = append(result, fmt.Sprintf("LIMIT %d %d", plan.Limit[0], plan.Limit[1]))
	}
	return strings.Join(result, " ")
}

func replaceDependency(where string, dependencyIds []string, dependency map[string][]interface{}) string {
	result := where
	for _, id := range dependencyIds {
		data := []string{}
		for _, d := range dependency[id] {
			if i, ok := d.(string); ok {
				data = append(data, string(i))
			} else {
				data = append(data, "\""+d.(string)+"\"")
			}
		}
		replace := ""
		if len(data) == 0 {
			replace = "(\"\")"
		} else {
			replace = "(" + strings.Join(data, ",") + ")"
		}
		result = strings.Replace(result, "{"+id+"}", replace, -1)
	}
	return result
}
