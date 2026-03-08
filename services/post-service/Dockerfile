FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ../../pkg /app/pkg
COPY services/post-service /app/services/post-service
WORKDIR /app/services/post-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/post-service ./cmd/...

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/post-service /usr/local/bin/post-service
EXPOSE 8003
CMD ["post-service"]
