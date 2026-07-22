# 빌드
FROM golang:1.24.3-bullseye AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app .

# 배포
FROM debian:bullseye-slim AS deploy

# TLS 인증서 이용
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=deploy-builder /app/app .

ENV TZ=Asia/Seoul

CMD ["./app"]