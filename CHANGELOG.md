# Changelog

## [0.3.0](https://github.com/IntegraDte/integradte-sdk-go/compare/integradte-sdk-go-v0.2.1...integradte-sdk-go-v0.3.0) (2026-09-12)


### ⚠ BREAKING CHANGES

* **api:** removed from httpintegra.Client, ports.IntegraDTEAPI and application.Service: RequeueOfflineDocument and RequestNumerations. Removed domain type: RequestNumerationsRequest. RequeueDocumentRequest stays, since RequeueDocument and RequeueOfflineDocumentStatus still use it.

### Features

* **api:** drop offline requeue and queued numeration requests ([ad04847](https://github.com/IntegraDte/integradte-sdk-go/commit/ad04847680f162cb05c2e2cefb42237dff36718f))

## [0.2.1](https://github.com/IntegraDte/integradte-sdk-go/compare/integradte-sdk-go-v0.2.0...integradte-sdk-go-v0.2.1) (2026-09-11)


### Bug Fixes

* **api:** point CreatePurchase and RequestNumbers at the current API routes ([9b2aa4c](https://github.com/IntegraDte/integradte-sdk-go/commit/9b2aa4c4ef5d3897992ab98a25234b5f22d0880c))
* **api:** point CreatePurchase and RequestNumbers at the current API routes ([6311e90](https://github.com/IntegraDte/integradte-sdk-go/commit/6311e90614bcdfed11917b4c081890250d273b1a))

## [0.2.0](https://github.com/IntegraDte/integradte-sdk-go/compare/integradte-sdk-go-v0.1.5...integradte-sdk-go-v0.2.0) (2026-09-11)


### ⚠ BREAKING CHANGES

* **api:** removed from httpintegra.Client, ports.IntegraDTEAPI and application.Service: CreateLicense, ListLicenses, GetLicense, ListLicenseDevices, EnableLicense, DisableLicense, RevokeLicense, ActivateLicense, RefreshLicense, SyncDocument and GetCurrentCertificate. Removed domain types: CreateLicenseRequest, LicenseActionRequest, LicenseBusiness, LicensePayload, SignedLicense, ActivateLicenseRequest, RefreshLicenseRequest and SyncDocumentRequest. GetCertificateInfo keeps its signature, but its data is now only {"has_valid_certificate": bool} (domain.CertificateInfo): the certificate details are gone, and a business without a certificate gets false with status 200 instead of a 400 error.

### Features

* **api:** sync with the public API, drop licenses, sync and current certificate ([1dacd8a](https://github.com/IntegraDte/integradte-sdk-go/commit/1dacd8aa997039b85efffbd3068299a4048ccd8f))

## [0.1.5](https://github.com/IntegraDte/integradte-sdk-go/compare/integradte-sdk-go-v0.1.4...integradte-sdk-go-v0.1.5) (2026-06-11)


### Features

* **api:** add new endpoints for document and business management ([a028d35](https://github.com/IntegraDte/integradte-sdk-go/commit/a028d35ccaf376a31d41fac29a0411c7b661eb3f))

## [0.1.4](https://github.com/IntegraDte/integradte-sdk-go/compare/integradte-sdk-go-v0.1.3...integradte-sdk-go-v0.1.4) (2026-03-18)

### Features

- **application:** create service layer to orchestrate API calls for document management. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **domain:** add DTE document structures and request builders for various DTE types. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **domain:** define common structures for DTE payloads and implement request creation functions. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **httpintegra:** implement HTTP client for Integra Facturacion API with document management features. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **ports:** define IntegraFacturacionAPI interface for application layer interaction. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))

### Bug Fixes

- align module path and docs with IntegraDte namespace ([d69f25d](https://github.com/IntegraDte/integradte-sdk-go/commit/d69f25d6148597f82d91c5f3e37f2e2bcc3d88b7))

## [0.1.2](https://github.com/IntegraDte/integradte-sdk-go/compare/fulldte-sdk-go-v0.1.1...fulldte-sdk-go-v0.1.2) (2026-02-27)

## [0.1.3](https://github.com/IntegraDte/integradte-sdk-go/compare/integradte-sdk-go-v0.1.2...integradte-sdk-go-v0.1.3) (2026-02-27)

### Features

- **application:** create service layer to orchestrate API calls for document management. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **domain:** add DTE document structures and request builders for various DTE types. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **domain:** define common structures for DTE payloads and implement request creation functions. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **httpintegra:** implement HTTP client for Integra Facturacion API with document management features. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **ports:** define integradteAPI interface for application layer interaction. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))

### Bug Fixes

- align module path and docs with IntegraDte namespace ([d69f25d](https://github.com/IntegraDte/integradte-sdk-go/commit/d69f25d6148597f82d91c5f3e37f2e2bcc3d88b7))

## [0.1.2](https://github.com/IntegraDte/integradte-sdk-go/compare/fulldte-sdk-go-v0.1.1...fulldte-sdk-go-v0.1.2) (2026-02-27)

### Bug Fixes

- align module path and docs with IntegraDte namespace ([d69f25d](https://github.com/IntegraDte/integradte-sdk-go/commit/d69f25d6148597f82d91c5f3e37f2e2bcc3d88b7))

## [0.1.1](https://github.com/IntegraDte/integradte-sdk-go/compare/fulldte-sdk-go-v0.1.0...fulldte-sdk-go-v0.1.1) (2026-02-26)

### Features

- **application:** create service layer to orchestrate API calls for document management. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **domain:** add DTE document structures and request builders for various DTE types. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **domain:** define common structures for DTE payloads and implement request creation functions. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **httpintegra:** implement HTTP client for Integra Facturacion API with document management features. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
- **ports:** define IntegraDTEAPI interface for application layer interaction. ([5a5b6df](https://github.com/IntegraDte/integradte-sdk-go/commit/5a5b6df515904491742a264660709e1514e1ff5c))
