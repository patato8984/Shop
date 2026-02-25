package cart_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	cart_model "github.com/patato8984/Shop/internal/modules/cart/model"
	cart_usescase "github.com/patato8984/Shop/internal/modules/cart/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type CartHandler struct {
	service *cart_usescase.CartService
}

func NewCartHandler(service *cart_usescase.CartService) *CartHandler {
	return &CartHandler{service: service}
}

var userErrorMapping = map[error]int{
	cart_model.ErrJson:             http.StatusBadRequest,
	cart_model.ErrUrl:              http.StatusBadRequest,
	cart_model.ErrCartItemsNotFund: http.StatusNotFound,
	cart_model.ErrCartNotFound:     http.StatusNotFound,
	cart_model.ErrMaxQuantity:      http.StatusBadRequest,
	cart_model.ErrMinQuantity:      http.StatusBadRequest,
	cart_model.ErrQuantityCatalog:  http.StatusBadRequest,
}

func (h *CartHandler) ResponseUserError(w http.ResponseWriter, err error, ctx context.Context, data any) {
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
			log.Error("BadRequest",
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
func (h *CartHandler) UpdateQuantityProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var productQuantity cart_model.RequestDeltaAndIdSkus
	if err := json.NewDecoder(r.Body).Decode(&productQuantity); err != nil {
		h.ResponseUserError(w, cart_model.ErrJson, ctx, nil)
		return
	}
	id, err := h.service.AddOrUpdateQuantityProduct(ctx, productQuantity.IdSkus, productQuantity.Delta)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, id)
}
func (h *CartHandler) GetSumCost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cost, err := h.service.GetSumCost(ctx)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, cost)
	return
}

func (h *CartHandler) DelProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	idSkus, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseUserError(w, cart_model.ErrUrl, ctx, nil)
		return
	}
	id, err := h.service.DelProduct(ctx, idSkus)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, id)
	return
}
func (h *CartHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var product cart_model.Items
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.ResponseUserError(w, cart_model.ErrJson, ctx, nil)
		return
	}
	id, err := h.service.AddProduct(ctx, product.Id_skus, product.Quantity)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, id)
	return
}
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cart, err := h.service.GetCart(ctx)
	if err != nil {
		h.ResponseUserError(w, err, ctx, cart)
		return
	}
	dto.Response(w, "success", http.StatusOK, cart)
	return
}
