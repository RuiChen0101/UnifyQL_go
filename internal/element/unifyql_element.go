package element

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/RuiChen0101/unfiyql/internal/utility"
)

var UnifyQLOperation = struct {
	Query int
	Count int
	Sum   int
}{
	Query: 0,
	Count: 1,
	Sum:   2,
}

type UnifyQLElement struct {
	Operation   int
	QueryTarget string
	QueryField  string
	With        []string
	Link        []string
	Where       string
	OrderBy     []string
	Limit       []int
}

func ExtractElement(queryStr string) (*UnifyQLElement, error) {
	result := UnifyQLElement{
		Operation:   UnifyQLOperation.Query,
		QueryTarget: "",
		With:        []string{},
		Link:        []string{},
		Where:       "",
	}

	operationRegexp := regexp.MustCompile(`(QUERY|SUM|COUNT)\s*([^\s]+)\s*(.*)`)
	capturedGroups := operationRegexp.FindStringSubmatch(queryStr)
	if len(capturedGroups) == 0 {
		return nil, errors.New("Invalid format")
	}

	switch capturedGroups[1] {
	case "QUERY":
		result.Operation = UnifyQLOperation.Query
	case "COUNT":
		result.Operation = UnifyQLOperation.Count
	case "SUM":
		result.Operation = UnifyQLOperation.Sum
	}
	splitQueryTarget := strings.Split(capturedGroups[2], ".")
	result.QueryTarget = splitQueryTarget[0]
	if len(splitQueryTarget) == 2 {
		result.QueryField = splitQueryTarget[1]
	}

	splitQueryStr := utility.RegSplit(capturedGroups[3], `\s*(WITH|LINK|WHERE|ORDER BY|LIMIT)\s*`)

	filteredQueryStr := []string{}
	for _, n := range splitQueryStr {
		if n != "" {
			filteredQueryStr = append(filteredQueryStr, n)
		}
	}

	dotRegexp := regexp.MustCompile(`\s*,\s*`)
	for i := 0; i < len(filteredQueryStr); i += 2 {
		keyword := filteredQueryStr[i]
		value := filteredQueryStr[i+1]
		switch keyword {
		case "WITH":
			result.With = dotRegexp.Split(value, -1)
		case "LINK":
			result.Link = dotRegexp.Split(value, -1)
		case "WHERE":
			result.Where = value
		case "ORDER BY":
			result.OrderBy = dotRegexp.Split(value, -1)
		case "LIMIT":
			result.Limit = []int{}
			for _, n := range dotRegexp.Split(value, -1) {
				val, err := strconv.Atoi(n)
				if err != nil {
					return nil, err
				}
				result.Limit = append(result.Limit, val)
			}
		}
	}

	return &result, nil
}
