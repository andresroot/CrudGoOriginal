package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)


func fechaInicial(periodo string,mes string) string {
	var miMes = mes
	i, _ := strconv.ParseInt(miMes, 0, 64)
	miMes= fmt.Sprintf("%02d", i)
	return periodo+"-"+miMes+"-01"
}

type ciudad struct {
	Codigo string `json:"Codigo"`
	CodigoCiudad string
	CodigoDepartamento string
	Nombre string `json:"Nombre"`
	NombreCiudad string
	NombreDepartamento string
}
type medioPago struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type Unidaddemedida struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type formaPago struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type tipoOrganizacion struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type FormaDePago struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type MedioDePago struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type documentoIdentificacion struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type regimenFiscal struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

type responsabilidadFiscal struct {
	Codigo string `json:"Codigo"`
	Nombre string `json:"nombre"`
}

// NOMBRE DE UNA FILA
func Nombre(tabla string,valor string) string{
	db := dbConn()
	var q = "select nombre from "+tabla+" where codigo=$1"
	selDB, err := db.Query(q,valor)
	var Nombre string
	err = selDB.Scan(&Nombre)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	return Nombre
}

func NombreModulo() string{

	d := empresa{}
	err := db.Get( &d,"SELECT * FROM empresa")
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	return d.Modulo

}
// NOMBRE FINANCIERO
func FinancieroNombre(codigo string) string{
	if codigo==""{
		return ""
	}else{

	db := dbConn()
		d := financiero{}
		err := db.Get( &d,"SELECT * FROM financiero where Codigo = $1",codigo)
		if err != nil {
			//log.Println("Error ciudad 2")
			panic(err.Error())
		}
			return d.Nombre
	}

}

// NOMBRE DOCUMENTO
func DocumentoNombre(codigo string) string{
	if codigo==""{
		return ""
	}else{

		db := dbConn()
		d := documento{}
		err := db.Get( &d,"SELECT * FROM documento where Codigo = $1",codigo)
		if err != nil {
			//log.Println("Error ciudad 2")
			panic(err.Error())
		}
		return d.Nombre
	}

}

// NOMBRE MODULO
func ModuloNombre(codigo string) string{
	if codigo==""{
		return ""
	}else{

		db := dbConn()
		d := financiero{}
		err := db.Get( &d,"SELECT * FROM financiero where Codigo = $1",codigo)
		if err != nil {
			//log.Println("Error ciudad 2")
			panic(err.Error())
		}
		return d.Nombre
	}

}

// LISTAD DOCUMENTO
func ListaDocumento() []documento{
	var lista= []documento{}
	db := dbConn()
	var q = "select codigo,nombre,consecutivo,inicial from documento order by cast(codigo as integer)"
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		var Consecutivo string
		var Inicial string
		err = selDB.Scan(&Codigo,&Nombre,&Consecutivo,&Inicial)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, documento{
			Codigo,Nombre,Consecutivo,Inicial})
	}
	return lista
}

func ListaDocumentoBanco() []documento{
	var lista= []documento{}
	db := dbConn()
	var q = "select codigo,nombre,consecutivo,inicial from documento where codigo='1' or codigo='2' order by codigo "
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		var Consecutivo string
		var Inicial string
		err = selDB.Scan(&Codigo,&Nombre,&Consecutivo,&Inicial)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, documento{
			Codigo,Nombre,Consecutivo,Inicial})
	}
	return lista
}

// LISTA TIPO ORGANIZACION
func ListaTipoOrganizacion() []tipoOrganizacion{
	var lista= []tipoOrganizacion{}
	db := dbConn()
	var q = "select codigo,nombre from tipoorganizacion"
	selDB, err := db.Query(q)
	if err != nil {
			//log.Println("Error ciudad 2")
			panic(err.Error())
		}
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		err = selDB.Scan(
			&Codigo,
			&Nombre)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, tipoOrganizacion{
			Codigo,Nombre})
	}
	return lista
}

