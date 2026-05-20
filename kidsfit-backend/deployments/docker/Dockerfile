# ============================================================
# KidsFit 后端服务多阶段构建 Dockerfile
# 支持通过 BUILD_ARG 指定要构建的服务名
# 使用方式: docker build --build-arg SERVICE=user-svc -t kidsfit-user-svc .
# ============================================================

# -------------------- 第一阶段：编译 --------------------
FROM golang:1.22-alpine AS builder

# 安装编译依赖（git用于拉取go模块，ca-certificates用于HTTPS请求）
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 先复制依赖文件，利用Docker缓存层优化构建速度
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建参数：指定要编译的服务名（user-svc / training-svc / vision-svc / reward-svc / gateway）
ARG SERVICE=user-svc

# 编译指定服务，启用静态链接和裁剪符号表以减小镜像体积
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /app/bin/${SERVICE} \
    github.com/kidsfit/api/cmd/${SERVICE}

# -------------------- 第二阶段：运行 --------------------
FROM alpine:3.19

# 安装运行时依赖（ca-certificates用于HTTPS，tzdata用于时区，curl用于健康检查）
RUN apk add --no-cache ca-certificates tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 创建非root用户运行服务，提升安全性
RUN addgroup -S kidsfit && adduser -S kidsfit -G kidsfit

# 设置工作目录
WORKDIR /app

# 从编译阶段复制二进制文件
ARG SERVICE=user-svc
COPY --from=builder /app/bin/${SERVICE} /app/server

# 复制配置文件目录
COPY configs/ /app/configs/

# 修改文件属主为非root用户
RUN chown -R kidsfit:kidsfit /app

# 切换到非root用户
USER kidsfit

# 暴露服务端口（各服务端口：user-svc:8001, training-svc:8002, vision-svc:8003, reward-svc:8004, gateway:8080）
EXPOSE 8001

# 健康检查：每30秒检查一次服务是否存活，超时5秒，重试3次
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8001/health || exit 1

# 启动服务
ENTRYPOINT ["/app/server"]
