module github.com/linksphere/feed-service

go 1.22.0

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/linksphere/pkg v0.0.0
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/linksphere/pkg => ../../pkg
