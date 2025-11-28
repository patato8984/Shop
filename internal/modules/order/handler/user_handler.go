package order_handler

import order_usescase "github.com/patato8984/Shop/internal/modules/order/usescase"

type OrderUserHandler struct {
	servise *order_usescase.OrderUserServise
}

func NewOrderUserHandler(servise *order_usescase.OrderUserServise) *OrderUserHandler {
	return &OrderUserHandler{servise: servise}
}
