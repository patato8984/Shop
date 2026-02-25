package cart_usescase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	cart_model "github.com/patato8984/Shop/internal/modules/cart/model"
	cart_repo "github.com/patato8984/Shop/internal/modules/cart/repo"
	"github.com/patato8984/Shop/internal/shared/cache"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"github.com/patato8984/Shop/internal/shared/outbox"
	"golang.org/x/sync/errgroup"
)

type CatalogProvider interface {
	GetStock(ctx context.Context, idSkus int) (int, error)
	WorkerGetPrice(ctx context.Context, idSkus int) (float64, error)
	SearchIdProduct(ctx context.Context, idSkus int) error
}
type CartService struct {
	provider CatalogProvider
	repo     cart_repo.CartRepo
	cache    cache.Cache
	kp       shared_events.EventPublisher
	tx       cart_repo.TxManager
	outbox   outbox.OutboxRepo
}

func NewCartService(provider CatalogProvider, repo *cart_repo.CartRepo, cache *cache.Cache, kp shared_events.EventPublisher, tx *cart_repo.TxManager, outbox *outbox.OutboxRepo) *CartService {
	return &CartService{provider: provider, repo: *repo, cache: *cache, kp: kp, tx: *tx, outbox: *outbox}
}
func (s *CartService) GetCart(ctx context.Context) (*cart_model.Cart, error) {
	var cart cart_model.Cart
	idCart, err := s.repo.SearchCart(ctx)
	if err != nil {
		return &cart, err
	}
	return cache.GetOrSet(ctx, s.cache, fmt.Sprintf("cart:%d", idCart), time.Minute*1, &cart, func() (*cart_model.Cart, error) {
		return s.repo.GetCart(ctx, idCart)
	})
}

func (s *CartService) AddProduct(ctx context.Context, idSkus, quantity int) (int, error) {
	id, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		var items cart_model.Items
		if err := items.ChangeQuantity(quantity); err != nil {
			return 0, err
		}
		idCart, err := s.repo.SearchCart(ctx)
		if err != nil {
			return 0, err
		}
		g, groupCtx := errgroup.WithContext(ctx)
		var (
			quantityCatalogs int
			priceCatalog     float64
		)
		g.Go(func() error {
			err := s.provider.SearchIdProduct(groupCtx, idSkus)
			return err
		})
		g.Go(func() error {
			quantityCatalog, err := s.provider.GetStock(groupCtx, idSkus)
			if err == nil {
				quantityCatalogs = quantityCatalog
			}
			return err
		})
		g.Go(func() error {
			price, err := s.provider.WorkerGetPrice(groupCtx, idSkus)
			if err == nil {
				priceCatalog = price
			}
			return err
		})
		if err := g.Wait(); err != nil {
			return 0, err
		}
		items.Quantity = quantityCatalogs
		items.Price_snapshot = priceCatalog
		price_snapshot, err := items.ChangeQuantityCatalogAndGetPrice(quantity)
		if err != nil {
			return 0, err
		}
		id, err := s.repo.AddProduct(ctx, idCart, idSkus, quantity, price_snapshot)
		if err != nil {
			return 0, err
		}
		event := cart_model.CartEvent{
			EventType: "product_add",
			PayLoad: cart_model.CartUpdate{
				IdCart: idCart,
			},
		}
		err = s.outbox.Add(ctx, "cart_event", event.EventType, event.PayLoad)
		return id, err
	})
	if err != nil {
		return 0, err
	}
	p, ok := id.(int)
	if !ok {
		return 0, errors.New("error convert type")
	}
	return p, err
}
func (s *CartService) ClearAllItems(ctx context.Context) (time.Time, error) {
	timeUpdate_at, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		var timeUpdate_at time.Time
		idCart, err := s.repo.SearchCart(ctx)
		if err != nil {
			return timeUpdate_at, err
		}
		err = s.repo.ClearAllItems(ctx, idCart)
		if err != nil {
			return timeUpdate_at, err
		}
		timeUpdate_at, err = s.repo.AddUpdate_at(ctx, idCart)
		if err != nil {
			return timeUpdate_at, err
		}
		event := cart_model.CartEvent{
			EventType: "clear_all_item",
			PayLoad: cart_model.CartUpdate{
				IdCart: idCart,
			},
		}
		err = s.outbox.Add(ctx, "cart_event", event.EventType, event.PayLoad)
		return timeUpdate_at, err
	})
	if err != nil {
		return time.Time{}, err
	}
	p, ok := timeUpdate_at.(time.Time)
	if !ok {
		return time.Time{}, errors.New("error convert type")
	}
	return p, err
}
func (s *CartService) AddOrUpdateQuantityProduct(ctx context.Context, idSkus, delta int) (int, error) {
	quantity, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		idCart, err := s.repo.SearchCart(ctx)
		if err != nil {
			return 0, err
		}
		idCartItems, err := s.repo.GetIdCartItems(ctx, idCart, idSkus)
		if err != nil {
			return 0, err
		}
		var items cart_model.Items
		items.Quantity, err = s.repo.GetQuantity(ctx, idCartItems, idSkus)
		if err != nil {
			return 0, err
		}
		err = items.ChangeQuantity(delta)
		if err != nil {
			return 0, err
		}
		err = s.repo.UpdateQuantity(ctx, idCart, idSkus, delta)
		if err != nil {
			return 0, err
		}
		timeUpdate_at, err := s.repo.AddUpdate_at(ctx, idCart)
		if err != nil {
			return 0, err
		}
		log.Print(timeUpdate_at)
		event := cart_model.CartEvent{
			EventType: "update_quantity_product",
			PayLoad: cart_model.CartUpdate{
				IdCart: idCart,
			},
		}
		err = s.outbox.Add(ctx, "cart_event", event.EventType, event.PayLoad)
		return items.Quantity, err
	})
	if err != nil {
		return 0, err
	}
	p, ok := quantity.(int)
	if !ok {
		return 0, errors.New("error convert type")
	}
	return p, err
}
func (s *CartService) GetSumCost(ctx context.Context) (*float64, error) {
	var sumCost float64
	idCart, err := s.repo.SearchCart(ctx)
	if err != nil {
		return &sumCost, err
	}
	return cache.GetOrSet(ctx, s.cache, fmt.Sprintf("cartSumCost:%d", idCart), time.Minute*1, &sumCost, func() (*float64, error) {
		return s.repo.GetAllCostProduct(ctx, idCart)
	})
}
func (s *CartService) DelProduct(ctx context.Context, idSkus int) (int, error) {
	idCart, err := s.repo.SearchCart(ctx)
	if err != nil {
		return 0, err
	}
	id, err := s.repo.DelProduct(ctx, idCart, idSkus)
	if err != nil {
		return 0, err
	}
	_ = s.cache.Del(ctx, fmt.Sprintf("cart:%d", idCart))
	return id, nil
}
