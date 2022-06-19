package converter_test

import (
	"testing"

	"github.com/RuiChen0101/UnifyQL_go/pkg/converter"
	"github.com/stretchr/testify/assert"
)

func TestConvertStandardQuery(t *testing.T) {
	queryStr := "QUERY tableA"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.Nil(t, err)
	assert.Equal(t, "SELECT tableA.* FROM tableA", sql)
}

func TestConvertStandardQueryWithField(t *testing.T) {
	queryStr := "QUERY tableA.fieldA"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.Nil(t, err)
	assert.Equal(t, "SELECT tableA.fieldA fieldA FROM tableA", sql)
}

func TestConvertCountQuery(t *testing.T) {
	queryStr := "COUNT tableA"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.Nil(t, err)
	assert.Equal(t, "SELECT count(tableA.*) count FROM tableA", sql)
}

func TestConvertSumQuery(t *testing.T) {
	queryStr := "SUM tableA.fieldA"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.Nil(t, err)
	assert.Equal(t, "SELECT sum(tableA.fieldA) sum FROM tableA", sql)
}

func TestConvertComplexQuery(t *testing.T) {
	queryStr := "QUERY tableA WITH tableB, tableC, tableD LINK tableC.fieldC=tableB.fieldB1,tableD.fieldD=tableA.fieldA1,tableA.fieldA2=tableB.fieldB2 WHERE tableD.fieldD1 != 0 ORDER BY tableA.fieldA4 DESC LIMIT 0,100"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.Nil(t, err)
	assert.Equal(t, "SELECT tableA.* FROM tableB,tableC,tableD,tableA WHERE tableC.fieldC=tableB.fieldB1 AND tableD.fieldD=tableA.fieldA1 AND tableA.fieldA2=tableB.fieldB2 AND tableD.fieldD1 != 0 ORDER BY tableA.fieldA4 DESC LIMIT 0, 100", sql)
}

func TestExtractorError(t *testing.T) {
	queryStr := "tableA"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.EqualError(t, err, "Invalid format")
	assert.Equal(t, "", sql)
}

func TestConvertInvalidSumQueryError(t *testing.T) {
	queryStr := "SUM tableA"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.EqualError(t, err, "Invalid format")
	assert.Equal(t, "", sql)
}

func TestConvertInvalidLimitError(t *testing.T) {
	queryStr := "QUERY tableA LIMIT 0,100,1000"
	sql, err := converter.ConvertToSQL(queryStr)
	assert.EqualError(t, err, "Invalid format")
	assert.Equal(t, "", sql)
}
