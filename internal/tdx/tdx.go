package tdx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

const (
	tokenURL = "https://tdx.transportdata.tw/auth/realms/TDXConnect/protocol/openid-connect/token"
	baseURL  = "https://tdx.transportdata.tw/api/basic/v2/Bus"
)

type Provider struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client

	mu          sync.RWMutex
	accessToken string
	expiresAt   time.Time
}

func NewProvider(clientID, clientSecret string) *Provider {
	return &Provider{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (p *Provider) refreshToken(ctx context.Context) error {
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request returned %d: %s", resp.StatusCode, body)
	}

	var tok tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	p.mu.Lock()
	p.accessToken = tok.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(tok.ExpiresIn)*time.Second - 5*time.Minute)
	p.mu.Unlock()

	return nil
}

func (p *Provider) getToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	token := p.accessToken
	expires := p.expiresAt
	p.mu.RUnlock()

	if token != "" && time.Now().Before(expires) {
		return token, nil
	}

	if err := p.refreshToken(ctx); err != nil {
		return "", err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.accessToken, nil
}

func (p *Provider) doGet(ctx context.Context, url string) ([]byte, error) {
	token, err := p.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TDX API returned %d: %s", resp.StatusCode, body)
	}

	return io.ReadAll(resp.Body)
}

func (p *Provider) SearchRoutes(ctx context.Context, city, keyword string) ([]model.Route, error) {
	filter := fmt.Sprintf("contains(RouteName/Zh_tw,'%s')", keyword)
	u := buildURL(baseURL+"/Route/City/"+city, filter)
	body, err := p.doGet(ctx, u)
	if err != nil {
		return nil, err
	}

	var tdxRoutes []TDXRoute
	if err := json.Unmarshal(body, &tdxRoutes); err != nil {
		return nil, err
	}
	return convertRoutes(tdxRoutes), nil
}

func (p *Provider) GetStops(ctx context.Context, city, routeID string, direction int) ([]model.Stop, error) {
	filter := fmt.Sprintf("RouteUID eq '%s' and Direction eq %d", routeID, direction)
	u := buildURL(baseURL+"/StopOfRoute/City/"+city, filter)
	body, err := p.doGet(ctx, u)
	if err != nil {
		return nil, err
	}

	var stopOfRoutes []TDXStopOfRoute
	if err := json.Unmarshal(body, &stopOfRoutes); err != nil {
		return nil, err
	}
	return convertStops(stopOfRoutes), nil
}

func (p *Provider) GetETA(ctx context.Context, city, routeID string, direction int) ([]model.StopETA, error) {
	filter := fmt.Sprintf("RouteUID eq '%s' and Direction eq %d", routeID, direction)
	u := buildURL(baseURL+"/EstimatedTimeOfArrival/City/"+city, filter)
	body, err := p.doGet(ctx, u)
	if err != nil {
		return nil, err
	}

	var tdxETAs []TDXETA
	if err := json.Unmarshal(body, &tdxETAs); err != nil {
		return nil, err
	}
	return convertETAs(tdxETAs), nil
}

// buildURL constructs a TDX API URL with properly encoded query parameters.
func buildURL(base, filter string) string {
	params := url.Values{}
	params.Set("$filter", filter)
	params.Set("$format", "JSON")
	return base + "?" + params.Encode()
}
