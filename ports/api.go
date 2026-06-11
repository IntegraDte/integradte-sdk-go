package ports

import (
	"context"

	"github.com/IntegraDte/integradte-sdk-go/domain"
)

// IntegraDTEAPI defines the outbound port.
type IntegraDTEAPI interface {
	CreateDocument(ctx context.Context, req domain.CreateDocumentRequest) (domain.APIResponse, error)
	ListDocuments(ctx context.Context, filter domain.DocumentFilter) (domain.APIResponse, error)
	GetDocument(ctx context.Context, id string) (domain.APIResponse, error)
	GetDocumentStats(ctx context.Context) (domain.APIResponse, error)
	GetDocumentStatsWithFilter(ctx context.Context, filter domain.DocumentFilter) (domain.APIResponse, error)
	CreateCession(ctx context.Context, req domain.CreateCessionRequest) (domain.APIResponse, error)
	GeneratePDF(ctx context.Context, req domain.GeneratePDFRequest, cedible bool) (domain.APIResponse, error)
	CreateBusiness(ctx context.Context, req domain.CreateBusinessRequest) (domain.APIResponse, error)
	ListBusinesses(ctx context.Context) (domain.APIResponse, error)
	GetBusiness(ctx context.Context, id string) (domain.APIResponse, error)
	UpdateBusiness(ctx context.Context, id string, req domain.UpdateBusinessRequest) (domain.APIResponse, error)
	EnableProductionMode(ctx context.Context, req domain.ProductionModeRequest) (domain.APIResponse, error)
	EnableCertificationMode(ctx context.Context) (domain.APIResponse, error)
	UploadCertificate(ctx context.Context, businessID string, req domain.UploadCertificateRequest) (domain.APIResponse, error)
	GetCertificateInfo(ctx context.Context) (domain.APIResponse, error)
	GetCurrentCertificate(ctx context.Context) (domain.APIResponse, error)
	GetMe(ctx context.Context) (domain.APIResponse, error)
	GetBillingBalance(ctx context.Context) (domain.APIResponse, error)
	ListBillingPayments(ctx context.Context, filter domain.PaymentFilter) (domain.APIResponse, error)
	CreatePurchase(ctx context.Context, req domain.CreatePurchaseRequest) (domain.APIResponse, error)
	ListPurchaseAcknowledgments(ctx context.Context, filter domain.PurchaseAcknowledgmentFilter) (domain.APIResponse, error)
	GetNumerationSummary(ctx context.Context) (domain.APIResponse, error)
	GetLastUsedFolio(ctx context.Context, codeSII string) (domain.APIResponse, error)
	UploadNumeration(ctx context.Context, req domain.UploadNumerationRequest) (domain.APIResponse, error)
	DeleteNumeration(ctx context.Context, id string) (domain.APIResponse, error)
	CreateLicense(ctx context.Context, req domain.CreateLicenseRequest) (domain.APIResponse, error)
	ListLicenses(ctx context.Context) (domain.APIResponse, error)
	GetLicense(ctx context.Context, id string) (domain.APIResponse, error)
	ListLicenseDevices(ctx context.Context, id string) (domain.APIResponse, error)
	EnableLicense(ctx context.Context, id string, req domain.LicenseActionRequest) (domain.APIResponse, error)
	DisableLicense(ctx context.Context, id string, req domain.LicenseActionRequest) (domain.APIResponse, error)
	RevokeLicense(ctx context.Context, id string, req domain.LicenseActionRequest) (domain.APIResponse, error)
	ActivateLicense(ctx context.Context, req domain.ActivateLicenseRequest) (domain.APIResponse, error)
	RefreshLicense(ctx context.Context, req domain.RefreshLicenseRequest) (domain.APIResponse, error)
	RequestNumbers(ctx context.Context, req domain.RequestNumbersRequest) ([]domain.NumberRange, error)
	RequestNumerations(ctx context.Context, req domain.RequestNumerationsRequest) (domain.APIResponse, error)
	SyncDocument(ctx context.Context, req domain.SyncDocumentRequest) (domain.APIResponse, error)
	RequeueDocument(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error)
	RequeueOfflineDocument(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error)
	RequeueOfflineDocumentStatus(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error)
}
