package httpintegra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/IntegraDte/integradte-sdk-go/domain"
)

func (c *Client) CreateDocument(ctx context.Context, req domain.CreateDocumentRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/documents", nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) GetDocument(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/documents/%s", id), nil, nil, nil)
}

func (c *Client) GetDocumentStats(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/documents/stats", nil, nil, nil)
}

func (c *Client) CreateCession(ctx context.Context, req domain.CreateCessionRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/cessions", nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) GeneratePDF(ctx context.Context, req domain.GeneratePDFRequest, cedible bool) (domain.APIResponse, error) {
	query := url.Values{}
	query.Set("cedible", fmt.Sprintf("%t", cedible))
	return c.doJSON(ctx, http.MethodPost, "/api/v1/pdfs/generate", query, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) CreateBusiness(ctx context.Context, req domain.CreateBusinessRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/businesses", nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) UpdateBusiness(ctx context.Context, id string, req domain.UpdateBusinessRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPut, fmt.Sprintf("/api/v1/businesses/%s", id), nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) UploadCertificate(ctx context.Context, businessID string, req domain.UploadCertificateRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPut, fmt.Sprintf("/api/v1/business/%s/certificate", businessID), nil, req, nil)
}

func (c *Client) GetCertificateInfo(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/business/certificate-info", nil, nil, nil)
}

func (c *Client) GetMe(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/users/me", nil, nil, nil)
}

func (c *Client) CreatePurchase(ctx context.Context, req domain.CreatePurchaseRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/purchases", nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) GetNumerationSummary(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/numerations/summary", nil, nil, nil)
}

func (c *Client) GetLastUsedFolio(ctx context.Context, codeSII string) (domain.APIResponse, error) {
	query := url.Values{}
	query.Set("code_sii", codeSII)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/numerations/last-used-number", query, nil, nil)
}

func (c *Client) UploadNumeration(ctx context.Context, req domain.UploadNumerationRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPut, "/api/v1/numerations", nil, req, nil)
}

func (c *Client) DeleteNumeration(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/numerations/%s", id), nil, nil, nil)
}
