package ports

import (
	"context"

	"github.com/IntegraDte/integradte-sdk-go/domain"
)

// IntegraDTEAPI defines the outbound port.
type IntegraDTEAPI interface {
	CreateDocument(ctx context.Context, req domain.CreateDocumentRequest) (domain.APIResponse, error)
	GetDocument(ctx context.Context, id string) (domain.APIResponse, error)
	GetDocumentStats(ctx context.Context) (domain.APIResponse, error)
	CreateCession(ctx context.Context, req domain.CreateCessionRequest) (domain.APIResponse, error)
	GeneratePDF(ctx context.Context, req domain.GeneratePDFRequest, cedible bool) (domain.APIResponse, error)
	CreateBusiness(ctx context.Context, req domain.CreateBusinessRequest) (domain.APIResponse, error)
	UpdateBusiness(ctx context.Context, id string, req domain.UpdateBusinessRequest) (domain.APIResponse, error)
	UploadCertificate(ctx context.Context, businessID string, req domain.UploadCertificateRequest) (domain.APIResponse, error)
	GetCertificateInfo(ctx context.Context) (domain.APIResponse, error)
	GetMe(ctx context.Context) (domain.APIResponse, error)
	CreatePurchase(ctx context.Context, req domain.CreatePurchaseRequest) (domain.APIResponse, error)
	GetNumerationSummary(ctx context.Context) (domain.APIResponse, error)
	GetLastUsedFolio(ctx context.Context, codeSII string) (domain.APIResponse, error)
	UploadNumeration(ctx context.Context, req domain.UploadNumerationRequest) (domain.APIResponse, error)
	DeleteNumeration(ctx context.Context, id string) (domain.APIResponse, error)
}
