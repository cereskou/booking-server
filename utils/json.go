package utils

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

//JSON -
var (
	JSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

//StructToJSONTagMap -
func StructToJSONTagMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Tag.Get("json")
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

//StructToMap -
func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Name
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

//UnmarshalStruct -
func UnmarshalStruct(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	b, _ := JSON.Marshal(data)
	JSON.Unmarshal(b, &result)

	return result
}
