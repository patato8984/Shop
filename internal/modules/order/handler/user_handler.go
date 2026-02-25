package order_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	order_model "github.com/patato8984/Shop/internal/modules/order/model"
	order_usescase "github.com/patato8984/Shop/internal/modules/order/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type OrderUserHandler struct {
	service *order_usescase.OrderService
}

func NewOrderUserHandler(service *order_usescase.OrderService) *OrderUserHandler {
	return &OrderUserHandler{service: service}
}

var userErrorMapping = map[error]int{
	order_model.ErrQuantityItemsCart: http.StatusUnprocessableEntity,
	order_model.ErrUrl:               http.StatusBadRequest,
	order_model.ErrHash:              http.StatusBadRequest,
	order_model.ErrStatusPey:         http.StatusOK,
	order_model.ErrStatusCancelled:   http.StatusOK,
	order_model.ErrOrderNotFound:     http.StatusNotFound,
	order_model.ErrStock:             http.StatusConflict,
	order_model.ErrCreateTransaction: http.StatusBadGateway,
}

func (h *OrderUserHandler) ResponseUserError(w http.ResponseWriter, err error, ctx context.Context, data any) {
	log := logger.FromContext(ctx)
	log.With(zap.Error(err))
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
	dto.Response(w, "server error", http.StatusInternalServerError, nil)

}
func (h *OrderUserHandler) AddNewOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var orderInfo order_model.Order
	err := json.NewDecoder(r.Body).Decode(&orderInfo)
	if err != nil {
		h.ResponseUserError(w, order_model.ErrJson, ctx, nil)
		return
	}
	idOrder, err := h.service.AddNowOrder(ctx, orderInfo.Address, orderInfo.Messages)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	orderInfo.Id = idOrder
	dto.Response(w, "success", http.StatusOK, orderInfo)
}
func (h *OrderUserHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allOrders, err := h.service.GetAllOrders(ctx)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, allOrders)
}
func (h *OrderUserHandler) GetOrderAndAllItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	idOrder, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseUserError(w, order_model.ErrUrl, ctx, nil)
		return
	}
	order, err := h.service.GetOrderAndAllItems(ctx, idOrder)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, order)
}
func (h *OrderUserHandler) Payment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	idOrder, err := strconv.Atoi(idStr)
	if err != nil {
		h.ResponseUserError(w, order_model.ErrUrl, ctx, nil)
		return
	}
	bankRequest, err := h.service.GetPayment(ctx, idOrder)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "paid", http.StatusOK, bankRequest)
	return
}
