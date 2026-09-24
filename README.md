# GoWallet — Digital Wallet Service

A digital wallet system built with **Go + PostgreSQL + React**, focusing on backend engineering concepts such as concurrency, transaction safety, authentication, and database performance.

## Tech Stack

| Layer          | Technology         |
| -------------- | ------------------ |
| Backend        | Go, Gin             |
| Database       | PostgreSQL 16       |
| DB Driver      | pgx v5 / pgxpool     |
| Authentication | JWT + bcrypt         |
| API Docs       | Swagger / OpenAPI    |
| Frontend       | React + Vite         |
| Infrastructure | Docker Compose       |
| CI             | GitHub Actions        |

## Key Features

- RESTful API with Go + Gin
- Repository pattern with interface-based design
- PostgreSQL transactions and `SELECT ... FOR UPDATE` row locking
- Connection pooling with `pgxpool`
- Goroutine + channel-based Worker Pool for bulk transfers
- Monthly table partitioning and database indexes
- JWT authentication with ownership checks (IDOR protection)
- bcrypt password hashing
- Transaction audit trail, including failed transfers
- Swagger API documentation
- GitHub Actions CI for build, vet, and test

## Project Structure

```text
e-wallet/
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── api/
│   │   ├── auth/
│   │   ├── db/
│   │   ├── domain/
│   │   ├── repository/
│   │   └── worker/
│   ├── schema/
│   ├── docker-compose.yml
│   └── .env.example
├── frontend/
│   └── src/
└── .github/workflows/ci.yml
```

## Getting Started

### Backend

```bash
cd backend
cp .env.example .env
docker-compose up -d
go mod download
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go
go run cmd/main.go
```

Backend: `http://localhost:8080`
Swagger: `http://localhost:8080/swagger/index.html`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend: `http://localhost:5173`

## API

| Method | Endpoint               | Auth | Description           |
| ------ | ----------------------- | ---- | ---------------------- |
| POST   | `/register`              | ❌    | Register account       |
| POST   | `/login`                 | ❌    | Login and receive JWT  |
| GET    | `/accounts/me`           | ✅    | Get own account        |
| POST   | `/deposit`               | ✅    | Deposit funds          |
| POST   | `/transfer`               | ✅    | Transfer funds          |
| POST   | `/transfer/bulk`          | ✅    | Bulk transfer           |
| GET    | `/transactions`           | ✅    | Transaction history      |
| GET    | `/transactions/{code}`    | ✅    | Get transaction          |

Protected endpoints derive account ownership from the authenticated JWT.

## Transfer Flow

```text
POST /transfer
      ↓
BEGIN TRANSACTION
      ↓
SELECT sender FOR UPDATE
      ↓
Check balance
      ↓
Debit sender
      ↓
Credit receiver
      ↓
Record transaction
      ↓
COMMIT
```

## Reset Database

```bash
cd backend
docker-compose down -v
docker-compose up -d
```
