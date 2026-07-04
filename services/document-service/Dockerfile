FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/document-service ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates curl
WORKDIR /app
COPY --from=builder /out/document-service .
EXPOSE 8082
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8082/health || exit 1
ENTRYPOINT ["./document-service"]
