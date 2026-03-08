module github.com/linksphere/post-service

go 1.22.0

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/linksphere/pkg v0.0.0
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.47 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/linksphere/pkg => ../../pkg
