package converter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RuiChen0101/unfiyql/pkg/element"
)

func ConvertToSQL(unifyQl string) (string, error) {
	el, err := element.ExtractElement(unifyQl)
	if err != nil {
		return "", err
	}
	return ConvertToSQLByElement(el)
}

func ConvertToSQLByElement(el *element.UnifyQLElement) (string, error) {
	result := []string{}
	el.With = append(el.With, el.QueryTarget)
	switch el.Operation {
	case element.UnifyQLOperation.Query:
		if el.QueryField != "" {
			result = append(result, fmt.Sprintf("SELECT %s.%s %s", el.QueryTarget, el.QueryField, el.QueryField))
		} else {
			result = append(result, fmt.Sprintf("SELECT %s.*", el.QueryTarget))
		}
	case element.UnifyQLOperation.Count:
		result = append(result, "SELECT count(*) count")
	case element.UnifyQLOperation.Sum:
		if el.QueryField == "" {
			return "", errors.New("Invalid format")
		}
		result = append(result, fmt.Sprintf("SELECT sum(%s.%s) sum", el.QueryTarget, el.QueryField))
	}
	result = append(result, fmt.Sprintf("FROM %s", strings.Join(el.With, ",")))
	if len(el.Link) != 0 {
		result = append(result, fmt.Sprintf("WHERE %s", strings.Join(el.Link, " AND ")))
	}
	if el.Where != "" {
		if len(el.Link) != 0 {
			result = append(result, "AND")
		} else {
			result = append(result, "WHERE")
		}
		result = append(result, el.Where)
	}
	if len(el.OrderBy) != 0 {
		result = append(result, fmt.Sprintf("ORDER BY %s", strings.Join(el.OrderBy, ",")))
	}
	if len(el.Limit) == 2 {
		result = append(result, fmt.Sprintf("LIMIT %d, %d", el.Limit[0], el.Limit[1]))
	} else if len(el.Limit) != 0 {
		return "", errors.New("Invalid format")
	}
	return strings.Join(result, " "), nil
}
