package httpintegra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/IntegraDte/integradte-sdk-go/domain"
)

const (
	// DefaultBaseURL is the production API host.
	DefaultBaseURL = "https://api.integradte.cl"
	defaultTimeout = 30 * time.Second
)

// Config defines adapter initialization settings.
type Config struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// APIError wraps non-2xx responses.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("integradte: status=%d body=%s", e.StatusCode, e.Body)
}

// Client is an HTTP adapter implementing the outbound port.
type Client struct {
	apiKey     string
	baseURL    *url.URL
	httpClient *http.Client
	userAgent  string
}

// New creates a new HTTP adapter client.
func New(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("integradte: API key is required")
	}

	base := cfg.BaseURL
	if strings.TrimSpace(base) == "" {
		base = DefaultBaseURL
	}

	parsedBase, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("integradte: invalid base URL: %w", err)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	userAgent := strings.TrimSpace(cfg.UserAgent)
	if userAgent == "" {
		userAgent = "integradte-sdk-go/0.2.0"
	}

	return &Client{
		apiKey:     cfg.APIKey,
		baseURL:    parsedBase,
		httpClient: httpClient,
		userAgent:  userAgent,
	}, nil
}

func (c *Client) buildURL(route string, query url.Values) string {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, route)
	u.RawQuery = query.Encode()
	return u.String()
}

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	route string,
	query url.Values,
	body any,
	extraHeaders map[string]string,
) (domain.APIResponse, error) {
	var payload io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("integradte: marshal request: %w", err)
		}
		payload = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(route, query), payload)
	if err != nil {
		return nil, fmt.Errorf("integradte: create request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range extraHeaders {
		if strings.TrimSpace(v) != "" {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("integradte: do request: %w", err)
	}
	defer resp.Body.Close()

	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("integradte: read response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(rawResp)}
	}

	if len(rawResp) == 0 {
		return domain.APIResponse{}, nil
	}

	var out domain.APIResponse
	if err := json.Unmarshal(rawResp, &out); err != nil {
		return nil, fmt.Errorf("integradte: decode response: %w", err)
	}

	return out, nil
}

func withIdempotency(idempotencyKey string) map[string]string {
	if strings.TrimSpace(idempotencyKey) == "" {
		return nil
	}
	return map[string]string{"idempotency-key": idempotencyKey}
}

func EncodeDataDTE(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
