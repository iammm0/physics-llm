package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/iammm0/physics-llm/internal/config"
)

// Point 表示要写入 Qdrant 的单个向量点
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

// Client 用于通过 Qdrant 的 HTTP API 做向量检索
type Client struct {
	client     *resty.Client
	collection string
}

// NewClient 初始化 Resty 客户端，BaseURL 即 cfg.QdrantURL（例如 "http://localhost:6333"）
func NewClient(cfg *config.Config) *Client {
	cli := resty.New().
		SetBaseURL(cfg.QdrantURL). // e.g. "http://localhost:6333"
		SetHeader("Content-Type", "application/json")
	return &Client{
		client:     cli,
		collection: cfg.QdrantCol, // e.g. "physics"
	}
}

// Search 调用 Qdrant 的 /collections/{collection}/points/query 接口，返回 payload["text"]
func (c *Client) Search(ctx context.Context, vector []float32, topK int) ([]string, error) {
	// 调用 /collections/{col}/points/query
	url := fmt.Sprintf("/collections/%s/points/query", c.collection)
	body := map[string]interface{}{
		"query":        vector,
		"limit":        topK,
		"with_payload": true,
	}

	var resp struct {
		Result struct {
			Points []struct {
				Payload map[string]interface{} `json:"payload"`
			} `json:"points"`
		} `json:"result"`
	}

	r, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&resp).
		Post(url)
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		return nil, fmt.Errorf("qdrant search error: %s", r.Status())
	}

	var texts []string
	for _, pt := range resp.Result.Points {
		if txt, ok := pt.Payload["text"].(string); ok {
			texts = append(texts, txt)
		}
	}
	return texts, nil
}

// EnsureCollection internal/store/qdrant.go  片段
func (c *Client) EnsureCollection(dim int) error {
	url := fmt.Sprintf("/collections/%s", c.collection)

	// 先尝试 GET
	r, err := c.client.R().Get(url)
	if err == nil && r.StatusCode() == http.StatusOK {
		return nil // 已存在
	}

	// 不存在就创建
	body := map[string]any{
		"vectors": map[string]any{
			"size":     dim,
			"distance": "Cosine",
		},
	}
	r, err = c.client.R().
		SetBody(body).
		Put(url)
	if err != nil {
		return err
	}
	if r.IsError() {
		return fmt.Errorf("create collection: %s", r.Status())
	}
	return nil
}

// Upsert 批量写入或更新向量点到 Qdrant（注意：用 PUT）
func (c *Client) Upsert(ctx context.Context, points []Point) error {
	// Endpoint 必须是 PUT /collections/{col}/points
	url := fmt.Sprintf("/collections/%s/points", c.collection)
	body := map[string]interface{}{"points": points}

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		Put(url) // ← 这里改成 Put
	if err != nil {
		return fmt.Errorf("qdrant upsert request failed: %w", err)
	}
	if resp.IsError() {
		// 打印一下返回体便于调试
		return fmt.Errorf("qdrant upsert error: %s — %s", resp.Status(), resp.String())
	}
	return nil
}
