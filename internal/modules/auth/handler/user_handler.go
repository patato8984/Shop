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

type UserHandler struct {
	service *usescase_user.UserService
}

func NewUserHandler(service *usescase_user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

var userErrorMapping = map[error]int{
	model.ErrJson:                    http.StatusBadRequest,
	model.ErrMailBusy:                http.StatusConflict,
	model.ErrCheckPassword:           http.StatusConflict,
	model.ErrShortPasswordOrNickname: http.StatusBadRequest,
	model.ErrUserNotFound:            http.StatusBadRequest,
	model.ErrCheckPassword:           http.StatusBadRequest,
}

func (h *UserHandler) ResponseUserError(w http.ResponseWriter, err error, ctx context.Context, data any) {
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
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.ResponseUserError(w, model.ErrJson, ctx, nil)
		return
	}
	err := h.service.RegisterCase(ctx, user)
	if err != nil {
		h.ResponseUserError(w, err, ctx, nil)
		return
	}
	dto.Response(w, "the user has been created", http.StatusOK, nil)
}
func (h *UserHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var passwordAndName model.User
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&passwordAndName); err != nil {
		h.ResponseUserError(w, model.ErrJson, ctx, nil)
		return
	}
	user, err := h.service.GetToken(ctx, passwordAndName.Nickname, passwordAndName.Password)
	if err != nil {
		h.ResponseUserError(w, err, ctx, user)
		return
	}
	dto.Response(w, "the user has been created", http.StatusOK, user)
}
