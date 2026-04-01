package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

// Handler holds the dependencies for the HTTP handlers.
type Handler struct {
	Store *store.MemoryStore
}

// NewHandler creates a new Handler.
func NewHandler(s *store.MemoryStore) *Handler {
	return &Handler{Store: s}
}

// RegisterRoutes sets up the API routes on the given router.
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/products", h.GetProducts).Methods("GET")
	r.HandleFunc("/products", h.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", h.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", h.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}", h.DeleteProduct).Methods("DELETE")
}

// Health returns a simple health check response.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GetProducts returns all products.
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products := h.Store.GetAll()
	respondJSON(w, http.StatusOK, products)
}

// GetProduct returns a single product by ID.
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	p, err := h.Store.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}
	respondJSON(w, http.StatusOK, p)
}

// CreateProduct creates a new product.
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if !p.Validate() {
		respondError(w, http.StatusBadRequest, "Invalid product: name required, price must be >= 0")
		return
	}

	created := h.Store.Create(p)
	respondJSON(w, http.StatusCreated, created)
}

// UpdateProduct updates an existing product.
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	updated, err := h.Store.Update(id, p)
	if err != nil {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}
	respondJSON(w, http.StatusOK, updated)
}

// DeleteProduct deletes a product by ID.
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.Store.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
