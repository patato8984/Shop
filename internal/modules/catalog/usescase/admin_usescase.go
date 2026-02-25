package catalog_usescase

import (
	"context"
	"errors"
	"log"

	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	catalog_admin_repo "github.com/patato8984/Shop/internal/modules/catalog/repo"
	"github.com/patato8984/Shop/internal/shared/cache"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"github.com/patato8984/Shop/internal/shared/outbox"
	"github.com/patato8984/Shop/pkg/ctxkey"
)

type CatalogAdminService struct {
	repo   *catalog_admin_repo.CatalogAdminRepo
	cache  cache.Cache
	kp     *shared_events.EventPublisher
	tx     *catalog_admin_repo.TxManager
	outbox *outbox.OutboxRepo
}

func NewCatalogAdminService(repo *catalog_admin_repo.CatalogAdminRepo, cache *cache.Cache, kp *shared_events.EventPublisher, tx *catalog_admin_repo.TxManager, outbox *outbox.OutboxRepo) *CatalogAdminService {
	return &CatalogAdminService{repo: repo, cache: *cache, kp: kp, tx: tx, outbox: outbox}
}

func (s CatalogAdminService) CreateNewProduct(ctx context.Context, nameProduct string) (catalog_model.Product, error) {
	var product catalog_model.Product
	if len(nameProduct) < 4 {
		return product, catalog_model.ErrShortName
	}
	products, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		id, err := s.repo.CreateOrSearchProduct(ctx, nameProduct)
		if err != nil {
			return product, err
		}
		product.Id = id
		product.Name = nameProduct
		event := catalog_model.CatalogEvent{
			EventType: "product_create",
			MetaDate: catalog_model.MetaDate{
				IDUser:   ctx.Value(ctxkey.UserIDKey).(int),
				RoleUser: ctx.Value(ctxkey.Role).(string),
			},
			PayLoad: product,
		}
		err = s.outbox.Add(ctx, "catalog_event", event.EventType, event)
		return product, err
	})
	if err != nil {
		return product, err
	}
	p, ok := products.(catalog_model.Product)
	if !ok {
		return catalog_model.Product{}, errors.New("error convert type")
	}
	return p, err
}
func (s CatalogAdminService) DelProduct(ctx context.Context, id int) (int, error) {
	if id <= 0 {
		return 0, catalog_model.ErrShortID
	}
	iD, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		id, err := s.repo.DelProduct(ctx, id)
		if err != nil {
			return id, err
		}
		event := catalog_model.CatalogEvent{
			EventType: "product_delete",
			MetaDate: catalog_model.MetaDate{
				IDUser:   ctx.Value(ctxkey.UserIDKey).(int),
				RoleUser: ctx.Value(ctxkey.Role).(string),
			},
			PayLoad: catalog_model.Product{
				Id: id,
			},
		}
		err = s.outbox.Add(ctx, "catalog_event", event.EventType, event)
		return id, err
	})
	p, ok := iD.(int)
	if !ok {
		return 0, errors.New("error convert type")
	}
	return p, err
}
func (s CatalogAdminService) CreateNewSkus(ctx context.Context, id int, sku catalog_model.SKU) (catalog_model.Product, error) {
	if id <= 0 {
		var mod catalog_model.Product
		return mod, catalog_model.ErrShortID
	}
	product, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		product, err := s.repo.CreateSkus(ctx, id, sku)
		if err != nil {
			return product, err
		}
		event := catalog_model.CatalogEvent{
			EventType: "skus_created",
			MetaDate: catalog_model.MetaDate{
				IDUser:   ctx.Value(ctxkey.UserIDKey).(int),
				RoleUser: ctx.Value(ctxkey.Role).(string),
			},
			PayLoad: product,
		}
		err = s.outbox.Add(ctx, "catalog_event", event.EventType, event)
		return product, err
	})
	if err != nil {
		return catalog_model.Product{}, err
	}
	p, ok := product.(catalog_model.Product)
	if !ok {
		return catalog_model.Product{}, errors.New("error convert type")
	}
	return p, err
}
func (s CatalogAdminService) AddStockToSkus(ctx context.Context, stock, id int) (int, error) {
	if id <= 0 {
		return 0, catalog_model.ErrShortID
	}
	newStock, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		idProduct, newStock, update_at, err := s.repo.UpdateStock(ctx, stock, id)
		if err != nil {
			return 0, err
		}
		stock, err = s.repo.GetStock(ctx, id)
		if err != nil {
			return 0, err
		}
		event := catalog_model.CatalogEvent{
			EventType: "skus_addStock",
			MetaDate: catalog_model.MetaDate{
				IDUser:   ctx.Value(ctxkey.UserIDKey).(int),
				RoleUser: ctx.Value(ctxkey.Role).(string),
			},
			PayLoad: catalog_model.StockUpdatedLoad{
				ProductID: idProduct,
				SkusID:    id,
				NewStock:  stock,
			},
		}
		err = s.kp.Publisher(ctx, "catalog_event", event.EventType, event)
		log.Print(update_at)
		return newStock, err
	})
	if err != nil {
		return 0, err
	}
	p, ok := newStock.(int)
	if !ok {
		return 0, errors.New("error convert type")
	}
	return p, err
}
func (s CatalogAdminService) UpdatePriceToSkus(ctx context.Context, idSkus int, price float64) (float64, error) {
	newPrice, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		idProduct, price, err := s.repo.UpdatePriceToSkus(ctx, idSkus, price)
		if err != nil {
			return price, err
		}
		event := catalog_model.CatalogEvent{
			EventType: "skus_price_update",
			MetaDate: catalog_model.MetaDate{
				IDUser:   ctx.Value(ctxkey.UserIDKey).(int),
				RoleUser: ctx.Value(ctxkey.Role).(string),
			},
			PayLoad: catalog_model.PriceUpdatedLoad{
				ProductID: idProduct,
				SkusID:    idSkus,
				NewPrice:  price,
			},
		}
		err = s.outbox.Add(ctx, "catalog_event", event.EventType, event)
		return price, err
	})
	if err != nil {
		return 0, err
	}
	p, ok := newPrice.(float64)
	if !ok {
		return 0, errors.New("error convert type")
	}
	return p, err
}
func (s CatalogAdminService) DelSkus(ctx context.Context, idSkus int) (int, error) {
	if idSkus <= 0 {
		return 0, catalog_model.ErrShortID
	}
	idSkuse, err := s.tx.WithinTransaction(ctx, func(ctx context.Context) (any, error) {
		idProduct, idSkus, err := s.repo.DelSkus(ctx, idSkus)
		if err != nil {
			return idSkus, err
		}
		event := catalog_model.CatalogEvent{
			EventType: "skus_deleted",
			MetaDate: catalog_model.MetaDate{
				IDUser:   ctx.Value(ctxkey.UserIDKey).(int),
				RoleUser: ctx.Value(ctxkey.Role).(string),
			},
			PayLoad: catalog_model.SkusDeleteLoad{
				ProductID: idProduct,
				SkusID:    idSkus,
			},
		}
		err = s.outbox.Add(ctx, "catalog_event", event.EventType, event)
		return idProduct, err
	})
	if err != nil {
		return 0, err
	}
	p, ok := idSkuse.(int)
	if !ok {
		return 0, errors.New("error convert type")
	}
	return p, err
}
