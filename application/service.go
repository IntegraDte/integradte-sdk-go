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

func (s *Service) GetHealth(ctx context.Context) (domain.APIResponse, error) {
	return s.api.GetHealth(ctx)
}

func (s *Service) Login(ctx context.Context, req domain.LoginRequest) (domain.APIResponse, error) {
	return s.api.Login(ctx, req)
}

func (s *Service) CreateFirstBusiness(ctx context.Context, userKey string, req domain.CreateFirstBusinessRequest) (domain.APIResponse, error) {
	return s.api.CreateFirstBusiness(ctx, userKey, req)
}

func (s *Service) UpdateDocument(ctx context.Context, id string, req domain.UpdateDocumentRequest) (domain.APIResponse, error) {
	return s.api.UpdateDocument(ctx, id, req)
}

func (s *Service) UpdateNumerationNextNumber(ctx context.Context, numerationID string, req domain.UpdateNumerationNextNumberRequest) (domain.APIResponse, error) {
	return s.api.UpdateNumerationNextNumber(ctx, numerationID, req)
}

func (s *Service) UpdateLowStockConfig(ctx context.Context, req domain.UpdateLowStockConfigRequest) (domain.APIResponse, error) {
	return s.api.UpdateLowStockConfig(ctx, req)
}

func (s *Service) ListNumerationRanges(ctx context.Context, filter domain.NumerationRangeFilter) (domain.APIResponse, error) {
	return s.api.ListNumerationRanges(ctx, filter)
}

func (s *Service) RequeuePurchase(ctx context.Context, req domain.RequeuePurchaseRequest) (domain.APIResponse, error) {
	return s.api.RequeuePurchase(ctx, req)
}

func (s *Service) ListBillingCharges(ctx context.Context, filter domain.ChargeFilter) (domain.APIResponse, error) {
	return s.api.ListBillingCharges(ctx, filter)
}

func (s *Service) ListBillingPlans(ctx context.Context) (domain.APIResponse, error) {
	return s.api.ListBillingPlans(ctx)
}

func (s *Service) ListBillingInvoices(ctx context.Context, filter domain.InvoiceFilter) (domain.APIResponse, error) {
	return s.api.ListBillingInvoices(ctx, filter)
}

func (s *Service) PreviewSubscriptionUpgrade(ctx context.Context, planID string) (domain.APIResponse, error) {
	return s.api.PreviewSubscriptionUpgrade(ctx, planID)
}

func (s *Service) GetConsumption(ctx context.Context) (domain.APIResponse, error) {
	return s.api.GetConsumption(ctx)
}

func (s *Service) ListConsumptionOverages(ctx context.Context, filter domain.ConsumptionOverageFilter) (domain.APIResponse, error) {
	return s.api.ListConsumptionOverages(ctx, filter)
}

func (s *Service) ListConsumptionOperations(ctx context.Context, filter domain.ConsumptionOperationFilter) (domain.APIResponse, error) {
	return s.api.ListConsumptionOperations(ctx, filter)
}

func (s *Service) RequeueCession(ctx context.Context, req domain.RequeueCessionRequest) (domain.APIResponse, error) {
	return s.api.RequeueCession(ctx, req)
}

func (s *Service) ListCessions(ctx context.Context, filter domain.CessionFilter) (domain.APIResponse, error) {
	return s.api.ListCessions(ctx, filter)
}

func (s *Service) GetCession(ctx context.Context, id string) (domain.APIResponse, error) {
	return s.api.GetCession(ctx, id)
}
