package handler

import (
	"encoding/json"
	"net/http"

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
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allproduct)
}
