package httpintegra

import (
	"bytes"
	"context"
	"crypto/rand"
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

	headerAPIKey         = "x-api-key"
	headerUserKey        = "x-user-key"
	headerIdempotencyKey = "idempotency-key"
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

// New creates a new HTTP adapter client. cfg.APIKey is required.
func New(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("integradte: API key is required")
	}
	return newClient(cfg)
}

// NewWithoutAPIKey creates a client that does not require cfg.APIKey, for the calls made
// before the account has an API key: GetHealth, Login and CreateFirstBusiness. If
// cfg.APIKey is empty, the methods that need x-api-key return an error without calling
// the API.
func NewWithoutAPIKey(cfg Config) (*Client, error) {
	return newClient(cfg)
}

func newClient(cfg Config) (*Client, error) {
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

// credential is the auth header a route expects. The zero value sends no auth header.
type credential struct {
	header string
	value  string
}

func (c *Client) apiKeyCredential() credential {
	return credential{header: headerAPIKey, value: c.apiKey}
}

func (c *Client) buildURL(route string, query url.Values) string {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, route)
	u.RawQuery = query.Encode()
	return u.String()
}

// doJSON sends a request authenticated with x-api-key and decodes the JSON response.
func (c *Client) doJSON(
	ctx context.Context,
	method string,
	route string,
	query url.Values,
	body any,
	extraHeaders map[string]string,
) (domain.APIResponse, error) {
	return c.doJSONAs(ctx, c.apiKeyCredential(), method, route, query, body, extraHeaders)
}

// doJSONAs is doJSON with an explicit credential, for the routes that do not use x-api-key.
func (c *Client) doJSONAs(
	ctx context.Context,
	cred credential,
	method string,
	route string,
	query url.Values,
	body any,
	extraHeaders map[string]string,
) (domain.APIResponse, error) {
	var out domain.APIResponse
	if err := c.send(ctx, cred, method, route, query, body, extraHeaders, &out); err != nil {
		return nil, err
	}
	if out == nil {
		return domain.APIResponse{}, nil
	}
	return out, nil
}

// doJSONInto sends a request authenticated with x-api-key and decodes the response into out.
func (c *Client) doJSONInto(
	ctx context.Context,
	method string,
	route string,
	body any,
	out any,
) error {
	return c.send(ctx, c.apiKeyCredential(), method, route, nil, body, nil, out)
}

// doIdempotent sends an x-api-key request to a route that mounts the API's
// IdempotencyMiddleware, which answers 400 when idempotency-key is missing.
func (c *Client) doIdempotent(
	ctx context.Context,
	method string,
	route string,
	body any,
	idempotencyKey string,
) (domain.APIResponse, error) {
	headers, err := idempotencyHeaders(ctx, idempotencyKey)
	if err != nil {
		return nil, err
	}
	return c.doJSON(ctx, method, route, nil, body, headers)
}

func (c *Client) send(
	ctx context.Context,
	cred credential,
	method string,
	route string,
	query url.Values,
	body any,
	extraHeaders map[string]string,
	out any,
) error {
	if cred.header != "" && strings.TrimSpace(cred.value) == "" {
		return fmt.Errorf("integradte: %s is required for %s %s", cred.header, method, route)
	}

	var payload io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("integradte: marshal request: %w", err)
		}
		payload = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(route, query), payload)
	if err != nil {
		return fmt.Errorf("integradte: create request: %w", err)
	}

	if cred.header != "" {
		req.Header.Set(cred.header, cred.value)
	}
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
		return fmt.Errorf("integradte: do request: %w", err)
	}
	defer resp.Body.Close()

	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("integradte: read response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &APIError{StatusCode: resp.StatusCode, Body: string(rawResp)}
	}
	if len(rawResp) == 0 {
		return nil
	}
	if err := json.Unmarshal(rawResp, out); err != nil {
		return fmt.Errorf("integradte: decode response: %w", err)
	}

	return nil
}

// withIdempotency forwards the caller's key on routes where the header is optional.
func withIdempotency(idempotencyKey string) map[string]string {
	if strings.TrimSpace(idempotencyKey) == "" {
		return nil
	}
	return map[string]string{headerIdempotencyKey: idempotencyKey}
}

type idempotencyKeyContextKey struct{}

// WithIdempotencyKey returns a copy of ctx carrying key as the idempotency-key for the
// routes that require one. Use it with methods that take no request struct, such as
// DeleteNumeration. An IdempotencyKey set on the request struct wins over the one in ctx.
// The key must be a UUID. The API replays the first response for the same key and path
// for 24 hours, so use a new key for each logical operation.
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKeyContextKey{}, key)
}

// idempotencyHeaders returns the idempotency-key header for a route that requires it: the
// caller's key if set, then the key stored with WithIdempotencyKey, and otherwise a new
// random UUID for this call.
func idempotencyHeaders(ctx context.Context, idempotencyKey string) (map[string]string, error) {
	key := strings.TrimSpace(idempotencyKey)
	if key == "" {
		if fromCtx, ok := ctx.Value(idempotencyKeyContextKey{}).(string); ok {
			key = strings.TrimSpace(fromCtx)
		}
	}
	if key == "" {
		generated, err := newUUID()
		if err != nil {
			return nil, fmt.Errorf("integradte: generate idempotency-key: %w", err)
		}
		key = generated
	}
	return map[string]string{headerIdempotencyKey: key}, nil
}

// newUUID returns a random version 4 UUID. The API accepts any UUID version.
func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func EncodeDataDTE(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
