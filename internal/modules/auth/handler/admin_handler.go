package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/patato8984/Shop/internal/modules/auth/model"
	usescase_user "github.com/patato8984/Shop/internal/modules/auth/usescase"
	"github.com/patato8984/Shop/internal/shared/dto"
)

type AdminHandler struct {
	service *usescase_user.AdminServise
}

func NewAdminHandler(servise *usescase_user.AdminServise) AdminHandler {
	return AdminHandler{service: servise}
}
func (h *AdminHandler) CreateNewAdmin(w http.ResponseWriter, r http.Request) {
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("error json", http.StatusBadRequest, nil))
		return
	}
	if err := h.service.CreateNewAdmin(user); err != nil {
		switch err.Error() {
		case "short password or nickname":
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusBadRequest, nil))
			return
		case "nickname busy":
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusConflict, nil))
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusBadRequest, nil))
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response("the admin has been created", http.StatusBadRequest, nil))
}
