package catalog_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	catalog_usescase "github.com/patato8984/Shop/internal/modules/catalog/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type CatalogAdminHandler struct {
	service *catalog_usescase.CatalogAdminService
}

func NewCatalogAdminHandler(service *catalog_usescase.CatalogAdminService) CatalogAdminHandler {
	return CatalogAdminHandler{service: service}
}

var adminErrorMapping = map[error]int{
	catalog_model.ErrJson:            http.StatusBadRequest,
	catalog_model.ErrUrl:             http.StatusBadRequest,
	catalog_model.ErrProductNotFound: http.StatusNotFound,
	catalog_model.ErrShortID:         http.StatusBadRequest,
	catalog_model.ErrSkuNotFound:     http.StatusNotFound,
	catalog_model.ErrShortName:       http.StatusBadRequest,
}

func (h CatalogAdminHandler) ResponseAdminError(w http.ResponseWriter, err error, ctx context.Context, data any) {
	log := logger.FromContext(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		log.Error("server timeout",
			zap.Error(err),
		)
		dto.Response(w, "timeout", http.StatusGatewayTimeout, nil)
		return
	}
	if errors.Is(err, context.Canceled) {
		return
	}
	for target, status := range adminErrorMapping {
		if errors.Is(err, target) {
			log.Warn("BadRequest",
				zap.Error(err),
			)
			dto.Response(w, target.Error(), status, data)
			return
		}
	}
	log.Error("server error",
		zap.Error(err),
	)
	dto.Response(w, "server error", http.StatusInternalServerError, nil)
}
func (h CatalogAdminHandler) CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var nameProduct string
	if err := json.NewDecoder(r.Body).Decode(&nameProduct); err != nil {
		h.ResponseAdminError(w, catalog_model.ErrJson, ctx, nil)
		return
	}
	product, err := h.service.CreateNewProduct(ctx, nameProduct)
	if err != nil {
		h.ResponseAdminError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, product)
}
func (h CatalogAdminHandler) DelProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseAdminError(w, catalog_model.ErrUrl, ctx, nil)
		return
	}
	deletedID, err := h.service.DelProduct(ctx, id)
	if err != nil {
		h.ResponseAdminError(w, err, ctx, nil)
	}
	dto.Response(w, "success delete", http.StatusOK, deletedID)
}
func (h CatalogAdminHandler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseAdminError(w, catalog_model.ErrUrl, ctx, nil)
		return
	}
	var price catalog_model.SKU
	if err := json.NewDecoder(r.Body).Decode(&price); err != nil {
		h.ResponseAdminError(w, catalog_model.ErrJson, ctx, nil)
		return
	}
	newPrice, err := h.service.UpdatePriceToSkus(ctx, id, price.Price)
	if err != nil {
		h.ResponseAdminError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success update", http.StatusOK, newPrice)
}
func (h CatalogAdminHandler) CreateNewSkus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseAdminError(w, catalog_model.ErrUrl, ctx, nil)
		return
	}
	var skus catalog_model.SKU
	if err := json.NewDecoder(r.Body).Decode(&skus); err != nil {
		h.ResponseAdminError(w, catalog_model.ErrJson, ctx, nil)
		return
	}
	data, errors := h.service.CreateNewSkus(ctx, id, skus)
	if errors != nil {
		h.ResponseAdminError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, data)
}

func (h CatalogAdminHandler) AddStockToSkus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseAdminError(w, catalog_model.ErrUrl, ctx, nil)
		return
	}
	var stock int
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		h.ResponseAdminError(w, catalog_model.ErrJson, ctx, nil)
		return
	}
	sku, err := h.service.AddStockToSkus(ctx, stock, id)
	if err != nil {
		h.ResponseAdminError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, sku)
}
func (h CatalogAdminHandler) DelSkus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseAdminError(w, catalog_model.ErrUrl, ctx, nil)
		return
	}
	deletedID, err := h.service.DelSkus(ctx, id)
	if err != nil {
		h.ResponseAdminError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success delete", http.StatusOK, deletedID)
}
