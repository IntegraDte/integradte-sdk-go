package domain

// Common SII DTE shared structures used to build data_dte JSON payloads.

type Emisor struct {
	RUTEmisor    string `json:"RUTEmisor"`
	RznSoc       string `json:"RznSoc"`
	GiroEmis     string `json:"GiroEmis"`
	Telefono     string `json:"Telefono,omitempty"`
	CorreoEmisor string `json:"CorreoEmisor,omitempty"`
	Acteco       []int  `json:"Acteco,omitempty"`
	CdgSIISucur  int    `json:"CdgSIISucur,omitempty"`
	DirOrigen    string `json:"DirOrigen"`
	CmnaOrigen   string `json:"CmnaOrigen"`
	CiudadOrigen string `json:"CiudadOrigen,omitempty"`
}

type EmisorBoleta struct {
	RUTEmisor    string `json:"RUTEmisor"`
	RznSocEmisor string `json:"RznSocEmisor"`
	GiroEmisor   string `json:"GiroEmisor"`
	Telefono     string `json:"Telefono,omitempty"`
	CorreoEmisor string `json:"CorreoEmisor,omitempty"`
	CdgSIISucur  int    `json:"CdgSIISucur,omitempty"`
	DirOrigen    string `json:"DirOrigen"`
	CmnaOrigen   string `json:"CmnaOrigen"`
	CiudadOrigen string `json:"CiudadOrigen,omitempty"`
}

type Receptor struct {
	RUTRecep     string `json:"RUTRecep"`
	CdgIntRecep  string `json:"CdgIntRecep,omitempty"`
	RznSocRecep  string `json:"RznSocRecep"`
	GiroRecep    string `json:"GiroRecep,omitempty"`
	Contacto     string `json:"Contacto,omitempty"`
	CorreoRecep  string `json:"CorreoRecep,omitempty"`
	DirRecep     string `json:"DirRecep,omitempty"`
	CmnaRecep    string `json:"CmnaRecep,omitempty"`
	CiudadRecep  string `json:"CiudadRecep,omitempty"`
	DirPostal    string `json:"DirPostal,omitempty"`
	CmnaPostal   string `json:"CmnaPostal,omitempty"`
	CiudadPostal string `json:"CiudadPostal,omitempty"`
}

type ReceptorBoleta struct {
	RUTRecep    string `json:"RUTRecep,omitempty"`
	RznSocRecep string `json:"RznSocRecep,omitempty"`
	GiroRecep   string `json:"GiroRecep,omitempty"`
	Contacto    string `json:"Contacto,omitempty"`
	CorreoRecep string `json:"CorreoRecep,omitempty"`
	DirRecep    string `json:"DirRecep,omitempty"`
	CmnaRecep   string `json:"CmnaRecep,omitempty"`
	CiudadRecep string `json:"CiudadRecep,omitempty"`
}

type ImptoReten struct {
	TipoImp  string  `json:"TipoImp,omitempty"`
	TasaImp  float64 `json:"TasaImp,omitempty"`
	MontoImp int     `json:"MontoImp,omitempty"`
}

type Totales struct {
	MntNeto       int          `json:"MntNeto,omitempty"`
	MntExe        int          `json:"MntExe,omitempty"`
	TasaIVA       float64      `json:"TasaIVA,omitempty"`
	IVA           int          `json:"IVA,omitempty"`
	IVAProp       int          `json:"IVAProp,omitempty"`
	IVANoRet      int          `json:"IVANoRet,omitempty"`
	ImptoReten    []ImptoReten `json:"ImptoReten,omitempty"`
	MntTotal      int          `json:"MntTotal"`
	MontoNF       int          `json:"MontoNF,omitempty"`
	TotalPeriodo  int          `json:"TotalPeriodo,omitempty"`
	SaldoAnterior int          `json:"SaldoAnterior,omitempty"`
	VlrPagar      int          `json:"VlrPagar,omitempty"`
}

type TotalesGuia struct {
	MntNeto    int          `json:"MntNeto,omitempty"`
	MntExe     int          `json:"MntExe,omitempty"`
	TasaIVA    float64      `json:"TasaIVA,omitempty"`
	IVA        int          `json:"IVA,omitempty"`
	ImptoReten []ImptoReten `json:"ImptoReten,omitempty"`
	MntTotal   int          `json:"MntTotal,omitempty"`
}

