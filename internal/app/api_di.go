package app

import (
	"database/sql"

	auth_handler "github.com/patato8984/Shop/internal/modules/auth/handler"
	auth_repo "github.com/patato8984/Shop/internal/modules/auth/repo"
	auth_usescase "github.com/patato8984/Shop/internal/modules/auth/usescase"
	cart_handler "github.com/patato8984/Shop/internal/modules/cart/handler"
	cart_repo "github.com/patato8984/Shop/internal/modules/cart/repo"
	cart_usescase "github.com/patato8984/Shop/internal/modules/cart/usescase"
	catalog_handler "github.com/patato8984/Shop/internal/modules/catalog/handler"
	catalog_repo "github.com/patato8984/Shop/internal/modules/catalog/repo"
	catalog_usescase "github.com/patato8984/Shop/internal/modules/catalog/usescase"
	order_handler "github.com/patato8984/Shop/internal/modules/order/handler"
	order_repo "github.com/patato8984/Shop/internal/modules/order/repo"
	order_usescase "github.com/patato8984/Shop/internal/modules/order/usescase"
	"github.com/patato8984/Shop/internal/shared/cache"
	"github.com/patato8984/Shop/internal/shared/config"
	"github.com/patato8984/Shop/internal/shared/dto"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"github.com/patato8984/Shop/internal/shared/middleware"
	"github.com/patato8984/Shop/internal/shared/outbox"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

func DependencyInitiationApi(db *sql.DB, r *redis.Client, config config.ApiConfig, log *zap.Logger, kafka *shared_events.EventPublisher) (DiAndRouterMethods, StandaloneHandlers) {
	sg := singleflight.Group{}
	cache := cache.NewRedisCache(r, log)
	outboxRepo := outbox.NewOutboxRepo(db)
	catalogUserRepo := catalog_repo.NewCatalogRepo(db)
	catalogAdminRepo := catalog_repo.NewCatalogAdminRepo(db)
	catalogUserTx := catalog_repo.NewCatalogTxManager(db)
	catalogUserService := catalog_usescase.NewCatalogService(catalogUserRepo, cache, kafka, &sg)
	catalogAdminService := catalog_usescase.NewCatalogAdminService(catalogAdminRepo, &cache, kafka, catalogUserTx, outboxRepo)
	catalogUserHandler := catalog_handler.NewCatalogHandler(catalogUserService)
	catalogAdminHandler := catalog_handler.NewCatalogAdminHandler(catalogAdminService)
	cartUserRepo := cart_repo.NewCartRepo(db)
	cartUserTx := cart_repo.NewTxManager(db)
	cartUserService := cart_usescase.NewCartService(catalogUserRepo, cartUserRepo, &cache, *kafka, cartUserTx, outboxRepo)
	cartUserHandler := cart_handler.NewCartHandler(cartUserService)
	authUserRepo := auth_repo.NewUserRepo(db)
	authAdminRepo := auth_repo.NewAdminRepo(db)
	authUserService := auth_usescase.NewUserService(cartUserRepo, authUserRepo, config.JwtKey, *kafka)
	authAdminService := auth_usescase.NewAdminService(authAdminRepo)
	authAdminSeed := auth_usescase.NewSeedService(config.Admins, authAdminRepo)
	authUserHandler := auth_handler.NewUserHandler(authUserService)
	authAdminHandler := auth_handler.NewAdminHandler(authAdminService)
	bankApi := dto.NewBankApi()
	orderUserTx := order_repo.NewOrderTxManager(db)
	orderUserRepo := order_repo.NewOrderUserRepo(db)
	orderUserService := order_usescase.NewOrderUserService(orderUserRepo, *orderUserTx, cartUserRepo, catalogUserRepo, catalogAdminRepo, bankApi, config.SharedApiBankKey)
	orderUserHandler := order_handler.NewOrderUserHandler(orderUserService)
	orderWebHook := order_handler.NewWeBHookService(orderUserService)
	authorizationMiddleware := middleware.NewAuthorization(config.JwtKey)
	createdCtx := middleware.NewCreatedCtx(config.TimeOut)
	logger := middleware.NewLoggerMiddleware(log)
	workerCart := cart_usescase.NewWorkerService(catalogUserRepo, cartUserRepo, *kafka)
	workerOutbox := outbox.NewOutboxWorker(db, kafka, log)
	var diAndRouterMethods = DiAndRouterMethods{
		Logger:        *logger,
		Authorization: *authorizationMiddleware,
		CreatedCtx:    *createdCtx,
		AuthAdminH:    &authAdminHandler,
		AuthUserH:     *authUserHandler,
		CatalogAdminH: catalogAdminHandler,
		CatalogUserH:  *catalogUserHandler,
		CartUserH:     *cartUserHandler,
		OrderUserH:    *orderUserHandler,
	}
	var StandaloneHandlers = StandaloneHandlers{
		WebHookBank:       *orderWebHook,
		WorkerCartService: *workerCart,
		SeedAdmins:        *authAdminSeed,
		WorkerOutbox:      *workerOutbox,
	}
	return diAndRouterMethods, StandaloneHandlers
}
