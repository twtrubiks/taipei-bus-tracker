package ebus

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

// Minimum HTML length to suspect a scraping failure rather than empty results.
const minHTMLLen = 1024

// Regex patterns for parsing eBus HTML responses.
var (
	// Search results: extract routeId, route name, and place (start - end).
	searchRouteRe = regexp.MustCompile(
		`(?s)routeid=([^&"\s]+).*?findroute-bus">\s*(.*?)\s*</p>.*?findroute-place">\s*(.*?)\s*</p>`,
	)

	// Stop list: extract sequence, stop name, and UniStopId.
	stopItemRe = regexp.MustCompile(
		`(?s)stationlist-number">\s*(\d+)\s*</span>\s*<span[^>]*stationlist-place">\s*(.*?)\s*</span>.*?UniStopId[^>]*value="(\d+)"`,
	)
)

// parseSearchRoutes parses the HTML returned by POST /Query/QBusRoute.
func parseSearchRoutes(html string) ([]model.Route, error) {
	matches := searchRouteRe.FindAllStringSubmatch(html, -1)
	routes := make([]model.Route, 0, len(matches))
	for _, m := range matches {
		routeID := strings.TrimSpace(m[1])
		name := strings.TrimSpace(m[2])
		place := strings.TrimSpace(m[3])

		startStop, endStop := splitPlace(place)
		routes = append(routes, model.Route{
			RouteID:   routeID,
			Name:      name,
			StartStop: startStop,
			EndStop:   endStop,
			Source:    "ebus",
		})
	}
	if len(routes) == 0 && len(html) > minHTMLLen {
		return nil, fmt.Errorf("ebus scrape: searchRoutes matched 0 routes from %d bytes HTML, possible format change", len(html))
	}
	return routes, nil
}

// splitPlace splits "新莊 - 永春高中" into ("新莊", "永春高中").
func splitPlace(place string) (string, string) {
	parts := strings.SplitN(place, " - ", 2)
	start := strings.TrimSpace(parts[0])
	end := ""
	if len(parts) == 2 {
		end = strings.TrimSpace(parts[1])
	}
	return start, end
}

// parseStopsHTML parses the full /Route/StopsOfRoute page HTML.
// direction 0 = GoDirectionRoute, 1 = BackDirectionRoute.
func parseStopsHTML(html string, direction int) ([]model.Stop, error) {
	sectionID := "GoDirectionRoute"
	if direction == 1 {
		sectionID = "BackDirectionRoute"
	}

	// Find the section for the requested direction.
	idx := strings.Index(html, `id="`+sectionID+`"`)
	if idx < 0 {
		if len(html) > minHTMLLen {
			return nil, fmt.Errorf("ebus scrape: section %s not found in %d bytes HTML, possible format change", sectionID, len(html))
		}
		return nil, nil
	}
	section := html[idx:]

	// Limit to this section only — find the next section or end.
	nextSection := ""
	if direction == 0 {
		nextSection = "BackDirectionRoute"
	}
	if nextSection != "" {
		endIdx := strings.Index(section, `id="`+nextSection+`"`)
		if endIdx > 0 {
			section = section[:endIdx]
		}
	}

	matches := stopItemRe.FindAllStringSubmatch(section, -1)
	stops := make([]model.Stop, 0, len(matches))
	for _, m := range matches {
		seq, _ := strconv.Atoi(strings.TrimSpace(m[1]))
		stops = append(stops, model.Stop{
			StopID:   strings.TrimSpace(m[3]),
			Name:     strings.TrimSpace(m[2]),
			Sequence: seq,
			Source:   "ebus",
		})
	}
	if len(stops) == 0 && len(section) > minHTMLLen {
		return nil, fmt.Errorf("ebus scrape: parseStops matched 0 stops from %d bytes section, possible format change", len(section))
	}
	return stops, nil
}
