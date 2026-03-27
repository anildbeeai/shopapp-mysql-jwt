package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"jwtapp/models"
)

// ─── JWT ──────────────────────────────────────────────────────────────────────

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given user.
func GenerateToken(user models.User, secret string) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken validates and parses a JWT string.
func ParseToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// ─── HTTP helpers ─────────────────────────────────────────────────────────────

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"success":false,"message":"encoding error"}`, http.StatusInternalServerError)
	}
}

// OK writes a 200 JSON response.
func OK(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: data})
}

// Created writes a 201 JSON response.
func Created(w http.ResponseWriter, msg string, data interface{}) {
	WriteJSON(w, http.StatusCreated, models.APIResponse{Success: true, Message: msg, Data: data})
}

// Fail writes an error JSON response.
func Fail(w http.ResponseWriter, code int, msg string) {
	WriteJSON(w, code, models.APIResponse{Success: false, Message: msg})
}
