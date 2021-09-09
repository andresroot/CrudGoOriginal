package main

// INICIA VENTASERVICIO IMPORTAR PAQUETES
import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "github.com/bitly/go-simplejson"
	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// INICIA VENTASERVICIO ESTRUCTURA JSON
type ventaservicioJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA VENTASERVICIO ESTRUCTURA
type ventaservicioLista struct {
	Codigo        string
	Fecha         time.Time
	Total         string
	Tercero       string
	TerceroNombre string
	CentroNombre  string
	VendedorNombre string
}

// INICIA VENTASERVICIO ESTRUCTURA
type ventaservicio struct {
	Codigo                    string
	Fecha                     time.Time
	Vence                     time.Time
	Hora                      string
	Descuento                 string
	Subtotaldescuento19       string
	Subtotaldescuento5        string
	Subtotaldescuento0        string
	Subtotal                  string
	Subtotal19                string
	Subtotal5                 string
	Subtotal0                 string
	Subtotaliva19             string
	Subtotaliva5              string
	Subtotaliva0              string
	Subtotalbase19            string
	Subtotalbase5             string
	Subtotalbase0             string
	TotalIva                  string
	Total                     string
	Neto                      string
	Items                     string
	Formadepago               string
	Mediodepago               string
	Resolucionventa           string
	Tercero                   string
	Vendedor                  string
	Accion                    string
	Detalle                   []ventaserviciodetalle `json:"Detalle"`
	DetalleEditar			  []ventaserviciodetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Cotizacionservicio		  string
	Tipo                      string
	Ret2201					  string
	Centro					  string
}

// INICIA VENTADETALLE ESTRUCTURA
type ventaserviciodetalle struct {
	Id                string
	Codigo            string
	Fila              string
	Cantidad          string
	Precio            string
	Descuento         string
	Montodescuento    string
	Sigratis          string
	Subtotal          string
	Subtotaldescuento string
	Pagina            string
	Nombreservicio    string
	Unidadservicio    string
	Codigoservicio    string
	Ivaservicio       string
	Tipo              string
	Fecha             time.Time
}

// INICIA COMPRA DETALLE EDITARr
type ventaserviciodetalleeditar struct {
	Id                string
	Codigo            string
	Fila              string
	Cantidad          string
	Precio            string
	Descuento         string
	Montodescuento    string
	Sigratis          string
	Subtotal          string
	Subtotaldescuento string
	Pagina            string
	Nombreservicio    string
	Unidadservicio    string
	Codigoservicio    string
	Ivaservicio       string
	Tipo              string
	Fecha             time.Time
}

// INICIA VENTASERVICIO CONSULTA DETALLE
func VentaservicioConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "ventaserviciodetalle.Id as id ,"
	consulta += "ventaserviciodetalle.Codigo as codigo,"
	consulta += "ventaserviciodetalle.Fila as fila,"
	consulta += "ventaserviciodetalle.Cantidad as cantidad,"
	consulta += "ventaserviciodetalle.Precio as precio,"
	consulta += "ventaserviciodetalle.Descuento as descuento,"
	consulta += "ventaserviciodetalle.Montodescuento as montodescuento,"
	consulta += "ventaserviciodetalle.Sigratis as sigratis,"
	consulta += "ventaserviciodetalle.Subtotal as subtotal,"
	consulta += "ventaserviciodetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "ventaserviciodetalle.Pagina as pagina ,"
	consulta += "ventaserviciodetalle.Fecha as fecha,"
	consulta += "ventaserviciodetalle.Nombreservicio as Nombreservicio, "
	consulta += "ventaserviciodetalle.Unidadservicio as Unidadservicio, "
	consulta += "ventaserviciodetalle.Codigoservicio as Codigoservicio, "
	consulta += "ventaserviciodetalle.Ivaservicio as Ivaservicio "
	consulta += "from ventaserviciodetalle "
	consulta += " where ventaserviciodetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA VENTASERVICIO LISTA
func VentaservicioLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/ventaservicio/ventaservicioLista.html")
	log.Println("Error ventaservicio 0")
	var consulta string

	consulta = "  SELECT vendedor.nombre as VendedorNombre,centro.nombre as CentroNombre,total,ventaservicio.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM ventaservicio "
	consulta += " inner join tercero on tercero.codigo=ventaservicio.tercero "
	consulta += " inner join centro on centro.codigo=ventaservicio.centro "
	consulta += " inner join vendedor on vendedor.codigo=ventaservicio.vendedor "
	consulta += " ORDER BY ventaservicio.codigo ASC"

	db := dbConn()
	res := []ventaservicioLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error ventaservicio888")
	tmp.Execute(w, varmap)
}

// INICIA VENTASERVICIO NUEVO
func VentaservicioNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio ventaservicio editar" + Codigo)

	db := dbConn()
	v := ventaservicio{}
	tc := tercero{}
	det := []ventaserviciodetalleeditar{}
	if Codigo == "False" {

	} else {

	err := db.Get(&v, "SELECT * FROM ventaservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	err2 := db.Select(&det, VentaservicioConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	err1 := db.Get(&tc, "SELECT * FROM tercero where codigo=$1", v.Tercero)
	if err1 != nil {
		log.Fatalln(err1)
	}
	}

	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"ventaservicio":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"resolucionventa": ListaResolucionventa(),
		"centro" :ListaCentro(),
		"ventaserviciotipoiva":TraerParametrosInventario().Ventaserviciotipoiva,
		"autoret2201":TraerParametrosInventario().Ventaserviciocuentaporcentajeret2201,
		"codigo": Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	t, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/ventaservicio/ventaservicioNuevo.html", "vista/ventaservicio/ventaservicioScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error ventaservicio nuevo 3")
	t.Execute(w, parametros)
}

// INICIA INSERTAR COMPROBANTE DE VENTASERVICIO
func InsertaDetalleComprobanteVentaservicio(miFilaComprobante comprobantedetalle, miComprobante comprobante, miVentaservicio ventaservicio){
	db := dbConn()

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miVentaservicio.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	var q string
	q = "insert into comprobantedetalle ("
	q += "Fila,"
	q += "Cuenta,"
	q += "Tercero,"
	q += "Centro,"
	q += "Concepto,"
	q += "Factura,"
	q += "Debito,"
	q += "Credito,"
	q += "Documento,"
	q += "Numero,"
	q += "Fecha,"
	q += "Fechaconsignacion"
	q += " ) values"
	q += "("
	q += parametros(12)
	q += ")"

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	if len(miFilaComprobante.Debito)>0 {
		miFilaComprobante.Debito=	miFilaComprobante.Debito + ".00"
	}

	if len(miFilaComprobante.Credito)>0 {
		miFilaComprobante.Credito=	miFilaComprobante.Credito + ".00"
	}

	// TERMINA COMPROBANTE GRABAR INSERTAR
	_, err = insForm.Exec(
		miFilaComprobante.Fila,
		miFilaComprobante.Cuenta  ,
		miVentaservicio.Tercero,
		miVentaservicio.Centro,
		miTercero.Nombre,
		"",
		Flotantedatabase(miFilaComprobante.Debito) ,
		Flotantedatabase(miFilaComprobante.Credito) ,
		miComprobante.Documento,
		miComprobante.Numero,
		miComprobante.Fecha,
		miComprobante.Fechaconsignacion	)
	if err != nil {
		panic(err)
	}
}

