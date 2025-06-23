package extractor

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/ledongthuc/pdf"
)

// textPDF：普通 PDF 文本提取
type textPDF struct{}

func (textPDF) Extract(p string) (string, error) {
	f, reader, err := pdf.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// scanPDF：扫描 PDF 调用 tesseract CLI 做 OCR
type scanPDF struct{}

func (scanPDF) Extract(p string) (string, error) {
	// -l eng+chi_sim 根据需要切换语言
	out, err := exec.Command("tesseract", p, "stdout", "-l", "eng+chi_sim").CombinedOutput()
	return string(out), err
}

func init() {
	Register(".pdf", textPDF{})
	Register(".scan.pdf", scanPDF{})
}
