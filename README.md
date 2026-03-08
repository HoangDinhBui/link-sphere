# 🔗 LinkSphere — Social Network API

A polyglot microservices-based social network API built with **Go** and **Python (FastAPI)**.

## 🏗 Architecture

```
                    ┌─────────────────┐
                    │   API Gateway   │
                    │    (NGINX)      │
                    └────────┬────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
  ┌─────┴─────┐   ┌────────┴────────┐   ┌──────┴──────┐
  │   User    │   │     Post       │   │   Comment   │
  │  Service  │   │    Service     │   │   Service   │
  │   (Go)    │   │     (Go)       │   │    (Go)     │
  └───────────┘   └────────────────┘   └─────────────┘
        │                    │                    │
  ┌─────┴─────┐   ┌────────┴────────┐   ┌──────┴──────┐
  │   Auth    │   │     Feed       │   │ Notification│
  │  Service  │   │    Service     │   │   Service   │
  │   (Go)    │   │     (Go)       │   │  (Python)   │
  └───────────┘   └────────────────┘   └─────────────┘
                                              │
                  ┌─────────────────┐   ┌─────┴───────┐
                  │     Search     │   │    Kafka    │
                  │    Service     │   │   Broker    │
                  │   (Python)     │   └─────────────┘
                  └────────┬───────┘
                           │
                  ┌────────┴────────┐
                  │   OpenSearch   │
                  └────────────────┘
```

## 🛠 Tech Stack

| Layer               | Technology                           |
| ------------------- | ------------------------------------ |
| **Go Services**     | Chi router, sqlx, go-redis, kafka-go |
| **Python Services** | FastAPI, aiokafka, opensearch-py     |
| **Database**        | PostgreSQL                           |
| **Cache**           | Redis                                |
| **Message Broker**  | Apache Kafka                         |
| **Search Engine**   | OpenSearch                           |
| **API Gateway**     | NGINX                                |
| **Container**       | Docker + Kubernetes                  |
| **CI/CD**           | GitHub Actions                       |

## 📁 Services

| Service              | Language | Port | Description                            |
| -------------------- | -------- | ---- | -------------------------------------- |
| User Service         | Go       | 8001 | Registration, profile, follow/unfollow |
| Auth Service         | Go       | 8002 | Login, JWT token                       |
| Post Service         | Go       | 8003 | Create/list/like posts                 |
| Comment Service      | Go       | 8004 | Create/list comments                   |
| Feed Service         | Go       | 8005 | News feed generation                   |
| Notification Service | Python   | 8006 | Realtime notifications (WebSocket)     |
| Search Service       | Python   | 8007 | Full-text search via OpenSearch        |

## 🚀 Getting Started

### Prerequisites

- Go 1.22+
- Python 3.12+
- Docker & Docker Compose

### Run with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Run locally (development)

```bash
# Start infrastructure only
make dev-infra

# Build Go services
make build

# Run individual service
./bin/user-service
```

### Run tests

```bash
make test
```

## 📬 API Endpoints

All endpoints are accessed through the API Gateway at `http://localhost:80`.

| Endpoint                     | Method | Service      | Description        |
| ---------------------------- | ------ | ------------ | ------------------ |
| `/api/v1/users/register`     | POST   | User         | Register account   |
| `/api/v1/auth/login`         | POST   | Auth         | Login              |
| `/api/v1/users/follow`       | POST   | User         | Follow user        |
| `/api/v1/users/unfollow`     | POST   | User         | Unfollow user      |
| `/api/v1/posts`              | POST   | Post         | Create post        |
| `/api/v1/posts`              | GET    | Post         | List posts         |
| `/api/v1/posts/{id}/like`    | POST   | Post         | Like post          |
| `/api/v1/posts/comment`      | POST   | Comment      | Create comment     |
| `/api/v1/posts/comments`     | POST   | Comment      | List comments      |
| `/api/v1/feed/get`           | POST   | Feed         | Get news feed      |
| `/api/v1/notifications/list` | POST   | Notification | List notifications |
| `/api/v1/search/posts`       | POST   | Search       | Search posts       |

## 📄 License

MIT
