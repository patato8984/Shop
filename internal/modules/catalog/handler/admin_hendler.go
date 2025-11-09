package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
	"github.com/patato8984/Shop/internal/modules/catalog/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
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
		json.NewEncoder(w).Encode(dto.Response("error json", http.StatusBadRequest, nil))
		return
	}
	product, err := h.service.CreateNewProduct(nameProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.Response("error server", http.StatusBadRequest, nil))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response("success", http.StatusOK, product))
}
func (h CatalogAdminHandler) CreateNewSkus(w http.ResponseWriter, r http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/catalog/")
	id, err := strconv.Atoi(idStr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("incorrect URL", http.StatusBadRequest, nil))
		return
	}
	var skus model.SKU
	if err := json.NewDecoder(r.Body).Decode(&skus); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("error json", http.StatusBadRequest, nil))
		return
	}
	data, errors := h.service.CreateNewSkus(id, skus)
	if errors != nil {
		switch errors.Error() {
		case "product not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.Response(errors.Error(), http.StatusNotFound, nil))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Response("error server", http.StatusInternalServerError, nil))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response("success", http.StatusOK, data))
}
func (h CatalogAdminHandler) AddStockToSkus(w http.ResponseWriter, r http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/catalog/skus/")
	id, err := strconv.Atoi(idStr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("incorrect URL", http.StatusBadRequest, nil))
		return
	}
	var stock int
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("error json", http.StatusBadRequest, nil))
		return
	}
	sku, err := h.service.AddStockToSkus(stock, id)
	if err != nil {
		switch err.Error() {
		case "sku not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusNotFound, nil))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Response("error server", http.StatusInternalServerError, nil))
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response("success", http.StatusOK, sku))
}
