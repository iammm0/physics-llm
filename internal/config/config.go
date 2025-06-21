package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	APIAddr          string
	OllamaURL        string
	OllamaModel      string
	OllamaEmbedModel string
	QdrantURL        string
	QdrantCol        string
	EmbedDim         int
	DocsDir          string
	KnowledgeDir     string
	ChunkSize        int
	ChunkOverlap     int
}

func LoadConfig() *Config {
	// 加载 .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("API_ADDR", ":8080")
	viper.SetDefault("OLLAMA_BASE_URL", "http://localhost:11434")
	viper.SetDefault("OLLAMA_MODEL", "deepseek-r1:14b")
	viper.SetDefault("OLLAMA_EMBED_MODEL", "mxbai-embed-large")
	viper.SetDefault("QDRANT_URL", "http://localhost:6333")
	viper.SetDefault("QDRANT_COLLECTION", "physics")
	viper.SetDefault("EMBED_DIM", 1024)
	viper.SetDefault("KNOWLEDGE_DIR", "./knowledge")
	viper.SetDefault("DOCS_DIR", "./docs")
	viper.SetDefault("CHUNK_SIZE", 500)
	viper.SetDefault("CHUNK_OVERLAP", 50)

	return &Config{
		APIAddr:          viper.GetString("API_ADDR"),
		OllamaURL:        viper.GetString("OLLAMA_BASE_URL"),
		OllamaModel:      viper.GetString("OLLAMA_MODEL"),
		OllamaEmbedModel: viper.GetString("OLLAMA_EMBED_MODEL"),
		QdrantURL:        viper.GetString("QDRANT_URL"),
		QdrantCol:        viper.GetString("QDRANT_COLLECTION"),
		EmbedDim:         viper.GetInt("EMBED_DIM"),
		DocsDir:          viper.GetString("DOCS_DIR"),
		KnowledgeDir:     viper.GetString("KNOWLEDGE_DIR"),
		ChunkSize:        viper.GetInt("CHUNK_SIZE"),
		ChunkOverlap:     viper.GetInt("CHUNK_OVERLAP"),
	}
}
