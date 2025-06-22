package ingest

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/iammm0/physics-llm/internal/config"
	"github.com/iammm0/physics-llm/internal/ollama"
	"github.com/iammm0/physics-llm/internal/store"
	"github.com/ledongthuc/pdf"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// extractText 根据后缀读取文本，支持 .txt/.md/.pdf
func extractText(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".txt", ".md":
		b, err := os.ReadFile(path)
		return string(b), err
	case ".pdf":
		return extractFromPDF(path)
	default:
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}
}

// extractFromPDF 用 ledongthuc/pdf 库提取整份 PDF 的文本
func extractFromPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("打开 PDF 失败: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	totalPages := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		page := r.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("读取第 %d 页文本失败: %w", pageIndex, err)
		}
		sb.WriteString(content)
		sb.WriteString("\n\n")
	}
	return sb.String(), nil
}

// chunkText 按指定长度 + 重叠切分文本
func chunkText(text string, size, overlap int) []string {
	var chunks []string
	for start := 0; start < len(text); start += size - overlap {
		end := start + size
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, strings.TrimSpace(text[start:end]))
		if end == len(text) {
			break
		}
	}
	return chunks
}

// Run 扫描 cfg.KnowledgeDir 下所有文件，切片、Embedding 并 Upsert 到 Qdrant
func Run(ctx context.Context, cfg *config.Config) error {
	llmClient := ollama.NewClient(cfg)
	dbClient := store.NewClient(cfg)

	// 1. 列出所有知识文件
	pattern := filepath.Join(cfg.KnowledgeDir, "*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("扫描知识库目录失败: %w", err)
	}
	log.Printf("发现 %d 个知识文件\n", len(files))

	// 2. 对每个文件处理
	for _, file := range files {
		text, err := extractText(file)
		if err != nil {
			log.Printf("跳过 %s: %v\n", file, err)
			continue
		}

		// 3. 文本切片
		chunks := chunkText(text, cfg.ChunkSize, cfg.ChunkOverlap)
		log.Printf("文件 %s 切成 %d 段\n", filepath.Base(file), len(chunks))

		// 4. Embedding + 构造 Point
		var points []store.Point
		for idx, chunk := range chunks {
			vec, err := llmClient.Embeddings(chunk)
			if err != nil {
				return fmt.Errorf("生成 Embedding 失败 (%s 段 %d): %w", file, idx, err)
			}
			id := uuid.New().String()
			points = append(points, store.Point{
				ID:     id,
				Vector: vec,
				Payload: map[string]interface{}{
					"text":   chunk,
					"source": filepath.Base(file),
					"index":  idx,
				},
			})
		}

		// 5. 批量 Upsert
		if err := dbClient.Upsert(ctx, points); err != nil {
			return fmt.Errorf("upsert 到 Qdrant 失败 (%s): %w", file, err)
		}
	}

	log.Println("知识库导入完成")
	return nil
}
