# go-service/Dockerfile
FROM golang:latest AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
