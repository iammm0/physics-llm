package config

import (
	"log"
	"os"
)

type Config struct {
	Addr       string // API 监听地址
	OllamaURL  string // http://ollama:11434
	QdrantURL  string // http://qdrant:6333
	Model      string // physics-phi
	Collection string // physics
}

func Load() *Config {
	// 可选：本地调试自动加载 .env
	_ = godotenv.Load(".env")

	cfg := &Config{
		Addr:       getEnv("API_ADDR", ":8080"),
		OllamaURL:  getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		QdrantURL:  getEnv("QDRANT_URL", "http://localhost:6333"),
		Model:      getEnv("OLLAMA_MODEL", "physics-phi"),
		Collection: getEnv("QDRANT_COLLECTION", "physics"),
	}
	log.Printf("[config] %+v", cfg)
	return cfg
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
