package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	APIAddr     string
	OllamaURL   string
	OllamaModel string
	QdrantURL   string
	QdrantCol   string
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
	viper.SetDefault("OLLAMA_MODEL", "physics-phi")
	viper.SetDefault("QDRANT_URL", "http://localhost:6333")
	viper.SetDefault("QDRANT_COLLECTION", "physics")

	return &Config{
		APIAddr:     viper.GetString("API_ADDR"),
		OllamaURL:   viper.GetString("OLLAMA_BASE_URL"),
		OllamaModel: viper.GetString("OLLAMA_MODEL"),
		QdrantURL:   viper.GetString("QDRANT_URL"),
		QdrantCol:   viper.GetString("QDRANT_COLLECTION"),
	}
}