// INICIA VENTASERVICIO INSERTAR AJAX
func VentaservicioAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempVentaservicio ventaservicio

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la ventaservicio
	err = json.Unmarshal(b, &tempVentaservicio)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	var Codigoactual string
	if tempVentaservicio.Accion == "Nuevo" {
		log.Println("Resolucion " + tempVentaservicio.Resolucionventa)
		Codigoactual=Numeroventa(tempVentaservicio.Resolucionventa)
		tempVentaservicio.Codigo=Codigoactual
	}else{
		Codigoactual=tempVentaservicio.Codigo
	}



	if tempVentaservicio.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from ventaserviciodetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempVentaservicio.Codigo)

		// borra detalle inventario
		delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo='Ventaservicio'")
		if err != nil {
			panic(err.Error())
		}
		delForm2.Exec(tempVentaservicio.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from ventaservicio WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempVentaservicio.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempVentaservicio.Detalle {
		var a = i
		var q string
		q = "insert into ventaserviciodetalle ("
		q += "Id,"
		q += "Codigo,"
		q += "Fila,"
		q += "Cantidad,"
		q += "Precio,"
		q += "Subtotal,"
		q += "Pagina,"
		q += "Nombreservicio,"
		q += "Unidadservicio,"
		q += "Codigoservicio,"
		q += "Ivaservicio,"
		q += "Descuento,"
		q += "Montodescuento,"
		q += "Sigratis,"
		q += "Subtotaldescuento,"
		q += "Tipo,"
		q += "Fecha"
		q += " ) values("
		q += parametros(17)
		q += ")"

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA VENTASERVICIO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Id,
			Codigoactual,
			x.Fila,
			x.Cantidad,
			x.Precio,
			x.Subtotal,
			x.Pagina,
			x.Nombreservicio,
			x.Unidadservicio,
			x.Codigoservicio,
			x.Ivaservicio,
			x.Descuento,
			x.Montodescuento,
			x.Sigratis,
			x.Subtotaldescuento,
			x.Tipo,
			x.Fecha)
		if err != nil {
			panic(err)
		}

		log.Println("Insertar Producto \n", x.Nombreservicio, a)
	}

	// INICIA INSERTAR VENTAS
	log.Println("Got %s age %s club %s\n", tempVentaservicio.Codigo, tempVentaservicio.Tercero, tempVentaservicio.Total)
	var q string
	q = "insert into ventaservicio ("
	q += "Codigo,"
	q += "Fecha,"
	q += "Vence,"
	q += "Hora,"
	q += "Descuento,"
	q += "Subtotaldescuento19,"
	q += "Subtotaldescuento5,"
	q += "Subtotaldescuento0,"
	q += "Subtotal,"
	q += "Subtotal19,"
	q += "Subtotal5,"
	q += "Subtotal0,"
	q += "Subtotaliva19,"
	q += "Subtotaliva5,"
	q += "Subtotaliva0,"
	q += "Subtotalbase19,"
	q += "Subtotalbase5,"
	q += "Subtotalbase0,"
	q += "TotalIva,"
	q += "Ret2201,"
	q += "Total,"
	q += "Neto,"
	q += "Items,"
	q += "Formadepago,"
	q += "Mediodepago,"
	q += "Resolucionventa,"
	q += "Tercero,"
	q += "Vendedor,"
	q += "Cotizacionservicio,"
	q += "Centro,"
	q += "Tipo"
	q += " ) values("
	q+=parametros(31)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempVentaservicio.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempVentaservicio.Codigo,
		tempVentaservicio.Fecha.Format(layout),
		tempVentaservicio.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempVentaservicio.Descuento,
		tempVentaservicio.Subtotaldescuento19,
		tempVentaservicio.Subtotaldescuento5,
		tempVentaservicio.Subtotaldescuento0,
		tempVentaservicio.Subtotal,
		tempVentaservicio.Subtotal19,
		tempVentaservicio.Subtotal5,
		tempVentaservicio.Subtotal0,
		tempVentaservicio.Subtotaliva19,
		tempVentaservicio.Subtotaliva5,
		tempVentaservicio.Subtotaliva0,
		tempVentaservicio.Subtotalbase19,
		tempVentaservicio.Subtotalbase5,
		tempVentaservicio.Subtotalbase0,
		tempVentaservicio.TotalIva,
		tempVentaservicio.Ret2201,
		tempVentaservicio.Total,
		tempVentaservicio.Neto,
		tempVentaservicio.Items,
		tempVentaservicio.Formadepago,
		tempVentaservicio.Mediodepago,
		tempVentaservicio.Resolucionventa,
		tempVentaservicio.Tercero,
		tempVentaservicio.Vendedor,
		tempVentaservicio.Cotizacionservicio,
		tempVentaservicio.Centro,
		tempVentaservicio.Tipo)

	if err != nil {
		panic(err)
	}

	// INSERTAR COMPROBANTE CONTABILIDAD
	var tempComprobante comprobante
	var tempComprobanteDetalle comprobantedetalle
	tempComprobante.Documento="7"
	tempComprobante.Numero=tempVentaservicio.Codigo
	tempComprobante.Fecha =tempVentaservicio.Fecha
	tempComprobante.Fechaconsignacion =tempVentaservicio.Fecha
	tempComprobante.Debito = tempVentaservicio.Neto + ".00"
	tempComprobante.Credito	= tempVentaservicio.Neto + ".00"
	tempComprobante.Periodo	= ""
	tempComprobante.Licencia = ""
	tempComprobante.Usuario	= ""
	tempComprobante.Estado	= ""

	// borra detalle anterior
	delForm, err := db.Prepare("DELETE from comprobantedetalle WHERE documento=$1 and numero=$2")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(tempComprobante.Documento, tempComprobante.Numero)

	// borra cabecera anterior

	delForm1, err := db.Prepare("DELETE from comprobante WHERE documento=$1 and numero=$2 ")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(tempComprobante.Documento, tempComprobante.Numero)

	// INSERTAR CABECERA COMPROBANTE

	log.Println("Got %s age %s club %s\n", tempComprobante.Documento, tempComprobante.Numero)

	var totalDebito float64
	var totalCredito float64
	var fila int
	fila=0
	totalDebito=0
	totalCredito=0

	// INSERTAR CUENTA DEBITO CLIENTE
	if (tempVentaservicio.Neto!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempVentaservicio.Neto)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaserviciocuentacliente
		tempComprobanteDetalle.Debito = tempVentaservicio.Neto
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO RET. 2201
	if (tempVentaservicio.Ret2201!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVentaservicio.Ret2201)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaserviciocuentaret2201
		tempComprobanteDetalle.Debito = tempVentaservicio.Ret2201
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}


	// INSERTAR CUENTA DEBITO DESCUENTO
	//if (tempVentaservicio.Descuento!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalDebito+=Flotante(tempVentaservicio.Descuento)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaserviciocuentadescuento
	//	tempComprobanteDetalle.Debito = tempVentaservicio.Descuento
	//	tempComprobanteDetalle.Credito = ""
	//	InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
	//	log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	//}

	// INSERTAR CUENTA DEBITO RET. 2201
	if (tempVentaservicio.Ret2201!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempVentaservicio.Ret2201)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaserviciocontracuentaret2201
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempVentaservicio.Ret2201
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO IVA 19%
	if (tempVentaservicio.Subtotaliva19!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVentaservicio.Subtotaliva19)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaservicioiva19
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempVentaservicio.Subtotaliva19
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO IVA 5%
	if (tempVentaservicio.Subtotaliva5!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVentaservicio.Subtotaliva5)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaservicioiva5
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempVentaservicio.Subtotaliva5
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO VENTASERVICIO IVA 19%
	if (tempVentaservicio.Subtotalbase19!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVentaservicio.Subtotalbase19)+Flotante(tempVentaservicio.Subtotaldescuento19)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaserviciocuenta19
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = humanize.Comma(int64(Flotante(tempVentaservicio.Subtotalbase19)))
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO VENTASERVICIO IVA 5%
	if (tempVentaservicio.Subtotalbase5!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVentaservicio.Subtotalbase5)+Flotante(tempVentaservicio.Subtotaldescuento5)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaserviciocuenta5
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = humanize.Comma(int64(Flotante(tempVentaservicio.Subtotalbase5)))
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO VENTASERVICIO IVA 0%
	if (tempVentaservicio.Subtotalbase0!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVentaservicio.Subtotalbase0)+Flotante(tempVentaservicio.Subtotaldescuento0)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaserviciocuenta0
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = humanize.Comma(int64(Flotante(tempVentaservicio.Subtotalbase0)))
		InsertaDetalleComprobanteVentaservicio(tempComprobanteDetalle,tempComprobante,tempVentaservicio)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	var cadenaDebito=FormatoFlotante(totalDebito)
	var cadenaCredito=FormatoFlotante(totalCredito)

	q = "insert into comprobante ("
	q += "Documento,"
	q += "Numero,"
	q += "Fecha,"
	q += "Fechaconsignacion,"
	q += "Debito,"
	q += "Credito,"
	q += "Periodo,"
	q += "Licencia,"
	q += "Usuario,"
	q += "Estado)"
	q += " values("
	q += parametros(10)
	q += ")"

	log.Println("Cadena SQL " + q)
	insForm, err = db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	_, err = insForm.Exec(
		tempComprobante.Documento,
		tempComprobante.Numero,
		tempComprobante.Fecha.Format(layout),
		tempComprobante.Fechaconsignacion.Format(layout),
		cadenaDebito,
		cadenaCredito,
		tempComprobante.Periodo,
		tempComprobante.Licencia,
		tempComprobante.Usuario,
		tempComprobante.Estado)
	if err != nil {
		panic(err)
	}

	var resultado bool
	resultado = true

	js, err := json.Marshal(SomeStruct{resultado})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

