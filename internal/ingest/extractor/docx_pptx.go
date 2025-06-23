package extractor

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// docx 读取 word/document.xml 并抽取文本节点
type docx struct{}

func (docx) Extract(p string) (string, error) {
	zr, err := zip.OpenReader(p)
	if err != nil {
		return "", err
	}
	defer zr.Close()

	var raw []byte
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			raw, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", err
			}
			break
		}
	}
	if raw == nil {
		return "", fmt.Errorf("docx: 找不到 word/document.xml")
	}
	return extractTextFromXML(raw, "w:t", "w:p"), nil
}

// pptx 读取 ppt/slides/slideN.xml 并抽取文本节点
type pptx struct{}

func (pptx) Extract(p string) (string, error) {
	zr, err := zip.OpenReader(p)
	if err != nil {
		return "", err
	}
	defer zr.Close()

	var allText []string
	for _, f := range zr.File {
		// slide 文件位于 ppt/slides/slideX.xml
		if filepath.Dir(f.Name) == "ppt/slides" && filepath.Ext(f.Name) == ".xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			raw, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", err
			}
			// pptx 文本节点是 <a:t>，段落结束用 <a:p>
			txt := extractTextFromXML(raw, "a:t", "a:p")
			allText = append(allText, txt)
		}
	}
	return strings.Join(allText, "\n"), nil
}

// extractTextFromXML 从 raw XML 中提取所有 <textTag> 节点文本，
// 并在每遇到一次 endTag 时插入换行。
func extractTextFromXML(raw []byte, textTag, endTag string) string {
	dec := xml.NewDecoder(bytes.NewReader(raw))
	var sb strings.Builder

	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			// 出错就先返回已累积文本
			return sb.String()
		}
		switch t := tok.(type) {
		case xml.CharData:
			// 直接写入文本节点内容
			sb.WriteString(string(t))
		case xml.EndElement:
			if t.Name.Local == endTag {
				sb.WriteString("\n")
			}
		}
	}
	return sb.String()
}

func init() {
	Register(".docx", docx{})
	Register(".pptx", pptx{})
}
