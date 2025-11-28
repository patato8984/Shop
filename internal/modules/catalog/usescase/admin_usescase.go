package catalog_usescase

import (
	"errors"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
	catalog_admin_repo "github.com/patato8984/Shop/internal/modules/catalog/repo"
)

type CatalogAdminServise struct {
	repo *catalog_admin_repo.CatalogAdminRepo
}

func NewCatalogAdminServise(repo *catalog_admin_repo.CatalogAdminRepo) CatalogAdminServise {
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
func (s CatalogAdminServise) DelProduct(id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("short id")
	}
	return s.repo.DelProduct(id)
}
func (s CatalogAdminServise) CreateNewSkus(id int, sku model.SKU) (model.Product, error) {
	if id <= 0 {
		var model model.Product
		return model, errors.New("short id")
	}
	return s.repo.CreateSkus(id, sku)
}
func (s CatalogAdminServise) AddStockToSkus(stock, id int) (model.SKU, error) {
	if id <= 0 {
		var model model.SKU
		return model, errors.New("short id")
	}
	return s.repo.AddStock(stock, id)
}

func (s CatalogAdminServise) DelSkus(id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("short id")
	}
	return s.repo.DelSkus(id)
}
