package main

import (
	"context"
	"github.com/yourname/physics-llm/internal/config"
	"github.com/yourname/physics-llm/internal/handler"
	"github.com/yourname/physics-llm/internal/ollama"
	"github.com/yourname/physics-llm/internal/store"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 初始化客户端
	ollamaClient := ollama.New(cfg.OllamaURL)
	qdrantClient := store.New(ctx, cfg.QdrantURL) // 目前未使用检索，可先保留
	defer qdrantClient.Close()

	// HTTP server
	r := gin.New()
	r.Use(gin.Recovery())

	h := &handler.ChatHandler{
		Ollama: ollamaClient,
		Model:  cfg.Model,
	}

	r.POST("/v1/chat", h.Chat)

	log.Printf("Listening on %s ...", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("server exit: %v", err)
	}
}
