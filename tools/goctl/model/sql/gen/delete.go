package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genDelete(table Table, withCache bool) (string, error) {
	keySet := collection.NewSet()
	keyVariableSet := collection.NewSet()
	for fieldName, key := range table.CacheKey {
		if fieldName == table.PrimaryKey.Name.Source() {
			keySet.AddStr(key.KeyExpression)
		}
		if fieldName != table.PrimaryKey.Name.Source() {
			keySet.AddStr(key.DataKeyExpression)
		}
		keyVariableSet.AddStr(key.Variable)
	}
	var containsIndexCache = false
	for _, item := range table.Fields {
		if item.IsKey {
			containsIndexCache = true
			break
		}
	}
	camel := table.Name.Snake2Camel()
	output, err := templatex.With("delete").
		Parse(template.Delete).
		Execute(map[string]interface{}{
			"upperStartCamelObject":     camel,
			"withCache":                 withCache,
			"containsIndexCache":        containsIndexCache,
			"lowerStartCamelPrimaryKey": stringx.From(table.PrimaryKey.Name.Snake2Camel()).LowerStart(),
			"dataType":                  table.PrimaryKey.DataType,
			"keys":                      strings.Join(keySet.KeysStr(), "\n"),
			"originalPrimaryKey":        table.PrimaryKey.Name.Source(),
			"keyValues":                 strings.Join(keyVariableSet.KeysStr(), ", "),
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
