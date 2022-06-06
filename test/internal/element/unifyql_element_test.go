package element_test

import (
	"testing"

	"github.com/RuiChen0101/unfiyql/internal/element"
	"github.com/stretchr/testify/assert"
)

func TestExtractStandardQuery(t *testing.T) {
	queryStr := "QUERY tableA"
	el, err := element.ExtractElement(queryStr)
	actual := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		QueryField:  "",
		With:        []string{},
		Link:        []string{},
		Where:       "",
	}
	assert.Nil(t, err)
	assert.EqualValues(t, *el, actual)
}

func TestExtractStandardQueryWithField(t *testing.T) {
	queryStr := "QUERY tableA.fieldA"
	el, err := element.ExtractElement(queryStr)
	actual := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		QueryField:  "fieldA",
		With:        []string{},
		Link:        []string{},
		Where:       "",
	}
	assert.Nil(t, err)
	assert.EqualValues(t, *el, actual)
}

func TestExtractCountQuery(t *testing.T) {
	queryStr := "COUNT tableA"
	el, err := element.ExtractElement(queryStr)
	actual := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Count,
		QueryTarget: "tableA",
		With:        []string{},
		Link:        []string{},
		Where:       "",
	}
	assert.Nil(t, err)
	assert.EqualValues(t, *el, actual)
}

func TestExtractSumQuery(t *testing.T) {
	queryStr := "SUM tableA.fieldA"
	el, err := element.ExtractElement(queryStr)
	actual := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Sum,
		QueryTarget: "tableA",
		QueryField:  "fieldA",
		With:        []string{},
		Link:        []string{},
		Where:       "",
	}
	assert.Nil(t, err)
	assert.EqualValues(t, *el, actual)
}

func TestExtractComplexQuery(t *testing.T) {
	queryStr := "QUERY tableA WITH tableB, tableC, tableD LINK tableC.fieldC=tableB.fieldB1,tableD.fieldD=tableA.fieldA1,tableA.fieldA2=tableB.fieldB2 WHERE tableD.fieldD1 != 0 ORDER BY tableA.fieldA4 DESC LIMIT 0,100"
	el, err := element.ExtractElement(queryStr)
	actual := element.UnifyQLElement{
		Operation:   element.UnifyQLOperation.Query,
		QueryTarget: "tableA",
		With:        []string{"tableB", "tableC", "tableD"},
		Link:        []string{"tableC.fieldC=tableB.fieldB1", "tableD.fieldD=tableA.fieldA1", "tableA.fieldA2=tableB.fieldB2"},
		Where:       "tableD.fieldD1 != 0",
		OrderBy:     []string{"tableA.fieldA4 DESC"},
		Limit:       []int{0, 100},
	}
	assert.Nil(t, err)
	assert.EqualValues(t, *el, actual)
}

func TestExtractInvalidFormatError(t *testing.T) {
	queryStr := "tableA"
	el, err := element.ExtractElement(queryStr)
	assert.Equal(t, err.Error(), "Invalid format")
	assert.Nil(t, el)
}

func TestExtractInvalidLimitError(t *testing.T) {
	queryStr := "QUERY tableA LIMIT \"aaa\",\"bbb\""
	el, err := element.ExtractElement(queryStr)
	assert.Equal(t, err.Error(), "strconv.Atoi: parsing \"\\\"aaa\\\"\": invalid syntax")
	assert.Nil(t, el)
}
