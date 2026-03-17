package ebus

import (
	"testing"
)

const searchResultHTML = `
<ul class="auto-list-pool findroute-list-pool">
                <li>

                            <a href=/Route/StopsOfRoute?routeid=0100029900 class="auto-list-link auto-list-findroute-link">
                                <span class="auto-list auto-list-findroute">
                                    <span class="auto-list-findroute-c">
                                        <p class="auto-list-findroute-bus">299  </p>
                                        <p class="auto-list-findroute-place">新莊 - 永春高中</p>

                                    </span>
                                </span>
                            </a>

</li>
                <li>

                            <a href=/Route/StopsOfRoute?routeid=0100029901 class="auto-list-link auto-list-findroute-link">
                                <span class="auto-list auto-list-findroute">
                                    <span class="auto-list-findroute-c">
                                        <p class="auto-list-findroute-bus">299區  </p>
                                        <p class="auto-list-findroute-place">新莊 - 台北車站</p>

                                    </span>
                                </span>
                            </a>

</li>
</ul>
`

func TestParseSearchRoutes(t *testing.T) {
	routes, err := parseSearchRoutes(searchResultHTML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}

	// First route
	if routes[0].RouteID != "0100029900" {
		t.Errorf("expected routeId 0100029900, got %s", routes[0].RouteID)
	}
	if routes[0].Name != "299" {
		t.Errorf("expected name '299', got '%s'", routes[0].Name)
	}
	if routes[0].StartStop != "新莊" {
		t.Errorf("expected startStop '新莊', got '%s'", routes[0].StartStop)
	}
	if routes[0].EndStop != "永春高中" {
		t.Errorf("expected endStop '永春高中', got '%s'", routes[0].EndStop)
	}

	// Second route
	if routes[1].RouteID != "0100029901" {
		t.Errorf("expected routeId 0100029901, got %s", routes[1].RouteID)
	}
	if routes[1].Name != "299區" {
		t.Errorf("expected name '299區', got '%s'", routes[1].Name)
	}
}

func TestParseSearchRoutes_Empty(t *testing.T) {
	routes, err := parseSearchRoutes(`<ul class="auto-list-pool"></ul>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(routes) != 0 {
		t.Errorf("expected 0 routes, got %d", len(routes))
	}
}

func TestParseSearchRoutes_NoSeparator(t *testing.T) {
	html := `<ul>
<li>
<a href=/Route/StopsOfRoute?routeid=0100000100 class="auto-list-link">
<p class="auto-list-findroute-bus">1</p>
<p class="auto-list-findroute-place">萬華</p>
</a>
</li>
</ul>`
	routes, err := parseSearchRoutes(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	if routes[0].StartStop != "萬華" {
		t.Errorf("expected startStop '萬華', got '%s'", routes[0].StartStop)
	}
	if routes[0].EndStop != "" {
		t.Errorf("expected empty endStop, got '%s'", routes[0].EndStop)
	}
}

const stopsPageHTML = `<html><body>
<div id="GoDirectionRoute" class="auto-list-pool-c stationlist-list-pool-c">
<ul class="auto-list-pool stationlist-list-pool">
<li>
<a href="javascript:void(0);">
<span class="auto-list auto-list-stationlist">
<span class="auto-list-stationlist-number"> 1</span>
<span class="auto-list-stationlist-place">三重客運新莊站</span>
<input id="item_UniStopId" name="item.UniStopId" type="hidden" value="2387801040" />
</span>
</a>
</li>
<li>
<a href="javascript:void(0);">
<span class="auto-list auto-list-stationlist">
<span class="auto-list-stationlist-number"> 2</span>
<span class="auto-list-stationlist-place">大安變電所(大安路)</span>
<input id="item_UniStopId" name="item.UniStopId" type="hidden" value="2426201720" />
</span>
</a>
</li>
</ul>
</div>
<div id="BackDirectionRoute" class="auto-list-pool-c stationlist-list-pool-c">
<ul class="auto-list-pool stationlist-list-pool">
<li>
<a href="javascript:void(0);">
<span class="auto-list auto-list-stationlist">
<span class="auto-list-stationlist-number"> 1</span>
<span class="auto-list-stationlist-place">永春高中</span>
<input id="item_UniStopId" name="item.UniStopId" type="hidden" value="1234567890" />
</span>
</a>
</li>
</ul>
</div>
</body></html>`

func TestParseStops_GoDirection(t *testing.T) {
	stops, err := parseStopsHTML(stopsPageHTML, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stops) != 2 {
		t.Fatalf("expected 2 stops for direction 0, got %d", len(stops))
	}

	if stops[0].StopID != "2387801040" {
		t.Errorf("expected stopId 2387801040, got %s", stops[0].StopID)
	}
	if stops[0].Name != "三重客運新莊站" {
		t.Errorf("expected name '三重客運新莊站', got '%s'", stops[0].Name)
	}
	if stops[0].Sequence != 1 {
		t.Errorf("expected sequence 1, got %d", stops[0].Sequence)
	}

	if stops[1].StopID != "2426201720" {
		t.Errorf("expected stopId 2426201720, got %s", stops[1].StopID)
	}
	if stops[1].Name != "大安變電所(大安路)" {
		t.Errorf("expected name '大安變電所(大安路)', got '%s'", stops[1].Name)
	}
	if stops[1].Sequence != 2 {
		t.Errorf("expected sequence 2, got %d", stops[1].Sequence)
	}
}

func TestParseStops_BackDirection(t *testing.T) {
	stops, err := parseStopsHTML(stopsPageHTML, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stops) != 1 {
		t.Fatalf("expected 1 stop for direction 1, got %d", len(stops))
	}

	if stops[0].Name != "永春高中" {
		t.Errorf("expected name '永春高中', got '%s'", stops[0].Name)
	}
	if stops[0].StopID != "1234567890" {
		t.Errorf("expected stopId 1234567890, got %s", stops[0].StopID)
	}
}

func TestParseStops_Empty(t *testing.T) {
	stops, err := parseStopsHTML("<html></html>", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stops) != 0 {
		t.Errorf("expected 0 stops, got %d", len(stops))
	}
}
