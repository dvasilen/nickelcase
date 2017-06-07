package template

import (
	"encoding/json"
	"fmt"
	"os"
)

func getenv(key string, v ...string) string {
	defaultValue := ""
	if len(v) > 0 {
		defaultValue = v[0]
	}

	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func createMap(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid map call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("map keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func unmarshalJsonObject(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func unmarshalJsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func fileExist(fpath string) bool {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return false
	}
	return true
}
