FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ../../pkg /app/pkg
COPY services/feed-service /app/services/feed-service
WORKDIR /app/services/feed-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/feed-service ./cmd/...

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/feed-service /usr/local/bin/feed-service
EXPOSE 8005
CMD ["feed-service"]
