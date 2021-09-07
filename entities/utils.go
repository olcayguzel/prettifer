package entities

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

func getUndelyingValue(v interface{}) reflect.Value {
	t := reflect.ValueOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func getUndelyingType(v interface{}) reflect.Type {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func formatValue(v interface{}, c Column) string {
	var result string
	switch c.Type {
	case DT_INTEGER, DT_COMPLEX:
		result = fmt.Sprintf("%d", v)
	case DT_FLOAT:
		result = fmt.Sprintf("%f", v)
	case DT_DATETIME:
		result = v.(time.Time).Format(c.Format)
	case DT_STRING:
		result = v.(string)
	}
	return result
}

func getAlignment(t DataType) Alignment {
	alignment := AL_LEFT
	switch t {
	case DT_COMPLEX, DT_FLOAT, DT_INTEGER:
		alignment = AL_RIGTH
	default:
		alignment = AL_LEFT
	}
	return alignment
}

func getType(v interface{}, field string) DataType {
	var dataType DataType = DT_STRUCT
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	fieldValue := f.Interface()

	switch fieldValue.(type) {
	case int, int16, int32, int64,
		uint, uint16, uint32, uint64:
		dataType = DT_INTEGER
	case float32, float64:
		dataType = DT_FLOAT
	case complex128, complex64:
		dataType = DT_COMPLEX
	case bool:
		dataType = DT_BOOLEAN
	case time.Time:
		dataType = DT_DATETIME
	case string:
		dataType = DT_STRING
	default:
		dataType = DT_STRUCT
	}
	return dataType
}

func getDataType(t reflect.Type) DataType {
	var dataType DataType = DT_STRUCT
	switch t.Kind() {
	case reflect.Complex128, reflect.Complex64:
		dataType = DT_COMPLEX
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dataType = DT_INTEGER
	case reflect.Float32, reflect.Float64:
		dataType = DT_FLOAT
	case reflect.ValueOf(time.Time{}).Kind():
		dataType = DT_DATETIME
	case reflect.String:
		dataType = DT_STRING
	}
	return dataType
}

func getDefaultDataFormat(t DataType) string {
	format := ""
	switch t {
	case DT_DATETIME:
		format = "2006-01-02 15:04:05"
	}
	return format
}

func padding(value string, width int, alignment Alignment) string {
	var result string
	len := len(value)
	if len >= width {
		result = value[0:(width-3)] + "..."
	} else {
		offset := (width - len)
		if alignment == AL_LEFT {
			result = value + strings.Repeat(" ", offset)
		} else if alignment == AL_RIGTH {
			result = strings.Repeat(" ", offset) + value
		} else {
			result = strings.Repeat(" ", offset/2) + value + strings.Repeat(" ", width-len-(offset/2))
		}
	}
	return result
}

func obtainStringValue(v interface{}, column Column) string {
	value := "<NULL>"
	if v != nil {
		if column.CalculateFunc != nil {
			value = column.CalculateFunc(v)
		} else {
			value = formatValue(v, column)
		}
	}
	return value
}

func getAutoGenerateColumnName(name string) string {
	columnName := ""
	if len(name) > 1 {
		columnName = strings.ToUpper(string(name[0]))
		prvIsUpperChar := unicode.IsUpper(rune(name[0]))
		nextIsUpperChar := unicode.IsUpper(rune(name[1]))
		name = name[1:]
		for len(name) > 0 {
			firstChar := rune(name[0])
			if len(name) > 1 {
				nextIsUpperChar = unicode.IsUpper(rune(name[1]))
			}
			name = name[1:]
			if unicode.IsUpper(firstChar) {
				if prvIsUpperChar && nextIsUpperChar {
					columnName += fmt.Sprintf("%v", strings.ToUpper(string(firstChar)))
				} else if prvIsUpperChar && !nextIsUpperChar {
					columnName += fmt.Sprintf(" %v", strings.ToUpper(string(firstChar)))
				} else {
					columnName += fmt.Sprintf(" %v", string(firstChar))
				}
				prvIsUpperChar = true
			} else if unicode.IsPunct(firstChar) {
				columnName += fmt.Sprintf(" %v", strings.ToUpper(string(name[0])))
			} else {
				columnName += fmt.Sprintf("%v", string(firstChar))
			}
		}
	} else {
		columnName = strings.ToUpper(name)
	}
	return columnName
}

func getValue(value interface{}, fieldName string) interface{} {
	var result interface{}
	val := getUndelyingValue(value)
	if val.Kind() == reflect.Struct {
		val = val.FieldByName(fieldName)
		result = val.Interface()
	}
	return result
}

func detectDataType(v interface{}, f string) DataType {
	var t DataType = DT_STRING
	value := getUndelyingType(v)
	field, ok := value.FieldByName(f)
	if ok {
		t = getDataType(field.Type)
	}
	return t
}

func stringSliceIndexOf(arr []string, s string) int {
	idx := -1
	if len(arr) > 0 {
		for i, v := range arr {
			if strings.Compare(s, v) == 0 {
				idx = i
				break
			}
		}
	}
	return idx
}
