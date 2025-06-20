package handler

import (
	"context"
	"fmt"
	"github.com/iammm0/physics-llm/internal/store"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iammm0/physics-llm/internal/config"
	"github.com/iammm0/physics-llm/internal/ollama"
)

const (
	DefaultTopK      = 3
	ContextSeparator = "\n---\n"
	PromptTemplate   = "你是物理领域专家，以下是与用户问题相关的文档片段：\n%s\n\n请结合上述内容，回答用户问题：%s"
)

type ChatRequest struct {
	Query string `json:"query" binding:"required"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	llm := ollama.NewClient(cfg)
	db := store.NewClient(cfg)

	r.POST("/v1/chat", func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()

		// 1) 对用户 query 生成向量
		vec, err := llm.Embeddings(req.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Embedding 失败: " + err.Error()})
			return
		}

		// 2) 向量检索 topK 文档片段
		docs, err := db.Search(ctx, vec, DefaultTopK)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "检索文档失败: " + err.Error()})
			return
		}

		// 3) 拼接 prompt
		combined := strings.Join(docs, ContextSeparator)
		prompt := fmt.Sprintf(PromptTemplate, combined, req.Query)

		// 4) 调用 Ollama 完成生成
		answer, err := llm.Complete(prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "调用模型失败: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, ChatResponse{Response: answer})
	})
}
