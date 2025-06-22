# ---- 构建阶段 ----
FROM golang:1.22-alpine AS builder

# 安装必要证书以支持 HTTPS（Resty 调用 Qdrant HTTP API）
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# 只复制 go.mod 和 go.sum，加速依赖下载
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码并编译
COPY . .
# 关闭 CGO，编译静态二进制，输出到 physics-llm
RUN CGO_ENABLED=0 GOOS=linux go build -o physics-llm ./cmd/api

# ---- 运行阶段 ----
FROM alpine:latest

# 安装证书
RUN apk add --no-cache ca-certificates

WORKDIR /app


# 从构建阶段拷贝可执行文件
COPY --from=builder /app/physics-llm .

# 暴露服务监听端口（与 API_ADDR 对应，默认 :8080）
EXPOSE 8080

# 启动程序
ENTRYPOINT ["./physics-llm"]