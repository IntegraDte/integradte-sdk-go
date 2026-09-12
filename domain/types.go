package domain

import "time"

// APIResponse is a generic response container.
type APIResponse map[string]any

// CreateDocumentRequest creates any DTE document through /documents.
type CreateDocumentRequest struct {
	UserID         string `json:"user_id,omitempty"`
	BusinessID     string `json:"business_id,omitempty"`
	CodeSII        string `json:"code_sii"`
	DataDTE        string `json:"data_dte"`
	IdempotencyKey string `json:"-"`
}

// UpdateDocumentRequest replaces a document's DTE (PUT /documents/:id). Set DataDTE (the
// DTE serialized as a JSON string, see httpintegra.EncodeDataDTE) or DataDTEJSON (the DTE
// as a value, or a string holding JSON). DataDTE wins when both are set.
type UpdateDocumentRequest struct {
	DataDTE        string `json:"data_dte,omitempty"`
	DataDTEJSON    any    `json:"data_dte_json,omitempty"`
	IdempotencyKey string `json:"-"`
}

// CreateCessionRequest creates a cession document.
type CreateCessionRequest struct {
	DocumentID       string `json:"document_id"`
	FactoringCode    string `json:"factoring_code"`
	FactoringName    string `json:"factoring_name"`
	FactoringAddress string `json:"factoring_address"`
	FactoringEmail   string `json:"factoring_email"`
	IdempotencyKey   string `json:"-"`
}

// RequeueCessionRequest requeues a cession (POST /cessions/requeue).
type RequeueCessionRequest struct {
	CessionID string `json:"cession_id"`
}

// CessionFilter controls cession queries. DocumentID lists the cessions of one document.
type CessionFilter struct {
	DocumentID string
	Page       int
	Limit      int
}

// GeneratePDFRequest requests a PDF generation.
type GeneratePDFRequest struct {
	DocumentID     string `json:"document_id"`
	Formato        string `json:"formato,omitempty"`
	CopiaCedible   bool   `json:"copia_cedible,omitempty"`
	IdempotencyKey string `json:"-"`
}

// LoginRequest exchanges a user's email and password for the x-user-key (POST /auth/login).
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateBusinessRequest creates a business profile.
type CreateBusinessRequest struct {
	BusinessName           string `json:"businessName"`
	RUT                    string `json:"rut"`
	Activity               string `json:"activity"`
	Address                string `json:"address"`
	Commune                string `json:"commune"`
	City                   string `json:"city"`
	EmailDTE               string `json:"emailDte"`
	EmailContact           string `json:"emailContact"`
	RUTLegalAgent          string `json:"rutLegalAgent"`
	FullNameLegalAgent     string `json:"fullNameLegalAgent"`
	ResolutionNumberDTE    string `json:"resolutionNumberDte"`
	ResolutionDateDTE      string `json:"resolutionDateDte"`
	ResolutionNumberTicket string `json:"resolutionNumberTicket"`
	ResolutionTicketDate   string `json:"resolutionTicketDate"`
	IdempotencyKey         string `json:"-"`
}

// UpdateBusinessRequest updates a business profile.
type UpdateBusinessRequest = CreateBusinessRequest

// CreateFirstBusinessRequest creates the account's first business with the x-user-key
// (POST /onboarding/businesses). Region or City is required. The resolution dates accept
// RFC 3339 or YYYY-MM-DD. Logo is raw base64 without a data: prefix.
type CreateFirstBusinessRequest struct {
	BusinessName           string `json:"businessName"`
	RUT                    string `json:"rut"`
	Activity               string `json:"activity"`
	Address                string `json:"address"`
	Commune                string `json:"commune"`
	Region                 string `json:"region,omitempty"`
	City                   string `json:"city,omitempty"`
	EmailDTE               string `json:"emailDte"`
	EmailContact           string `json:"emailContact"`
	RUTLegalAgent          string `json:"rutLegalAgent"`
	FullNameLegalAgent     string `json:"fullNameLegalAgent"`
	ResolutionNumberDTE    string `json:"resolutionNumberDte"`
	ResolutionDateDTE      string `json:"resolutionDateDte"`
	ResolutionNumberTicket string `json:"resolutionNumberTicket"`
	ResolutionTicketDate   string `json:"resolutionTicketDate"`
	Logo                   string `json:"logo,omitempty"`
	LogoContentType        string `json:"logoContentType,omitempty"`
}

// ProductionModeRequest moves the authenticated business to production.
type ProductionModeRequest struct {
	ResolutionNumberDTE    string `json:"resolution_number_dte"`
	ResolutionDateDTE      string `json:"resolution_date_dte"`
	ResolutionNumberTicket string `json:"resolution_number_ticket"`
	ResolutionTicketDate   string `json:"resolution_ticket_date"`
}

// DocumentFilter controls document listing and statistics queries.
type DocumentFilter struct {
	CodeSII  string
	Status   string
	FromDate string
	ToDate   string
	Page     int
	Limit    int
}

