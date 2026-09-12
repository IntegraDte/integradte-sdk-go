# integradte-sdk-go

SDK en Go para consumir la API de [IntegraDTE](https://api.integradte.cl), con arquitectura hexagonal.

## Instalacion

```bash
go get github.com/IntegraDte/integradte-sdk-go
```

Para instalar una version especifica:

```bash
go get github.com/IntegraDte/integradte-sdk-go@v0.1.0
```

## Estructura hexagonal

- `domain`: modelos de negocio (requests/responses)
- `ports`: puertos (interfaces)
- `application`: casos de uso
- `adapters/httpintegra`: adaptador HTTP concreto para IntegraDTE

## Uso recomendado (hexagonal)

```go
package main

import (
	"context"
	"fmt"

	"github.com/IntegraDte/integradte-sdk-go/adapters/httpintegra"
	"github.com/IntegraDte/integradte-sdk-go/application"
	"github.com/IntegraDte/integradte-sdk-go/domain"
)

func main() {
	adapter, err := httpintegra.New(httpintegra.Config{
		APIKey: "TU_X_API_KEY",
		// BaseURL: "https://api.integradte.cl", // opcional
	})
	if err != nil {
		panic(err)
	}

	service := application.NewService(adapter)

	dataDTE, err := httpintegra.EncodeDataDTE(map[string]any{
		"Encabezado": map[string]any{
			"IdDoc": map[string]any{
				"TipoDTE": 33,
				"FchEmis": "2026-02-03",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	resp, err := service.CreateDocument(context.Background(), domain.CreateDocumentRequest{
		CodeSII: "33",
		DataDTE: dataDTE,
		// Opcional. Debe ser un UUID; si lo omites, el SDK genera uno por llamada.
		IdempotencyKey: "0190d7a4-6c3e-7b8e-9f21-3a4b5c6d7e8f",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
```

## Crear `data_dte` desde structs tipados

Puedes construir el DTE con structs (`Dte33Data`, `Dte34Data`, `Dte39Data`, `Dte41Data`, `Dte46Data`, `Dte52Data`, `Dte56Data`, `Dte61Data`) y convertirlo a `CreateDocumentRequest`.

Estos tipos **no incluyen** `TED` ni `TmstFirma`.

```go
dte := domain.Dte33Data{
	Encabezado: domain.Encabezado33{
		IdDoc: domain.IdDocBase{TipoDTE: 33, FchEmis: "2026-02-03"},
		Emisor: domain.Emisor{
			RUTEmisor:  "12345689-3",
			RznSoc:     "EMPRESA DE PRUEBA",
			GiroEmis:   "Servicios de desarrollo de software",
			DirOrigen:  "Av. Apoquindo 3000",
			CmnaOrigen: "Las Condes",
		},
		Receptor: domain.Receptor{
			RUTRecep:    "12236547-6",
			RznSocRecep: "Cliente de Prueba Ltda",
		},
		Totales: domain.Totales{MntNeto: 100000, IVA: 19000, MntTotal: 119000},
	},
	Detalle: []domain.Detalle{{NroLinDet: 1, NmbItem: "Servicio", MontoItem: 100000}},
}

// El tercer argumento es el idempotency-key: un UUID, o "" para que el SDK genere uno.
req, err := domain.DTE33ToRequest("user_id", "business_id", "", dte)
if err != nil {
	panic(err)
}

resp, err := service.CreateDocument(context.Background(), req)
```

## Endpoints implementados

- Salud y acceso (sin `x-api-key`): `GetHealth`, `Login`, `CreateFirstBusiness`
- Usuarios y empresas: `GetMe`, `CreateBusiness`, `ListBusinesses`, `GetBusiness`, `UpdateBusiness`, `EnableProductionMode`, `EnableCertificationMode`
- Documentos: `CreateDocument`, `UpdateDocument`, `ListDocuments`, `GetDocument`, `GetDocumentStats`, `GetDocumentStatsWithFilter`, `RequeueDocument`, `RequeueOfflineDocumentStatus`
- Cesiones y PDF: `CreateCession`, `ListCessions`, `GetCession`, `RequeueCession`, `GeneratePDF`
- Certificados: `UploadCertificate`, `GetCertificateInfo`
- Billing: `GetBillingBalance`, `ListBillingPayments`, `ListBillingCharges`, `ListBillingPlans`, `ListBillingInvoices`, `PreviewSubscriptionUpgrade`
- Consumo: `GetConsumption`, `ListConsumptionOverages`, `ListConsumptionOperations`
- Compras: `CreatePurchase`, `ListPurchaseAcknowledgments`, `RequeuePurchase`
- Numeraciones: `GetNumerationSummary`, `GetLastUsedFolio`, `ListNumerationRanges`, `UploadNumeration`, `DeleteNumeration`, `UpdateNumerationNextNumber`, `UpdateLowStockConfig`, `RequestNumbers`

Todos devuelven `domain.APIResponse` (el JSON como mapa), salvo `RequestNumbers`, que devuelve `[]domain.NumberRange`.

| Metodo | Ruta | Parametros y notas |
|---|---|---|
| `GetHealth` | `GET /api/v1/health` | Sin credenciales. JSON crudo, sin el sobre `{success, message, data}`: `resp["service"]`, `resp["uptime_seconds"]` |
| `Login` | `POST /api/v1/auth/login` | `domain.LoginRequest`. Sin credenciales. Devuelve `data.xUserKey` |
| `CreateFirstBusiness` | `POST /api/v1/onboarding/businesses` | `userKey`, `domain.CreateFirstBusinessRequest`. Usa `x-user-key` en vez de `x-api-key`. Devuelve la API key en `data.apiToken.xApiKey` |
| `UpdateDocument` | `PUT /api/v1/documents/:id` | `domain.UpdateDocumentRequest` (`DataDTE` o `DataDTEJSON`) |
| `ListNumerationRanges` | `GET /api/v1/numerations/ranges` | `domain.NumerationRangeFilter` (`code_sii`) |
| `UpdateNumerationNextNumber` | `PATCH /api/v1/numerations/:numerationId/next-number` | `numerationID` es `ranges[].id` de `ListNumerationRanges`; `domain.UpdateNumerationNextNumberRequest` |
| `UpdateLowStockConfig` | `PATCH /api/v1/numerations/low-stock` | `domain.UpdateLowStockConfigRequest`. Se combina por `code_sii` y devuelve la configuracion completa |
| `RequeuePurchase` | `POST /api/v1/purchase-acknowledgments/requeue` | `domain.RequeuePurchaseRequest` |
| `ListBillingCharges` | `GET /api/v1/billing/charges` | `domain.ChargeFilter` (`status`, `pricing_key`, `from_date`, `to_date`, `page`, `limit`) |
| `ListBillingPlans` | `GET /api/v1/billing/plans` | `data` es un arreglo |
| `ListBillingInvoices` | `GET /api/v1/billing/invoices` | `domain.InvoiceFilter` (`status`). `data` es un arreglo |
| `PreviewSubscriptionUpgrade` | `GET /api/v1/billing/subscription/upgrade/preview` | `planID` (id o `code` del plan). Solo cotiza, no cobra |
| `GetConsumption` | `GET /api/v1/consumption` | Consumo del ciclo actual |
| `ListConsumptionOverages` | `GET /api/v1/consumption/overages` | `domain.ConsumptionOverageFilter` (`page`, `limit`). El total viene en `data.total` |
| `ListConsumptionOperations` | `GET /api/v1/consumption/operations` | `domain.ConsumptionOperationFilter` (`period`, `YYYY-MM`). Sin paginar |
| `RequeueCession` | `POST /api/v1/cessions/requeue` | `domain.RequeueCessionRequest` |
| `ListCessions` | `GET /api/v1/cessions` | `domain.CessionFilter` (`document_id`, `page`, `limit`). La lista viene en `data.cessions` y el total en `data.total` |
| `GetCession` | `GET /api/v1/cessions/:id` | |

## Idempotencia

Estas rutas de la API exigen el header `idempotency-key` y responden 400 sin el. El valor debe ser un UUID (sirve cualquier version):

| Metodo | Ruta |
|---|---|
| `CreateDocument` | `POST /api/v1/documents` |
| `UpdateDocument` | `PUT /api/v1/documents/:id` |
| `CreateBusiness` | `POST /api/v1/businesses` |
| `UpdateBusiness` | `PUT /api/v1/businesses/:id` |
| `UploadCertificate` | `PUT /api/v1/business/:id/certificate` |
| `UploadNumeration` | `PUT /api/v1/numerations` |
| `DeleteNumeration` | `DELETE /api/v1/numerations/:id` |
| `UpdateNumerationNextNumber` | `PATCH /api/v1/numerations/:numerationId/next-number` |
| `UpdateLowStockConfig` | `PATCH /api/v1/numerations/low-stock` |
| `CreatePurchase` | `POST /api/v1/purchase-acknowledgments` |
| `CreateCession` | `POST /api/v1/cessions` |

En esas rutas el SDK siempre envia el header. Usa, en este orden:

1. El campo `IdempotencyKey` del request, si lo indicas.
2. La clave guardada en el contexto con `httpintegra.WithIdempotencyKey(ctx, key)`. Sirve para los metodos que no reciben un request, como `DeleteNumeration`.
3. Si no hay ninguna, un UUID v4 nuevo en cada llamada.

```go
ctx := httpintegra.WithIdempotencyKey(context.Background(), "0190d7a4-6c3e-7b8e-9f21-3a4b5c6d7e8f")
resp, err := service.DeleteNumeration(ctx, "id-del-rango")
```

La API guarda cada clave por usuario y ruta durante 24 horas y no compara el body: si repites una clave, devuelve la primera respuesta. Por eso:

- Usa una clave nueva para cada operacion. Repite una clave solo para reintentar una operacion cuya respuesta no alcanzaste a recibir.
- Si la API respondio un error, reintenta con una clave nueva. La API no guarda algunas respuestas (por ejemplo, los errores de validacion) y un reintento con la misma clave responde 500 `failed to parse cached response`.
- `UpdateDocument` nunca guarda su respuesta, asi que cada intento necesita una clave nueva. Si no indicas clave, el SDK ya genera una por llamada.

`GeneratePDF` no es una ruta idempotente: solo envia el header si indicas `IdempotencyKey`.

## Onboarding sin API key

`Login` y `CreateFirstBusiness` se usan antes de tener API key. Como `httpintegra.New` exige `APIKey`, para ese paso usa `httpintegra.NewWithoutAPIKey`. Con un cliente sin API key, los metodos que necesitan `x-api-key` devuelven un error sin llamar a la API.

```go
ctx := context.Background()

onboarding, err := httpintegra.NewWithoutAPIKey(httpintegra.Config{})
if err != nil {
	panic(err)
}

login, err := onboarding.Login(ctx, domain.LoginRequest{Email: "tu@empresa.cl", Password: "tu-clave"})
if err != nil {
	panic(err)
}
loginData, _ := login["data"].(map[string]any)
userKey, _ := loginData["xUserKey"].(string)

business, err := onboarding.CreateFirstBusiness(ctx, userKey, domain.CreateFirstBusinessRequest{
	BusinessName:           "EMPRESA DE PRUEBA SpA",
	RUT:                    "76000000-0",
	Activity:               "Servicios de desarrollo de software",
	Address:                "Av. Apoquindo 3000",
	Commune:                "Las Condes",
	Region:                 "Metropolitana", // Region o City
	EmailDTE:               "dte@empresa.cl",
	EmailContact:           "contacto@empresa.cl",
	RUTLegalAgent:          "11111111-1",
	FullNameLegalAgent:     "Nombre Representante",
	ResolutionNumberDTE:    "0",
	ResolutionDateDTE:      "2014-08-22",
	ResolutionNumberTicket: "0",
	ResolutionTicketDate:   "2014-08-22",
})
if err != nil {
	panic(err)
}
businessData, _ := business["data"].(map[string]any)
apiToken, _ := businessData["apiToken"].(map[string]any)
apiKey, _ := apiToken["xApiKey"].(string)

adapter, err := httpintegra.New(httpintegra.Config{APIKey: apiKey})
```

## Certificado digital

`GetCertificateInfo` solo indica si la empresa puede firmar: `data.has_valid_certificate` es `true` cuando la empresa tiene certificado, abre con su clave y no esta vencido (la misma validacion que usa la emision). Si la empresa no tiene certificado responde `false`, no un error. La API no devuelve datos del certificado ni su clave privada; el tipo `domain.CertificateInfo` describe el contenido de `data`.

```go
resp, err := service.GetCertificateInfo(context.Background())
if err != nil {
	panic(err)
}

data, _ := resp["data"].(map[string]any)
valid, _ := data["has_valid_certificate"].(bool)
fmt.Println("certificado valido:", valid)
```

## Versionado y releases automaticos

El repo incluye dos workflows:

- `CI` (`.github/workflows/ci.yml`): ejecuta `go mod tidy`, `go vet` y `go test ./...` en push/PR.
- `Release Please` (`.github/workflows/release-please.yml`): genera PR de release, crea tag semantico (`vX.Y.Z`) y GitHub Release.

### Como disparar una nueva version

Usa Conventional Commits en `main`:

- `fix: ...` => patch
- `feat: ...` => minor
- `feat!: ...` o `BREAKING CHANGE:` => major

Cuando haya cambios, Release Please abrira/actualizara un PR de release. Al mergearlo, crea el tag y la release automaticamente.
