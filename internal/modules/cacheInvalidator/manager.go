package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	cart_model "github.com/patato8984/Shop/internal/modules/cart/model"
	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	"github.com/patato8984/Shop/internal/shared/cache"
)

type EventHandler func(ctx context.Context, payload []byte) error
type InvalidatorManager struct {
	Handler map[string]EventHandler
	Cache   cache.Cache
}

func NewInvalidationManager(c cache.Cache) *InvalidatorManager {
	m := &InvalidatorManager{
		Cache:   c,
		Handler: make(map[string]EventHandler),
	}
	m.Handler["product_add"] = m.handleUpdateCart
	m.Handler["cartUpdateWorker"] = m.handleUpdateCart
	m.Handler["clear_all_item"] = m.handleUpdateCart
	m.Handler["update_quantity_product"] = m.handleUpdateCart
	m.Handler["product_create"] = m.HandlerProductUpdate
	m.Handler["product_delete"] = m.HandlerProductUpdate
	m.Handler["skus_created"] = m.HandlerSkusCreate
	m.Handler["skus_addStock"] = m.HandlerSkusStockUpdate
	m.Handler["skus_price_update"] = m.HandlerSkusPriceUpdate
	m.Handler["skus_deleted"] = m.HandlerSkusDeleted
	return m
}
func (m *InvalidatorManager) HandlerCartUpdate(ctx context.Context, data []byte) error {
	var cart cart_model.CartUpdate
	if err := json.Unmarshal(data, &cart); err != nil {
		return err
	}
	var key []string
	key = append(key, fmt.Sprintf("cart:%d", cart.IdCart))
	key = append(key, fmt.Sprintf("cartSumCost:%d", cart.IdCart))
	if err := m.Cache.PipelineDelCart(ctx, key); err != nil {
		return err
	}
	return nil
}

func (m *InvalidatorManager) handleUpdateCart(ctx context.Context, data []byte) error {
	var p cart_model.CartUpdate
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	var key []string
	key = append(key, fmt.Sprintf("cart:%d", p.IdCart))
	key = append(key, fmt.Sprintf("cartSumCost:%d", p.IdCart))
	if err := m.Cache.PipelineDelCart(ctx, key); err != nil {
		return err
	}
	return nil
}
func (m *InvalidatorManager) HandlerProductUpdate(ctx context.Context, data []byte) error {
	var idProduct catalog_model.Product
	if err := json.Unmarshal(data, &idProduct); err != nil {
		return err
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("allSkus:%d", idProduct.Id)); !ok {
		return errors.New(fmt.Sprintf("error deleted allSkus:%d", idProduct.Id))
	}
	if ok := m.Cache.Del(ctx, "AllProducts"); !ok {
		return errors.New(fmt.Sprintf("error deleted allProduct"))
	}
	return nil
}
func (m *InvalidatorManager) HandlerSkusCreate(ctx context.Context, data []byte) error {
	var product catalog_model.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return err
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("Skus:%d", product.SKUs[0].Id)); !ok {
		return errors.New(fmt.Sprintf("error deleted Skus:%d", product.SKUs[0].Id))
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("allSkus:%d", product.Id)); !ok {
		return errors.New(fmt.Sprintf("error deleted allProduct"))
	}
	return nil
}
func (m *InvalidatorManager) HandlerSkusPriceUpdate(ctx context.Context, data []byte) error {
	var priceUpdate catalog_model.PriceUpdatedLoad
	if err := json.Unmarshal(data, &priceUpdate); err != nil {
		return err
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("Skus:%d", priceUpdate.ProductID)); !ok {
		return errors.New(fmt.Sprintf("error deleted Skus:%d", priceUpdate.ProductID))
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("allSkus:%d", priceUpdate.ProductID)); !ok {
		return errors.New(fmt.Sprintf("error deleted allSkus:%d", priceUpdate.SkusID))
	}
	return nil
}
func (m *InvalidatorManager) HandlerSkusDeleted(ctx context.Context, data []byte) error {
	var skusDelete catalog_model.SkusDeleteLoad
	if err := json.Unmarshal(data, &skusDelete); err != nil {
		return err
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("Skus:%d", skusDelete.ProductID)); !ok {
		return errors.New(fmt.Sprintf("error deleted Skus:%d", skusDelete.ProductID))
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("allSkus:%d", skusDelete.ProductID)); !ok {
		return errors.New(fmt.Sprintf("error deleted allSkus:%d", skusDelete.SkusID))
	}
	return nil
}
func (m *InvalidatorManager) HandlerSkusStockUpdate(ctx context.Context, data []byte) error {
	var stockUpdate catalog_model.StockUpdatedLoad
	if err := json.Unmarshal(data, &stockUpdate); err != nil {
		return err
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("Skus:%d", stockUpdate.ProductID)); !ok {
		return errors.New(fmt.Sprintf("error deleted Skus:%d", stockUpdate.ProductID))
	}
	if ok := m.Cache.Del(ctx, fmt.Sprintf("allSkus:%d", stockUpdate.ProductID)); !ok {
		return errors.New(fmt.Sprintf("error deleted allSkus:%d", stockUpdate.SkusID))
	}
	return nil
}
