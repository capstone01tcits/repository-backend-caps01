# 🔐 Go Auth API

REST API Authentication dengan Fiber + PostgreSQL + JWT

## Tech Stack
- **Framework**: Fiber v2
- **Database**: PostgreSQL + GORM
- **Auth**: JWT (Access + Refresh Token)
- **Password**: bcrypt

## Struktur Folder
```
go-auth/
├── cmd/
│   └── main.go           # Entry point
├── config/
│   └── config.go         # Config & DB connection
├── internal/
│   ├── handler/          # HTTP handlers
│   ├── middleware/        # JWT middleware
│   ├── model/            # Struct & request types
│   ├── repository/       # Database queries
│   └── service/          # Business logic
├── pkg/utils/            # JWT & response helpers
├── .env.example
├── Dockerfile
└── docker-compose.yml
```

## Setup & Run

### 1. Clone & Setup
```bash
cp .env.example .env
# Edit .env sesuai konfigurasi kamu
```

### 2. Run dengan Docker (Recommended)
```bash
docker-compose up -d
```

### 3. Run Manual
```bash
go mod tidy
go run cmd/main.go
```

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | ✗ | Register user baru |
| POST | `/api/auth/login` | ✗ | Login |
| POST | `/api/auth/refresh` | ✗ | Refresh access token |
| GET  | `/api/auth/me` | ✓ | Get profile |
| GET  | `/health` | ✗ | Health check |

## Contoh Request

### Register
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret123"}'
```

### Login
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'
```

### Get Profile (dengan token)
```bash
curl http://localhost:3000/api/auth/me \
  -H "Authorization: Bearer <access_token>"
```

### Refresh Token
```bash
curl -X POST http://localhost:3000/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

## Response Format
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 1234567890,
    "user": {
      "id": "uuid",
      "name": "John",
      "email": "john@example.com"
    }
  }
}
```
