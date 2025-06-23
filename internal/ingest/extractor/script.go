package extractor

import (
	"os"
	"strings"
)

// codeExt 简单去注释，保留代码与文档字符串
type codeExt struct{}

func (codeExt) Extract(p string) (string, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(b), "\n")
	for i, l := range lines {
		// 去掉 Python/R/Matlab 中以 # 或 % 开头的注释
		if idx := strings.IndexAny(l, "#%"); idx >= 0 {
			l = l[:idx]
		}
		lines[i] = strings.TrimSpace(l)
	}
	return strings.Join(lines, "\n"), nil
}

func init() {
	Register(".py", codeExt{})
	Register(".r", codeExt{})
	Register(".m", codeExt{}) // MATLAB 脚本
}
