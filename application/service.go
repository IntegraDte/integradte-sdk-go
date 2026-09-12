package application

import (
	"context"

	"github.com/IntegraDte/integradte-sdk-go/domain"
	"github.com/IntegraDte/integradte-sdk-go/ports"
)

// Service is the application layer orchestrating use-cases.
type Service struct {
	api ports.IntegraDTEAPI
}

func NewService(api ports.IntegraDTEAPI) *Service {
	return &Service{api: api}
}

func (s *Service) CreateDocument(ctx context.Context, req domain.CreateDocumentRequest) (domain.APIResponse, error) {
	return s.api.CreateDocument(ctx, req)
}

func (s *Service) GetDocument(ctx context.Context, id string) (domain.APIResponse, error) {
	return s.api.GetDocument(ctx, id)
}

func (s *Service) ListDocuments(ctx context.Context, filter domain.DocumentFilter) (domain.APIResponse, error) {
	return s.api.ListDocuments(ctx, filter)
}

func (s *Service) GetDocumentStats(ctx context.Context) (domain.APIResponse, error) {
	return s.api.GetDocumentStats(ctx)
}

func (s *Service) GetDocumentStatsWithFilter(ctx context.Context, filter domain.DocumentFilter) (domain.APIResponse, error) {
	return s.api.GetDocumentStatsWithFilter(ctx, filter)
}

func (s *Service) CreateCession(ctx context.Context, req domain.CreateCessionRequest) (domain.APIResponse, error) {
	return s.api.CreateCession(ctx, req)
}

func (s *Service) GeneratePDF(ctx context.Context, req domain.GeneratePDFRequest, cedible bool) (domain.APIResponse, error) {
	return s.api.GeneratePDF(ctx, req, cedible)
}

func (s *Service) CreateBusiness(ctx context.Context, req domain.CreateBusinessRequest) (domain.APIResponse, error) {
	return s.api.CreateBusiness(ctx, req)
}

func (s *Service) ListBusinesses(ctx context.Context) (domain.APIResponse, error) {
	return s.api.ListBusinesses(ctx)
}

func (s *Service) GetBusiness(ctx context.Context, id string) (domain.APIResponse, error) {
	return s.api.GetBusiness(ctx, id)
}

func (s *Service) UpdateBusiness(ctx context.Context, id string, req domain.UpdateBusinessRequest) (domain.APIResponse, error) {
	return s.api.UpdateBusiness(ctx, id, req)
}

func (s *Service) EnableProductionMode(ctx context.Context, req domain.ProductionModeRequest) (domain.APIResponse, error) {
	return s.api.EnableProductionMode(ctx, req)
}

func (s *Service) EnableCertificationMode(ctx context.Context) (domain.APIResponse, error) {
	return s.api.EnableCertificationMode(ctx)
}

func (s *Service) UploadCertificate(ctx context.Context, businessID string, req domain.UploadCertificateRequest) (domain.APIResponse, error) {
	return s.api.UploadCertificate(ctx, businessID, req)
}

func (s *Service) GetCertificateInfo(ctx context.Context) (domain.APIResponse, error) {
	return s.api.GetCertificateInfo(ctx)
}

func (s *Service) GetMe(ctx context.Context) (domain.APIResponse, error) {
	return s.api.GetMe(ctx)
}

func (s *Service) GetBillingBalance(ctx context.Context) (domain.APIResponse, error) {
	return s.api.GetBillingBalance(ctx)
}

func (s *Service) ListBillingPayments(ctx context.Context, filter domain.PaymentFilter) (domain.APIResponse, error) {
	return s.api.ListBillingPayments(ctx, filter)
}

func (s *Service) CreatePurchase(ctx context.Context, req domain.CreatePurchaseRequest) (domain.APIResponse, error) {
	return s.api.CreatePurchase(ctx, req)
}

func (s *Service) ListPurchaseAcknowledgments(ctx context.Context, filter domain.PurchaseAcknowledgmentFilter) (domain.APIResponse, error) {
	return s.api.ListPurchaseAcknowledgments(ctx, filter)
}

func (s *Service) GetNumerationSummary(ctx context.Context) (domain.APIResponse, error) {
	return s.api.GetNumerationSummary(ctx)
}

func (s *Service) GetLastUsedFolio(ctx context.Context, codeSII string) (domain.APIResponse, error) {
	return s.api.GetLastUsedFolio(ctx, codeSII)
}

func (s *Service) UploadNumeration(ctx context.Context, req domain.UploadNumerationRequest) (domain.APIResponse, error) {
	return s.api.UploadNumeration(ctx, req)
}

func (s *Service) DeleteNumeration(ctx context.Context, id string) (domain.APIResponse, error) {
	return s.api.DeleteNumeration(ctx, id)
}

func (s *Service) RequestNumbers(ctx context.Context, req domain.RequestNumbersRequest) ([]domain.NumberRange, error) {
	return s.api.RequestNumbers(ctx, req)
}

func (s *Service) RequeueDocument(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error) {
	return s.api.RequeueDocument(ctx, req)
}

func (s *Service) RequeueOfflineDocumentStatus(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error) {
	return s.api.RequeueOfflineDocumentStatus(ctx, req)
}
