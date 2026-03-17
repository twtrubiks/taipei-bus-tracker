package ebus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

const (
	pageURL      = "https://ebus.gov.taipei"
	stopDynsURL  = "https://ebus.gov.taipei/EBus/GetStopDyns"
	searchURL    = "https://ebus.gov.taipei/Query/QBusRoute"
	stopsPageURL = "https://ebus.gov.taipei/Route/StopsOfRoute"
)

var csrfTokenRe = regexp.MustCompile(`name="__RequestVerificationToken".*?value="([^"]+)"`)

type Provider struct {
	httpClient *http.Client

	mu        sync.RWMutex
	csrfToken string
	cookies   []*http.Cookie
}

func NewProvider() *Provider {
	return &Provider{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *Provider) fetchCSRFToken(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch ebus page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	matches := csrfTokenRe.FindSubmatch(body)
	if len(matches) < 2 {
		return fmt.Errorf("CSRF token not found in ebus page")
	}

	p.mu.Lock()
	p.csrfToken = string(matches[1])
	p.cookies = resp.Cookies()
	p.mu.Unlock()

	return nil
}

func (p *Provider) getCSRFToken(ctx context.Context) (string, []*http.Cookie, error) {
	p.mu.RLock()
	token := p.csrfToken
	cookies := p.cookies
	p.mu.RUnlock()

	if token != "" {
		return token, cookies, nil
	}

	if err := p.fetchCSRFToken(ctx); err != nil {
		return "", nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.csrfToken, p.cookies, nil
}

// doPostWithCSRF performs a POST request with CSRF token and cookies.
// On non-200 status, it clears the cached token for retry on next call.
func (p *Provider) doPostWithCSRF(ctx context.Context, targetURL string, extra url.Values, headers map[string]string) ([]byte, error) {
	token, cookies, err := p.getCSRFToken(ctx)
	if err != nil {
		return nil, err
	}

	data := url.Values{"__RequestVerificationToken": {token}}
	for k, vs := range extra {
		data[k] = vs
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		p.mu.Lock()
		p.csrfToken = ""
		p.mu.Unlock()
		return nil, fmt.Errorf("ebus %s returned status %d", targetURL, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *Provider) getStopDyns(ctx context.Context, routeID string, direction int) ([]EBusStopDynRaw, error) {
	body, err := p.doPostWithCSRF(ctx, stopDynsURL, url.Values{
		"routeId": {routeID},
		"gb":      {fmt.Sprintf("%d", direction)},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("ebus GetStopDyns: %w", err)
	}

	var stops []EBusStopDynRaw
	if err := json.Unmarshal(body, &stops); err != nil {
		return nil, fmt.Errorf("failed to parse ebus response: %w", err)
	}
	return stops, nil
}

// SearchRoutes searches routes by keyword via eBus HTML scraping.
func (p *Provider) SearchRoutes(ctx context.Context, _ string, keyword string) ([]model.Route, error) {
	body, err := p.doPostWithCSRF(ctx, searchURL, url.Values{
		"QueryModel.QueryString": {keyword},
	}, map[string]string{"X-Requested-With": "XMLHttpRequest"})
	if err != nil {
		return nil, fmt.Errorf("ebus SearchRoutes: %w", err)
	}
	return parseSearchRoutes(string(body))
}

// GetStops returns stops for a route by scraping the eBus route page.
func (p *Provider) GetStops(ctx context.Context, _ string, routeID string, direction int) ([]model.Stop, error) {
	reqURL := fmt.Sprintf("%s?routeId=%s", stopsPageURL, routeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ebus GetStops failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ebus stops page returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseStopsHTML(string(body), direction)
}

// GetETA returns estimated arrival times from eBus.
func (p *Provider) GetETA(ctx context.Context, _ string, routeID string, direction int) ([]model.StopETA, error) {
	stops, err := p.getStopDyns(ctx, routeID, direction)
	if err != nil {
		return nil, err
	}
	return convertETAs(stops), nil
}
