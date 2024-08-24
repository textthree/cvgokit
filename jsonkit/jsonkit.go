package jsonkit

import (
	"encoding/json"
	"fmt"
)

// 序列化成json
func JsonEncode(val interface{}) []byte {
	json, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Json encode error")
	}
	return json
}

// JSON反序列化
func JsonDecodeToMap(str string) (thisMap map[string]any, err error) {
	if err = json.Unmarshal([]byte(str), &thisMap); err != nil {
		fmt.Printf("Json反序列化为Map出错: %s\n", err.Error())
	}
	return
}

// JSON 数组字符串反序列化成 map 数组
func JsonDecodeToMapArray(str string) []map[string]interface{} {
	var thisMap []map[string]interface{}
	var err error
	if err = json.Unmarshal([]byte(str), &thisMap); err != nil {
		fmt.Printf("Json反序列化为Map出错: %s\n", err.Error(), str)
		return nil
	}
	return thisMap
}

// JSON 反序列化
// example: strkit.JsonDecode(jsonStr, &v)
func JsonDecode(jsonString string, v any) error {
	return json.Unmarshal([]byte(jsonString), v)
}

func Json_decode_map(str string) (thisMap map[string]interface{}, err error) {
	if err = json.Unmarshal([]byte(str), &thisMap); err != nil {
		fmt.Printf("Json反序列化为Map出错: %s\n", err.Error())
	}
	return
}
