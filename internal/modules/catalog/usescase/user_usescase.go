package usescase

import (
	modelscatalog "github.com/patato8984/Shop/internal/modules/catalog/model"
	"github.com/patato8984/Shop/internal/modules/catalog/repo"
)

type CatalogService struct {
	repo repo.CatalogRepo
}

func NewCatalogService(repo *repo.CatalogRepo) *CatalogService {
	return &CatalogService{repo: *repo}
}
func (s CatalogService) GetAllProducts() ([]modelscatalog.Product, error) {
	return s.repo.GetAll()
}
