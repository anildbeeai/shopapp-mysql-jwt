package middleware

import (
	"net/http"
	"strings"

	"jwtapp/utils"
)

// setCORSHeaders writes all necessary CORS headers onto the response.
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// CORS middleware — sets headers on every response AND short-circuits OPTIONS.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204 — no body needed
			return
		}
		next.ServeHTTP(w, r)
	})
}

// PreflightHandler handles OPTIONS requests Gorilla Mux would reject with 405.
// Register BEFORE all routes: r.Methods("OPTIONS").PathPrefix("/").HandlerFunc(...)
func PreflightHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}

// Auth validates the JWT bearer token and injects user claims into request headers.
func Auth(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				utils.Fail(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ParseToken(tokenStr, secret)
			if err != nil {
				utils.Fail(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}
			r.Header.Set("X-User-ID", intToStr(claims.UserID))
			r.Header.Set("X-User-Email", claims.Email)
			r.Header.Set("X-User-Role", claims.Role)
			next(w, r)
		}
	}
}

// AdminOnly ensures the authenticated user has role "admin".
func AdminOnly(secret string) func(http.HandlerFunc) http.HandlerFunc {
	authMW := Auth(secret)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return authMW(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-User-Role") != "admin" {
				utils.Fail(w, http.StatusForbidden, "Admin access required")
				return
			}
			next(w, r)
		})
	}
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	b := make([]byte, 0, 10)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
