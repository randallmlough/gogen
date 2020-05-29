package template

import (
	"reflect"
	"strings"
	"text/template"
)

type FuncMap = template.FuncMap

var StandardTemplateFuncs = FuncMap{
	"titleCase":         strings.Title,
	"valueOf":           reflectValue,
	"typeOf":            reflectType,
	"isCustomType":      isCustomType,
	"getUnderlyingType": getUnderlyingType,
}

func reflectValue(v interface{}) reflect.Value {
	return reflect.ValueOf(v)
}
func reflectType(v interface{}) reflect.Type {
	return reflect.TypeOf(v)
}
func isCustomType(v interface{}) bool {
	typ := reflect.TypeOf(v)
	switch typ.String() {
	case "bool",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"string":
		return false
	}
	return true
}

func getUnderlyingType(typ reflect.Type) string {
	return typ.Kind().String()
}
