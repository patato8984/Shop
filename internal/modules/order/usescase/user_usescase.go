package order_usescase

import (
	"context"
	"log"
	"time"

	order_model "github.com/patato8984/Shop/internal/modules/order/model"
	order_repo "github.com/patato8984/Shop/internal/modules/order/repo"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type CatalogProvider interface {
	WorkerGetPrice(ctx context.Context, idSkus int) (float64, error)
	GetStock(ctx context.Context, idSkus int) (int, error)
}
type CatalogAdminProvider interface {
	UpdateStock(ctx context.Context, delta, id int) (int, int, time.Time, error)
}
type CartProvider interface {
	GetNumberItemsSkus(ctx context.Context, idCart int) (int, error)
	GetAllIdItemsCart(ctx context.Context, idCart int) ([]int, error)
	SearchCart(ctx context.Context) (int, error)
	GetAllCostProduct(ctx context.Context, idCart int) (*float64, error)
	GetAllIdSkusAndQuantity(ctx context.Context, idCart int) (map[int]int, error)
	ClearAllItems(ctx context.Context, idCart int) error
}
type EmulatorBankApiProvider interface {
	CreateTransaction(ctx context.Context, idOrder int, amount float64) (string, string, error)
}
type OrderService struct {
	catalogAdminProvider CatalogAdminProvider
	catalogProvider      CatalogProvider
	cartProvider         CartProvider
	apiBankProvider      EmulatorBankApiProvider
	repo                 *order_repo.OrderUserRepo
	TxManager            *order_repo.TxManager
	sharedSecret         string
}

func NewOrderUserService(repo *order_repo.OrderUserRepo, txManager order_repo.TxManager, cartProvider CartProvider, catalogProvider CatalogProvider, catalogAdminProvider CatalogAdminProvider, apiBankProvider EmulatorBankApiProvider, sharedSecret string) *OrderService {
	return &OrderService{
		repo:                 repo,
		TxManager:            &txManager,
		cartProvider:         cartProvider,
		catalogProvider:      catalogProvider,
		catalogAdminProvider: catalogAdminProvider,
		apiBankProvider:      apiBankProvider,
		sharedSecret:         sharedSecret,
	}
}

// requires refactoring
func (s *OrderService) AddNowOrder(ctx context.Context, address, message string) (int, error) {
	_, err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		idCart, err := s.cartProvider.SearchCart(ctx)
		if err != nil {
			return nil, err
		}
		quantity, err := s.cartProvider.GetNumberItemsSkus(ctx, idCart)
		if err != nil {
			return nil, err
		}
		if quantity <= 0 {
			return nil, order_model.ErrQuantityItemsCart
		}
		idOrder, _, err := s.repo.AddOrder(ctx, address, message) // ligation
		if err != nil {
			return nil, err
		}
		allCost, err := s.cartProvider.GetAllCostProduct(ctx, idCart)
		if err != nil {
			return nil, err
		}
		idSkusAndQuantity, err := s.cartProvider.GetAllIdSkusAndQuantity(ctx, idCart)
		if err != nil {
			return nil, err
		}
		var allQuantity int
		for idSkus, quantity := range idSkusAndQuantity {
			stock, err := s.catalogProvider.GetStock(ctx, idSkus)
			if err != nil {
				return nil, err
			}
			if stock < quantity {
				return nil, order_model.ErrStock
			}
			price, err := s.catalogProvider.WorkerGetPrice(ctx, idSkus)
			if err != nil {
				return nil, err
			}
			err = s.repo.AddOrderItems(ctx, idOrder, idSkus, quantity, price)
			if err != nil {
				return nil, err
			}
			allQuantity += quantity
			_, _, _, err = s.catalogAdminProvider.UpdateStock(ctx, -quantity, idSkus)
			if err != nil {
				return nil, err
			}
		}
		_, err = s.repo.FinalUpdate(ctx, allQuantity, *allCost, "CREATED", idOrder)
		if err != nil {
			return nil, err
		}
		err = s.cartProvider.ClearAllItems(ctx, idCart)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return 0, err
	}
	idCart, err := s.cartProvider.SearchCart(ctx)
	if err != nil {
		return 0, err
	}
	return idCart, nil
}
func (s *OrderService) GetAllOrders(ctx context.Context) ([]order_model.Order, error) {
	orders, err := s.repo.GetAllOrdersBaseInfo(ctx)
	if err != nil {
		return orders, err
	}
	return orders, nil
}
func (s *OrderService) GetOrderAndAllItems(ctx context.Context, idOrder int) (order_model.Order, error) {
	var (
		order order_model.Order
		items []order_model.Order_items
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		orders, err := s.repo.GetOrder(ctx, idOrder)
		if err == nil {
			order = orders
		}
		return err
	})
	g.Go(func() error {
		item, err := s.repo.GetAllOrdersItems(ctx, idOrder)
		if err == nil {
			items = item
		}
		return err
	})
	if err := g.Wait(); err != nil {
		return order, err
	}
	order.Items = items
	return order, nil
}

