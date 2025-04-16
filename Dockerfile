FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /ollama-cors-proxy

FROM alpine:latest
WORKDIR /
COPY --from=builder /ollama-cors-proxy /ollama-cors-proxy

EXPOSE ${PROXY_PORT}
ENTRYPOINT ["/ollama-cors-proxy"]
