package app

import (
	manager "github.com/patato8984/Shop/internal/modules/cacheInvalidator"
	"github.com/patato8984/Shop/internal/shared/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func DependencyInitiationCacheInvalidator(r *redis.Client, log *zap.Logger) manager.InvalidatorManager {
	cache := cache.NewRedisCache(r, log)
	invalidationManager := manager.NewInvalidationManager(cache)
	return *invalidationManager
}
