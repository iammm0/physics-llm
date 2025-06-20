# Physics-LLM

> 在私人物理服务器上部署的 **离线物理专业大模型 + 私有知识库**。  
> - **Ollama** 提供 LLM 与 Embeddings  
> - **Qdrant** 存储向量检索  
> - **Go (Gin)** 暴露 REST / SSE API  
> - 可选前端（Next.js 或原生 HTMX）

---

## ✨ 功能速览

| 功能 | 描述 |
|------|------|
| **LLM 对话** | `/v1/chat` 单接口即可与物理专业模型交流 |
| **私有知识库 (RAG)** | 支持 PDF / LaTeX / Markdown 一键向量化并检索（后续 CLI `ingest`） |
| **纯 Go 服务** | 自动加载环境变量，二进制可直接 `docker run` 或裸机跑 |
| **GPU 加速** | 全栈兼容 NVIDIA CUDA；显存越大可用的模型越大 |
| **可扩展** | 支持多模型、WebSocket 流式输出、Prometheus 监控 |

---

## 目录结构

```
physics-llm/
├─ build/               # Dockerfiles & compose
├─ cmd/
│  ├─ api/              # HTTP 主入口 (main.go)
│  └─ ingest/           # ⬆️ 待实现：知识库导入 CLI
├─ internal/
│  ├─ config/           # 读取 env / .env
│  ├─ handler/          # Gin 路由 & 业务
│  ├─ ollama/           # Ollama SDK 封装
│  ├─ store/            # Qdrant gRPC 客户端
│  └─ rag/              # ⬆️ 待实现：检索增强生成
├─ web/                 # ⬆️ 可选前端（Next.js / 静态）
└─ README.md
```

---

## 核心代码（MVP）

1. **`go.mod`**

```go
module github.com/yourname/physics-llm

go 1.22

require (
    github.com/gin-gonic/gin v1.11.0
    github.com/joho/godotenv v1.5.1
    github.com/qdrant/go-client v0.0.26
    github.com/ollama/ollama/api v0.0.0-20250610-d0b1e3f
)
```

2. **关键包**

| 路径 | 作用 |
|------|------|
| `internal/config/config.go` | 加载 `API_ADDR` / `OLLAMA_BASE_URL` / `QDRANT_URL` 等 |
| `internal/ollama/ollama.go` | 同步问答 / 生成 Embeddings |
| `internal/store/qdrant.go`  | gRPC 直连 Qdrant |
| `internal/handler/chat.go`  | `POST /v1/chat` 路由 |
| `cmd/api/main.go`           | 监听、优雅关机、日志 |

3. **运行指令**

```bash
# 拉依赖
go mod tidy

# 导出环境变量（裸机调试示例）
export API_ADDR=":8080"
export OLLAMA_BASE_URL="http://localhost:11434"
export QDRANT_URL="localhost:6333"
export OLLAMA_MODEL="physics-phi"
export QDRANT_COLLECTION="physics"

# 启动服务
go run ./cmd/api
```

4. **测试请求**

```bash
curl -X POST http://localhost:8080/v1/chat      -H "Content-Type: application/json"      -d '{"query":"请解释玻色–爱因斯坦凝聚"}'
```

---

## 一键部署 (Docker Compose)

```yaml
version: "3.9"
services:
  ollama:
    image: ollama/ollama:latest-cuda
    volumes:
      - /opt/ollama:/root/.ollama
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]
    ports: ["3000:11434"]

  qdrant:
    image: qdrant/qdrant:v1.9.1
    volumes:
      - /opt/vector-db:/qdrant/storage
    ports: ["6333:6333"]

  api:
    build: ./build/api
    environment:
      - OLLAMA_BASE_URL=http://ollama:11434
      - QDRANT_URL=http://qdrant:6333
      - OLLAMA_MODEL=physics-phi
    depends_on: [ollama, qdrant]
    ports: ["8080:8080"]
```

```bash
docker compose pull           # 首次拉取镜像
docker compose up -d          # 后台启动
docker compose exec ollama ollama pull physics-phi   # 拉物理模型
```

---

## 环境变量一览

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `API_ADDR` | `:8080` | API 监听地址 |
| `OLLAMA_BASE_URL` | `http://localhost:11434` | Ollama 服务地址 |
| `OLLAMA_MODEL` | `physics-phi` | 默认推理模型 |
| `QDRANT_URL` | `localhost:6333` | Qdrant gRPC 端口 |
| `QDRANT_COLLECTION` | `physics` | 向量库集合名 |

可在生产中写入 `.env`，开发时自动读取（由 `godotenv` 实现）。

---

## 后续路线图

| 阶段 | 内容 |
|------|------|
| **📚 Ingest CLI** | `cmd/ingest`：PDF → 分块 → Embeddings → Qdrant |
| **🔍 RAG Pipeline** | `internal/rag`：搜索 top-k + Prompt 拼接 |
| **🚰 SSE / WS** | 改造 `/v1/chat` 支持流式推送 |
| **🪝 Prometheus** | `/metrics` 暴露 Go & Ollama 指标 |
| **🔐 Auth** | JWT / OIDC 中间件，支持多用户隔离 |
| **🖥️ 前端** | Next.js + tRPC or HTMX + Tailwind |

---

## 贡献与反馈

- 提交 PR 前请运行 `go test ./...`  
- 如需新增特性或修改目录，请先在 Issue 中沟通  
- 发现物理公式回答不完整？  
  - 检查知识库是否已包含该章节 PDF  
  - 提升检索 `topK` 或更新 Prompt 模板  

---

> **License**  
> MIT © 2025 赵明俊
