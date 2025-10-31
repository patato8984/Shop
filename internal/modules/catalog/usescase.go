package catalog

type CatalogService struct {
	repo CatalogRepo
}

func NewCatalogService(repo *CatalogRepo) *CatalogService {
	return &CatalogService{repo: *repo}
}
