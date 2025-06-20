package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourname/physics-llm/internal/ollama"
)

type ChatHandler struct {
	Ollama *ollama.Client
	Model  string
}

type chatReq struct {
	Query string `json:"query" binding:"required"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req chatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
	defer cancel()

	answer, err := h.Ollama.Ask(ctx, h.Model, req.Query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"answer": answer})
}
