package order_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	order_model "github.com/patato8984/Shop/internal/modules/order/model"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type WeBHookService interface {
	WeBhookBankPayment(ctx context.Context, bankApiData order_model.BankApiRequest) error
}
type WeBHookHandler struct {
	weBHookService WeBHookService
}

func NewWeBHookService(weBHookService WeBHookService) *WeBHookHandler {
	return &WeBHookHandler{weBHookService: weBHookService}
}

var vebHookErrorMapping = map[error]int{
	order_model.ErrHeader:        http.StatusBadRequest,
	order_model.ErrJson:          http.StatusBadRequest,
	order_model.ErrUrl:           http.StatusBadRequest,
	order_model.ErrHash:          http.StatusBadRequest,
	order_model.ErrOrderNotFound: http.StatusNotFound,
}

func (h *WeBHookHandler) ResponseUserError(w http.ResponseWriter, err error, ctx context.Context, data any) {
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
	for target, status := range vebHookErrorMapping {
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
func (h *WeBHookHandler) WebHookPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var bankApiData order_model.BankApiRequest
	strSignature := r.Header.Get("X-Signature")
	err := json.NewDecoder(r.Body).Decode(&bankApiData)
	if err != nil {
		h.ResponseUserError(w, order_model.ErrJson, ctx, nil)
		return
	}
	bankApiData.Signature = strSignature
	err = h.weBHookService.WeBhookBankPayment(ctx, bankApiData)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "success", http.StatusOK, nil)
	return
}
