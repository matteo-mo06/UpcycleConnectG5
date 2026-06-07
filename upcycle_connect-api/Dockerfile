FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/seed .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY db/migrations /app/db/migrations
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 8084
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./api"]
