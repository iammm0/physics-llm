package extractor

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// yamlExt 和 JSON 类似，解析后抽取所有字符串
type yamlExt struct{}

func (yamlExt) Extract(p string) (string, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	var v interface{}
	if err := yaml.Unmarshal(b, &v); err != nil {
		return "", err
	}
	var parts []string
	extractJSONStrings(v, &parts) // 复用 JSON 中的抽取函数
	return strings.Join(parts, " "), nil
}

func init() {
	Register(".yaml", yamlExt{})
	Register(".yml", yamlExt{})
}
