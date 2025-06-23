package ingest

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/iammm0/physics-llm/internal/config"
	"github.com/iammm0/physics-llm/internal/ingest/extractor"
	"github.com/iammm0/physics-llm/internal/ollama"
	"github.com/iammm0/physics-llm/internal/store"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func extractText(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ex, ok := extractor.Get(ext); ok {
		return ex.Extract(path)
	}
	// 回退：当普通 UTF-8 文本
	b, err := os.ReadFile(path)
	return string(b), err
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