// LISTA FINANCIERO
func ListaFinanciero() []financiero{
	var lista= []financiero{}
	db := dbConn()

	var q = "select codigo,nombre from financiero"
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		err = selDB.Scan(
			&Codigo,
			&Nombre)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, financiero{
			Codigo,Nombre})
	}
	return lista
}

// LISTA RESPONSABILIDAD FISCAL
func ListaResponsabilidadFiscal() []responsabilidadFiscal{
	var lista= []responsabilidadFiscal{}
	db := dbConn()
	var q = "select codigo,nombre from responsabilidadFiscal"
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		err = selDB.Scan(
			&Codigo, &Nombre)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, responsabilidadFiscal{
			Codigo,Nombre})
	}
	return lista
}

// LISTA REGIMEN FISCAL
func ListaRegimenFiscal() []regimenFiscal{
	var lista= []regimenFiscal{}
	db := dbConn()
	var q = "select codigo,nombre from regimenfiscal"
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		err = selDB.Scan(&Codigo, &Nombre)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, regimenFiscal{
			Codigo, Nombre})
	}
	return lista
}

// LISTA DOCUMENTO IDENTIFICACION
func ListaDocumentoIdentificacion() []documentoIdentificacion{
	var lista= []documentoIdentificacion{}
	db := dbConn()
	var q = "select codigo,nombre from documentoidentificacion"
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		err = selDB.Scan(&Codigo, &Nombre)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, documentoIdentificacion{
			Codigo,Nombre})
	}

	return lista
}

// LISTA RESOLUCION VENTA
func ListaResolucionventa() []Resolucionventa{
	log.Println("resolucionventa")
	db := dbConn()
	res := []Resolucionventa{}
	db.Select(&res, "SELECT * FROM resolucionventa ORDER BY codigo ASC")
	log.Println("resolucionventa2")
	return res
}

// LISTA RESOLUCION SOPORTE
func ListaResolucionsoporte() []Resolucionsoporte{
	log.Println("resolucionsoporte")
	db := dbConn()
	res := []Resolucionsoporte{}
	db.Select(&res, "SELECT * FROM resolucionsoporte ORDER BY codigo ASC")
	log.Println("resolucionsoporte2")
	return res
}




