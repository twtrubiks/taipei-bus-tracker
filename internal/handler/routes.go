package handler

import (
	"net/http"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/cache"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

const defaultCity = "Taipei"

// Handlers holds dependencies for HTTP handlers.
type Handlers struct {
	Fallback *FallbackService
	Primary  model.BusDataSource
	Cache    *cache.Cache
}

func NewHandlers(primary model.BusDataSource, fallback model.BusDataSource, c *cache.Cache) *Handlers {
	return &Handlers{
		Fallback: NewFallbackService(primary, fallback, c),
		Primary:  primary,
		Cache:    c,
	}
}

// SearchRoutes handles GET /api/routes/search?q={keyword}&city={city}
func (h *Handlers) SearchRoutes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "missing query parameter: q")
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		city = defaultCity
	}

	cacheKey := "routes:" + city + ":" + q
	if cached, ok := h.Cache.Get(cacheKey); ok {
		if routes, ok := cached.([]model.Route); ok {
			writeJSON(w, http.StatusOK, routes)
			return
		}
	}

	ctx := r.Context()
	routes, err := h.Primary.SearchRoutes(ctx, city, q)
	if err != nil {
		writeError(w, http.StatusBadGateway, "upstream unavailable")
		return
	}

	h.Cache.Set(cacheKey, routes)
	writeJSON(w, http.StatusOK, routes)
}

// GetStops handles GET /api/routes/{routeId}/stops?gb={direction}
func (h *Handlers) GetStops(w http.ResponseWriter, r *http.Request) {
	routeID := r.PathValue("routeId")
	if routeID == "" {
		writeError(w, http.StatusBadRequest, "missing routeId")
		return
	}

	direction := parseDirection(r.URL.Query().Get("gb"))
	city := r.URL.Query().Get("city")
	if city == "" {
		city = defaultCity
	}

	cacheKey := "stops:" + city + ":" + routeID + ":" + r.URL.Query().Get("gb")
	if cached, ok := h.Cache.Get(cacheKey); ok {
		if stops, ok := cached.([]model.Stop); ok {
			writeJSON(w, http.StatusOK, stops)
			return
		}
	}

	ctx := r.Context()
	stops, err := h.Primary.GetStops(ctx, city, routeID, direction)
	if err != nil {
		writeError(w, http.StatusBadGateway, "upstream unavailable")
		return
	}

	h.Cache.Set(cacheKey, stops)
	writeJSON(w, http.StatusOK, stops)
}

// ETAResponse is the response format for the ETA endpoint.
type ETAResponse struct {
	Route     string          `json:"route"`
	Direction int             `json:"direction"`
	Source    string          `json:"source"`
	UpdatedAt string          `json:"updatedAt"`
	Stops     []model.StopETA `json:"stops"`
}

// GetETA handles GET /api/routes/{routeId}/eta?gb={direction}
func (h *Handlers) GetETA(w http.ResponseWriter, r *http.Request) {
	routeID := r.PathValue("routeId")
	if routeID == "" {
		writeError(w, http.StatusBadRequest, "missing routeId")
		return
	}

	direction := parseDirection(r.URL.Query().Get("gb"))
	city := r.URL.Query().Get("city")
	if city == "" {
		city = defaultCity
	}

	ctx := r.Context()
	etas, err := h.Fallback.GetETA(ctx, city, routeID, direction)
	if err != nil {
		writeError(w, http.StatusBadGateway, "upstream unavailable")
		return
	}

	// Fill in status strings
	for i := range etas {
		etas[i].Status = model.ETAStatus(etas[i].ETA)
	}

	source := "tdx"
	if len(etas) > 0 && etas[0].Source != "" {
		source = etas[0].Source
	}

	resp := ETAResponse{
		Route:     routeID,
		Direction: direction,
		Source:    source,
		UpdatedAt: time.Now().Format(time.RFC3339),
		Stops:     etas,
	}
	writeJSON(w, http.StatusOK, resp)
}

func parseDirection(s string) int {
	if s == "1" {
		return 1
	}
	return 0
}
