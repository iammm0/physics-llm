// cmd/api/main.go

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iammm0/physics-llm/internal/config"
	"github.com/iammm0/physics-llm/internal/handler"
	"github.com/iammm0/physics-llm/internal/ingest"
	"github.com/iammm0/physics-llm/internal/store"
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()

	// 2. 初始化 Qdrant 客户端并确保 collection 存在
	db := store.NewClient(cfg)
	if err := db.EnsureCollection(cfg.EmbedDim); err != nil {
		log.Fatalf("qdrant 初始化失败: %v", err)
	}

	// 3. 批量导入知识库文件到 Qdrant
	//    10 分钟超时，导入过程中会自动分片、生成 embedding 并 upsert
	ingestCtx, ingestCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer ingestCancel()
	if err := ingest.Run(ingestCtx, cfg); err != nil {
		log.Fatalf("文档导入失败: %v", err)
	}

	// 4. 设置 Gin 路由
	router := gin.Default()
	handler.RegisterRoutes(router, cfg)

	// 5. 启动 HTTP 服务
	srv := &http.Server{
		Addr:    cfg.APIAddr,
		Handler: router,
	}

	go func() {
		log.Printf("开始监听 %s ...", cfg.APIAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	// 6. 捕获系统信号，优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("接收到关闭信号，正在优雅退出...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("服务器优雅关闭失败: %v", err)
	}

	log.Println("服务器已退出")
}
