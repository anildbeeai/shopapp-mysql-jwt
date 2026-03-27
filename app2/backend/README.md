PS C:\@lab_pro-2\shopapp-mysql-jwt\app2\backend> go mod tidy
verifying github.com/golang-jwt/jwt/v5@v5.2.0: checksum mismatch
        downloaded: h1:d/ix8ftRUorsN+5eMIlF4T6J8CAt9rch3My2winC1Jw=
        go.sum:     h1:d/alXHFMECZGKZHQgwd5FnAKNQbMsVZIlsDZDZ6SXNE=

SECURITY ERROR
This download does NOT match an earlier download recorded in go.sum.
The bits may have been replaced on the origin server, or an attacker may
have intercepted the download attempt.

For more information, see 'go help module-auth'.
PS C:\@lab_pro-2\shopapp-mysql-jwt\app2\backend> go run main.go
verifying github.com/golang-jwt/jwt/v5@v5.2.0: checksum mismatch
        downloaded: h1:d/ix8ftRUorsN+5eMIlF4T6J8CAt9rch3My2winC1Jw=
        go.sum:     h1:d/alXHFMECZGKZHQgwd5FnAKNQbMsVZIlsDZDZ6SXNE=

SECURITY ERROR
This download does NOT match an earlier download recorded in go.sum.
The bits may have been replaced on the origin server, or an attacker may
have intercepted the download attempt.

For more information, see 'go help module-auth'.
PS C:\@lab_pro-2\shopapp-mysql-jwt\app2\backend>

go mod tidy
go run main.go
npx serve .
====================================================================================

# Backend — Go JWT API with MySQL

## Prerequisites
- Go 1.21+
- MySQL 8.0+ (or MariaDB 10.6+)

---

## 1. Create the MySQL Database

Option A — let the Go app auto-migrate (recommended):
```sql
-- Just create the empty database first:
CREATE DATABASE IF NOT EXISTS shopapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Option B — run the full schema manually:
```bash
mysql -u root -p < schema.sql
```

---

## 2. Configure Environment

```bash
cp .env.example .env
# Edit .env and set your MySQL password + JWT secret
```

**Available env variables:**

| Variable      | Default                    | Description              |
|---------------|----------------------------|--------------------------|
| SERVER_PORT   | 8080                       | HTTP port                |
| DB_HOST       | 127.0.0.1                  | MySQL hostname           |
| DB_PORT       | 3306                       | MySQL port               |
| DB_USER       | root                       | MySQL username           |
| DB_PASSWORD   | (empty)                    | MySQL password           |
| DB_NAME       | shopapp                    | Database name            |
| JWT_SECRET    | super-secret-…             | JWT signing key          |

---

## 3. Run the Server

```bash
# Install dependencies
go mod tidy
go run main.go

npx serve .

# Set env vars (Linux/macOS)
export DB_PASSWORD=your_password
export JWT_SECRET=your_jwt_secret

# Or source .env
export $(grep -v '^#' .env | xargs)

# Start
go run main.go
# → http://localhost:8080
```

The server will:
1. Connect to MySQL
2. Auto-create tables (if not present)
3. Seed default users + products (if tables are empty)

---

## Default Credentials

| Role  | Email                | Password |
|-------|----------------------|----------|
| Admin | admin@example.com    | admin123 |
| User  | user@example.com     | user123  |

---

## Database Schema

```
users
├── id          INT AUTO_INCREMENT PK
├── name        VARCHAR(120)
├── email       VARCHAR(255) UNIQUE
├── password    VARCHAR(255)  ← bcrypt hash
├── role        ENUM('user','admin')
├── created_at  DATETIME
└── updated_at  DATETIME

products
├── id          INT AUTO_INCREMENT PK
├── name        VARCHAR(255)
├── description TEXT
├── price       DECIMAL(10,2)
├── stock       INT
├── created_at  DATETIME
└── updated_at  DATETIME

refresh_tokens  (reserved for token revocation)
├── id          INT AUTO_INCREMENT PK
├── user_id     INT FK→users.id
├── token_hash  VARCHAR(255)
├── expires_at  DATETIME
└── created_at  DATETIME
```

---

## API Reference

### Auth
| Method | Endpoint            | Auth     | Body                          |
|--------|---------------------|----------|-------------------------------|
| POST   | /api/auth/register  | –        | `{name, email, password}`     |
| POST   | /api/auth/login     | –        | `{email, password}`           |
| GET    | /api/auth/profile   | Bearer   | –                             |

### Products
| Method | Endpoint              | Auth          |
|--------|-----------------------|---------------|
| GET    | /api/products         | –             |
| GET    | /api/products/:id     | –             |
| POST   | /api/products         | Bearer (any)  |
| PUT    | /api/products/:id     | Bearer (any)  |
| DELETE | /api/products/:id     | Bearer (any)  |

### Admin
| Method | Endpoint                  | Auth          |
|--------|---------------------------|---------------|
| GET    | /api/admin/stats          | Bearer admin  |
| GET    | /api/admin/users          | Bearer admin  |
| DELETE | /api/admin/users/:id      | Bearer admin  |
| POST   | /api/admin/products       | Bearer admin  |
| PUT    | /api/admin/products/:id   | Bearer admin  |
| DELETE | /api/admin/products/:id   | Bearer admin  |

### JWT Token Usage
```
Authorization: Bearer <token>
```

---

## Project Structure

```
backend/
├── main.go              ← Entry point, router setup
├── go.mod / go.sum      ← Module dependencies
├── schema.sql           ← Full MySQL schema (manual setup option)
├── .env.example         ← Environment config template
├── config/
│   └── config.go        ← Reads env vars with defaults
├── db/
│   └── db.go            ← MySQL connection pool + auto-migration + seeder
├── models/
│   └── models.go        ← Structs: User, Product, DTOs
├── utils/
│   └── utils.go         ← JWT generate/parse + HTTP response helpers
├── middleware/
│   └── middleware.go    ← CORS, Auth, AdminOnly middleware
└── handlers/
    ├── auth.go          ← Register, Login, Profile
    ├── product.go       ← GetAll, GetOne, Create, Update, Delete
    └── admin.go         ← GetUsers, DeleteUser, Stats
```
