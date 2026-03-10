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
	pageURL     = "https://ebus.gov.taipei"
	stopDynsURL = "https://ebus.gov.taipei/EBus/GetStopDyns"
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

func (p *Provider) getStopDyns(ctx context.Context, routeID string, direction int) ([]EBusStopDynRaw, error) {
	token, cookies, err := p.getCSRFToken(ctx)
	if err != nil {
		return nil, err
	}

	data := url.Values{
		"__RequestVerificationToken": {token},
		"routeId": {routeID},
		"gb":      {fmt.Sprintf("%d", direction)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, stopDynsURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ebus GetStopDyns failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// Token might be expired, clear and retry once
		p.mu.Lock()
		p.csrfToken = ""
		p.mu.Unlock()
		return nil, fmt.Errorf("ebus returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stops []EBusStopDynRaw
	if err := json.Unmarshal(body, &stops); err != nil {
		return nil, fmt.Errorf("failed to parse ebus response: %w", err)
	}

	return stops, nil
}

// SearchRoutes is not supported by eBus provider.
func (p *Provider) SearchRoutes(_ context.Context, _, _ string) ([]model.Route, error) {
	return nil, fmt.Errorf("eBus provider does not support SearchRoutes")
}

// GetStops is not supported by eBus provider.
func (p *Provider) GetStops(_ context.Context, _, _ string, _ int) ([]model.Stop, error) {
	return nil, fmt.Errorf("eBus provider does not support GetStops")
}

// GetETA returns estimated arrival times from eBus.
func (p *Provider) GetETA(ctx context.Context, _ string, routeID string, direction int) ([]model.StopETA, error) {
	stops, err := p.getStopDyns(ctx, routeID, direction)
	if err != nil {
		return nil, err
	}
	return convertETAs(stops), nil
}
