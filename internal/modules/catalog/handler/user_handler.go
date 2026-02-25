package catalog_handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	catalog_usescase "github.com/patato8984/Shop/internal/modules/catalog/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type CatalogHandler struct {
	service *catalog_usescase.CatalogService
}

func NewCatalogHandler(service *catalog_usescase.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

var userErrorMapping = map[error]int{
	catalog_model.ErrJson:            http.StatusBadRequest,
	catalog_model.ErrUrl:             http.StatusBadRequest,
	catalog_model.ErrProductNotFound: http.StatusNotFound,
	catalog_model.ErrShortID:         http.StatusBadRequest,
	catalog_model.ErrSkuNotFound:     http.StatusNotFound,
	catalog_model.ErrShortName:       http.StatusBadRequest,
}

func (h CatalogHandler) ResponseUserError(w http.ResponseWriter, err error, ctx context.Context, data any) {
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
	for target, status := range userErrorMapping {
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
func (h CatalogHandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allproduct, err := h.service.GetAllProducts(ctx)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, allproduct)
}
func (h CatalogHandler) GetAllSkus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseUserError(w, catalog_model.ErrUrl, ctx, nil)
		return
	}
	product, err := h.service.GetAllSkus(ctx, id)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, product)
}

func (h CatalogHandler) GetSkus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseUserError(w, catalog_model.ErrUrl, ctx, nil)
		return
	}
	skus, err := h.service.GetSkus(ctx, id)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, skus)
}
