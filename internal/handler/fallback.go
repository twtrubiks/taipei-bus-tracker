package handler

import (
	"context"
	"fmt"
	"log"
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

// fallbackFetch is a generic cache-first fallback fetcher.
func fallbackFetch[T any](
	c *cache.Cache,
	cacheKey string,
	label string,
	primaryFn func() ([]T, error),
	fallbackFn func() ([]T, error),
) ([]T, error) {
	// 1. Fresh cache
	if cached, ok := c.Get(cacheKey); ok {
		if items, ok := cached.([]T); ok && len(items) > 0 {
			return items, nil
		}
	}

	// 2. Primary
	items, err := primaryFn()
	if err == nil && len(items) > 0 {
		c.Set(cacheKey, items)
		return items, nil
	}
	if err != nil {
		log.Printf("[fallback] %s primary failed: %v", label, err)
	}

	// 3. Fallback
	if fallbackFn != nil {
		items, err = fallbackFn()
		if err == nil && len(items) > 0 {
			c.Set(cacheKey, items)
			return items, nil
		}
		if err != nil {
			log.Printf("[fallback] %s fallback failed: %v", label, err)
		}
	}

	// 4. Stale cache
	if cached, ok := c.GetStale(cacheKey); ok {
		if items, ok := cached.([]T); ok && len(items) > 0 {
			log.Printf("[fallback] %s: using stale cache", label)
			return items, nil
		}
	}

	return items, err
}

func (f *FallbackService) SearchRoutes(ctx context.Context, city, keyword string) ([]model.Route, error) {
	cacheKey := "routes:" + city + ":" + keyword
	var fb func() ([]model.Route, error)
	if f.Fallback != nil {
		fb = func() ([]model.Route, error) { return f.Fallback.SearchRoutes(ctx, city, keyword) }
	}
	return fallbackFetch(f.Cache, cacheKey, fmt.Sprintf("SearchRoutes(%s, %s)", city, keyword),
		func() ([]model.Route, error) { return f.Primary.SearchRoutes(ctx, city, keyword) },
		fb,
	)
}

func (f *FallbackService) GetStops(ctx context.Context, city, routeID string, direction int) ([]model.Stop, error) {
	cacheKey := "stops:" + city + ":" + routeID + ":" + strconv.Itoa(direction)
	var fb func() ([]model.Stop, error)
	if f.Fallback != nil {
		fb = func() ([]model.Stop, error) { return f.Fallback.GetStops(ctx, city, routeID, direction) }
	}
	return fallbackFetch(f.Cache, cacheKey, fmt.Sprintf("GetStops(%s, %s, %d)", city, routeID, direction),
		func() ([]model.Stop, error) { return f.Primary.GetStops(ctx, city, routeID, direction) },
		fb,
	)
}

func (f *FallbackService) GetETA(ctx context.Context, city, routeID string, direction int) ([]model.StopETA, error) {
	cacheKey := "eta:" + city + ":" + routeID + ":" + strconv.Itoa(direction)
	var fb func() ([]model.StopETA, error)
	if f.Fallback != nil {
		fb = func() ([]model.StopETA, error) { return f.Fallback.GetETA(ctx, city, routeID, direction) }
	}
	return fallbackFetch(f.Cache, cacheKey, fmt.Sprintf("GetETA(%s, %s, %d)", city, routeID, direction),
		func() ([]model.StopETA, error) { return f.Primary.GetETA(ctx, city, routeID, direction) },
		fb,
	)
}
