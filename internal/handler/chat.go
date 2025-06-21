package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iammm0/physics-llm/internal/config"
	"github.com/iammm0/physics-llm/internal/ollama"
	"github.com/iammm0/physics-llm/internal/store"
)

const (
	DefaultTopK      = 3
	ContextSeparator = "\n---\n"
	systemPrompt     = "你是经验丰富的物理学家，请使用严谨、准确的语言回答用户问题。"
	userPromptTmpl   = "以下是与用户问题相关的文档片段：\n%s\n\n请基于这些内容回答：%s"
)

type ChatRequest struct {
	Query string `json:"query" binding:"required"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// RegisterRoutes 挂载 /v1/chat
func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	llm := ollama.NewClient(cfg)
	db := store.NewClient(cfg)

	r.POST("/v1/chat", func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 超时控制
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		// 1) 生成用户 Query 的向量
		vec, err := llm.Embeddings(req.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Embedding 失败: " + err.Error()})
			return
		}

		// 2) 检索 topK 文档片段
		docs, err := db.Search(ctx, vec, DefaultTopK)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "检索文档失败: " + err.Error()})
			return
		}

		// 3) 组装用户 prompt
		combined := strings.Join(docs, ContextSeparator)
		userPrompt := fmt.Sprintf(userPromptTmpl, combined, req.Query)

		// 4) 调用 Ollama
		answer, err := llm.Complete(userPrompt, systemPrompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "调用模型失败: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, ChatResponse{Response: answer})
	})
}
