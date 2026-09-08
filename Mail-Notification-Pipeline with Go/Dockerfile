# --- 1단계: 빌드 ---
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# --- 2단계: 실행 (버전 고정 및 Non-root) ---
FROM alpine:3.19
WORKDIR /app

# 보안을 위해 appuser 생성
RUN adduser -D -u 10001 appuser
COPY --from=builder --chown=appuser:appuser /app/server .

USER appuser
EXPOSE 8080

CMD ["./server"]