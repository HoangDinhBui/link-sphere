FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ../../pkg /app/pkg
COPY services/comment-service /app/services/comment-service
WORKDIR /app/services/comment-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/comment-service ./cmd/...

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/comment-service /usr/local/bin/comment-service
EXPOSE 8004
CMD ["comment-service"]
