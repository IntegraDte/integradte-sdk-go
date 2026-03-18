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
		CodeSII:        "33",
		DataDTE:        dataDTE,
		IdempotencyKey: "mi-idempotency-key-1",
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

req, err := domain.DTE33ToRequest("user_id", "business_id", "idempotency-key", dte)
if err != nil {
	panic(err)
}

resp, err := service.CreateDocument(context.Background(), req)
```

## Endpoints implementados

- `CreateDocument`
- `GetDocument`
- `GetDocumentStats`
- `CreateCession`
- `GeneratePDF`
- `CreateBusiness`
- `UpdateBusiness`
- `UploadCertificate`
- `GetCertificateInfo`
- `GetMe`
- `CreatePurchase`
- `GetNumerationSummary`
- `GetLastUsedFolio`
- `UploadNumeration`
- `DeleteNumeration`

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
