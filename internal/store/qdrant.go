package store

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/iammm0/physics-llm/internal/config"
)

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
	// 构建请求 URL 和 Body
	url := fmt.Sprintf("/collections/%s/points/query", c.collection)
	body := map[string]interface{}{
		"query":        vector,
		"limit":        topK,
		"with_payload": true,
	}

	// 定义响应结构（只关心 payload 字段）
	var resp struct {
		Result []struct {
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
	}

	// 发起请求并解析
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

	// 提取文本字段
	var results []string
	for _, item := range resp.Result {
		if txt, ok := item.Payload["text"].(string); ok {
			results = append(results, txt)
		}
	}
	return results, nil
}
