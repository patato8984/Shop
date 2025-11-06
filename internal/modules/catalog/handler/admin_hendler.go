package handler

import (
	"encoding/json"
	"net/http"

	"github.com/patato8984/Shop/internal/modules/catalog/usescase"
)

type CatalogAdminHandler struct {
	service *usescase.CatalogAdminServise
}

func NewCatalogAdminHandler(service *usescase.CatalogAdminServise) CatalogAdminHandler {
	return CatalogAdminHandler{service: service}
}
func (h CatalogAdminHandler) CreateNewProduct(w http.ResponseWriter, r http.Request) {
	var nameProduct string
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&nameProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid json"})
		return
	}
	product, err := h.service.CreateNewProduct(nameProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error server"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
