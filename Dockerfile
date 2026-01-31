# 构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装git（用于获取依赖）
RUN apk add --no-cache git

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ai-bridge ./cmd/server

# 运行阶段
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/ai-bridge .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./ai-bridge"]