// INICIA VENTASERVICIO EXISTE
func VentaservicioExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM ventaservicio  WHERE codigo=$1", Codigo)
	err := row.Scan(&total)
	if err != nil {
		log.Fatal(err)
	}
	var resultado bool
	if total > 0 {
		resultado = true
	} else {
		resultado = false
	}
	js, err := json.Marshal(SomeStruct{resultado})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// INICIA VENTASERVICIO EDITAR
func VentaservicioEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio ventaservicio editar" + Codigo)

	db := dbConn()

	// traer ventaservicio
	v := ventaservicio{}
	err := db.Get(&v, "SELECT * FROM ventaservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []ventaserviciodetalleeditar{}

	err2 := db.Select(&det, VentaservicioConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	t := tercero{}
	err1 := db.Get(&t, "SELECT * FROM tercero where codigo=$1", v.Tercero)
	if err1 != nil {
		log.Fatalln(err1)
	}
	log.Println("codigo nombre99" + t.Codigo + t.Nombre)

	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"ventaservicio":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"resolucionventa": ListaResolucionventa(),
		"centro" :ListaCentro(),
		"ventaserviciotipoiva":TraerParametrosInventario().Ventaserviciotipoiva,
		"autoret2201":TraerParametrosInventario().Ventaserviciocuentaporcentajeret2201,
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/ventaservicio/ventaservicioEditar.html", "vista/ventaservicio/ventaservicioScript.html")
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error ventaservicio nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA VENTASERVICIO BORRAR
func VentaservicioBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio ventaservicio editar" + Codigo)

	db := dbConn()

	// traer ventaservicio
	v := ventaservicio{}
	err := db.Get(&v, "SELECT * FROM ventaservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []ventaserviciodetalleeditar{}
	err2 := db.Select(&det, VentaservicioConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	t := tercero{}
	err1 := db.Get(&t, "SELECT * FROM tercero where codigo=$1", v.Tercero)
	if err1 != nil {
		log.Fatalln(err1)
	}
	log.Println("codigo nombre99" + t.Codigo + t.Nombre)

	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"ventaservicio":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"resolucionventa":  ListaResolucionventa(),
		"ventaserviciotipoiva":TraerParametrosInventario().Ventaserviciotipoiva,
		"autoret2201":TraerParametrosInventario().Ventaserviciocuentaporcentajeret2201,
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/ventaservicio/ventaservicioBorrar.html", "vista/ventaservicio/ventaservicioScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error ventaservicio nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA VENTASERVICIO ELIMINAR
func VentaservicioEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar ventaservicio
	delForm, err := db.Prepare("DELETE from ventaservicio WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from ventaserviciodetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borar detalle Inventario
	Borrarinventario(codigo,"Ventaservicio")

	// borra detalle anterior
	delForm, err = db.Prepare("DELETE from comprobantedetalle WHERE documento=$1 and numero=$2")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec("7", codigo)

	// borra cabecera anterior

	delForm1, err = db.Prepare("DELETE from comprobante WHERE documento=$1 and numero=$2 ")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec("7", codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/VentaservicioLista", 301)
}

// TRAER COTIZACION
func DatosCotizacionservicio(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacionservicio editar" + Codigo)
	db := dbConn()
	var res []cotizacionservicio

	// traer COTIZACION
	v := cotizacionservicio{}
	err := db.Get(&v, "SELECT * FROM cotizacionservicio where codigo=$1", Codigo)
	var valida bool
	valida=true

	switch err {
	case nil:
		log.Printf("cotizacionservicio existe: %+v\n", v)
	case sql.ErrNoRows:
		log.Println("cotizacionservicio NO Existe")
		valida=false
	default:
		log.Printf("cotizacionservicio error: %s\n", err)
	}
	det := []cotizacionserviciodetalleeditar{}
	t := tercero{}

	// trae datos si existe cotizacionservicio
	if valida==true {
		err2 := db.Select(&det, CotizacionservicioConsultaDetalle(), Codigo)
		if err2 != nil {
			fmt.Println(err2)
		}
		// traer tercero
		err1 := db.Get(&t, "SELECT * FROM tercero where codigo=$1", v.Tercero)
		if err1 != nil {
			log.Fatalln(err1)
		}
		v.TerceroDetalle=t;
		v.DetalleEditar=det;
		res = append(res, v)
	}

	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
	data, err := json.Marshal(res)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

// INICIA PDF
func VentaservicioPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer ventaservicio
	miVentaservicio := ventaservicio{}
	err := db.Get(&miVentaservicio, "SELECT * FROM ventaservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []ventaserviciodetalleeditar{}
	err2 := db.Select(&miDetalle, VentaservicioConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miVentaservicio.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer Vendedor
	miVendedor := vendedor{}
	err4 := db.Get(&miVendedor, "SELECT * FROM vendedor where codigo=$1", miVentaservicio.Vendedor)
	if err4 != nil {
		log.Fatalln(err4)
	}

	var e empresa = ListaEmpresa()
	var c ciudad = TraerCiudad(e.Ciudad)
	var re Resolucionventa = TraerResolucionventa(miVentaservicio.Resolucionventa)
	var buf bytes.Buffer
	var err1 error

	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)
	ene := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetHeaderFunc(func() {
		// LOGO
		pdf.SetY(5)
		pdf.Image(imageFile("logo.png"), 20, 30, 40, 0, false,
			"", 0, "")

		// EMPRESA
		pdf.SetFont("Arial", "", 9)
		pdf.SetY(20)
		pdf.CellFormat(190, 10, Mayuscula(e.Nombre), "0", 0,
			"C", false, 0, "")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, "Nit No. "+Coma(e.Codigo)+" - "+e.Dv, "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, e.Iva+" - "+e.ReteIva, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, "Actividad Ica - "+e.ActividadIca, "0",
			0, "C", false, 0, "")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, e.Direccion, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+" - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")

		// RESOLUCION
		pdf.SetY(20)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "Resolucion No. "+re.Numero, "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "Del Numero "+re.Prefijo+" "+re.NumeroInicial+" al "+re.Prefijo+" "+Coma(re.NumeroFinal), "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "Vigencia del "+re.FechaInicial.Format("02/01/2006")+" al "+re.FechaFinal.Format("02/01/2006"), "0", 0, "C",
			false, 0, "")

		// FACTURA NUMERO
		pdf.SetFont("Arial", "", 10)
		pdf.Ln(5)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "FACTURA ELECTRONICA", "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, " DE VENTA No. "+re.Prefijo+" "+Codigo, "0", 0, "C",
			false, 0, "")
		pdf.Ln(6)

	})

	pdf.SetFooterFunc(func() {
		pdf.SetFont("Arial", "", 8)
		pdf.SetY(256)
		pdf.SetX(20)
		pdf.CellFormat(80, 10, "Andres Eduardo Ojeda Medina Nit." +
			" 80.853.536-7 SADCONF Derechos de Autor 13-16-230 de 30-06-2006  www.Sadconf.com.co", "",
			0, "L", false, 0, "")
		pdf.SetX(130)
		pdf.CellFormat(78, 10, fmt.Sprintf(" %d de {nb}", pdf.PageNo()), "",
			0, "R", false, 0, "")

	})

	// inicia suma de paginas
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.Ln(1)
	VentaservicioCabecera(pdf,miTercero,miVentaservicio, miVendedor)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			VentaservicioFilaDetalle(pdf,miFila,a)
		}
		VentaservicioPieDePagina(pdf,miTercero,miVentaservicio)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					VentaservicioFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			VentaservicioCabecera(pdf,miTercero,miVentaservicio,miVendedor)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					VentaservicioFilaDetalle(pdf,miFila,a)
				}
			}

			VentaservicioPieDePagina(pdf,miTercero,miVentaservicio)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					VentaservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			VentaservicioCabecera(pdf,miTercero,miVentaservicio,miVendedor)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					VentaservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			VentaservicioCabecera(pdf,miTercero,miVentaservicio,miVendedor)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					VentaservicioFilaDetalle(pdf,miFila,a)
				}
			}

			VentaservicioPieDePagina(pdf,miTercero,miVentaservicio)
		}
	}

	// genera pdf
	err1 = pdf.Output(&buf)
	if err1 != nil {
		panic(err1.Error())
	}
	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	w.Write(buf.Bytes())
}

func VentaservicioCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miVentaservicio ventaservicio, miVendedor vendedor ){
	pdf.SetFont("Arial", "", 10)
	pdf.Ln(4)
	pdf.SetX(20)
	pdf.CellFormat(90, 5, "DATOS DEL ADQUIRIENTE", "1", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(94, 5, "LUGAR DE ENTREGA O SERVICIO", "1", 0,
		"L", false, 0, "")
	pdf.Ln(8)

	// DETALLE ADQUIRIENTE
	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Nit. No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(40)
	pdf.CellFormat(40, 4, Coma(miTercero.Codigo)+" - "+miTercero.Dv, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Nit. No. ", "", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, Coma(miTercero.Codigo)+" - "+miTercero.Dv, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Nombre", "", 0,
		"L", false, 0, "")
	pdf.SetX(40)
	pdf.CellFormat(40, 4, miTercero.Nombre, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Nombre", "", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, miTercero.Nombre, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Direccion", "", 0,
		"L", false, 0, "")
	pdf.SetX(40)
	pdf.CellFormat(40, 4, miTercero.Direccion, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Direccion", "", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, miTercero.Direccion, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Forma de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miVentaservicio.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miVentaservicio.Fecha.Format("02/01/2006")+" "+Titulo(miVentaservicio.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miVentaservicio.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miVentaservicio.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Cotizacion Servicio No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, miVentaservicio.Cotizacionservicio, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Vendedor", "", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, miVendedor.Nombre, "", 0,
		"L", false, 0, "")

	// CABECERA PRODUCTO
	pdf.SetFont("Arial", "", 10)
	pdf.Ln(6)
	pdf.SetX(20)
	pdf.CellFormat(184, 5, "ITEM", "1", 0,
		"L", false, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(190, 5, "CODIGO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(190, 5, "PRODUCTO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(105)
	pdf.CellFormat(190, 5, "UNIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(121)
	pdf.CellFormat(190, 5, "IVA", "0", 0,
		"L", false, 0, "")
	pdf.SetX(129)
	pdf.CellFormat(190, 5, "DES.", "0", 0,
		"L", false, 0, "")
	pdf.SetX(140)
	pdf.CellFormat(190, 5, "CANTIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(160)
	pdf.CellFormat(190, 5, "P. UNITARIO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(190)
	pdf.CellFormat(190, 5, "TOTAL", "0", 0,
		"L", false, 0, "")
	pdf.Ln(8)
}
func VentaservicioFilaDetalle(pdf *gofpdf.Fpdf,miFila ventaserviciodetalleeditar, a int ){
	pdf.SetFont("Arial", "", 9)
	pdf.SetX(20)
	var yinicial  float64
	yinicial=pdf.GetY()
	pdf.CellFormat(40, 4, strconv.Itoa(a), "", 0,
		"L", false, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 4, Subcadena(miFila.Codigoservicio,0,6), "", 0,
		"L", false, 0, "")
	var y float64
	y=pdf.GetY()
	pdf.SetX(42)
	pdf.MultiCell(68,4, Mayuscula(miFila.Nombreservicio), "","L", false)
	var yfinal float64
	yfinal=pdf.GetY()
	pdf.SetY(y)
	pdf.SetX(81)
	pdf.CellFormat(40, 4, Subcadena(miFila.Unidadservicio, 0,6), "", 0,
		"R", false, 0, "")
	pdf.SetX(89)
	pdf.CellFormat(40, 4, miFila.Ivaservicio, "", 0,
		"R", false, 0, "")
	pdf.SetX(98)
	pdf.CellFormat(40, 4, miFila.Descuento, "", 0,
		"R", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(40, 4, miFila.Cantidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(143)
	pdf.CellFormat(40, 4, miFila.Precio, "", 0,
		"R", false, 0, "")
	pdf.SetX(165)
	pdf.CellFormat(40, 4, miFila.Subtotal, "", 0,
		"R", false, 0, "")
	pdf.Ln(yfinal-yinicial+3)

}

func VentaservicioPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miVentaservicio ventaservicio){

	Totalletras,err := IntLetra(Cadenaentero(miVentaservicio.Total))
	if err!= nil{
		fmt.Println(err)
	}

	pdf.SetFont("Arial", "", 8)
	pdf.SetY(222)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "SON: " +Mayuscula(Totalletras)+" PESOS MDA. CTE.", "0", 0,
		"L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetY(229)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "Esta factura es un titulo valor para su emisor o poseedor en caso de", "0", 0,
		"C", false, 0, "")
	pdf.Ln(3)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "endoso, presta merito ejecutivo y cumple con los requisitos del art.", "0", 0,
		"C", false, 0, "")
	pdf.Ln(3)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "617 del E. T. y art. 773 y 774 del Codigo de Comercio", "0", 0,
		"C", false, 0, "")
	pdf.Ln(10)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "__________________________________________", "0", 0,
		"C", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "A C E P T A D A ", "0", 0, "C",
		false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetY(229)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "SUBTOTAL", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "DESCUENTO", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "IVA 19%", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "IVA 5%", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "TOTAL", "0", 0, "L",
		false, 0, "")

	pdf.SetY(229)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVentaservicio.Subtotal, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVentaservicio.Descuento, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVentaservicio.Subtotaliva19, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVentaservicio.Subtotaliva5, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVentaservicio.Total, "0", 0, "R",
		false, 0, "")
	pdf.Image(imageFile("QR.jpg"), 20, 229, 25, 0, false,
		"", 0, "")
	pdf.SetFont("Arial", "", 8)
	pdf.SetY(249)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "Cufexxxxxxxxxx", "", 0,
		"L", false, 0, "")

}