// PaymentFilter controls billing payment history queries.
type PaymentFilter struct {
	Status   string
	FromDate string
	ToDate   string
	Page     int
	Limit    int
}

// ChargeFilter controls billing charge queries (GET /billing/charges). Dates use YYYY-MM-DD.
type ChargeFilter struct {
	Status     string
	PricingKey string
	FromDate   string
	ToDate     string
	Page       int
	Limit      int
}

// InvoiceFilter controls billing invoice queries (GET /billing/invoices). Status is open,
// paid or void; empty returns every invoice.
type InvoiceFilter struct {
	Status string
}

// ConsumptionOverageFilter paginates the current cycle's overages (GET /consumption/overages).
type ConsumptionOverageFilter struct {
	Page  int
	Limit int
}

// ConsumptionOperationFilter selects the month of the operations detail
// (GET /consumption/operations). Period is YYYY-MM in UTC; empty means the current month.
type ConsumptionOperationFilter struct {
	Period string
}

// PurchaseAcknowledgmentFilter controls received purchase queries.
type PurchaseAcknowledgmentFilter struct {
	DocumentType string
	Action       string
	FromDate     string
	ToDate       string
	Page         int
	Limit        int
}

// UploadCertificateRequest uploads a digital certificate.
type UploadCertificateRequest struct {
	Certificate    string    `json:"certificate"`
	Password       string    `json:"password"`
	ExpiredDate    time.Time `json:"expired_date"`
	IdempotencyKey string    `json:"-"`
}

// CertificateInfo is the data returned by /business/certificate-info. HasValidCertificate is
// true when the business has a certificate that opens with its stored password and has not
// expired; a business without a certificate gets false, not an error.
type CertificateInfo struct {
	HasValidCertificate bool `json:"has_valid_certificate"`
}

// CreatePurchaseRequest registers a purchase document.
type CreatePurchaseRequest struct {
	XMLBase64         string `json:"xml_base64"`
	RUTEmisor         string `json:"rut_emisor"`
	RazonSocialEmisor string `json:"razon_social_emisor"`
	TipoDTE           string `json:"tipo_dte"`
	Folio             int    `json:"folio"`
	MntTotal          string `json:"mnt_total"`
	FechaEmision      string `json:"fecha_emision"`
	EmailEmisor       string `json:"email_emisor"`
	AccionDoc         string `json:"accion_doc"`
	IdempotencyKey    string `json:"-"`
}

// RequeuePurchaseRequest requeues a purchase acknowledgment
// (POST /purchase-acknowledgments/requeue).
type RequeuePurchaseRequest struct {
	PurchaseID string `json:"purchase_id"`
}

// UploadNumerationRequest uploads CAF/folios.
type UploadNumerationRequest struct {
	CodeSII        string `json:"code_sii"`
	StartNumber    int    `json:"start_number"`
	EndNumber      int    `json:"end_number"`
	CAFBase64      string `json:"caf_base64"`
	CreationDate   string `json:"creation_date"`
	DueDate        string `json:"due_date"`
	IdempotencyKey string `json:"-"`
}

// NumerationRangeFilter controls GET /numerations/ranges. An empty CodeSII returns the
// ranges of every document type.
type NumerationRangeFilter struct {
	CodeSII string
}

// UpdateNumerationNextNumberRequest sets the folio the next document of a CAF range gets
// (PATCH /numerations/:numerationId/next-number).
type UpdateNumerationNextNumberRequest struct {
	NextNumber     int    `json:"next_number"`
	IdempotencyKey string `json:"-"`
}

// LowStockConfigItem is the low-folio alert config of one document type. CodeSII is a
// string ("33"). Threshold is always sent, so 0 is a valid threshold.
type LowStockConfigItem struct {
	CodeSII         string `json:"code_sii"`
	Threshold       int    `json:"threshold"`
	RequestQuantity int    `json:"request_quantity"`
}

// UpdateLowStockConfigRequest merges the low-folio config by code_sii
// (PATCH /numerations/low-stock). Codes left out keep their current config.
type UpdateLowStockConfigRequest struct {
	Items          []LowStockConfigItem `json:"items"`
	IdempotencyKey string               `json:"-"`
}

// RequestNumbersRequest reserves available folios for offline use.
type RequestNumbersRequest struct {
	DocumentType int `json:"document_type"`
	Quantity     int `json:"quantity"`
}

// NumberRange is a reserved folio range returned by /numerations/request.
type NumberRange struct {
	DocumentType   int    `json:"document_type"`
	InitialFolio   int    `json:"folio_inicial"`
	FinalFolio     int    `json:"folio_final"`
	FolioXMLBase64 string `json:"folio_xml_base64"`
}

// RequeueDocumentRequest requeues a normal or offline document.
type RequeueDocumentRequest struct {
	DocumentID string `json:"document_id"`
}
