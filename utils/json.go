package utils

import "encoding/json"

var (
	EmptyJsonArray  = StringToBytes("[]")
	EmptyJsonObject = StringToBytes("{}")
)

// JsonArray 将slice转换成[]byte数据，如果slice为nil或空则返回空json array bytes
func JsonArray[T any](args ...[]T) []byte {
	if len(args) == 0 || len(args[0]) == 0 {
		return EmptyJsonArray
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}

// JsonObject 将object转换成[]byte数据，如果object为nil或空则返回空json object bytes
func JsonObject(args ...any) []byte {
	if len(args) == 0 {
		return EmptyJsonObject
	}

	if args[0] == nil {
		return EmptyJsonObject
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}