func (s *OrderService) GetPayment(ctx context.Context, idOrder int) (order_model.BankTransactionResponse, error) {
	var plug order_model.BankTransactionResponse
	statuse, err := s.repo.GetStatuseOrder(ctx, idOrder)
	plug.Status = statuse
	if err != nil {
		return plug, err
	}
	switch plug.Status {
	case "PAID":
		return plug, order_model.ErrStatusPey
	case "CANCELLED":
		return plug, order_model.ErrStatusCancelled
	case "WAITING_PAYMENT":
		url, err := s.repo.GetUrl(ctx, idOrder)
		log.Print(url)
		log.Print(err)
		if err == nil {
			plug.BankURL = url
			return plug, nil
		}
	}
	amount, err := s.repo.GetPriceAllProduct(ctx, idOrder)
	if err != nil {
		return plug, err
	}
	bankRequest, err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		var bankRequest order_model.BankTransactionResponse
		url, idPayment, err := s.apiBankProvider.CreateTransaction(ctx, idOrder, amount)
		if err != nil {
			return bankRequest, order_model.ErrCreateTransaction
		}
		log.Print(idPayment)
		err = s.repo.AddBankPaymentId(ctx, idOrder, idPayment)
		if err != nil {
			return bankRequest, err
		}
		err = s.repo.AddOrUpdatedUrl(ctx, url, idOrder)
		if err != nil {
			return bankRequest, err
		}
		statuse, err := s.repo.UpdateStatusOrder(ctx, idOrder, "WAITING_PAYMENT")
		if err != nil {
			return bankRequest, err
		}
		bankRequest.BankURL = url
		bankRequest.Status = statuse
		return bankRequest, nil
	})
	res, ok := bankRequest.(order_model.BankTransactionResponse)
	if !ok {
		return plug, err
	}
	return res, nil
}

func (s *OrderService) WeBhookBankPayment(ctx context.Context, bankApiData order_model.BankApiRequest) error {
	log := logger.FromContext(ctx)
	if bankApiData.PaymentId == "" {
		return order_model.ErrHash
	}
	idOrder, err := s.repo.GetIdOrderFromByBankId(ctx, bankApiData.PaymentId)
	if err != nil {
		return err
	}
	status, err := s.repo.GetStatuseOrder(ctx, idOrder)
	if err != nil {
		return err
	}
	if status == "PAID" {
		return nil
	}
	paymentId, err := s.repo.GetBankPaymentId(ctx, idOrder)
	if err != nil {
		return err
	}
	err = bankApiData.HashComparison(paymentId, s.sharedSecret)
	if err != nil {
		log.Error("error hash Comparison",
			zap.Error(err),
		)
		return order_model.ErrHash
	}
	_, err = s.repo.UpdateStatusOrder(ctx, idOrder, "PAID")
	if err != nil {
		return err
	}
	return nil
}
