.PHONY: build test clean docker-up docker-down docker-build lint

# ========================
# Go Services
# ========================
GO_SERVICES = user-service auth-service post-service comment-service feed-service

build:
	@echo "Building all Go services..."
	@for service in $(GO_SERVICES); do \
		echo "  Building $$service..."; \
		cd services/$$service && go build -o ../../bin/$$service ./cmd/... && cd ../..; \
	done
	@echo "Done."

test:
	@echo "Running tests..."
	@for service in $(GO_SERVICES); do \
		echo "  Testing $$service..."; \
		cd services/$$service && go test ./... && cd ../..; \
	done
	@echo "Done."

lint:
	@echo "Linting Go services..."
	@for service in $(GO_SERVICES); do \
		cd services/$$service && golangci-lint run ./... && cd ../..; \
	done

clean:
	@rm -rf bin/
	@echo "Cleaned."

# ========================
# Docker
# ========================
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# ========================
# Development
# ========================
dev-infra:
	docker-compose up -d postgres redis kafka opensearch

go-work-sync:
	go work sync
