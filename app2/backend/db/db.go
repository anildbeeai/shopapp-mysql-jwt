package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"jwtapp/config"
)

// Connect opens a MySQL connection pool and verifies connectivity.
func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Retry connection up to 10 times (useful when MySQL is starting up)
	for i := 1; i <= 10; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		log.Printf("⏳ Waiting for MySQL (%d/10): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}

// Migrate creates all tables and seeds default data if they don't exist.
func Migrate(db *sql.DB) error {
	log.Println("🔄 Running database migrations...")

	statements := []string{
		// ── users ──────────────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS users (
			id         INT          NOT NULL AUTO_INCREMENT,
			name       VARCHAR(120) NOT NULL,
			email      VARCHAR(255) NOT NULL,
			password   VARCHAR(255) NOT NULL,
			role       ENUM('user','admin') NOT NULL DEFAULT 'user',
			created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uq_email (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		// ── products ───────────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS products (
			id          INT            NOT NULL AUTO_INCREMENT,
			name        VARCHAR(255)   NOT NULL,
			description TEXT,
			price       DECIMAL(10,2)  NOT NULL DEFAULT 0.00,
			stock       INT            NOT NULL DEFAULT 0,
			created_at  DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		// ── refresh_tokens (optional — for future token revocation) ────────
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id         INT          NOT NULL AUTO_INCREMENT,
			user_id    INT          NOT NULL,
			token_hash VARCHAR(255) NOT NULL,
			expires_at DATETIME     NOT NULL,
			created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_user_id (user_id),
			CONSTRAINT fk_rt_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}

	if err := seed(db); err != nil {
		return fmt.Errorf("seed: %w", err)
	}

	log.Println("✅ Database migrations complete.")
	return nil
}

// seed inserts default admin, sample user, and sample products
// only if no rows exist (idempotent).
func seed(db *sql.DB) error {
	// seed users
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		log.Println("🌱 Seeding default users...")
		// Pre-hashed passwords (bcrypt cost 10):
		//   admin123 → below hash
		//   user123  → below hash
		rows := []struct{ name, email, hash, role string }{
			{
				"Admin User", "admin@example.com",
				"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
				"admin",
			},
			{
				"John Doe", "user@example.com",
				"$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
				"user",
			},
		}
		for _, r := range rows {
			_, err := db.Exec(
				`INSERT INTO users (name, email, password, role) VALUES (?,?,?,?)`,
				r.name, r.email, r.hash, r.role,
			)
			if err != nil {
				return err
			}
		}
	}

	// seed products
	if err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		log.Println("🌱 Seeding sample products...")
		products := []struct {
			name, desc string
			price      float64
			stock      int
		}{
			{"Laptop Pro", "High performance laptop for professionals", 1299.99, 15},
			{"Wireless Mouse", "Ergonomic wireless mouse with long battery life", 29.99, 100},
			{"Mechanical Keyboard", "RGB backlit mechanical gaming keyboard", 89.99, 50},
			{"4K Monitor", "27-inch 4K UHD IPS display", 399.99, 25},
			{"USB-C Hub", "7-in-1 USB-C hub with HDMI, SD card, USB-A ports", 49.99, 75},
		}
		for _, p := range products {
			_, err := db.Exec(
				`INSERT INTO products (name, description, price, stock) VALUES (?,?,?,?)`,
				p.name, p.desc, p.price, p.stock,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
