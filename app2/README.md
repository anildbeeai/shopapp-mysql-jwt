# 🛒 ShopApp — AngularJS + Go + MySQL + JWT

A production-ready full-stack application with:

| Layer      | Technology                          |
|------------|-------------------------------------|
| Frontend   | AngularJS 1.8.3 SPA                 |
| Admin      | AngularJS 1.8.3 SPA                 |
| Backend    | Go (Golang) REST API                |
| Database   | MySQL 8.0 (persistent storage)      |
| Auth       | JWT (HS256, 24h expiry)             |
| Passwords  | bcrypt (cost 10)                    |

---

## 🚀 Quick Start

### Step 1 — Set up MySQL

```sql
-- In your MySQL client:
CREATE DATABASE IF NOT EXISTS shopapp
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;
```

Or run the full schema manually:
```bash
mysql -u root -p < backend/schema.sql
```

### Step 2 — Configure the backend

```bash
cp backend/.env.example backend/.env
# Edit backend/.env — set DB_PASSWORD and JWT_SECRET
```

### Step 3 — Start everything

**Terminal 1 — Backend (Go)**
```bash
chmod +x start-backend.sh
./start-backend.sh
# → http://localhost:8080
```

**Terminal 2 — Frontend (AngularJS)**
```bash
chmod +x serve-frontend.sh
./serve-frontend.sh
# → http://localhost:3000
```

**Terminal 3 — Admin Panel (AngularJS)**
```bash
chmod +x serve-admin.sh
./serve-admin.sh
# → http://localhost:4000
```

> **No Python?** Use any static file server:
> `npx serve frontend -p 3000`  or  `npx serve admin -p 4000`

---
go mod tidy
go run main.go

npx serve .

## 🔑 Default Credentials

| Role  | Email                | Password | URL                      |
|-------|----------------------|----------|--------------------------|
| Admin | admin@example.com    | admin123 | http://localhost:4000    |
| User  | user@example.com     | user123  | http://localhost:3000    |

---

## 🔐 JWT Authentication Flow

```
┌─────────────┐     POST /api/auth/login      ┌──────────────┐
│  AngularJS  │ ───────────────────────────▶  │  Go Backend  │
│  Frontend   │ ◀───────────────────────────  │  (port 8080) │
└─────────────┘     { token, user }           └──────┬───────┘
       │                                             │
       │  localStorage.setItem('jwt_token', token)  │
       │                                             │
       │  $http interceptor adds:                   │
       │  Authorization: Bearer <token>  ──────────▶│
       │                                             │
       │                               middleware.Auth()
       │                               → ParseToken()
       │                               → inject X-User-* headers
       │                               → call handler
```

---

## 📁 Project Structure

```
shopapp-mysql/
├── README.md
├── start-backend.sh          ← Start Go server
├── serve-frontend.sh         ← Serve frontend on :3000
├── serve-admin.sh            ← Serve admin panel on :4000
│
├── backend/
│   ├── main.go               ← Entry point + router
│   ├── go.mod / go.sum       ← Go module
│   ├── schema.sql            ← Full MySQL DDL + seed
│   ├── .env.example          ← Environment config template
│   ├── README.md             ← Backend-specific docs
│   ├── config/
│   │   └── config.go         ← Env-based configuration
│   ├── db/
│   │   └── db.go             ← MySQL connect + migrate + seed
│   ├── models/
│   │   └── models.go         ← User, Product, DTOs
│   ├── utils/
│   │   └── utils.go          ← JWT helpers + HTTP response helpers
│   ├── middleware/
│   │   └── middleware.go     ← CORS, Auth, AdminOnly
│   └── handlers/
│       ├── auth.go           ← Register, Login, Profile
│       ├── product.go        ← Full product CRUD
│       └── admin.go          ← User management + dashboard stats
│
├── frontend/                 ← User-facing store (port 3000)
│   ├── index.html
│   ├── css/style.css
│   ├── js/
│   │   ├── app.js            ← Routes + JWT $http interceptor
│   │   ├── services/
│   │   │   ├── auth.service.js
│   │   │   └── product.service.js
│   │   └── controllers/
│   │       ├── nav, toast, home, auth, product, profile controllers
│   └── views/
│       └── home, products, product-detail, login, register, profile
│
└── admin/                    ← Admin panel (port 4000)
    ├── index.html
    ├── css/admin.css
    ├── js/
    │   ├── app.js            ← Routes + JWT $http interceptor
    │   ├── services/
    │   │   ├── auth.service.js
    │   │   └── api.service.js
    │   └── controllers/
    │       ├── root, toast, login, dashboard, products, users controllers
    └── views/
        └── login, dashboard, products, users
```

---

## 🗄️ Database Tables

| Table           | Purpose                                  |
|-----------------|------------------------------------------|
| users           | All registered users (role: user/admin)  |
| products        | Product catalog                          |
| refresh_tokens  | Token revocation (reserved for future)   |

---

## ⚙️ Environment Variables

| Variable    | Default           | Description              |
|-------------|-------------------|--------------------------|
| SERVER_PORT | 8080              | HTTP listen port         |
| DB_HOST     | 127.0.0.1         | MySQL host               |
| DB_PORT     | 3306              | MySQL port               |
| DB_USER     | root              | MySQL username           |
| DB_PASSWORD | (empty)           | MySQL password           |
| DB_NAME     | shopapp           | MySQL database name      |
| JWT_SECRET  | super-secret-…    | JWT signing secret       |
