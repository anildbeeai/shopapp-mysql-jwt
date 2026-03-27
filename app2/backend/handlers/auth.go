package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"jwtapp/models"
	"jwtapp/utils"
)

// AuthHandler holds DB + JWT secret dependency.
type AuthHandler struct {
	DB        *sql.DB
	JWTSecret string
}

// Register POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.Fail(w, http.StatusBadRequest, "Name, email and password are required")
		return
	}
	if len(req.Password) < 6 {
		utils.Fail(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Check if email already exists
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=?)", req.Email).Scan(&exists)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		utils.Fail(w, http.StatusConflict, "Email address is already registered")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Insert user
	result, err := h.DB.Exec(
		"INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, 'user')",
		req.Name, req.Email, string(hash),
	)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	lastID, _ := result.LastInsertId()

	user := models.User{
		ID:    int(lastID),
		Name:  req.Name,
		Email: req.Email,
		Role:  "user",
	}

	token, err := utils.GenerateToken(user, h.JWTSecret)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.Created(w, "Registration successful", models.AuthResponse{
		Token: token,
		User:  user.Safe(),
	})
}

// Login POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		utils.Fail(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	var user models.User
	err := h.DB.QueryRow(
		"SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email=?",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		utils.Fail(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Fail(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(user, h.JWTSecret)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.OK(w, models.AuthResponse{
		Token: token,
		User:  user.Safe(),
	})
}

// Profile GET /api/auth/profile  (requires Auth middleware)
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("X-User-Email")

	var user models.User
	err := h.DB.QueryRow(
		"SELECT id, name, email, role, created_at, updated_at FROM users WHERE email=?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		utils.Fail(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.OK(w, user.Safe())
}
