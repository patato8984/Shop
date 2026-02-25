package catalog_usescase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	catalog_repo "github.com/patato8984/Shop/internal/modules/catalog/repo"
	"github.com/patato8984/Shop/internal/shared/cache"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"golang.org/x/sync/singleflight"
)

type CatalogService struct {
	repo    catalog_repo.CatalogRepo
	cache   cache.Cache
	sfGroup *singleflight.Group
}

func NewCatalogService(repo *catalog_repo.CatalogRepo, cache cache.Cache, kp *shared_events.EventPublisher, sg *singleflight.Group) *CatalogService {
	return &CatalogService{repo: *repo, sfGroup: sg, cache: cache}
}
func (s *CatalogService) GetAllProducts(ctx context.Context) (*[]catalog_model.Product, error) {
	var product []catalog_model.Product
	data, err, _ := s.sfGroup.Do("AllProducts", func() (interface{}, error) {
		return cache.GetOrSet(ctx, s.cache, "AllProducts", time.Minute*1, &product, func() (*[]catalog_model.Product, error) {
			return s.repo.GetAll(ctx)
		})
	})
	if err != nil {
		return &product, err
	}
	products, ok := data.(*[]catalog_model.Product)
	if !ok {
		return &product, errors.New("error type")
	}
	return products, nil
}
func (s *CatalogService) GetSkus(ctx context.Context, id int) (*catalog_model.SKU, error) {
	var skus catalog_model.SKU
	if id <= 0 {
		return &skus, catalog_model.ErrShortID
	}
	return cache.GetOrSet(ctx, s.cache, fmt.Sprintf("Skus:%d", id), time.Minute*3, &skus, func() (*catalog_model.SKU, error) {
		return s.repo.GetSkus(ctx, id)
	})
}
func (s *CatalogService) GetAllSkus(ctx context.Context, id int) (*catalog_model.Product, error) {
	if id <= 0 {
		return &catalog_model.Product{}, catalog_model.ErrShortID
	}
	data, err, _ := s.sfGroup.Do(strconv.Itoa(id), func() (interface{}, error) {
		var skus catalog_model.Product
		return cache.GetOrSet(ctx, s.cache, fmt.Sprintf("allSkus:%d", id), 1*time.Minute, &skus, func() (*catalog_model.Product, error) {
			return s.repo.GetAllSkus(ctx, id)
		})
	})
	if err != nil {
		return &catalog_model.Product{}, err
	}
	skus, ok := data.(*catalog_model.Product)
	if !ok {
		return skus, errors.New("error type")
	}
	return skus, nil
}
