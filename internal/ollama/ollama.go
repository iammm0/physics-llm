package ollama

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/iammm0/physics-llm/internal/config"
)

type Client struct {
	cli        *resty.Client
	model      string
	embedModel string // embeddings
}

// ChatMessage 与 Ollama /api/chat JSON 保持一致
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewClient === 对外构造器 ===
func NewClient(cfg *config.Config) *Client {
	c := resty.New().
		SetBaseURL(cfg.OllamaURL).  // 例: http://localhost:11434
		SetTimeout(60*time.Second). // 可按需调整
		SetHeader("Content-Type", "application/json")

	return &Client{
		cli:        c,
		model:      cfg.OllamaModel,
		embedModel: cfg.OllamaEmbedModel, // 新增字段
	}
}

/*
Complete 发送聊天请求，返回 assistant 的 content

prompt —— 用户问题，自动封装为 `{"role":"user", ...}`
system —— 可选系统提示词；留空则不发送 system 消息
*/
func (c *Client) Complete(prompt string, system string) (string, error) {
	var msgs []ChatMessage

	if system != "" {
		msgs = append(msgs, ChatMessage{Role: "system", Content: system})
	}
	msgs = append(msgs, ChatMessage{Role: "user", Content: prompt})

	reqBody := map[string]interface{}{
		"model":    c.model,
		"messages": msgs,
		"stream":   false,
	}

	var resp struct {
		Message ChatMessage `json:"message"` // 只关心 assistant 最终回复
	}

	r, err := c.cli.R().
		SetBody(reqBody).
		SetResult(&resp).
		Post("/api/chat")
	if err != nil {
		return "", err
	}
	if r.IsError() {
		return "", fmt.Errorf("ollama chat error: %s", r.Status())
	}
	return resp.Message.Content, nil
}

// Embeddings 调 /api/embeddings，返回 float32 切片
func (c *Client) Embeddings(text string) ([]float32, error) {
	reqBody := map[string]string{
		"model":  "mxbai-embed-large", // 你可写到 cfg 里
		"prompt": text,
	}

	var resp struct {
		Embedding []float32 `json:"embedding"`
	}

	r, err := c.cli.R().
		SetBody(reqBody).
		SetResult(&resp).
		Post("/api/embeddings")
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		return nil, fmt.Errorf("ollama embeddings error: %s", r.Status())
	}
	return resp.Embedding, nil
}