type TotalesBoleta struct {
	MntNeto    int          `json:"MntNeto,omitempty"`
	MntExe     int          `json:"MntExe,omitempty"`
	IVA        int          `json:"IVA,omitempty"`
	ImptoReten []ImptoReten `json:"ImptoReten,omitempty"`
	MntTotal   int          `json:"MntTotal"`
}

type IdDocBase struct {
	TipoDTE       int    `json:"TipoDTE"`
	Folio         int    `json:"Folio"`
	FchEmis       string `json:"FchEmis"`
	IndServicio   int    `json:"IndServicio,omitempty"`
	FmaPago       int    `json:"FmaPago,omitempty"`
	PeriodoDesde  string `json:"PeriodoDesde,omitempty"`
	PeriodoHasta  string `json:"PeriodoHasta,omitempty"`
	TermPagoGlosa string `json:"TermPagoGlosa,omitempty"`
	FchVenc       string `json:"FchVenc,omitempty"`
	MedioPago     string `json:"MedioPago,omitempty"`
	TpoCtaPago    string `json:"TpoCtaPago,omitempty"`
	NumCtaPago    string `json:"NumCtaPago,omitempty"`
	MntBruto      int    `json:"MntBruto,omitempty"`
	IndMntNeto    int    `json:"IndMntNeto,omitempty"`
}

type IdDocBoleta struct {
	TipoDTE     int    `json:"TipoDTE"`
	Folio       int    `json:"Folio"`
	FchEmis     string `json:"FchEmis"`
	IndServicio int    `json:"IndServicio,omitempty"`
	FmaPago     int    `json:"FmaPago,omitempty"`
	MedioPago   string `json:"MedioPago,omitempty"`
	MntBruto    int    `json:"MntBruto,omitempty"`
	IndMntNeto  int    `json:"IndMntNeto,omitempty"`
}

type IdDocGuia struct {
	TipoDTE       int    `json:"TipoDTE"`
	Folio         int    `json:"Folio"`
	FchEmis       string `json:"FchEmis"`
	TipoDespacho  int    `json:"TipoDespacho,omitempty"`
	IndTraslado   int    `json:"IndTraslado,omitempty"`
	IndServicio   int    `json:"IndServicio,omitempty"`
	FmaPago       int    `json:"FmaPago,omitempty"`
	FchVenc       string `json:"FchVenc,omitempty"`
	PeriodoDesde  string `json:"PeriodoDesde,omitempty"`
	PeriodoHasta  string `json:"PeriodoHasta,omitempty"`
	TermPagoGlosa string `json:"TermPagoGlosa,omitempty"`
	MntBruto      int    `json:"MntBruto,omitempty"`
	IndMntNeto    int    `json:"IndMntNeto,omitempty"`
}

type CdgItem struct {
	TpoCodigo string `json:"TpoCodigo,omitempty"`
	VlrCodigo string `json:"VlrCodigo,omitempty"`
}

type Detalle struct {
	NroLinDet      int       `json:"NroLinDet"`
	CdgItem        []CdgItem `json:"CdgItem,omitempty"`
	IndExe         int       `json:"IndExe,omitempty"`
	NmbItem        string    `json:"NmbItem"`
	DscItem        string    `json:"DscItem,omitempty"`
	QtyItem        float64   `json:"QtyItem,omitempty"`
	UnmdItem       string    `json:"UnmdItem,omitempty"`
	PrcItem        float64   `json:"PrcItem,omitempty"`
	DescuentoPct   float64   `json:"DescuentoPct,omitempty"`
	DescuentoMonto int       `json:"DescuentoMonto,omitempty"`
	RecargoPct     float64   `json:"RecargoPct,omitempty"`
	RecargoMonto   int       `json:"RecargoMonto,omitempty"`
	CodImpAdic     []string  `json:"CodImpAdic,omitempty"`
	MontoItem      int       `json:"MontoItem"`
}

type DscRcgGlobal struct {
	NroLinDR int     `json:"NroLinDR"`
	TpoMov   string  `json:"TpoMov"`
	GlosaDR  string  `json:"GlosaDR,omitempty"`
	TpoValor string  `json:"TpoValor"`
	ValorDR  float64 `json:"ValorDR"`
	IndExeDR int     `json:"IndExeDR,omitempty"`
}

