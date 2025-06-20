package ollama

import (
	"context"
	"github.com/ollama/ollama/api"
)

type Client struct{ *api.Client }

func New(base string) *Client {
	return &Client{Client: &api.Client{Addr: base}}
}

// 同步问答（MVP 版本）
func (c *Client) Ask(ctx context.Context, model, prompt string) (string, error) {
	resp, err := c.Client.Chat(ctx, api.ChatRequest{
		Model: model,
		Messages: []api.Message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	})
	if err != nil {
		return "", err
	}
	return resp.Message.Content, nil
}
