package extractor

import (
	"os"
	"strings"
)

// rmdExt 当纯文本读入，简单去掉代码块标记 ```{r ...}
type rmdExt struct{}

func (rmdExt) Extract(p string) (string, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	src := string(b)
	// 删掉 ```{r ...} 到 ``` 之间
	for {
		start := strings.Index(src, "```{r")
		if start < 0 {
			break
		}
		end := strings.Index(src[start+1:], "```")
		if end < 0 {
			src = src[:start]
			break
		}
		src = src[:start] + src[start+end+4:]
	}
	// 去掉其余 Markdown 标记
	lines := strings.Split(src, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimPrefix(strings.TrimSpace(l), "#")
	}
	return strings.Join(lines, "\n"), nil
}

func init() {
	Register(".rmd", rmdExt{})
}
