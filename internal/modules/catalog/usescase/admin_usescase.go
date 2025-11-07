package usescase

import (
	"errors"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
	"github.com/patato8984/Shop/internal/modules/catalog/repo"
)

type CatalogAdminServise struct {
	repo *repo.CatalogAdminRepo
}

func NewCatalogAdminServise(repo *repo.CatalogAdminRepo) CatalogAdminServise {
	return CatalogAdminServise{repo: repo}
}

func (s CatalogAdminServise) CreateNewProduct(nameProduct string) (model.Product, error) {
	var product model.Product
	if len(nameProduct) < 4 {
		return product, errors.New("short name")
	}
	id, err := s.repo.CreateOrSearchProduct(nameProduct)
	if err != nil {
		return product, err
	}
	product.Id = id
	product.Name = nameProduct
	return product, nil
}
func (s CatalogAdminServise) CreateNewSkus(id int, sku model.SKU) (model.Product, error) {
	return s.repo.CreateSkus(id, sku)
}
