package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

var (
	emisorPrueba = Emisor{
		RUTEmisor:  "76326028-3",
		RznSoc:     "SERVICIOS INFORMATICOS SIGANET LIMITADA",
		GiroEmis:   "Servicios de desarrollo de software",
		DirOrigen:  "Av. Apoquindo 3000",
		CmnaOrigen: "Las Condes",
	}
	receptorPrueba = Receptor{RUTRecep: "18923860-6", RznSocRecep: "Cliente"}
	detallePrueba  = []Detalle{{NroLinDet: 1, NmbItem: "Aceite", MontoItem: 100000}}
)

// El bloque completo que exige la Res. Ex. SII N 144 de 2026, en la forma del esquema.
func transporteCanonico() *Transporte {
	return &Transporte{
		Patente:      "UVSS30",
		PatenteCarro: "PPRR50",
		RUTTrans:     "11111111-1",
		Chofer:       &Chofer{RUTChofer: "12000000-0", NombreChofer: "Juan Perez"},
		DirDest:      "Av Destino 456",
		CmnaDest:     "Quilpue",
		CiudadDest:   "Quilpue",
		FchSalida:    "2026-11-05",
		HraSalida:    "08:30:00",
		FchLlegada:   "2026-11-06",
	}
}

// El mismo bloque como lo pedian las versiones anteriores: chofer plano.
func transporteLegado() *Transporte {
	return &Transporte{
		Patente:      "UVSS30",
		PatenteCarro: "PPRR50",
		RUTTrans:     "11111111-1",
		NombreChofer: "Juan Perez",
		RUTChofer:    "12000000-0",
		DirDest:      "Av Destino 456",
		CmnaDest:     "Quilpue",
		CiudadDest:   "Quilpue",
		FchSalida:    "2026-11-05",
		HraSalida:    "08:30:00",
		FchLlegada:   "2026-11-06",
	}
}

func guia52(enEncabezado, alTope *Transporte) Dte52Data {
	return Dte52Data{
		Encabezado: Encabezado52{
			IdDoc:      IdDocGuia{TipoDTE: 52, Folio: 1, FchEmis: "2026-11-05", TipoDespacho: 1, IndTraslado: 1},
			Emisor:     emisorPrueba,
			Receptor:   receptorPrueba,
			Transporte: enEncabezado,
			Totales:    TotalesGuia{MntNeto: 100000, IVA: 19000, MntTotal: 119000},
		},
		Detalle:          detallePrueba,
		TransporteCompat: alTope,
	}
}

func payload(t *testing.T, v any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return m
}

func bloqueTransporte(t *testing.T, v any) map[string]any {
	t.Helper()
	p := payload(t, v)
	enc, ok := p["Encabezado"].(map[string]any)
	if !ok {
		t.Fatal("no hay Encabezado en el payload")
	}
	bloque, ok := enc["Transporte"].(map[string]any)
	if !ok {
		t.Fatalf("no hay Transporte dentro del Encabezado: %v", enc)
	}
	return bloque
}

func TestTransporte_DelTopePasaAlEncabezado(t *testing.T) {
	p := payload(t, guia52(nil, transporteCanonico()))

	if _, hay := p["Transporte"]; hay {
		t.Error("el bloque quedo al tope del documento")
	}
	bloqueTransporte(t, guia52(nil, transporteCanonico()))
}

func TestTransporte_LasDosFormasProducenElMismoPayload(t *testing.T) {
	desdeEncabezado, err := json.Marshal(guia52(transporteCanonico(), nil))
	if err != nil {
		t.Fatal(err)
	}
	desdeElTope, err := json.Marshal(guia52(nil, transporteLegado()))
	if err != nil {
		t.Fatal(err)
	}

	if string(desdeEncabezado) != string(desdeElTope) {
		t.Errorf("los payloads difieren\nencabezado: %s\ntope:       %s", desdeEncabezado, desdeElTope)
	}
}

