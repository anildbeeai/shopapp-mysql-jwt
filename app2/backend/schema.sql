-- ============================================================
-- ShopApp MySQL Schema
-- Run this ONCE to create the database, then let the Go app
-- handle the rest via auto-migration on startup.
-- ============================================================

-- 1. Create database (if it doesn't exist yet)
CREATE DATABASE IF NOT EXISTS shopapp
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE shopapp;

-- 2. Users table
CREATE TABLE IF NOT EXISTS users (
  id         INT          NOT NULL AUTO_INCREMENT,
  name       VARCHAR(120) NOT NULL,
  email      VARCHAR(255) NOT NULL,
  password   VARCHAR(255) NOT NULL,
  role       ENUM('user','admin') NOT NULL DEFAULT 'user',
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Products table
CREATE TABLE IF NOT EXISTS products (
  id          INT            NOT NULL AUTO_INCREMENT,
  name        VARCHAR(255)   NOT NULL,
  description TEXT,
  price       DECIMAL(10,2)  NOT NULL DEFAULT 0.00,
  stock       INT            NOT NULL DEFAULT 0,
  created_at  DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Refresh tokens table (for future token revocation)
CREATE TABLE IF NOT EXISTS refresh_tokens (
  id         INT          NOT NULL AUTO_INCREMENT,
  user_id    INT          NOT NULL,
  token_hash VARCHAR(255) NOT NULL,
  expires_at DATETIME     NOT NULL,
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  CONSTRAINT fk_rt_user FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. Seed admin user  (password: admin123)
INSERT IGNORE INTO users (name, email, password, role) VALUES
  ('Admin User', 'admin@example.com',
   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
   'admin');

-- 6. Seed regular user  (password: user123)
INSERT IGNORE INTO users (name, email, password, role) VALUES
  ('John Doe', 'user@example.com',
   '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
   'user');

-- 7. Seed sample products
INSERT IGNORE INTO products (id, name, description, price, stock) VALUES
  (1, 'Laptop Pro',           'High performance laptop for professionals',         1299.99, 15),
  (2, 'Wireless Mouse',       'Ergonomic wireless mouse with long battery life',     29.99, 100),
  (3, 'Mechanical Keyboard',  'RGB backlit mechanical gaming keyboard',              89.99,  50),
  (4, '4K Monitor',           '27-inch 4K UHD IPS display',                        399.99,  25),
  (5, 'USB-C Hub',            '7-in-1 USB-C hub with HDMI, SD card, USB-A ports',   49.99,  75);

-- Done!
SELECT 'Schema created and seeded successfully.' AS status;
