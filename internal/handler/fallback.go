package handler

import (
	"context"
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

func (f *FallbackService) SearchRoutes(ctx context.Context, city, keyword string) ([]model.Route, error) {
	cacheKey := "routes:" + city + ":" + keyword

	routes, err := f.Primary.SearchRoutes(ctx, city, keyword)
	if err == nil && len(routes) > 0 {
		f.Cache.Set(cacheKey, routes)
		return routes, nil
	}
	if err != nil {
		log.Printf("[fallback] SearchRoutes primary failed: %v", err)
	}

	if f.Fallback != nil {
		routes, err = f.Fallback.SearchRoutes(ctx, city, keyword)
		if err == nil && len(routes) > 0 {
			f.Cache.Set(cacheKey, routes)
			return routes, nil
		}
		if err != nil {
			log.Printf("[fallback] SearchRoutes fallback failed: %v", err)
		}
	}

	if cached, ok := f.Cache.Get(cacheKey); ok {
		if routes, ok := cached.([]model.Route); ok && len(routes) > 0 {
			return routes, nil
		}
	}

	if cached, ok := f.Cache.GetStale(cacheKey); ok {
		if routes, ok := cached.([]model.Route); ok && len(routes) > 0 {
			log.Printf("[fallback] SearchRoutes(%s, %s): using stale cache", city, keyword)
			return routes, nil
		}
	}

	if err == nil {
		return routes, nil
	}

	return nil, err
}

func (f *FallbackService) GetStops(ctx context.Context, city, routeID string, direction int) ([]model.Stop, error) {
	cacheKey := "stops:" + city + ":" + routeID + ":" + strconv.Itoa(direction)

	stops, err := f.Primary.GetStops(ctx, city, routeID, direction)
	if err == nil && len(stops) > 0 {
		f.Cache.Set(cacheKey, stops)
		return stops, nil
	}
	if err != nil {
		log.Printf("[fallback] GetStops primary failed: %v", err)
	}

	if f.Fallback != nil {
		stops, err = f.Fallback.GetStops(ctx, city, routeID, direction)
		if err == nil && len(stops) > 0 {
			f.Cache.Set(cacheKey, stops)
			return stops, nil
		}
		if err != nil {
			log.Printf("[fallback] GetStops fallback failed: %v", err)
		}
	}

	if cached, ok := f.Cache.Get(cacheKey); ok {
		if stops, ok := cached.([]model.Stop); ok && len(stops) > 0 {
			return stops, nil
		}
	}

	if cached, ok := f.Cache.GetStale(cacheKey); ok {
		if stops, ok := cached.([]model.Stop); ok && len(stops) > 0 {
			log.Printf("[fallback] GetStops(%s, %s, %d): using stale cache", city, routeID, direction)
			return stops, nil
		}
	}

	if err == nil {
		return stops, nil
	}

	return nil, err
}

func (f *FallbackService) GetETA(ctx context.Context, city, routeID string, direction int) ([]model.StopETA, error) {
	cacheKey := "eta:" + city + ":" + routeID + ":" + strconv.Itoa(direction)

	etas, err := f.Primary.GetETA(ctx, city, routeID, direction)
	if err == nil && len(etas) > 0 {
		f.Cache.Set(cacheKey, etas)
		return etas, nil
	}
	if err != nil {
		log.Printf("[fallback] GetETA primary failed: %v", err)
	}

	if f.Fallback != nil {
		etas, err = f.Fallback.GetETA(ctx, city, routeID, direction)
		if err == nil && len(etas) > 0 {
			f.Cache.Set(cacheKey, etas)
			return etas, nil
		}
		if err != nil {
			log.Printf("[fallback] GetETA fallback failed: %v", err)
		}
	}

	if cached, ok := f.Cache.Get(cacheKey); ok {
		if etas, ok := cached.([]model.StopETA); ok && len(etas) > 0 {
			return etas, nil
		}
	}

	if cached, ok := f.Cache.GetStale(cacheKey); ok {
		if etas, ok := cached.([]model.StopETA); ok && len(etas) > 0 {
			log.Printf("[fallback] GetETA(%s, %s, %d): using stale cache", city, routeID, direction)
			return etas, nil
		}
	}

	if err == nil {
		return etas, nil
	}

	return nil, err
}
