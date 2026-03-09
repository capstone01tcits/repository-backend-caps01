# AI Video Gen - Backend Service

Go REST API dengan authentication, PostgreSQL, JWT, dan AI service proxy untuk video generation.

## Tech Stack

**Go Backend (Port 3000)**
- Framework: Fiber v2
- Database: PostgreSQL + GORM
- Authentication: JWT (Access + Refresh Token)
- Password Hashing: bcrypt
- API Gateway: Proxy ke Python AI Service

**Python AI Service (Port 8000)**
- Framework: FastAPI + Uvicorn
- Purpose: AI processing, video generation, image processing, ML inference

## Struktur Folder

```
Sevima-BackEnd Ai Video Gen/
├── cmd/
│   └── main.go                    # Go server entry point
├── config/
│   └── config.go                  # Config loader & DB connection
├── internal/
│   ├── handler/                   # HTTP request handlers
│   │   ├── ai_handler.go         # AI proxy handler
│   │   └── auth_handler.go       # Auth handlers
│   ├── middleware/
│   │   └── auth.go               # JWT validation middleware
│   ├── model/
│   │   └── user.go               # User struct & request types
│   ├── repository/
│   │   └── user_repository.go    # Database queries
│   └── service/
│       └── auth_service.go       # Business logic
├── pkg/utils/
│   ├── jwt.go                    # JWT token generation
│   └── response.go               # Standardized response format
├── ai-service/
│   ├── main.py                   # Python AI service entry
│   └── requirements.txt          # Python dependencies
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Environment Configuration

Create `.env` file:
```bash
cp .env.example .env
```

Configure in `.env`:
```
APP_PORT=3000
APP_ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=go_auth
JWT_SECRET=your_jwt_secret
JWT_EXPIRE_HOURS=24
JWT_REFRESH_SECRET=your_refresh_secret
JWT_REFRESH_EXPIRE_HOURS=168
AI_SERVICE_URL=http://localhost:8000
```

## Setup & Run

### Option 1: Run dengan Docker (Recommended)

```bash
docker-compose up -d
```

This runs both PostgreSQL and Go service. For Python AI service:
```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

### Option 2: Run Manual

**Setup Go Backend:**
```bash
go mod tidy
go run cmd/main.go
```

Server akan berjalan di `http://localhost:3000`

**Setup Python AI Service:**
```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

Service akan berjalan di `http://localhost:8000`

## API Endpoints

### Health Check (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check Go backend status |
| GET | `/api/ai/health` | Check Python AI service connectivity |

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login and get tokens |
| POST | `/api/auth/refresh` | Refresh access token |

### Profile (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/me` | Get current user profile |

### AI Proxy (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/ai/generate` | Generate video |
| GET | `/api/ai/status` | Get AI service status |

All AI endpoints require Authorization header with valid JWT token. User context (X-User-ID, X-User-Email) is automatically injected by the proxy.

## Usage Examples

### 1. Health Check (No Auth)

```bash
curl http://localhost:3000/health
```

Response:
```json
{
  "status": "ok"
}
```

### 2. AI Health Check (No Auth)

```bash
curl http://localhost:3000/api/ai/health
```

Response (if AI service is running):
```json
{
  "success": true,
  "status": "ok",
  "message": "AI service is running"
}
```

### 3. Register User

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secure123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### 4. Login

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secure123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

### 5. Get User Profile (Protected)

```bash
curl http://localhost:3000/api/auth/me \
  -H "Authorization: Bearer {access_token}"
```

Response:
```json
{
  "success": true,
  "message": "Profile retrieved",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### 6. Refresh Access Token

```bash
curl -X POST http://localhost:3000/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "{refresh_token}"}'
```

Response:
```json
{
  "success": true,
  "message": "Token refreshed",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

### 7. AI Service - Generate Video (Protected)

```bash
curl -X POST http://localhost:3000/api/ai/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Create a video about artificial intelligence",
    "duration": 30,
    "style": "professional"
  }'
```

User ID and email are automatically added as X-User-ID and X-User-Email headers.

## Response Format

All responses follow standardized format:

```json
{
  "success": boolean,
  "message": "Descriptive message",
  "data": {} or []     // Actual data (optional)
}
```

Error Response:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error" (optional)
}
```

## Database Schema

### User Table
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

## Project Architecture

```
Client/Frontend
    |
    v
Go Backend (Port 3000)
    |
    +-- Auth Service --> PostgreSQL (Users, Auth)
    |
    +-- Middleware (JWT Validation)
    |
    +-- AI Gateway (Proxy) --> Python AI Service (Port 8000)
                                   |
                                   v
                            AI Processing
                            (Video Gen, Image, ML)
```

## Authorization

Protected endpoints require Authorization header with Bearer token format:

```
Authorization: Bearer {access_token}
```

Token is obtained from `/api/auth/login` response. Access tokens expire after JWT_EXPIRE_HOURS configured in .env (default 24 hours). Use refresh token to get new access token via `/api/auth/refresh`.

## Troubleshooting

Problem: "AI service is unreachable"
- Ensure Python AI service is running on port 8000
- Check AI_SERVICE_URL in .env configuration
- Run: `cd ai-service && python main.py`

Problem: "Invalid token"
- Token may have expired, use refresh endpoint
- Ensure Authorization header format is correct: "Bearer {token}"
- Check JWT_SECRET matches between service and token

Problem: "Connection refused"
- Go backend listening on port 3000
- PostgreSQL running on configured host/port
- Both services have network connectivity
