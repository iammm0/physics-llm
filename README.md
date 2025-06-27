<table align="center">
  <tr>
    <td align="center" valign="middle" width="220">
      <img
        src="./web/src/assets/Physics-LLM.png"
        alt="Physics-LLM Logo"
      />
    </td>
    <td align="left" valign="middle">
      <h1 style="margin:0;">Physics-LLM</h1>
      <p style="margin:0;">
        <a href="https://github.com/iammm0/physics-llm/releases/tag/v0.0.1">
          <img
            src="https://img.shields.io/badge/version-v0.0.1-blue.svg"
            alt="version"
          />
        </a>
        <a href="LICENSE">
          <img
            src="https://img.shields.io/badge/license-MIT-green.svg"
            alt="license"
          />
        </a>
      </p>
    </td>
  </tr>
</table>


> **本地离线物理大模型 + 私有知识库 + 向量检索 RAG**
>
> * **Ollama**：本地 LLM（聊天 + Embedding）
> * **Qdrant**：向量数据库
> * **Go (Gin)**：REST `/v1/chat` + 自动知识库导入
> * **Vite + React (TS)**：前端聊天窗口

---

## 目录结构

```
physics-llm/
├─ cmd/
│  └─ api/                 # HTTP 服务入口 (main.go)
├─ internal/
│  ├─ config/              # 读取 .env / ENV
│  ├─ handler/             # Gin 路由 ( /v1/chat )
│  ├─ ingest/              # 启动时扫描 knowledge/ → Upsert Qdrant
│  ├─ ollama/              # Ollama REST 客户端
│  └─ store/               # Qdrant HTTP 客户端 (Search / Upsert / Ensure)
├─ knowledge/              # 放置 PDF / MD / TXT 等各种文件格式的物理资料
├─ web/               # React (TS) 前端聊天应用
├─ build-scripts/                  # Dockerfiles & compose
└─ README.md
```

---

## 环境变量（.env）

```dotenv
# API
API_ADDR=:8080

# Ollama
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=deepseek-r1:14b
OLLAMA_EMBED_MODEL=mxbai-embed-large

# Qdrant
QDRANT_URL=http://localhost:6333
QDRANT_COLLECTION=physics
EMBED_DIM=1024

# 知识库
KNOWLEDGE_DIR=./knowledge
CHUNK_SIZE=500
CHUNK_OVERLAP=50
```

---

## 启动步骤

```bash
# 1. 拉模型
ollama pull deepseek-r1:14b
ollama pull mxbai-embed-large

# 2. 本地启动 ollama 服务
ollama serve   # 默认在 http://localhost:11434

# 2. 启动 Qdrant
docker compose up -d

# 3. 访问前端
cd frontend && npm install && npm run dev   # http://localhost:5173
```

> **首启自动导入知识库**：
>
> `ingest.Run()` 会扫描 `knowledge/` 目录，将所有 PDF/DOCS/MD/TXT/RMarkDown/JSON/XML/YAML/HTML 提取文本 → 切片 → Embedding → `Upsert` 到 Qdrant。若文件更新，重启服务即可增量导入。

---

## API 快速测试

```bash
curl -H 'Content-Type: application/json' \
     -d '{"query":"解释量子隧穿"}' \
     http://localhost:8080/v1/chat
```

返回示例：

```json
{"response":"量子隧穿是一种…"}
```

---

## 示例关键实现

| 位置                          | 说明                                                        |
| --------------------------- | --------------------------------------------------------- |
| `internal/ingest/ingest.go` | 提取文本（PDF: `ledongthuc/pdf`），切片、生成 UUID、Embedding、`Upsert` |
| `internal/store/qdrant.go`  | `EnsureCollection` + `Search` + `Upsert (PUT)`            |
| `internal/handler/chat.go`  | Embedding → Search → Prompt → Chat (stream\:false)        |
| `internal/ollama/ollama.go` | `/api/embeddings` & `/api/chat` 封装                        |

---

## 时序流程
```mermaid
sequenceDiagram
participant Frontend as React 前端 (http://localhost:5173)
participant API as Go API /v1/chat
participant Ollama as Ollama LLM (Embeddings & Chat)
participant Qdrant as Qdrant 向量库

Frontend->>API: POST /v1/chat { query }
    API->>Ollama: /api/embeddings { prompt=query }
    Ollama-->>API: [向量]
    API->>Qdrant: /collections/physics/points/query<br/>query=[向量], limit=k
    Qdrant-->>API: payload.text (top-k 片段)
    API->>Ollama: /api/chat { system+user prompt }
    Ollama-->>API: assistant answer
    API-->>Frontend: { response: answer }
```
---

## TODO

* [ ] SSE / WebSocket 流式输出
* [ ] PDF 数学公式 OCR
* [ ] 文档增量更新检测
* [ ] Prometheus /metrics
* [ ] JWT / 角色权限

---

MIT © 2025 赵明俊
