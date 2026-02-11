# E-Wallet API

RESTful API untuk sistem e-wallet sederhana yang dibangun dengan Go dan Gin framework.

## Fitur

- User Management (Register, Login, Get Profile)
- Wallet Management (Top Up, Get Balance)
- Transaction Management (Transfer, Transaction History)
- JWT Authentication
- Password Hashing dengan bcrypt
- Database Transaction untuk memastikan atomicity
- Race condition handling untuk concurrent transactions
- **Swagger/OpenAPI Documentation**

## Tech Stack

- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- PostgreSQL
- JWT untuk authentication
- bcrypt untuk password hashing
- Swagger/OpenAPI untuk API documentation

## Struktur Proyek

```
ewallet/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── config/
│   ├── config.go                # Konfigurasi aplikasi
│   └── database.go              # Database setup & migration
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── wallet_handler.go
│   │   └── transaction_handler.go
│   ├── middleware/              # Middleware (Auth)
│   │   └── auth.go
│   ├── models/                  # Database models
│   │   ├── user.go
│   │   ├── wallet.go
│   │   └── transaction.go
│   ├── repository/              # Database operations
│   │   ├── user_repository.go
│   │   ├── wallet_repository.go
│   │   └── transaction_repository.go
│   └── service/                 # Business logic
│       ├── auth_service.go
│       ├── wallet_service.go
│       └── transaction_service.go
├── migrations/                  # SQL migrations
│   └── 001_create_tables.sql
├── pkg/
│   └── utils/                   # Utilities
│       ├── jwt.go
│       └── response.go
├── .env.example                 # Environment variables template
├── .gitignore
├── go.mod
└── README.md
```

## Instalasi

### Prerequisites

- Go 1.21 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi

### Setup

1. Clone repository:
```bash
git clone <repository-url>
cd ewallet
```

2. Install dependencies:
```bash
go mod download
```

3. Setup database PostgreSQL dan buat database:
```sql
CREATE DATABASE ewallet_db;
```

4. Copy file `.env.example` ke `.env` dan sesuaikan konfigurasi:
```bash
cp .env.example .env
```

5. Edit file `.env` dengan konfigurasi database Anda:
```env
SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ewallet_db
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-key-change-this
JWT_EXPIRY=24h
```

6. (Optional) Run migrations manually:
```bash
make migrate-up
```
Note: Migrations akan otomatis dijalankan saat aplikasi start.

7. Jalankan aplikasi:
```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080`

**Migrations akan otomatis dijalankan** saat aplikasi pertama kali start.

## Demo Data

Proyek ini menyediakan **demo data** untuk memudahkan testing:

### 5 Demo Users (Password: `password123`)

| Email | Name | Balance |
|-------|------|---------|
| alice@example.com | Alice Johnson | Rp 960,000 |
| bob@example.com | Bob Smith | Rp 700,000 |
| charlie@example.com | Charlie Brown | Rp 525,000 |
| diana@example.com | Diana Prince | Rp 300,000 |
| eve@example.com | Eve Davis | Rp 115,000 |

### Quick Test

1. Login with `alice@example.com` / `password123`
2. Check balance (should be Rp 960,000)
3. Transfer to Bob (user_id: 2)
4. View transaction history

Lihat [DEMO_DATA.md](DEMO_DATA.md) untuk detail lengkap.

### Remove Demo Data

```bash
make migrate-down  # Remove seed data
```

## API Documentation

### Swagger UI

Setelah aplikasi berjalan, akses dokumentasi API interaktif di:

**http://localhost:8080/swagger/index.html**

Swagger UI menyediakan:
- Daftar lengkap semua endpoints
- Detail request/response untuk setiap endpoint
- Kemampuan untuk test API langsung dari browser
- Schema definitions untuk semua models

### Generate Swagger Documentation

Jika Anda melakukan perubahan pada API annotations, regenerate documentation dengan:

```bash
make swagger
# atau
~/go/bin/swag init -g cmd/server/main.go -o docs
```

## API Endpoints

### Health Check

```
GET /health
```

### Authentication

#### Register
```
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login
```
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

### User Management

#### Get Profile
```
GET /api/users/profile
Authorization: Bearer <token>
```

### Wallet Management

#### Get Balance
```
GET /api/wallets/balance
Authorization: Bearer <token>
```

#### Top Up
```
POST /api/wallets/topup
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": 100000
}
```

### Transaction Management

#### Transfer
```
POST /api/transactions/transfer
Authorization: Bearer <token>
Content-Type: application/json

{
  "receiver_id": 2,
  "amount": 50000
}
```

#### Get Transaction History
```
GET /api/transactions/history?limit=50
Authorization: Bearer <token>
```

## Validasi Business Logic

1. **Transfer:**
   - Amount harus lebih besar dari 0
   - Tidak bisa transfer ke diri sendiri
   - Saldo pengirim harus mencukupi
   - Receiver harus ada di database

2. **Top Up:**
   - Amount harus lebih besar dari 0

3. **Register:**
   - Semua field wajib diisi
   - Email harus valid
   - Password minimal 6 karakter
   - Email harus unik

## Database Migrations

Proyek ini menggunakan [golang-migrate](https://github.com/golang-migrate/migrate) untuk database migrations.

### Migration Commands

```bash
make migrate-up          # Apply all pending migrations
make migrate-down        # Rollback last migration
make migrate-version     # Check current migration version
make migrate-create name=add_column  # Create new migration
```

Lihat [MIGRATION.md](MIGRATION.md) untuk panduan lengkap.

### Automatic Migrations

Migrations akan **otomatis dijalankan** saat aplikasi start. Jika ingin run manual:

```bash
make migrate-up
```

## Database Schema

### Users Table
- id (Primary Key)
- name
- email (Unique)
- password (hashed)
- created_at
- updated_at
- deleted_at

### Wallets Table
- id (Primary Key)
- user_id (Foreign Key, Unique)
- balance (Decimal, Default: 0)
- created_at
- updated_at
- deleted_at

### Transactions Table
- id (Primary Key)
- sender_id (Foreign Key, nullable)
- receiver_id (Foreign Key)
- amount (Decimal)
- type (topup/transfer)
- status (pending/success/failed)
- created_at
- updated_at
- deleted_at

## Security Features

1. **Password Hashing:** Password di-hash menggunakan bcrypt
2. **JWT Authentication:** Protected endpoints memerlukan valid JWT token
3. **Database Transactions:** Transfer menggunakan database transaction untuk memastikan atomicity
4. **Row Locking:** Menggunakan `FOR UPDATE` untuk mencegah race condition pada concurrent transactions
5. **Deadlock Prevention:** Wallet locking dilakukan dalam urutan konsisten (ID rendah terlebih dahulu)

## Testing dengan cURL

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'
```

### Get Profile
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Top Up
```bash
curl -X POST http://localhost:8080/api/wallets/topup \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount":100000}'
```

### Transfer
```bash
curl -X POST http://localhost:8080/api/transactions/transfer \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"receiver_id":2,"amount":50000}'
```

### Get Transaction History
```bash
curl -X GET http://localhost:8080/api/transactions/history \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Error Handling

API menggunakan format response yang konsisten:

### Success Response
```json
{
  "success": true,
  "message": "Success message",
  "data": { }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error"
}
```

### HTTP Status Codes
- 200: OK
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 404: Not Found
- 500: Internal Server Error

## Development

### Install dependencies
```bash
go mod download
```

### Run application
```bash
go run cmd/server/main.go
```

### Build application
```bash
go build -o bin/ewallet cmd/server/main.go
```

### Run built binary
```bash
./bin/ewallet
```

## License

MIT License
