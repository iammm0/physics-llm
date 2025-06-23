package extractor

type Extractor interface {
	Extract(path string) (string, error) // 返回纯文本
}

var registry = map[string]Extractor{}

func Register(ext string, ex Extractor) { registry[ext] = ex }
func Get(ext string) (Extractor, bool)  { ex, ok := registry[ext]; return ex, ok }
