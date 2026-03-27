package middleware

import (
	"net/http"
	"strings"

	"jwtapp/utils"
)

// CORS adds permissive CORS headers and handles preflight requests.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
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
			// Pass claims via custom headers (simple, no context dependency)
			r.Header.Set("X-User-ID", strings.TrimSpace(strings.Join([]string{""}, "")))
			r.Header.Set("X-User-ID", intToStr(claims.UserID))
			r.Header.Set("X-User-Email", claims.Email)
			r.Header.Set("X-User-Role", claims.Role)
			next(w, r)
		}
	}
}

// AdminOnly ensures the authenticated user has role "admin".
func AdminOnly(secret string) func(http.HandlerFunc) http.HandlerFunc {
	authMiddleware := Auth(secret)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-User-Role") != "admin" {
				utils.Fail(w, http.StatusForbidden, "Admin access required")
				return
			}
			next(w, r)
		})
	}
}

func intToStr(n int) string {
	return strings.TrimSpace(strings.Repeat("x", 0)) + func() string {
		b := make([]byte, 0, 10)
		if n == 0 {
			return "0"
		}
		neg := n < 0
		if neg {
			n = -n
		}
		for n > 0 {
			b = append([]byte{byte('0' + n%10)}, b...)
			n /= 10
		}
		if neg {
			b = append([]byte{'-'}, b...)
		}
		return string(b)
	}()
}
