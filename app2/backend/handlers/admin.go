package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jwtapp/models"
	"jwtapp/utils"
)

// AdminHandler holds DB dependency.
type AdminHandler struct {
	DB *sql.DB
}

// GetUsers GET /api/admin/users
func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(
		"SELECT id, name, email, role, created_at, updated_at FROM users ORDER BY id",
	)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer rows.Close()

	users := []models.UserPublic{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			utils.Fail(w, http.StatusInternalServerError, "Failed to scan user")
			return
		}
		users = append(users, u.Safe())
	}
	utils.OK(w, users)
}

// DeleteUser DELETE /api/admin/users/{id}
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Prevent deleting admin accounts
	var role string
	err = h.DB.QueryRow("SELECT role FROM users WHERE id=?", id).Scan(&role)
	if err == sql.ErrNoRows {
		utils.Fail(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Database error")
		return
	}
	if role == "admin" {
		utils.Fail(w, http.StatusForbidden, "Cannot delete admin accounts")
		return
	}

	// Prevent self-deletion
	selfEmail := r.Header.Get("X-User-Email")
	var userEmail string
	h.DB.QueryRow("SELECT email FROM users WHERE id=?", id).Scan(&userEmail)
	if userEmail == selfEmail {
		utils.Fail(w, http.StatusForbidden, "Cannot delete your own account")
		return
	}

	res, err := h.DB.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		utils.Fail(w, http.StatusNotFound, "User not found")
		return
	}
	utils.OK(w, map[string]string{"message": "User deleted successfully"})
}

// Stats GET /api/admin/stats
func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	var stats models.DashboardStats

	// total users
	h.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	// admin users
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role='admin'").Scan(&stats.AdminUsers)
	// total products
	h.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&stats.TotalProducts)
	// inventory value = SUM(price * stock)
	h.DB.QueryRow("SELECT COALESCE(SUM(price * stock), 0) FROM products").Scan(&stats.TotalRevenue)

	utils.OK(w, stats)
}
