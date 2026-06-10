# User Management API (Go)
 
A RESTful API built with Go for managing users, using MongoDB for persistence and JWT for authentication.
 
## Tech Stack
 
- **Go 1.26** — standard `net/http` with Go 1.26 routing (`r.PathValue`)
- **MongoDB** — official Go driver (`go.mongodb.org/mongo-driver`)
- **JWT** — `golang-jwt/jwt` with HS256 signing
- **Docker** — multi-stage build + docker-compose for API and MongoDB
## Project Structure
 
```
backend-go/
├── cmd/api/
│   └── main.go                  # Entry point, server setup, graceful shutdown
├── internal/
│   ├── handler/
│   │   └── user_handler.go      # HTTP handlers
│   ├── middleware/
│   │   └── middleware.go        # JWT auth + request logger
│   ├── model/
│   │   └── user.go              # User struct and request/response types
│   ├── repository/
│   │   └── user_repository.go   # UserRepository interface + MongoDB adapter
│   └── service/
│       ├── user_service.go      # Business logic
│       └── user_service_test.go # Unit tests with mocked repository
├── .env.example
├── docker-compose.yml
├── Dockerfile
└── go.mod
```
 
### Architecture decisions
 
- **Hexagonal-inspired** — `repository` defines an interface (port); `mongoUserRepository` is the adapter. Business logic in `service` depends only on the interface, not on MongoDB directly — making it easy to mock in tests.
- **Standard library routing** — Go 1.26 added path parameters (`{id}`) to `net/http`, so no external router is needed.
- **Background goroutine** — logs total user count every 10 seconds using `time.Ticker`.
- **Graceful shutdown** — listens for `SIGINT`/`SIGTERM` and gives in-flight requests 5 seconds to finish.
---
 
## Getting Started
 
### Prerequisites
 
- Go >= 1.26
- MongoDB running locally **or** Docker + Docker Compose
### Option A — Run with Docker (recommended)
 
```bash
cp .env.example .env          # set JWT_SECRET
docker-compose up --build
```
 
API is available at `http://localhost:8080`.
 
### Option B — Run locally
 
```bash
cp .env.example .env          # set MONGO_URI and JWT_SECRET
 
# load env vars, then run
export $(cat .env | xargs)
go run ./cmd/api
```
 
---
 
## Environment Variables
 
| Variable | Default | Description |
|----------|---------|-------------|
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `JWT_SECRET` | `changeme` | Secret key for signing JWT tokens |
 
---
 
## API Endpoints
 
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/auth/register` | — | Register a new user |
| `POST` | `/auth/login` | — | Login and receive a JWT token |
| `GET` | `/users` | ✅ | List all users |
| `GET` | `/users/{id}` | ✅ | Get a user by ID |
| `PUT` | `/users/{id}` | ✅ | Update a user's name or email |
| `DELETE` | `/users/{id}` | ✅ | Delete a user |
 
---
 
## JWT Guide
 
### 1. Register a user
 
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com", "password": "secret123"}'
```
 
Response:
```json
{
  "id": "665f1a2b3c4d5e6f7a8b9c0d",
  "name": "Alice",
  "email": "alice@example.com",
  "createdAt": "2024-06-01T10:00:00Z"
}
```
 
### 2. Login and get a token
 
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "password": "secret123"}'
```
 
Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
 
The token is a signed **HS256 JWT** containing:
 
| Claim | Value |
|-------|-------|
| `sub` | User ID (MongoDB ObjectID as string) |
| `iat` | Issued at (Unix timestamp) |
| `exp` | Expires at (24 hours from issue) |
 
### 3. Use the token on protected endpoints
 
Pass the token in the `Authorization` header as a Bearer token:
 
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
 
# List all users
curl http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN"
 
# Get user by ID
curl http://localhost:8080/users/665f1a2b3c4d5e6f7a8b9c0d \
  -H "Authorization: Bearer $TOKEN"
 
# Update user
curl -X PUT http://localhost:8080/users/665f1a2b3c4d5e6f7a8b9c0d \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Smith"}'
 
# Delete user
curl -X DELETE http://localhost:8080/users/665f1a2b3c4d5e6f7a8b9c0d \
  -H "Authorization: Bearer $TOKEN"
```
 
---
 
## Running Tests
 
```bash
go test ./...
```
 
Tests use a **mock repository** (`mockUserRepository`) that implements the `UserRepository` interface in-memory — no MongoDB connection required.
 
Test cases cover:
 
- `TestRegister` — user is created with a hashed password
- `TestLogin_Success` — valid credentials return a JWT token
- `TestLogin_WrongPassword` — wrong password returns an error
- `TestCount` — returns correct user count