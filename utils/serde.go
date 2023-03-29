package utils

import "encoding/json"

func ToJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func FromJSON[T any](jsonStr string) T {
	var result T
	_ = json.Unmarshal([]byte(jsonStr), &result)
	return result
}
