# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./
COPY static/ ./static/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /calculator .

# Stage 2: Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /calculator .

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./calculator"]