// LISTA FORMA DE PAGO
func ListaFormaDePago() []FormaDePago{
	log.Println("ciudad1")
	db := dbConn()
	res := []FormaDePago{}
	db.Select(&res, "SELECT * FROM formadepago ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}

// LISTA MEDIO DE PAGO
func ListaMedioDePago() []MedioDePago{
	log.Println("ciudad1")
	db := dbConn()
	res := []MedioDePago{}
	db.Select(&res, "SELECT * FROM mediodepago ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}

// LISTA BODEGA
func ListaBodega() []bodega{
	log.Println("ciudad1")
	db := dbConn()
	res := []bodega{}
	db.Select(&res, "SELECT * FROM bodega ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}

// LISTA PRODUCTO
func ListaProducto() []producto{
	log.Println("producto")
	db := dbConn()
	res := []producto{}
	db.Select(&res, "SELECT * FROM producto ORDER BY codigo ASC")
	log.Println("producto")
	return res
}

// LISTA VENDEDOR
func ListaVendedor() []vendedor{
	log.Println("ciudad1")
	db := dbConn()
	res := []vendedor{}
	db.Select(&res, "SELECT * FROM vendedor ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}

// LISTA ALMACENISTA
func ListaAlmacenista() []almacenista{
	log.Println("ciudad1")
	db := dbConn()
	res := []almacenista{}
	db.Select(&res, "SELECT * FROM almacenista ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}

// CIUDAD
func ListaCiudad() []ciudad{
	log.Println("ciudad1")
	db := dbConn()
	res := []ciudad{}
	db.Select(&res, "SELECT * FROM ciudad ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}

func ListaCuentaBanco() []plandecuentaempresa{
	log.Println("ciudad1")
	db := dbConn()
	res := []plandecuentaempresa{}
	db.Select(&res, "SELECT * FROM plandecuentaempresa where substring(codigo,1,4)='1110' and nivel='A' ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}

func ListaCuentaAuxiliar() []plandecuentaempresa{
	log.Println("ciudad1")
	db := dbConn()
	res := []plandecuentaempresa{}
	db.Select(&res, "SELECT * FROM plandecuentaempresa where nivel='A' ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}
func ListaCuentaCaja() []plandecuentaempresa{
	log.Println("ciudad1")
	db := dbConn()
	res := []plandecuentaempresa{}
	db.Select(&res, "SELECT * FROM plandecuentaempresa where substring(codigo,1,2)='11' and nivel='A' ORDER BY codigo ASC")
	log.Println("ciudad2")
	return res
}
// LISTA GRUPO
func ListaGrupo() []grupo{
	log.Println("grupo1")
	db := dbConn()
	res := []grupo{}
	db.Select(&res, "SELECT * FROM grupo ORDER BY codigo ASC")
	log.Println("grupo2")
	return res
}

// LISTA CENTRO
func ListaCentro() []centro{
	log.Println("centro")
	db := dbConn()
	res := []centro{}
	db.Select(&res, "SELECT * FROM centro ORDER BY codigo ASC")
	log.Println("centro2")
	return res
}

// LISTA EMPRESA
func ListaEmpresa() empresa {
	log.Println("inicio empresa unica111")
	db := dbConn()
	t := empresa{}
	err := db.Get( &t,"SELECT * FROM empresa limit 1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio empresa unica2222")
	log.Println("inicio empresa unica2222"+t.Nombre)
	return t
}


// TRAER CIUDAD
func TraerCiudad(codigoCiudad string) ciudad {
	log.Println("inicio empresa unica111")
	db := dbConn()
	t := ciudad{}
	err := db.Get( &t,"SELECT * FROM ciudad where Codigo = $1",codigoCiudad)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio empresa unica2222")
	return t
}

func TraerDocumento(codigoDocumento string) documento {
	log.Println("inicio empresa unica111")
	db := dbConn()
	d := documento{}
	err := db.Get( &d,"SELECT * FROM documento where Codigo = $1",codigoDocumento)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio empresa unica2222")
	return d
}

// TRAER TIPO
func TraerTipo(codigoTipo string) string {
	log.Println("inicio tipo unica111")
	db := dbConn()
	t := tipoOrganizacion{}
	err := db.Get( &t,"SELECT * FROM tipoorganizacion where Codigo = $1",codigoTipo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio tipoorganizacion2222")
	return t.Nombre
}

// TRAER DOCUMENTO
func TraerDocumentoIdentificacion(codigoDocumento string) string {
	log.Println("inicio documento unica111")
	db := dbConn()
	t := documentoIdentificacion{}
	err := db.Get( &t,"SELECT * FROM documentoidentificacion where Codigo = $1",codigoDocumento)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio documento2222")
	return t.Nombre
}

// TRAER REGIMEN FISCAL
func TraerRegimen(codigoRegimen string) string {
	log.Println("inicio regimen unica111")
	db := dbConn()
	t := regimenFiscal{}
	err := db.Get( &t,"SELECT * FROM regimenfiscal where Codigo = $1",codigoRegimen)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio regimen2222")
	return t.Nombre
}

// TRAER RESPONSABILIDAD FISCAL
func TraerFiscal(codigoFiscal string) string {
	log.Println("inicio fiscal unica111")
	db := dbConn()
	t := responsabilidadFiscal{}
	err := db.Get( &t,"SELECT * FROM responsabilidadfiscal where Codigo = $1",codigoFiscal)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio fiscal2222")
	return t.Nombre
}

// TRAER GRUPO
func TraerGrupo(codigoGrupo string) string {
	log.Println("inicio grupo unica111")
	db := dbConn()
	t := grupo{}
	err := db.Get( &t,"SELECT * FROM grupo where Codigo = $1",codigoGrupo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio grpo2222")
	return t.Nombre
}

// LISTA SUBGRUPO
func ListaSubgrupo() []Subgrupo{
	var lista= []Subgrupo{}
	db := dbConn()
	var q = "select codigo,nombre from subgrupo"
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	log.Println("Error subgrupo 3")
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		err = selDB.Scan(
			&Codigo, &Nombre)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, Subgrupo{
			Codigo, Nombre, ""})
	}
	log.Println("Error subgrupo777")
	return lista
}

// LISTA UNIDAD DE MEDIDA
func ListaUnidaddemedida() []Unidaddemedida{
	var lista= []Unidaddemedida{}
	db := dbConn()
	var q = "select codigo,nombre from unidaddemedida"
	selDB, err := db.Query(q)
	if err != nil {
		//log.Println("Error ciudad 2")
		panic(err.Error())
	}
	log.Println("Error unidaddemedida 3")
	for selDB.Next() {
		var Codigo  string
		var Nombre  string
		err = selDB.Scan(&Codigo, &Nombre)
		if err != nil {
			//log.Println("Error ciudad 6")
			panic(err.Error())
		}
		lista = append(lista, Unidaddemedida{
			Codigo, Nombre})
	}
	log.Println("Error unidaddemedida 777")
	return lista
}

// TRAER SUBGRUPO
func TraerSubgrupo(codigoSubgrupo string) string {
	log.Println("inicio subgrupo unica111")
	db := dbConn()
	t := Subgrupo{}
	err := db.Get( &t,"SELECT * FROM subgrupo where Codigo = $1",codigoSubgrupo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio subgrupo2222")
	return t.Nombre
}

// TRAE UNIDAD DE MEDIDA
func TraerUnidaddemedida(codigoUnidaddemedida string) string {
	log.Println("inicio unidaddemedida unica111")
	db := dbConn()
	t := Unidaddemedida{}
	err := db.Get( &t,"SELECT * FROM unidaddemedida where Codigo = $1",codigoUnidaddemedida)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio unidaddemedida2222")
	return t.Nombre
}

// TRAE RESOLUCION VENTA
func TraerResolucionventa(codigoResolucionventa string) Resolucionventa {
	log.Println("inicio resolucionventa unica111")
	db := dbConn()
	t := Resolucionventa{}
	err := db.Get( &t,"SELECT * FROM resolucionventa where Codigo = $1",codigoResolucionventa)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio resolucionventa unica2222")
	return t
}


 // TRAE RESOLUCION SOPORTE
func TraerResolucionsoporte(codigoResolucionsoporte string) Resolucionsoporte {
	log.Println("inicio resolucionsoporte unica111")
	db := dbConn()
	t := Resolucionsoporte{}
	err := db.Get( &t,"SELECT * FROM resolucionsoporte where Codigo = $1",codigoResolucionsoporte)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio resolucionsoporte unica2222")
	return t
}


// TRAE MEDIO DE PAGO
func TraerMediodepago(codigoMediodepago string) string {
	log.Println("inicio mediodepago unica111")
	db := dbConn()
	t := MedioDePago{}
	err := db.Get( &t,"SELECT * FROM mediodepago where Codigo = $1",codigoMediodepago)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio codigomediodepago2222")
	return t.Nombre
}

// TRAE FORMA DE PAGO
func TraerFormadepago(codigoFormadepago string) string {
	log.Println("inicio formadepago unica111")
	db := dbConn()
	t := FormaDePago{}
	err := db.Get( &t,"SELECT * FROM formadepago where Codigo = $1",codigoFormadepago)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("inicio codigoformadepago2222")
	return t.Nombre
}

// CORTAR CADENA
func Subcadena(s string, start, end int) string {
	counter, startIdx := 0, 0
	for i := range s {
		if counter == start {
			startIdx = i
		}
		if counter == end {
			return s[startIdx:i]
		}
		counter++
	}
	return s[startIdx:]
}

// NUMEROS A LETRAS
var (
	//ErrValorNoAdmitido error para valor no admitidos
	ErrValorNoAdmitido = errors.New("Valor no admitido")
	us                 = []string{"cero", "uno", "dos", "tres", "cuatro", "cinco", "seis", "siete", "ocho", "nueve"}
	ds                 = []string{"X", "y", "veinte", "treinta", "cuarenta", "cincuenta", "sesenta", "setenta", "ochenta", "noventa"}
	des                = []string{"diez", "once", "doce", "trece", "catorce", "quince", "dieciseis", "diecisiete", "dieciocho", "diecinueve"}
	cs                 = []string{"x", "cien", "doscientos", "trescientos", "cuatrocientos", "quinientos", "seiscientos", "setecientos", "ochocientos", "novecientos"}
)

// NUMEROS A LETRAS
func IntLetra(n int) (s string, err error) {
	var aux string
	sb := strings.Builder{}
	if n < 0 {
		sb.WriteString("menos")
		n *= -1
	}
	millones := quotient(n, 1000000)
	if millones > 999999 {
		return s, ErrValorNoAdmitido
	}
	millares := quotient(n, 1000) % 1000
	centenas := quotient(n, 100) % 10
	decenas := quotient(n, 10) % 10
	unidades := n % 10
	if millones == 1 {
		sb.WriteString(" un millón")
	} else if millones > 1 {
		aux, err = IntLetra(millones)
		sb.WriteString(aux)
		sb.WriteString(" millones")
	}
	if millares == 1 {
		sb.WriteString(" mil")
	} else if millares > 1 {
		aux, err = IntLetra(millares)
		sb.WriteRune(' ')
		sb.WriteString(aux)
		sb.WriteString(" mil")
	}
	if centenas == 1 {
		if n%100 == 0 {
			sb.WriteString(" cien")
		} else {
			sb.WriteString(" ciento")
		}
	} else if centenas > 0 {
		sb.WriteRune(' ')
		sb.WriteString(cs[centenas])
	}
	if decenas == 1 {
		sb.WriteRune(' ')
		sb.WriteString(des[n%10])
		unidades = 0
	} else if decenas == 2 && unidades > 0 {
		sb.WriteString(" veinti")
		sb.WriteString(us[unidades])
		unidades = 0
	} else if decenas > 1 {
		sb.WriteRune(' ')
		sb.WriteString(ds[decenas])
		if unidades > 0 {
			sb.WriteString(" y")
		}
	}
	if unidades > 0 {
		sb.WriteRune(' ')
		sb.WriteString(us[unidades])
	} else if n == 0 {
		sb.WriteString(us[0])
	}
	return strings.TrimSpace(sb.String()), err
}

func quotient(numerator, denominator int) int {
	return numerator / denominator
}

// CADENA CONVIERTE A NUMERO
func Cadenaentero(numero string) int{
	resultado := strings.ReplaceAll(numero,",","")
	i,err := strconv.Atoi(resultado)
	if err!= nil{
		fmt.Println(err)
	}
	return i
}


// CADENA CONVIERTE A NUMERO
func Cadenapunto(numero string) string{

	r,err:=strconv.ParseInt(numero, 10,64)
	if err != nil {
		log.Fatalln(err)
	}
	numero = humanize.Comma(r)
	numero= strings.Replace(numero, ",", ".", -1)

	return numero
}
// BODEGA LLENAR
func BodegaLlenar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM bodega ORDER BY codigo DESC")
	if err != nil {
		panic(err.Error())
	}
	res := []bodega{}
	for selDB.Next() {
		var Codigo string
		var Nombre string
		err = selDB.Scan(&Codigo, &Nombre)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, bodega{Codigo, Nombre})
	}
	data, _ := json.Marshal(res)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// PONER COMAS A LOS NUMEROS
func Coma(s string) string{
	var r int64
	var numero string
	if s==""{
		s="0"
	}
	s= strings.Replace(s, ".", "", -1)
	s= strings.Replace(s, ",", "", -1)

	r,err:=strconv.ParseInt(s, 10,64)
	if err != nil {
		log.Fatalln(err)
	}
	numero = humanize.Comma(r)
	if err != nil {
		log.Fatalln(err)
	}
	return numero
}

// TIPO DE LETRA ORACION
func Titulo(s string) string {
	return strings.Title(strings.ToLower(s))
}

// TIPO DE LETRA MAYUSCULAS
func Mayuscula(s string) string {
	return strings.ToUpper(s)
}

// TIPO DE LETRA MINUSCULAS
func Minuscula(s string) string {
	return strings.ToLower(s)
}
// QUITA PUNTOS
func Puntos(s string) string {
	return strings.Replace(s, ".", "", -1)
}
func Quitapuntos(s string) string {
	return strings.Replace(s, ".", "", -1)
}

func Quitacoma(s string) string {
	return strings.Replace(s, ",", "", -1)
}
// FLOTANTE A NUMERO
func Flotante(numero string) float64{
	if numero==""{
		numero="0"
	}
	resultado := strings.ReplaceAll(numero,",","")
	//s := fmt.Sprintf("%v", b)
	log.Println("string a flotante  " + resultado)
	f1, err := strconv.ParseFloat(resultado, 8)

	//i,err := strconv.Atoi(resultado)
	if err!= nil{
		fmt.Println("Error flotante")
		fmt.Println(err)
	}
	return f1
}

func Entero(numero string) int{
	if numero==""{
		numero="0"
	}
	i, err := strconv.Atoi(numero)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	return i
}

func Flotantedatabase(numero string) string{
	if numero==""{
		numero="0.00"
	}
	//resultado := strings.ReplaceAll(numero,",","")
	//resultado = strings.ReplaceAll(resultado,".",",")

	resultado := strings.ReplaceAll(numero,",","")
	//resultado = strings.ReplaceAll(resultado,".",",")

	//s := fmt.Sprintf("%v", b)
	//*f1, err := strconv.ParseFloat(resultado, 8)
	//
	////i,err := strconv.Atoi(resultado)
	//if err!= nil{
	//	fmt.Println("Error flotante")
	//	fmt.Println(err)
	//}
	//return f1
	return resultado
}

func mesLetras(miMes string) string{
	switch miMes{
	case "1":
		return "Enero"
	case "2":
		return "Febrero"
	case "3":
		return "Marzo"
	case "4":
		return "Abril"
	case "5":
		return "Mayo"
	case "6":
		return "Junio"
	case "7":
		return "Julio"
	case "8":
		return "Agosto"
	case "9":
		return "Septiembre"
	case "10":
		return "Octubre"
	case "11":
		return "Noviembre"
	case "12":
		return "Diciembre"
	case "13":
		return "De Cierre"
	default:
		return ""
	}
}

func parametros(totalParametro int) string{
	var miComa=""
	var cadenaNumero string
	var cadenaParametro string
	cadenaParametro=""
	for i := 1; i <= totalParametro; i++ {
		cadenaNumero=strconv.Itoa(i)
		cadenaParametro+=miComa+"$"+cadenaNumero
		if miComa=="" {
			miComa=","
		}

	}
	return cadenaParametro
}

// DE FLOTANTE A CADENA
func CadenaFlotante(NumeroFlotante float64) string{
	return fmt.Sprintf("%.2f",NumeroFlotante)
	//return fmt.Sprintf("%f",NumeroFlotante)

}

func TraerParametrosInventario()configuracioninventario{
// PARAMETROS INVENTARIO
	db := dbConn()
	parametrosinventario := configuracioninventario{}

	err := db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
	panic(err.Error())
	}
	return parametrosinventario
}

func TraerParametrosContabilidad()configuracioncontabilidad{
	// PARAMETROS INVENTARIO
	db := dbConn()
	parametroscontabilidad := configuracioncontabilidad{}

	err := db.Get(&parametroscontabilidad, "SELECT * FROM configuracioncontabilidad")
	if err != nil {
		panic(err.Error())
	}
	return parametroscontabilidad
}
func FormatoImprimir(numeroFlotante float64) string{
	var cadenaDebito=fmt.Sprintf("%.2f",numeroFlotante)
	log.Println("Cadena debito " + cadenaDebito)

	var enteroDebito string
	var decimalDebito string
	enteroDebito=Subcadena(cadenaDebito,0,strings.Index(cadenaDebito,"."))
	decimalDebito=Subcadena(cadenaDebito,strings.Index(cadenaDebito,".")+1, len(cadenaDebito))
	log.Println("Cadena entero " + enteroDebito)
	log.Println("Cadena decimal " + decimalDebito)
	var numero int64
	numero,err:=strconv.ParseInt(enteroDebito,10,64)
	if err != nil {
		panic(err.Error())
	}
	cadenaDebito=humanize.Comma(numero)+"."+decimalDebito

	if cadenaDebito=="0.00"{
		cadenaDebito=""
	}

	return cadenaDebito

}
func FormatoFlotante(numeroFlotante float64) string{
	var cadenaDebito=fmt.Sprintf("%.2f",numeroFlotante)
	log.Println("Cadena debito " + cadenaDebito)

	var enteroDebito string
	var decimalDebito string
	enteroDebito=Subcadena(cadenaDebito,0,strings.Index(cadenaDebito,"."))
	decimalDebito=Subcadena(cadenaDebito,strings.Index(cadenaDebito,".")+1, len(cadenaDebito))
	log.Println("Cadena entero " + enteroDebito)
	log.Println("Cadena decimal " + decimalDebito)
	var numero int64
	numero,err:=strconv.ParseInt(enteroDebito,10,64)
	if err != nil {
		panic(err.Error())
	}
	cadenaDebito=humanize.Comma(numero)+"."+decimalDebito

	log.Println("Cadena final " + cadenaDebito)
    return cadenaDebito

}

func FormatoFlotanteComprobante(numeroFlotante float64) string{
	var cadenaDebito=fmt.Sprintf("%.2f",numeroFlotante)
	log.Println("Cadena debito " + cadenaDebito)

	//var enteroDebito string
	//var decimalDebito string
	//enteroDebito=Subcadena(cadenaDebito,0,strings.Index(cadenaDebito,"."))
	//decimalDebito=Subcadena(cadenaDebito,strings.Index(cadenaDebito,".")+1, len(cadenaDebito))
	//log.Println("Cadena entero " + enteroDebito)
	//log.Println("Cadena decimal " + decimalDebito)
	//var numero int64
	//numero,err:=strconv.ParseInt(enteroDebito,10,64)
	//if err != nil {
	//	panic(err.Error())
	//}
	//cadenaDebito=(numero)+"."+decimalDebito

	log.Println("Cadena final " + cadenaDebito)
	return cadenaDebito

}
func FormatoNumeroComprobante(numeroFlotante string) string{
	var cadenaDebito=numeroFlotante
	log.Println("Cadena debito " + cadenaDebito)

	var enteroDebito string
	var decimalDebito string
	enteroDebito=Subcadena(cadenaDebito,0,strings.Index(cadenaDebito,"."))
	decimalDebito=Subcadena(cadenaDebito,strings.Index(cadenaDebito,".")+1, len(cadenaDebito))
	log.Println("Cadena entero " + enteroDebito)
	log.Println("Cadena decimal " + decimalDebito)
	var numero int64
	numero,err:=strconv.ParseInt(enteroDebito,10,64)
	if err != nil {
		panic(err.Error())
	}
	cadenaDebito=humanize.Comma(numero)
	cadenaDebito=strings.Replace(cadenaDebito,",",".",-1)
	cadenaDebito=cadenaDebito+","+decimalDebito
	if cadenaDebito=="0,00"{
		cadenaDebito=""
	}
	return cadenaDebito

}

func FormatoFlotanteEntero(numeroFlotante float64) string{
	var cadenaDebito=fmt.Sprintf("%.2f",numeroFlotante)
	var enteroDebito string
	var decimalDebito string
	enteroDebito=Subcadena(cadenaDebito,0,strings.Index(cadenaDebito,"."))
	decimalDebito=Subcadena(cadenaDebito,strings.Index(cadenaDebito,".")+1, len(cadenaDebito))
	log.Println("Cadena entero " + enteroDebito)
	log.Println("Cadena decimal " + decimalDebito)
	var numero int64
	numero,err:=strconv.ParseInt(enteroDebito,10,64)
	if err != nil {
		panic(err.Error())
	}
	cadenaDebito=humanize.Comma(numero)
	//+"."+decimalDebito
	return cadenaDebito

}
func Borrarinventario(CodigoDocumento string, TipoDocumento string ){

	delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo=$2")
	if err != nil {
		panic(err.Error())

	}
	delForm2.Exec(CodigoDocumento, TipoDocumento)
}

func Numeroventa(Codigoresolucion string) string{
	db := dbConn()
	miResolucion := Resolucionventa{}
	err := db.Get(&miResolucion, "SELECT * FROM resolucionventa where codigo=$1",Codigoresolucion)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Resolucion generador " + Codigoresolucion)
	log.Println("Resolucion numero actual " + miResolucion.NumeroActual)

	if miResolucion.NumeroActual == "0" {

		log.Println("Resolucion cero " )
		insForm, err := db.Prepare("UPDATE resolucionventa set	numeroactual=$2  " + " WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Codigoresolucion, miResolucion.NumeroInicial)
		return miResolucion.NumeroInicial
	} else {

		log.Println("numero nuevo " )
		var numeroactual int64
		var numerocadena string
	     numeroactual, err := strconv.ParseInt(miResolucion.NumeroActual, 10, 64)
		if err == nil {
		}
		numerocadena=strconv.FormatInt(numeroactual+1,10)
		insForm, err := db.Prepare("UPDATE resolucionventa set	numeroactual=$2  " + " WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Codigoresolucion, numerocadena)
		return numerocadena
	}
}
type ultimoDocumento struct {
	Ultimo string
}
func NumeroDocumento(Codigodocumento string) string{
	db := dbConn()
	miDocumento := ultimoDocumento{}
	err := db.Get(&miDocumento, "select max(cast(numero as integer))+1 as ultimo from comprobantedetalle where documento=$1",Codigodocumento)
	if err != nil {
		log.Fatalln(err)
	}
	return miDocumento.Ultimo
}

func Numerosoporte(Codigoresolucion string) string{
	db := dbConn()
	miResolucion := Resolucionsoporte{}
	err := db.Get(&miResolucion, "SELECT * FROM resolucionsoporte where codigo=$1",Codigoresolucion)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Resolucion generador " + Codigoresolucion)
	log.Println("Resolucion numero actual " + miResolucion.NumeroActual)

	if miResolucion.NumeroActual == "0" {

		log.Println("Resolucion cero " )
		insForm, err := db.Prepare("UPDATE resolucionsoporte set	numeroactual=$2  " + " WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Codigoresolucion, miResolucion.NumeroInicial)
		return miResolucion.NumeroInicial
	} else {

		log.Println("numero nuevo " )
		var numeroactual int64
		var numerocadena string
		numeroactual, err := strconv.ParseInt(miResolucion.NumeroActual, 10, 64)
		if err == nil {
		}
		numerocadena=strconv.FormatInt(numeroactual+1,10)
		insForm, err := db.Prepare("UPDATE resolucionsoporte set	numeroactual=$2  " + " WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Codigoresolucion, numerocadena)
		return numerocadena
	}
}
