package ollama

import (
	"github.com/go-resty/resty/v2"
	"github.com/iammm0/physics-llm/internal/config"
)

type Client struct {
	rest  *resty.Client
	model string
}

func NewClient(cfg *config.Config) *Client {
	client := resty.New().SetBaseURL(cfg.OllamaURL)
	return &Client{rest: client, model: cfg.OllamaModel}
}

// Complete 调用 Ollama Completion API
func (c *Client) Complete(prompt string) (string, error) {
	var resp struct {
		Completion string `json:"completion"`
	}
	_, err := c.rest.R().
		SetBody(map[string]interface{}{"model": c.model, "prompt": prompt}).
		SetResult(&resp).
		Post("/api/completions")
	if err != nil {
		return "", err
	}
	return resp.Completion, nil
}

// Embeddings 调用 Ollama Embedding API
func (c *Client) Embeddings(text string) ([]float32, error) {
	var resp struct {
		Embeddings []float32 `json:"embeddings"`
	}
	_, err := c.rest.R().
		SetBody(map[string]interface{}{"model": c.model, "input": text}).
		SetResult(&resp).
		Post("/api/embeddings")
	if err != nil {
		return nil, err
	}
	return resp.Embeddings, nil
}
