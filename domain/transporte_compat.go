package domain

import "encoding/json"

// Compatibilidad del bloque Transporte.
//
// El esquema del SII ubica Transporte DENTRO de Encabezado, entre Receptor y Totales, y exige
// que el chofer vaya en un bloque Chofer propio. Las versiones anteriores de este SDK lo
// pedian al tope del documento y con el chofer plano. Ninguna de esas dos formas validaba
// nunca contra el esquema, pero si hay codigo escrito asi.
//
// Estos MarshalJSON aceptan las dos formas y serializan siempre la que define el esquema, asi
// el payload sale correcto en el cable sin depender de que version de la API este desplegada.

// MarshalJSON mueve el chofer plano al bloque Chofer y descarta NroLicencia, que el SII no
// declara en ninguna version del esquema.
func (t Transporte) MarshalJSON() ([]byte, error) {
	type plano Transporte
	p := plano(t)

	if p.Chofer == nil && (p.NombreChofer != "" || p.RUTChofer != "") {
		p.Chofer = &Chofer{RUTChofer: p.RUTChofer, NombreChofer: p.NombreChofer}
	}
	p.NombreChofer = ""
	p.RUTChofer = ""
	p.NroLicencia = ""

	return json.Marshal(p)
}

// migrarTransporte mueve el bloque del tope del documento al Encabezado, que es el unico lugar
// desde el que se serializa. Si se setearon los dos, gana el del Encabezado: es la forma
// correcta y la que se eligio explicitamente.
func migrarTransporte(compat **Transporte, encabezado **Transporte) {
	if *compat == nil {
		return
	}
	if *encabezado == nil {
		*encabezado = *compat
	}
	*compat = nil
}

func (d Dte33Data) MarshalJSON() ([]byte, error) {
	type plano Dte33Data
	p := plano(d)
	migrarTransporte(&p.TransporteCompat, &p.Encabezado.Transporte)
	return json.Marshal(p)
}

func (d Dte34Data) MarshalJSON() ([]byte, error) {
	type plano Dte34Data
	p := plano(d)
	migrarTransporte(&p.TransporteCompat, &p.Encabezado.Transporte)
	return json.Marshal(p)
}

func (d Dte46Data) MarshalJSON() ([]byte, error) {
	type plano Dte46Data
	p := plano(d)
	migrarTransporte(&p.TransporteCompat, &p.Encabezado.Transporte)
	return json.Marshal(p)
}

func (d Dte52Data) MarshalJSON() ([]byte, error) {
	type plano Dte52Data
	p := plano(d)
	migrarTransporte(&p.TransporteCompat, &p.Encabezado.Transporte)
	return json.Marshal(p)
}

func (d Dte56Data) MarshalJSON() ([]byte, error) {
	type plano Dte56Data
	p := plano(d)
	migrarTransporte(&p.TransporteCompat, &p.Encabezado.Transporte)
	return json.Marshal(p)
}

func (d Dte61Data) MarshalJSON() ([]byte, error) {
	type plano Dte61Data
	p := plano(d)
	migrarTransporte(&p.TransporteCompat, &p.Encabezado.Transporte)
	return json.Marshal(p)
}
