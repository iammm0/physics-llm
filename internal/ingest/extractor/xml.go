package extractor

import (
	"encoding/xml"
	"io"
	"os"
	"strings"
)

// xmlExt 用 xml.Decoder 扫描所有 CharData，拼接成一句话
type xmlExt struct{}

func (xmlExt) Extract(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	var sb strings.Builder
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return sb.String(), nil
		}
		if cd, ok := tok.(xml.CharData); ok {
			txt := strings.TrimSpace(string(cd))
			if txt != "" {
				sb.WriteString(txt)
				sb.WriteString(" ")
			}
		}
	}
	return sb.String(), nil
}

func init() {
	Register(".xml", xmlExt{})
}
