package app

import (
	auth_handler "github.com/patato8984/Shop/internal/modules/auth/handler"
	auth_usescase "github.com/patato8984/Shop/internal/modules/auth/usescase"
	cart_handler "github.com/patato8984/Shop/internal/modules/cart/handler"
	cart_usescase "github.com/patato8984/Shop/internal/modules/cart/usescase"
	catalog_handler "github.com/patato8984/Shop/internal/modules/catalog/handler"
	order_handler "github.com/patato8984/Shop/internal/modules/order/handler"
	"github.com/patato8984/Shop/internal/shared/middleware"
	"github.com/patato8984/Shop/internal/shared/outbox"
)

type DiAndRouterMethods struct {
	Logger        middleware.Logger
	Authorization middleware.Authorization
	CreatedCtx    middleware.CreatedCtx
	AuthAdminH    *auth_handler.AdminHandler
	AuthUserH     auth_handler.UserHandler
	CatalogAdminH catalog_handler.CatalogAdminHandler
	CatalogUserH  catalog_handler.CatalogHandler
	CartUserH     cart_handler.CartHandler
	OrderUserH    order_handler.OrderUserHandler
}

type StandaloneHandlers struct {
	WebHookBank       order_handler.WeBHookHandler
	WorkerCartService cart_usescase.PriceUpdaterService
	SeedAdmins        auth_usescase.SeedService
	WorkerOutbox      outbox.Worker
}
