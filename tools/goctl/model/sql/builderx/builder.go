package builderx

import (
	"fmt"
	"reflect"

	"github.com/go-xorm/builder"
)

const dbTag = "db"

// NewEq wraps builder.Eq
func NewEq(in interface{}) builder.Eq {
	return builder.Eq(ToMap(in))
}

// NewGt wraps builder.Gt
func NewGt(in interface{}) builder.Gt {
	return builder.Gt(ToMap(in))
}

// ToMap converts interface into map
func ToMap(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(dbTag); tagv != "" {
			// set key of map to value in struct field
			val := v.Field(i)
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()

			if reflect.DeepEqual(current, zero) {
				continue
			}
			out[tagv] = current
		}
	}

	return out
}

// FieldNames returns field names from given in.
// deprecated: use RawFieldNames instead automatically while model generating after goctl version v1.1.0
func FieldNames(in interface{}) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(dbTag); tagv != "" {
			out = append(out, tagv)
		} else {
			out = append(out, fi.Name)
		}
	}

	return out
}

// RawFieldNames converts golang struct field into slice string
func RawFieldNames(in interface{}) []string {
	v := reflect.ValueOf(in)

	return nestFieldNames(v) // recursive processor
}

// recursive processor
func nestFieldNames(v reflect.Value) (out []string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		ft := fi.Type

		// fix for support nest struct fields
		if fi.Anonymous && ft.Kind() == reflect.Struct {
			out = append(out, nestFieldNames(v.Field(i))...)
		} else {
			if tagV := fi.Tag.Get(dbTag); tagV != "" {
				out = append(out, fmt.Sprintf("`%s`", tagV))
			} else {
				//out = append(out, fmt.Sprintf(`"%s"`, fi.Name))
			}
		}
	}
	return out
}
