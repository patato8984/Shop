package app

import (
	"net/http"
)

type App struct {
	RouterMethods DiAndRouterMethods
}

func NewApp(routerMethods DiAndRouterMethods) *App {
	return &App{RouterMethods: routerMethods}
}
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
func (a *App) RegisterAuthRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("POST /api/v1/admin/registration", Chain(http.HandlerFunc(a.RouterMethods.AuthAdminH.CreateNewAdmin), a.RouterMethods.Logger.Response, a.RouterMethods.Authorization.AuthenticationAdmin, a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("POST /api/v1/user/registration", Chain(http.HandlerFunc(a.RouterMethods.AuthUserH.RegisterUser), a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/authentication", Chain(http.HandlerFunc(a.RouterMethods.AuthUserH.Authentication), a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("POST /api/v1/admin/catalog", Chain(http.HandlerFunc(a.RouterMethods.CatalogAdminH.CreateNewProduct), a.RouterMethods.Authorization.AuthenticationAdmin, a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("DELETE /api/v1/admin/catalog/{id}", Chain(http.HandlerFunc(a.RouterMethods.CatalogAdminH.DelProduct), a.RouterMethods.Authorization.AuthenticationAdmin, a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("PUT /api/v1/admin/catalog/skus/{id}/price", Chain(http.HandlerFunc(a.RouterMethods.CatalogAdminH.UpdatePrice), a.RouterMethods.Authorization.AuthenticationAdmin, a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("POST /api/v1/admin/catalog/{id}/skus", Chain(http.HandlerFunc(a.RouterMethods.CatalogAdminH.CreateNewSkus), a.RouterMethods.Authorization.AuthenticationAdmin, a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("POST /api/v1/admin/catalog/skus/{id}/stock", Chain(http.HandlerFunc(a.RouterMethods.CatalogAdminH.AddStockToSkus), a.RouterMethods.Authorization.AuthenticationAdmin, a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("DELETE /api/v1/admin/catalog/skus/{id}", Chain(http.HandlerFunc(a.RouterMethods.CatalogAdminH.DelSkus), a.RouterMethods.Authorization.AuthenticationAdmin, a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/catalog", Chain(http.HandlerFunc(a.RouterMethods.CatalogUserH.GetAllProduct), a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/catalog/skus/{id}", Chain(http.HandlerFunc(a.RouterMethods.CatalogUserH.GetSkus), a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/catalog/skus/detail/{id}", Chain(http.HandlerFunc(a.RouterMethods.CatalogUserH.GetAllSkus), a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("POST /api/v1/user/cart", Chain(http.HandlerFunc(a.RouterMethods.CartUserH.AddProduct), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/cart", Chain(http.HandlerFunc(a.RouterMethods.CartUserH.GetCart), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("DELETE /api/v1/user/cart/product/{id}", Chain(http.HandlerFunc(a.RouterMethods.CartUserH.DelProduct), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/cart/product/{id}/sum", Chain(http.HandlerFunc(a.RouterMethods.CartUserH.GetSumCost), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("PUT /api/v1/user/cart/product/{id}/quantity", Chain(http.HandlerFunc(a.RouterMethods.CartUserH.UpdateQuantityProduct), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("POST /api/v1/user/order", Chain(http.HandlerFunc(a.RouterMethods.OrderUserH.AddNewOrder), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/order", Chain(http.HandlerFunc(a.RouterMethods.OrderUserH.GetAllOrders), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("GET /api/v1/user/order/{id}", Chain(http.HandlerFunc(a.RouterMethods.OrderUserH.GetOrderAndAllItems), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	mux.Handle("POST /api/v1/user/order/{id}/payment", Chain(http.HandlerFunc(a.RouterMethods.OrderUserH.Payment), a.RouterMethods.Authorization.Authorization, a.RouterMethods.Logger.Response, a.RouterMethods.Logger.RequestID, a.RouterMethods.CreatedCtx.CreatedCtxTimeout))
	return mux
}
