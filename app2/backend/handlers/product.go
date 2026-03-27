package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jwtapp/models"
	"jwtapp/utils"
)

// ProductHandler holds DB dependency.
type ProductHandler struct {
	DB *sql.DB
}

// GetAll GET /api/products
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(
		"SELECT id, name, description, price, stock, created_at, updated_at FROM products ORDER BY id",
	)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			utils.Fail(w, http.StatusInternalServerError, "Failed to scan product")
			return
		}
		products = append(products, p)
	}
	utils.OK(w, products)
}

// GetOne GET /api/products/{id}
func (h *ProductHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p models.Product
	err = h.DB.QueryRow(
		"SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id=?", id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		utils.Fail(w, http.StatusNotFound, "Product not found")
		return
	}
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Database error")
		return
	}
	utils.OK(w, p)
}

// Create POST /api/products  or  POST /api/admin/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Name == "" {
		utils.Fail(w, http.StatusBadRequest, "Product name is required")
		return
	}
	if req.Price < 0 {
		utils.Fail(w, http.StatusBadRequest, "Price must be non-negative")
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)",
		req.Name, req.Description, req.Price, req.Stock,
	)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to create product")
		return
	}
	lastID, _ := result.LastInsertId()

	// Return the freshly created row
	var p models.Product
	h.DB.QueryRow(
		"SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id=?", lastID,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	utils.Created(w, "Product created successfully", p)
}

// Update PUT /api/products/{id}  or  PUT /api/admin/products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req models.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Name == "" {
		utils.Fail(w, http.StatusBadRequest, "Product name is required")
		return
	}

	res, err := h.DB.Exec(
		"UPDATE products SET name=?, description=?, price=?, stock=? WHERE id=?",
		req.Name, req.Description, req.Price, req.Stock, id,
	)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		utils.Fail(w, http.StatusNotFound, "Product not found")
		return
	}

	var p models.Product
	h.DB.QueryRow(
		"SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id=?", id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	utils.OK(w, p)
}

// Delete DELETE /api/products/{id}  or  DELETE /api/admin/products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		utils.Fail(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	res, err := h.DB.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		utils.Fail(w, http.StatusNotFound, "Product not found")
		return
	}
	utils.OK(w, map[string]string{"message": "Product deleted successfully"})
}

// ─── helper ───────────────────────────────────────────────────────────────────

func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}
