package order_usescase

import order_repo "github.com/patato8984/Shop/internal/modules/order/repo"

type OrderUserServise struct {
	repo *order_repo.OrderUserRepo
}

func NewOrderUserServise(repo *order_repo.OrderUserRepo) *OrderUserServise {
	return &OrderUserServise{repo: repo}
}
