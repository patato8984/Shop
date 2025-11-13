package usescase

import (
	"errors"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
	"github.com/patato8984/Shop/internal/modules/catalog/repo"
)

type CatalogService struct {
	repo repo.CatalogRepo
}

func NewCatalogService(repo *repo.CatalogRepo) *CatalogService {
	return &CatalogService{repo: *repo}
}
func (s CatalogService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAll()
}
func (s CatalogService) GetSkus(id int) (model.SKU, error) {
	var skus model.SKU
	if id <= 0 {
		return skus, errors.New("short id")
	}
	return s.repo.GetSkus(id)
}
func (s CatalogService) GetAllSkus(id int) (model.Product, error) {
	var skus model.Product
	if id <= 0 {
		return skus, errors.New("short id")
	}
	return s.repo.GetAllSkus(id)
}
