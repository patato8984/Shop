package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/patato8984/Shop/internal/modules/catalog/usescase"
)

type CatalogHandler struct {
	service *usescase.CatalogService
}

func NewCatalogHandler(service *usescase.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}
func (h CatalogHandler) GetAllProduct(w http.ResponseWriter, r http.Request) {
	allproduct, err := h.service.GetAllProducts()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error server"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allproduct)
}
func (h CatalogHandler) GetProduct(w http.ResponseWriter, r http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/author/")
	id, err := strconv.Atoi(idStr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Incorrect URL"})
		return
	}
	product, erro := h.service.GetProduct(id)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid id"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
