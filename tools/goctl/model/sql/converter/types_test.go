package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/ddl-parser/parser"
)

func TestConvertDataType(t *testing.T) {
	v, err := ConvertDataType(parser.TinyInt, false)
	assert.Nil(t, err)
	assert.Equal(t, "int64", v)

	v, err = ConvertDataType(parser.TinyInt, true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullInt64", v)

	v, err = ConvertDataType(parser.Timestamp, false)
	assert.Nil(t, err)
	assert.Equal(t, "time.Time", v)

	v, err = ConvertDataType(parser.Timestamp, true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullTime", v)
}