func TestTransporte_ElChoferPlanoVaAlBloqueChofer(t *testing.T) {
	bloque := bloqueTransporte(t, guia52(nil, transporteLegado()))

	chofer, ok := bloque["Chofer"].(map[string]any)
	if !ok {
		t.Fatalf("no se armo el bloque Chofer: %v", bloque)
	}
	if chofer["NombreChofer"] != "Juan Perez" {
		t.Errorf("NombreChofer = %v", chofer["NombreChofer"])
	}
	if chofer["RUTChofer"] != "12000000-0" {
		t.Errorf("RUTChofer = %v", chofer["RUTChofer"])
	}
	if _, hay := bloque["NombreChofer"]; hay {
		t.Error("quedo NombreChofer plano en el bloque")
	}
}

func TestTransporte_NroLicenciaNoSeSerializa(t *testing.T) {
	raw, err := json.Marshal(guia52(nil, &Transporte{Patente: "UVSS30", NroLicencia: "A1234"}))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "NroLicencia") {
		t.Errorf("NroLicencia llego al payload: %s", raw)
	}
	if !strings.Contains(string(raw), "UVSS30") {
		t.Errorf("se perdio el resto del bloque: %s", raw)
	}
}

func TestTransporte_ElEncabezadoLePisaAlTope(t *testing.T) {
	bloque := bloqueTransporte(t, guia52(
		&Transporte{Patente: "CANONICA"},
		&Transporte{Patente: "LEGADA"},
	))

	if bloque["Patente"] != "CANONICA" {
		t.Errorf("no gano el del Encabezado: %v", bloque["Patente"])
	}
}

func TestTransporte_ConservaFechasYHoraDeTraslado(t *testing.T) {
	bloque := bloqueTransporte(t, guia52(nil, transporteLegado()))

	for campo, esperado := range map[string]string{
		"FchSalida":    "2026-11-05",
		"HraSalida":    "08:30:00",
		"FchLlegada":   "2026-11-06",
		"PatenteCarro": "PPRR50",
	} {
		if bloque[campo] != esperado {
			t.Errorf("%s = %v, esperaba %q", campo, bloque[campo], esperado)
		}
	}
}

// El arreglo no es solo de la guia: la Res. Ex. SII N 144 de 2026 cubre tambien las facturas
// que amparan traslado de bienes.
func TestTransporte_LaFacturaTambienLoMueve(t *testing.T) {
	dte := Dte33Data{
		Encabezado: Encabezado33{
			IdDoc:    IdDocBase{TipoDTE: 33, Folio: 10, FchEmis: "2026-11-05"},
			Emisor:   emisorPrueba,
			Receptor: receptorPrueba,
			Totales:  Totales{MntNeto: 100000, IVA: 19000, MntTotal: 119000},
		},
		Detalle:          detallePrueba,
		TransporteCompat: &Transporte{Patente: "UVSS30", NombreChofer: "Juan Perez", RUTChofer: "12000000-0"},
	}

	bloque := bloqueTransporte(t, dte)
	if _, ok := bloque["Chofer"].(map[string]any); !ok {
		t.Fatalf("no se armo el bloque Chofer: %v", bloque)
	}
	if _, hay := payload(t, dte)["Transporte"]; hay {
		t.Error("el bloque quedo al tope del documento")
	}
}

func TestTransporte_SinTransporteNoAparece(t *testing.T) {
	dte := Dte33Data{
		Encabezado: Encabezado33{
			IdDoc:    IdDocBase{TipoDTE: 33, Folio: 10, FchEmis: "2026-11-05"},
			Emisor:   emisorPrueba,
			Receptor: receptorPrueba,
			Totales:  Totales{MntNeto: 100000, IVA: 19000, MntTotal: 119000},
		},
		Detalle: detallePrueba,
	}

	p := payload(t, dte)
	if _, hay := p["Transporte"]; hay {
		t.Error("apareció Transporte al tope sin que se setee")
	}
	enc := p["Encabezado"].(map[string]any)
	if _, hay := enc["Transporte"]; hay {
		t.Error("apareció Transporte en el Encabezado sin que se setee")
	}
}

// Serializar no tiene que mutar lo que recibe: el receptor es por valor justamente para eso.
func TestTransporte_NoMutaElDocumento(t *testing.T) {
	dte := guia52(nil, transporteCanonico())
	if _, err := json.Marshal(dte); err != nil {
		t.Fatal(err)
	}

	if dte.TransporteCompat == nil {
		t.Error("se vacio TransporteCompat del documento original")
	}
	if dte.Encabezado.Transporte != nil {
		t.Error("se escribio Encabezado.Transporte en el documento original")
	}
}
