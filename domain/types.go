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

// CreateCessionRequest creates a cession document.
type CreateCessionRequest struct {
	DocumentID       string `json:"document_id"`
	FactoringCode    string `json:"factoring_code"`
	FactoringName    string `json:"factoring_name"`
	FactoringAddress string `json:"factoring_address"`
	FactoringEmail   string `json:"factoring_email"`
	IdempotencyKey   string `json:"-"`
}

// GeneratePDFRequest requests a PDF generation.
type GeneratePDFRequest struct {
	DocumentID     string `json:"document_id"`
	Formato        string `json:"formato,omitempty"`
	CopiaCedible   bool   `json:"copia_cedible,omitempty"`
	IdempotencyKey string `json:"-"`
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
	Certificate string    `json:"certificate"`
	Password    string    `json:"password"`
	ExpiredDate time.Time `json:"expired_date"`
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

// UploadNumerationRequest uploads CAF/folios.
type UploadNumerationRequest struct {
	CodeSII      string `json:"code_sii"`
	StartNumber  int    `json:"start_number"`
	EndNumber    int    `json:"end_number"`
	CAFBase64    string `json:"caf_base64"`
	CreationDate string `json:"creation_date"`
	DueDate      string `json:"due_date"`
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
