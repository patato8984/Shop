package catalog

import (
	"database/sql"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
)

type CatalogRepo struct {
	db *sql.DB
}

func NewCatalogRepo(db *sql.DB) *CatalogRepo {
	return &CatalogRepo{db: db}
}
func (c CatalogRepo) GetAllProducts() []model.Product {

}
