FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ../../pkg /app/pkg
COPY services/auth-service /app/services/auth-service
WORKDIR /app/services/auth-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service ./cmd/...

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/auth-service /usr/local/bin/auth-service
EXPOSE 8002
CMD ["auth-service"]
