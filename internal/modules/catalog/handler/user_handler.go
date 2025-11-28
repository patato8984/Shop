package catalog_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	catalog_usescase "github.com/patato8984/Shop/internal/modules/catalog/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
)

type CatalogHandler struct {
	service *catalog_usescase.CatalogService
}

func NewCatalogHandler(service *catalog_usescase.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}
func (h CatalogHandler) GetAllProduct(w http.ResponseWriter, r http.Request) {
	allproduct, err := h.service.GetAllProducts()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.Response("error server", http.StatusInternalServerError, nil))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allproduct)
}
func (h CatalogHandler) GetSkus(w http.ResponseWriter, r http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/catalog/skus/")
	id, err := strconv.Atoi(idStr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("incorrect URL", http.StatusBadRequest, nil))
		return
	}
	skus, err := h.service.GetSkus(id)
	if err != nil {
		switch err.Error() {
		case "skus not fund":
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusBadRequest, nil))
			return
		case "short id":
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusBadRequest, nil))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Response("error server", http.StatusBadRequest, nil))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response("success", http.StatusOK, skus))
}
func (h CatalogHandler) GetAllSkus(w http.ResponseWriter, r http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/catalog/")
	id, err := strconv.Atoi(idStr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("incorrect URL", http.StatusBadRequest, nil))
		return
	}
	product, erro := h.service.GetAllSkus(id)
	if erro != nil {
		switch erro.Error() {
		case "short id":
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Response("invalid id", http.StatusBadRequest, nil))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Response("error server", http.StatusBadRequest, nil))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response("success", http.StatusOK, product))
}
