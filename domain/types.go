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

// CreateLicenseRequest creates an offline license.
type CreateLicenseRequest struct {
	Name              string   `json:"name"`
	LicenseKey        string   `json:"license_key,omitempty"`
	DeviceID          string   `json:"device_id,omitempty"`
	DeviceFingerprint string   `json:"device_fingerprint"`
	Features          []string `json:"features,omitempty"`
	CLIMinVersion     string   `json:"cli_min_version,omitempty"`
	ValidityHours     int      `json:"validity_hours,omitempty"`
}

// LicenseActionRequest changes an offline license status.
type LicenseActionRequest struct {
	Reason string `json:"reason"`
}

// LicenseBusiness contains the business data embedded in a signed offline license.
type LicenseBusiness struct {
	BusinessName        string `json:"business_name"`
	RUT                 string `json:"rut"`
	Activity            string `json:"activity"`
	Address             string `json:"address"`
	Commune             string `json:"commune"`
	Region              string `json:"region"`
	EmailDTE            string `json:"email_dte"`
	EmailContact        string `json:"email_contact"`
	ResolutionNumberDTE string `json:"resolution_number_dte"`
	ResolutionDateDTE   string `json:"resolution_date_dte"`
	IsProd              bool   `json:"is_prod"`
}

// LicensePayload is the signed payload used by offline clients.
type LicensePayload struct {
	LicenseID         string          `json:"license_id"`
	BusinessID        string          `json:"business_id"`
	DeviceID          string          `json:"device_id"`
	DeviceFingerprint string          `json:"device_fingerprint"`
	Business          LicenseBusiness `json:"business"`
	Features          []string        `json:"features"`
	IssuedAt          string          `json:"issued_at"`
	ExpiresAt         string          `json:"expires_at"`
	LastValidatedAt   string          `json:"last_validated_at"`
	Status            string          `json:"status"`
	CLIMinVersion     string          `json:"cli_min_version"`
}

// SignedLicense contains an offline license payload and its Ed25519 signature.
type SignedLicense struct {
	Payload   LicensePayload `json:"payload"`
	Signature string         `json:"signature"`
}

// ActivateLicenseRequest activates an offline license for a device.
type ActivateLicenseRequest struct {
	LicenseKey         string `json:"license_key"`
	DeviceID           string `json:"device_id"`
	MachineFingerprint string `json:"machine_fingerprint"`
	Hostname           string `json:"hostname"`
	Platform           string `json:"platform"`
	Arch               string `json:"arch"`
	CLIVersion         string `json:"cli_version"`
}

// RefreshLicenseRequest refreshes an existing signed offline license.
type RefreshLicenseRequest struct {
	DeviceID           string        `json:"device_id"`
	MachineFingerprint string        `json:"machine_fingerprint"`
	CLIVersion         string        `json:"cli_version"`
	License            SignedLicense `json:"license"`
}

// RequestNumbersRequest reserves available folios for offline use.
type RequestNumbersRequest struct {
	DocumentType int `json:"document_type"`
	Quantity     int `json:"quantity"`
}

// NumberRange is a reserved folio range returned by /v1/numbers/request.
type NumberRange struct {
	DocumentType   int    `json:"document_type"`
	InitialFolio   int    `json:"folio_inicial"`
	FinalFolio     int    `json:"folio_final"`
	FolioXMLBase64 string `json:"folio_xml_base64"`
}

// RequestNumerationsRequest publishes a numeration request to RabbitMQ.
type RequestNumerationsRequest struct {
	CodeSII  string `json:"code_sii"`
	Quantity int    `json:"quantity"`
}

// SyncDocumentRequest uploads an already signed offline document.
type SyncDocumentRequest struct {
	DocumentID   string        `json:"document_id"`
	DocumentType int           `json:"document_type"`
	Folio        int           `json:"folio"`
	XMLBase64    string        `json:"xml_base64"`
	PDFBase64    string        `json:"pdf_base64,omitempty"`
	TEDXMLBase64 string        `json:"ted_xml_base64,omitempty"`
	GeneratedAt  string        `json:"generated_at"`
	RawPayload   any           `json:"raw_payload"`
	License      SignedLicense `json:"license"`
}

// RequeueDocumentRequest requeues a normal or offline document.
type RequeueDocumentRequest struct {
	DocumentID string `json:"document_id"`
}
