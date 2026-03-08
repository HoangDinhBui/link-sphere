FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ../../pkg /app/pkg
COPY services/user-service /app/services/user-service
WORKDIR /app/services/user-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user-service ./cmd/...

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/user-service /usr/local/bin/user-service
EXPOSE 8001
CMD ["user-service"]
