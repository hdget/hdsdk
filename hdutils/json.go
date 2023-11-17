package hdutils

import (
	"encoding/json"
	"reflect"
)

var (
	EmptyJsonArray  = StringToBytes("[]")
	EmptyJsonObject = StringToBytes("{}")
)

// JsonArray 将slice转换成[]byte数据，如果slice为nil或空则返回空json array bytes
func JsonArray(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return EmptyJsonArray
	}

	if reflect.TypeOf(args[0]).Kind() != reflect.Slice {
		return EmptyJsonArray
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}

// JsonObject 将object转换成[]byte数据，如果object为nil或空则返回空json object bytes
func JsonObject(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return EmptyJsonObject
	}

	// 如果传入值为slice,则返回empty object
	if reflect.TypeOf(args[0]).Kind() == reflect.Slice {
		return EmptyJsonObject
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}
