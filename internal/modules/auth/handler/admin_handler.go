package auth_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/patato8984/Shop/internal/modules/auth/model"
	usescase_user "github.com/patato8984/Shop/internal/modules/auth/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type AdminHandler struct {
	service *usescase_user.AdminService
}

func NewAdminHandler(service *usescase_user.AdminService) AdminHandler {
	return AdminHandler{service: service}
}

var adminErrorMapping = map[error]int{
	model.ErrJson:                    http.StatusBadRequest,
	model.ErrMailBusy:                http.StatusConflict,
	model.ErrNickNameBusy:            http.StatusBadRequest,
	model.ErrShortPasswordOrNickname: http.StatusBadRequest,
}

func (h *AdminHandler) ResponseAdminError(w http.ResponseWriter, err error, ctx context.Context, data any) {
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
	dto.Response(w, "server error", http.StatusInternalServerError, nil)

}

func (h *AdminHandler) CreateNewAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.ResponseAdminError(w, model.ErrJson, ctx, nil)
		return
	}
	if err := h.service.CreateNewAdmin(ctx, user); err != nil {
		h.ResponseAdminError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "the admin has been created", http.StatusOK, nil)
}