type Referencia struct {
	NroLinRef int    `json:"NroLinRef,omitempty"`
	TpoDocRef string `json:"TpoDocRef,omitempty"`
	FolioRef  string `json:"FolioRef,omitempty"`
	FchRef    string `json:"FchRef,omitempty"`
	CodRef    int    `json:"CodRef,omitempty"`
	RazonRef  string `json:"RazonRef,omitempty"`
}

// Chofer identifica a quien conduce. El esquema del SII lo exige como bloque propio dentro de
// Transporte, no como campos planos.
type Chofer struct {
	RUTChofer    string `json:"RUTChofer,omitempty"`
	NombreChofer string `json:"NombreChofer,omitempty"`
}

// Transporte es el bloque de transporte del documento, que va DENTRO de Encabezado, entre
// Receptor y Totales.
//
// PatenteCarro, FchSalida, HraSalida y FchLlegada los agrego el Formato DTE v2.5 (2026-02) y
// la Res. Ex. SII N 144 de 2026 los volvio obligatorios para los documentos que amparan
// traslado de bienes. Aca son opcionales porque la API todavia no los exige.
//
// El SII valida los formatos al caracter: FchSalida y FchLlegada en AAAA-MM-DD, y HraSalida en
// HH:MM:SS con segundos y ceros a la izquierda. Tanto "08:30" como "8:30:00" se rechazan.
type Transporte struct {
	Patente      string  `json:"Patente,omitempty"`
	PatenteCarro string  `json:"PatenteCarro,omitempty"`
	RUTTrans     string  `json:"RUTTrans,omitempty"`
	Chofer       *Chofer `json:"Chofer,omitempty"`
	DirDest      string  `json:"DirDest,omitempty"`
	CmnaDest     string  `json:"CmnaDest,omitempty"`
	CiudadDest   string  `json:"CiudadDest,omitempty"`
	FchSalida    string  `json:"FchSalida,omitempty"`
	HraSalida    string  `json:"HraSalida,omitempty"`
	FchLlegada   string  `json:"FchLlegada,omitempty"`

	// Deprecated: usar Chofer.NombreChofer. Se acepta y se mueve al serializar.
	NombreChofer string `json:"-"`
	// Deprecated: usar Chofer.RUTChofer. Se acepta y se mueve al serializar.
	RUTChofer string `json:"-"`
	// Deprecated: el SII no lo declara en ningun esquema. Se descarta al serializar.
	NroLicencia string `json:"-"`
}

// TransporteGuia era el bloque propio de la guia (52). El esquema del SII define un unico
// Transporte igual para todos los tipos.
//
// Deprecated: usar Transporte.
type TransporteGuia = Transporte

type SubTotInfo struct {
	NroSTI      int    `json:"NroSTI,omitempty"`
	GlosaSTI    string `json:"GlosaSTI,omitempty"`
	OrdenSTI    int    `json:"OrdenSTI,omitempty"`
	SubTotNeto  int    `json:"SubTotNeto,omitempty"`
	SubTotExe   int    `json:"SubTotExe,omitempty"`
	SubTotIVA   int    `json:"SubTotIVA,omitempty"`
	SubTotAdic  int    `json:"SubTotAdic,omitempty"`
	SubTotTotal int    `json:"SubTotTotal,omitempty"`
}

type Comisiones struct {
	NroLinCom    int     `json:"NroLinCom,omitempty"`
	TipoMovim    string  `json:"TipoMovim,omitempty"`
	Glosa        string  `json:"Glosa,omitempty"`
	TasaComision float64 `json:"TasaComision,omitempty"`
	ValComision  int     `json:"ValComision,omitempty"`
}

type OtraMoneda struct {
	TpoMoneda     string  `json:"TpoMoneda,omitempty"`
	TpoCambio     float64 `json:"TpoCambio,omitempty"`
	MntNetoOtrMnd int     `json:"MntNetoOtrMnd,omitempty"`
	MntExeOtrMnd  int     `json:"MntExeOtrMnd,omitempty"`
	IVAOtrMnd     int     `json:"IVAOtrMnd,omitempty"`
	MntTotOtrMnd  int     `json:"MntTotOtrMnd,omitempty"`
}
