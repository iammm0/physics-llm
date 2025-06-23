package extractor

import (
	"encoding/json"
	"os"
	"strings"
)

// jsonExt 解析 JSON，递归抽取所有键名和字符串型值
type jsonExt struct{}

func (jsonExt) Extract(p string) (string, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return "", err
	}
	var parts []string
	extractJSONStrings(v, &parts)
	return strings.Join(parts, " "), nil
}

// extractJSONStrings 递归扫 interface{}，把所有 string 类型（包括 map 键名）收集到 out
func extractJSONStrings(v interface{}, out *[]string) {
	switch x := v.(type) {
	case string:
		*out = append(*out, x)
	case []interface{}:
		for _, e := range x {
			extractJSONStrings(e, out)
		}
	case map[string]interface{}:
		for k, e := range x {
			*out = append(*out, k)
			extractJSONStrings(e, out)
		}
	}
}

func init() {
	Register(".json", jsonExt{})
}
