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
	// DefaultTopK TopK：从检索到的文档片段中拿前 N 条
	DefaultTopK = 3

	// ContextSeparator 文档片段间分隔符
	ContextSeparator = "\n---\n"

	// 系统指令：说明模型的部署背景、目标受众、维护团队、主导开发者等
	systemPrompt = `你是运行在天津城建大学私人服务器上的 Physics-LLM，基于 Deepseek 本地模型部署，
由天津城建大学理学院物理研究社研发并维护。主导开发者为 22 级应用物理学专业 1 班赵明俊。
你的使命是帮助天津城建大学范围内的本科生和研究生解答物理问题，检索并总结相关课程资料与文档，
以严谨、准确的专业语言输出。回答中必要时可引用文献、课程名称或具体章节。`

	// 用户模板：首先给出检索到的文档片段，再让模型作答
	userPromptTmpl = `以下是与用户问题相关的文档片段（已按相关度排序，最多取前 %d 条）：
%s

请基于上述内容，并结合你的物理学专业知识，详细回答下面的问题：
“%s”`
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
