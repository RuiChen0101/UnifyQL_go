package execution_plan

type ExecutionPlan struct {
	Operation  int                      `json:"operation"`
	Query      string                   `json:"query,omitempty"`
	With       []string                 `json:"with,omitempty"`
	Link       []string                 `json:"link,omitempty"`
	Where      string                   `json:"where,omitempty"`
	OrderBy    []string                 `json:"orderBy,omitempty"`
	Limit      []int                    `json:"limit,omitempty"`
	Dependency map[string]ExecutionPlan `json:"dependency,omitempty"`
}
