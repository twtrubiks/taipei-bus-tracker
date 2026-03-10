package handler

import (
	"context"
	"strconv"

	"github.com/twtrubiks/taipei-bus-tracker/internal/cache"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

// FallbackService wraps primary and fallback data sources with cache.
type FallbackService struct {
	Primary  model.BusDataSource
	Fallback model.BusDataSource
	Cache    *cache.Cache
}

func NewFallbackService(primary, fallback model.BusDataSource, c *cache.Cache) *FallbackService {
	return &FallbackService{
		Primary:  primary,
		Fallback: fallback,
		Cache:    c,
	}
}

func (f *FallbackService) GetETA(ctx context.Context, city, routeID string, direction int) ([]model.StopETA, error) {
	cacheKey := "eta:" + city + ":" + routeID + ":" + strconv.Itoa(direction)

	// Try primary
	etas, err := f.Primary.GetETA(ctx, city, routeID, direction)
	if err == nil {
		f.Cache.Set(cacheKey, etas)
		return etas, nil
	}

	// Try fallback
	if f.Fallback != nil {
		etas, err = f.Fallback.GetETA(ctx, city, routeID, direction)
		if err == nil {
			f.Cache.Set(cacheKey, etas)
			return etas, nil
		}
	}

	// Try cache
	if cached, ok := f.Cache.Get(cacheKey); ok {
		if etas, ok := cached.([]model.StopETA); ok {
			return etas, nil
		}
	}

	return nil, err
}
