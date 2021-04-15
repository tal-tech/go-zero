package sqlx

import (
	"errors"
	"reflect"
	"strings"

	"github.com/tal-tech/go-zero/core/mapping"
)

const tagName = "db"

var (
	// ErrNotMatchDestination is an error that indicates not matching destination to scan.
	ErrNotMatchDestination = errors.New("not matching destination to scan")
	// ErrNotReadableValue is an error that indicates value is not addressable or interfaceable.
	ErrNotReadableValue = errors.New("value not addressable or interfaceable")
	// ErrNotSettable is an error that indicates the passed in variable is not settable.
	ErrNotSettable = errors.New("passed in variable is not settable")
	// ErrUnsupportedValueType is an error that indicates unsupported unmarshal type.
	ErrUnsupportedValueType = errors.New("unsupported unmarshal type")
)

type rowsScanner interface {
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(v ...interface{}) error
}

func getTaggedFieldValueMap(v reflect.Value) (map[string]interface{}, error) {
	rt := mapping.Deref(v.Type())
	size := rt.NumField()
	result := make(map[string]interface{}, size)

	for i := 0; i < size; i++ {
		key := parseTagName(rt.Field(i))
		if len(key) == 0 {
			return nil, nil
		}

		valueField := reflect.Indirect(v).Field(i)
		switch valueField.Kind() {
		case reflect.Ptr:
			if !valueField.CanInterface() {
				return nil, ErrNotReadableValue
			}
			if valueField.IsNil() {
				baseValueType := mapping.Deref(valueField.Type())
				valueField.Set(reflect.New(baseValueType))
			}
			result[key] = valueField.Interface()
		default:
			if !valueField.CanAddr() || !valueField.Addr().CanInterface() {
				return nil, ErrNotReadableValue
			}
			result[key] = valueField.Addr().Interface()
		}
	}

	return result, nil
}

func mapStructFieldsIntoSlice(v reflect.Value, columns []string, strict bool) ([]interface{}, error) {
	fields := unwrapFields(v)
	if strict && len(columns) < len(fields) {
		return nil, ErrNotMatchDestination
	}

	taggedMap, err := getTaggedFieldValueMap(v)
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	if len(taggedMap) == 0 {
		for i := 0; i < len(values); i++ {
			valueField := fields[i]
			switch valueField.Kind() {
			case reflect.Ptr:
				if !valueField.CanInterface() {
					return nil, ErrNotReadableValue
				}
				if valueField.IsNil() {
					baseValueType := mapping.Deref(valueField.Type())
					valueField.Set(reflect.New(baseValueType))
				}
				values[i] = valueField.Interface()
			default:
				if !valueField.CanAddr() || !valueField.Addr().CanInterface() {
					return nil, ErrNotReadableValue
				}
				values[i] = valueField.Addr().Interface()
			}
		}
	} else {
		for i, column := range columns {
			if tagged, ok := taggedMap[column]; ok {
				values[i] = tagged
			} else {
				var anonymous interface{}
				values[i] = &anonymous
			}
		}
	}

	return values, nil
}

func parseTagName(field reflect.StructField) string {
	key := field.Tag.Get(tagName)
	if len(key) == 0 {
		return ""
	}

	options := strings.Split(key, ",")
	return options[0]
}

func unmarshalRow(v interface{}, scanner rowsScanner, strict bool) error {
	if !scanner.Next() {
		if err := scanner.Err(); err != nil {
			return err
		}
		return ErrNotFound
	}

	rv := reflect.ValueOf(v)
	if err := mapping.ValidatePtr(&rv); err != nil {
		return err
	}

	rte := reflect.TypeOf(v).Elem()
	rve := rv.Elem()
	switch rte.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		if rve.CanSet() {
			return scanner.Scan(v)
		}

		return ErrNotSettable
	case reflect.Struct:
		columns, err := scanner.Columns()
		if err != nil {
			return err
		}

		values, err := mapStructFieldsIntoSlice(rve, columns, strict)
		if err != nil {
			return err
		}

		return scanner.Scan(values...)
	default:
		return ErrUnsupportedValueType
	}
}

func unmarshalRows(v interface{}, scanner rowsScanner, strict bool) error {
	rv := reflect.ValueOf(v)
	if err := mapping.ValidatePtr(&rv); err != nil {
		return err
	}

	rt := reflect.TypeOf(v)
	rte := rt.Elem()
	rve := rv.Elem()
	switch rte.Kind() {
	case reflect.Slice:
		if rve.CanSet() {
			ptr := rte.Elem().Kind() == reflect.Ptr
			appendFn := func(item reflect.Value) {
				if ptr {
					rve.Set(reflect.Append(rve, item))
				} else {
					rve.Set(reflect.Append(rve, reflect.Indirect(item)))
				}
			}
			fillFn := func(value interface{}) error {
				if rve.CanSet() {
					if err := scanner.Scan(value); err != nil {
						return err
					}

					appendFn(reflect.ValueOf(value))
					return nil
				}
				return ErrNotSettable
			}

			base := mapping.Deref(rte.Elem())
			switch base.Kind() {
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64,
				reflect.String:
				for scanner.Next() {
					value := reflect.New(base)
					if err := fillFn(value.Interface()); err != nil {
						return err
					}
				}
			case reflect.Struct:
				columns, err := scanner.Columns()
				if err != nil {
					return err
				}

				for scanner.Next() {
					value := reflect.New(base)
					values, err := mapStructFieldsIntoSlice(value, columns, strict)
					if err != nil {
						return err
					}

					if err := scanner.Scan(values...); err != nil {
						return err
					}

					appendFn(value)
				}
			default:
				return ErrUnsupportedValueType
			}

			return nil
		}

		return ErrNotSettable
	default:
		return ErrUnsupportedValueType
	}
}

func unwrapFields(v reflect.Value) []reflect.Value {
	var fields []reflect.Value
	indirect := reflect.Indirect(v)

	for i := 0; i < indirect.NumField(); i++ {
		child := indirect.Field(i)
		if child.Kind() == reflect.Ptr && child.IsNil() {
			baseValueType := mapping.Deref(child.Type())
			child.Set(reflect.New(baseValueType))
		}

		child = reflect.Indirect(child)
		childType := indirect.Type().Field(i)
		if child.Kind() == reflect.Struct && childType.Anonymous {
			fields = append(fields, unwrapFields(child)...)
		} else {
			// fix for scan nest struct fields
			if tagV := childType.Tag.Get(tagName); tagV != "" {
				fields = append(fields, child)
			}
		}
	}

	return fields
}
