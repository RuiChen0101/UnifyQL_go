package execution_plan

type ExecutionPlan struct {
	Operation  int                      `json:"operation"`
	Query      string                   `json:"query"`
	With       []string                 `json:"with"`
	Link       []string                 `json:"link"`
	Where      string                   `json:"where"`
	OrderBy    []string                 `json:"orderBy"`
	Limit      []int                    `json:"limit"`
	Dependency map[string]ExecutionPlan `json:"dependency"`
}
