package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"jwtapp/config"
	"jwtapp/db"
	"jwtapp/handlers"
	"jwtapp/middleware"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MySQL: %v\n", err)
	}
	defer database.Close()
	log.Println("✅ Connected to MySQL successfully.")

	if err := db.Migrate(database); err != nil {
		log.Fatalf("❌ Migration failed: %v\n", err)
	}

	authH    := &handlers.AuthHandler{DB: database, JWTSecret: cfg.JWTSecret}
	productH := &handlers.ProductHandler{DB: database}
	adminH   := &handlers.AdminHandler{DB: database}

	auth      := middleware.Auth(cfg.JWTSecret)
	adminOnly := middleware.AdminOnly(cfg.JWTSecret)

	r := mux.NewRouter()

	// Handle ALL OPTIONS preflight requests inline — no separate function needed
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodOptions)

	r.Use(middleware.CORS)

	// Auth
	r.HandleFunc("/api/auth/register", authH.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login",    authH.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/profile",  auth(authH.Profile)).Methods(http.MethodGet)

	// Products
	r.HandleFunc("/api/products",      productH.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/api/products/{id}", productH.GetOne).Methods(http.MethodGet)
	r.HandleFunc("/api/products",      auth(productH.Create)).Methods(http.MethodPost)
	r.HandleFunc("/api/products/{id}", auth(productH.Update)).Methods(http.MethodPut)
	r.HandleFunc("/api/products/{id}", auth(productH.Delete)).Methods(http.MethodDelete)

	// Admin
	r.HandleFunc("/api/admin/stats",         adminOnly(adminH.Stats)).Methods(http.MethodGet)
	r.HandleFunc("/api/admin/users",         adminOnly(adminH.GetUsers)).Methods(http.MethodGet)
	r.HandleFunc("/api/admin/users/{id}",    adminOnly(adminH.DeleteUser)).Methods(http.MethodDelete)
	r.HandleFunc("/api/admin/products",      adminOnly(productH.Create)).Methods(http.MethodPost)
	r.HandleFunc("/api/admin/products/{id}", adminOnly(productH.Update)).Methods(http.MethodPut)
	r.HandleFunc("/api/admin/products/{id}", adminOnly(productH.Delete)).Methods(http.MethodDelete)

	addr := ":" + cfg.ServerPort
	fmt.Printf("\n🚀 ShopApp API running at http://localhost%s\n", addr)
	fmt.Println("   Admin : admin@example.com / admin123")
	fmt.Println("   User  : user@example.com  / user123")
	fmt.Printf("   DB    : %s@%s:%s/%s\n\n", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("❌ Server failed: %v\n", err)
	}
}
