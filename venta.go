package main

// INICIA VENTA IMPORTAR PAQUETES
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
	"math"
	"net/http"
	"strconv"
	"time"
)

// INICIA VENTA ESTRUCTURA JSON
type ventaJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA VENTA ESTRUCTURA
type ventaLista struct {
	Codigo        string
	Fecha         time.Time
	Total         string
	Tercero       string
	TerceroNombre string
	CentroNombre  string
	VendedorNombre string
}

// INICIA VENTA ESTRUCTURA
type venta struct {
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
	Resolucion                string
	Tercero                   string
	Vendedor                  string
	Accion                    string
	Detalle                   []ventadetalle `json:"Detalle"`
	DetalleEditar			  []ventadetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Cotizacion				  string
	Tipo                      string
	Ret2201					  string
	Centro					  string
}

// INICIA VENTADETALLE ESTRUCTURA
type ventadetalle struct {
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
	Bodega            string
	Producto          string
	Tipo              string
	Fecha             time.Time
}

// INICIA COMPRA DETALLE EDITARr
type ventadetalleeditar struct {
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
	Bodega            string
	Producto          string
	ProductoNombre    string
	ProductoIva       string
	ProductoUnidad    string
	Tipo              string
	Fecha             time.Time
}

// INICIA COMPRA CONSULTA DETALLE
func VentaConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "ventadetalle.Id as id ,"
	consulta += "ventadetalle.Codigo as codigo,"
	consulta += "ventadetalle.Fila as fila,"
	consulta += "ventadetalle.Cantidad as cantidad,"
	consulta += "ventadetalle.Precio as precio,"
	consulta += "ventadetalle.Descuento as descuento,"
	consulta += "ventadetalle.Montodescuento as montodescuento,"
	consulta += "ventadetalle.Sigratis as sigratis,"
	consulta += "ventadetalle.Subtotal as subtotal,"
	consulta += "ventadetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "ventadetalle.Pagina as pagina ,"
	consulta += "ventadetalle.Bodega as bodega,"
	consulta += "ventadetalle.Producto as producto,"
	consulta += "ventadetalle.Fecha as fecha,"
	consulta += "producto.nombre as ProductoNombre, "
	consulta += "producto.iva as ProductoIva, "
	consulta += "producto.unidad as ProductoUnidad "
	consulta += "from ventadetalle "
	consulta += "inner join producto on producto.codigo=ventadetalle.producto "
	consulta += " where ventadetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA VENTA LISTA
func VentaLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/venta/ventaLista.html")
	log.Println("Error venta 0")
	var consulta string

	consulta = "  SELECT vendedor.nombre as VendedorNombre,centro.nombre as CentroNombre,total,venta.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM venta "
	consulta += " inner join tercero on tercero.codigo=venta.tercero "
	consulta += " inner join centro on centro.codigo=venta.centro "
	consulta += " inner join vendedor on vendedor.codigo=venta.vendedor "
	consulta += " ORDER BY venta.codigo ASC"

	db := dbConn()
	res := []ventaLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error venta888")
	tmp.Execute(w, varmap)
}

// INICIA VENTA NUEVO
func VentaNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio venta editar" + Codigo)

	db := dbConn()
	v := venta{}
	tc := tercero{}
	det := []ventadetalleeditar{}
	if Codigo == "False" {

	} else {

	err := db.Get(&v, "SELECT * FROM venta where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	err2 := db.Select(&det, VentaConsultaDetalle(), Codigo)
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
		"venta":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"resolucionventa": ListaResolucionventa(),
		"centro" :ListaCentro(),
		"ventatipoiva":TraerParametrosInventario().Ventatipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
		"codigo": Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/venta/ventaNuevo.html", "vista/venta/ventaScript.html")
	fmt.Printf("%v, %v", tc, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error venta nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA INSERTAR COMPROBANTE DE VENTA
func InsertaDetalleComprobanteVenta(miFilaComprobante comprobantedetalle, miComprobante comprobante, miVenta venta){
	db := dbConn()

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miVenta.Tercero)
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
	q += " ) values("
	q += parametros(12)
	q += " ) "

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
	miVenta.Tercero,
	miVenta.Centro,
	miTercero.Nombre,
	"",
	Flotantedatabase(miFilaComprobante.Debito),
	Flotantedatabase(miFilaComprobante.Credito) ,
	miComprobante.Documento,
	miComprobante.Numero,
	miComprobante.Fecha,
	miComprobante.Fechaconsignacion	)
	if err != nil {
	panic(err)
	}
}

// INICIA VENTA INSERTAR AJAX
func VentaAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempVenta venta

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la venta
	err = json.Unmarshal(b, &tempVenta)
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
	if tempVenta.Accion == "Nuevo" {
		log.Println("Resolucion " + tempVenta.Resolucion)
		Codigoactual=Numeroventa(tempVenta.Resolucion)
		tempVenta.Codigo=Codigoactual
	}else{
		Codigoactual=tempVenta.Codigo
	}



	if tempVenta.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from ventadetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempVenta.Codigo)

		// borra detalle inventario
		delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo='Venta'")
		if err != nil {
			panic(err.Error())
		}
		delForm2.Exec(tempVenta.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from venta WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempVenta.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempVenta.Detalle {
		var a = i
		var q string
		q = "insert into ventadetalle ("
		q += "Id,"
		q += "Codigo,"
		q += "Fila,"
		q += "Cantidad,"
		q += "Precio,"
		q += "Subtotal,"
		q += "Pagina,"
		q += "Bodega,"
		q += "Producto,"
		q += "Descuento,"
		q += "Montodescuento,"
		q += "Sigratis,"
		q += "Subtotaldescuento,"
		q += "Tipo,"
		q += "Fecha"
		q += " ) values("
		q += parametros(15)
		q += " ) "

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA VENTA GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Id,
			Codigoactual,
			x.Fila,
			x.Cantidad,
			x.Precio,
			x.Subtotal,
			x.Pagina,
			x.Bodega,
			x.Producto,
			x.Descuento,
			x.Montodescuento,
			x.Sigratis,
			x.Subtotaldescuento,
			x.Tipo,
			x.Fecha)
		if err != nil {
			panic(err)
		}

		log.Println("Insertar Producto \n", x.Producto, a)
	}

	// INSERTA DETALLE INVENTARIO
	for i, x := range tempVenta.Detalle {
		var a = i
		var q string
		q = "insert into inventario ("
		q += "Fecha,"
		q += "Tipo,"
		q += "Codigo,"
		q += "Bodega,"
		q += "Producto,"
		q += "Cantidad,"
		q += "Precio,"
		q += "Operacion"
		q += " ) values("
		q += parametros(8)
		q += " ) "

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA VENTA GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Fecha,
			x.Tipo,
			Codigoactual,
			x.Bodega,
			x.Producto,
			x.Cantidad,
			x.Precio,
			operacionVenta)
		if err != nil {
			panic(err)
		}
		log.Println("Insertar Producto \n", x.Producto, a)
	}

	// INICIA INSERTAR VENTAS
	log.Println("Got %s age %s club %s\n", tempVenta.Codigo, tempVenta.Tercero, tempVenta.Total)
	var q string
	q = "insert into venta ("
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
	q += "Resolucion,"
	q += "Tercero,"
	q += "Vendedor,"
	q += "Cotizacion,"
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
	log.Println("Hora", tempVenta.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempVenta.Codigo,
		tempVenta.Fecha.Format(layout),
		tempVenta.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempVenta.Descuento,
		tempVenta.Subtotaldescuento19,
		tempVenta.Subtotaldescuento5,
		tempVenta.Subtotaldescuento0,
		tempVenta.Subtotal,
		tempVenta.Subtotal19,
		tempVenta.Subtotal5,
		tempVenta.Subtotal0,
		tempVenta.Subtotaliva19,
		tempVenta.Subtotaliva5,
		tempVenta.Subtotaliva0,
		tempVenta.Subtotalbase19,
		tempVenta.Subtotalbase5,
		tempVenta.Subtotalbase0,
		tempVenta.TotalIva,
		tempVenta.Ret2201,
		tempVenta.Total,
		tempVenta.Neto,
		tempVenta.Items,
		tempVenta.Formadepago,
		tempVenta.Mediodepago,
		tempVenta.Resolucion,
		tempVenta.Tercero,
		tempVenta.Vendedor,
		tempVenta.Cotizacion,
		tempVenta.Centro,
	    tempVenta.Tipo)

	if err != nil {
		panic(err)
	}

	// INSERTAR COMPROBANTE CONTABILIDAD
	var tempComprobante comprobante
	var tempComprobanteDetalle comprobantedetalle
	tempComprobante.Documento="7"
	tempComprobante.Numero=tempVenta.Codigo
	tempComprobante.Fecha =tempVenta.Fecha
	tempComprobante.Fechaconsignacion =tempVenta.Fecha
	tempComprobante.Debito = tempVenta.Neto + ".00"
	tempComprobante.Credito	= tempVenta.Neto + ".00"
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
	if (tempVenta.Neto!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempVenta.Neto)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventacuentacliente
		tempComprobanteDetalle.Debito = tempVenta.Neto
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA DEBITO RET. 2201
	if (tempVenta.Ret2201!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempVenta.Ret2201)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventacuentaret2201
		tempComprobanteDetalle.Debito = tempVenta.Ret2201
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA DEBITO DESCUENTO
	//if (tempVenta.Descuento!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalDebito+=Flotante(tempVenta.Descuento)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventacuentadescuento
	//	tempComprobanteDetalle.Debito = tempVenta.Descuento
	//	tempComprobanteDetalle.Credito = ""
	//	InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
	//	log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	//}

	// INSERTAR CUENTA CREDITO RET. 2201
	if (tempVenta.Ret2201!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVenta.Ret2201)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventacontracuentaret2201
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempVenta.Ret2201
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO IVA 19%
	if (tempVenta.Subtotaliva19!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVenta.Subtotaliva19)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaiva19
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempVenta.Subtotaliva19
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO IVA 5%
	if (tempVenta.Subtotaliva5!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVenta.Subtotaliva5)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventaiva5
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempVenta.Subtotaliva5
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO VENTA IVA 19%
	if (tempVenta.Subtotalbase19!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVenta.Subtotalbase19)+Flotante(tempVenta.Subtotaldescuento19)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventacuenta19
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = humanize.Comma(int64(Flotante(tempVenta.Subtotalbase19)))
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO VENTA IVA 5%
	if (tempVenta.Subtotalbase5!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVenta.Subtotalbase5)+Flotante(tempVenta.Subtotaldescuento5)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventacuenta5
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = humanize.Comma(int64(Flotante(tempVenta.Subtotalbase5)))
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO VENTA IVA 0%
	if (tempVenta.Subtotalbase0!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempVenta.Subtotalbase0)+Flotante(tempVenta.Subtotaldescuento0)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Ventacuenta0
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = humanize.Comma(int64(Flotante(tempVenta.Subtotalbase0)))
		InsertaDetalleComprobanteVenta(tempComprobanteDetalle,tempComprobante,tempVenta)
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
	q += "Estado"
	q += " ) values( "
	q += parametros(10)
	q += " ) "

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

// INICIA VENTA EXISTE
func VentaExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM venta  WHERE codigo=$1", Codigo)
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

// INICIA VENTA EDITAR
func VentaEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio venta editar" + Codigo)

	db := dbConn()

	// traer venta
	v := venta{}
	err := db.Get(&v, "SELECT * FROM venta where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []ventadetalleeditar{}

	err2 := db.Select(&det, VentaConsultaDetalle(), Codigo)
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
		"venta":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"resolucionventa": ListaResolucionventa(),
		"centro" :ListaCentro(),
		"ventatipoiva":TraerParametrosInventario().Ventatipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/venta/ventaEditar.html", "vista/venta/ventaScript.html")
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error venta nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA VENTA BORRAR
func VentaBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio venta editar" + Codigo)

	db := dbConn()

	// traer venta
	v := venta{}
	err := db.Get(&v, "SELECT * FROM venta where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []ventadetalleeditar{}
	err2 := db.Select(&det, VentaConsultaDetalle(), Codigo)
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
		"venta":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"resolucionventa":  ListaResolucionventa(),
		"ventatipoiva":TraerParametrosInventario().Ventatipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/venta/ventaBorrar.html", "vista/venta/ventaScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error venta nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA VENTA ELIMINAR
func VentaEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar venta
	delForm, err := db.Prepare("DELETE from venta WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from ventadetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borar detalle Inventario
	Borrarinventario(codigo,"Venta")

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
	http.Redirect(w, r, "/VentaLista", 301)
}

// TRAER COTIZACION
func DatosCotizacion(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacion editar" + Codigo)
	db := dbConn()
	var res []cotizacion

	// traer COTIZACION
	v := cotizacion{}
	err := db.Get(&v, "SELECT * FROM cotizacion where codigo=$1", Codigo)
	var valida bool
	valida=true

	switch err {
	case nil:
		log.Printf("cotizacion existe: %+v\n", v)
	case sql.ErrNoRows:
		log.Println("cotizacion NO Existe")
		valida=false
	default:
		log.Printf("cotizacion error: %s\n", err)
	}
	det := []cotizaciondetalleeditar{}
	t := tercero{}

	// trae datos si existe cotizacion
	if valida==true {
		err2 := db.Select(&det, CotizacionConsultaDetalle(), Codigo)
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


// INICIA VENTA PDF
func VentaPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer venta
	miVenta := venta{}
	err := db.Get(&miVenta, "SELECT * FROM venta where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []ventadetalleeditar{}
	err2 := db.Select(&miDetalle, VentaConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miVenta.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer Vendedor
	miVendedor := vendedor{}
	err4 := db.Get(&miVendedor, "SELECT * FROM vendedor where codigo=$1", miVenta.Vendedor)
	if err4 != nil {
		log.Fatalln(err4)
	}

	var buf bytes.Buffer
	var err1 error

	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)

	VentaHeader(pdf,miVenta)
	VentaFooter(pdf,miVenta)

	// inicia suma de paginas
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.Ln(1)
	VentaCabecera(pdf,miTercero,miVenta, miVendedor)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			VentaFilaDetalle(pdf,miFila,a)
		}
		VentaPieDePagina(pdf,miTercero,miVenta)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					VentaFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			VentaCabecera(pdf,miTercero,miVenta,miVendedor)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					VentaFilaDetalle(pdf,miFila,a)
				}
			}

			VentaPieDePagina(pdf,miTercero,miVenta)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					VentaFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			VentaCabecera(pdf,miTercero,miVenta,miVendedor)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					VentaFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			VentaCabecera(pdf,miTercero,miVenta,miVendedor)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					VentaFilaDetalle(pdf,miFila,a)
				}
			}

			VentaPieDePagina(pdf,miTercero,miVenta)
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
// TERMINA VENTA PDF

// INICIA VENTA HEADER
func VentaHeader(pdf *gofpdf.Fpdf, miVenta venta){
// ENCABEZADO
	var e empresa = ListaEmpresa()
	var c ciudad = TraerCiudad(e.Ciudad)
	var re Resolucionventa = TraerResolucionventa(miVenta.Resolucion)
	ene := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetHeaderFunc(func() {
		// LOGO
		pdf.SetY(5)
		pdf.Image(imageFile("logo.png"), 20, 30, 40, 0, false,
		"", 0, "")

		// EMPRESA
		pdf.SetFont("Arial", "", 9)

		// CUADRO EMPRESA
		pdf.SetY(20)
		pdf.SetX(20)
		pdf.SetDrawColor(119,134,153)
		pdf.CellFormat(184, 69, "", "1", 0, "C",
		false, 0, "")

		// DATOS VENTA
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
		pdf.SetFont("Arial", "", 12)
		pdf.Ln(7)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "FACTURA ELECTRONICA", "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, " DE VENTA No. "+re.Prefijo+" "+miVenta.Codigo, "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
	})
}
// TERMINA VENTA HEADER

// INICIA VENTA CABECERA
func VentaCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miVenta venta, miVendedor vendedor ){
	pdf.SetFont("Arial", "", 10)

	// RELLENO TITULO
	pdf.SetY(46)
	pdf.SetX(20)
	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)

	pdf.Ln(4)
	pdf.SetX(20)
	pdf.CellFormat(90, 5, "DATOS DEL ADQUIRIENTE", "1", 0,
		"L", true, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(94, 5, "LUGAR DE ENTREGA O SERVICIO", "1", 0,
		"L", true, 0, "")
	pdf.Ln(8)

	// DETALLE ADQUIRIENTE
	pdf.SetTextColor(0,0,0)
	pdf.SetX(21)
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

	pdf.SetX(21)
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

	pdf.SetX(21)
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

	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Forma de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miVenta.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miVenta.Fecha.Format("02/01/2006")+" "+Titulo(miVenta.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miVenta.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miVenta.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Cotizacion No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, miVenta.Cotizacion, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Vendedor", "", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, miVendedor.Nombre, "", 0,
		"L", false, 0, "")
	//pdf.Ln(-1)

	// CABECERA PRODUCTO
	pdf.SetFont("Arial", "", 10)

	// RELLENO TITULO
	pdf.SetY(78)
	pdf.SetX(20)
	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)

	pdf.Ln(6)
	pdf.SetX(20)
	pdf.CellFormat(183, 5, "ITEM", "1", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(190, 5, "CODIGO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(190, 5, "PRODUCTO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(106)
	pdf.CellFormat(190, 5, "UNIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(122)
	pdf.CellFormat(190, 5, "IVA", "0", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(190, 5, "DESC.", "0", 0,
		"L", false, 0, "")
	pdf.SetX(141)
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
// TERMINA VENTA CABECERA

// INICIA VENTA DETALLE
func VentaFilaDetalle(pdf *gofpdf.Fpdf,miFila ventadetalleeditar, a int ){
	pdf.SetFont("Arial", "", 9)
	if math.Mod(float64(a),2)==0 {
		pdf.SetFillColor(224,231,239)
		pdf.SetTextColor(0,0,0)
	} else{
		pdf.SetFillColor(255,255,255)
		pdf.SetTextColor(0,0,0)
	}

	pdf.SetTextColor(0,0,0)
	pdf.SetX(21)
	pdf.CellFormat(184, 4, strconv.Itoa(a), "", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 4, Subcadena(miFila.Producto,0,12), "", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(40, 4, Subcadena(miFila.ProductoNombre,0,32),  "", 0,
		"L", false, 0, "")
	pdf.SetX(80)
	pdf.CellFormat(40, 4, miFila.ProductoUnidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(89)
	pdf.CellFormat(40, 4, miFila.ProductoIva, "", 0,
		"R", false, 0, "")
	pdf.SetX(100)
	pdf.CellFormat(40, 4, miFila.Descuento, "", 0,
		"R", false, 0, "")
	pdf.SetX(119)
	pdf.CellFormat(40, 4, miFila.Cantidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(141)
	pdf.CellFormat(40, 4, miFila.Precio, "", 0,
		"R", false, 0, "")
	pdf.SetX(164)
	pdf.CellFormat(40, 4, miFila.Subtotal, "", 0,
		"R", false, 0, "")
	pdf.Ln(4)

}
// TERMINA VENTA DETALLE

// INICIA PIE DE PAGINA
func VentaPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miVenta venta){

	Totalletras,err := IntLetra(Cadenaentero(miVenta.Total))
	if err!= nil{
		fmt.Println(err)
	}
	ene := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 8)
	pdf.SetY(222)
	pdf.SetX(21)
	pdf.CellFormat(190, 10, "SON: " +ene(Mayuscula(Totalletras))+" PESOS MDA. CTE.", "0", 0,
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
	pdf.CellFormat(190, 10, "TOTAL", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "DESCUENTO", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "SUBTOTAL", "0", 0, "L",
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
	pdf.CellFormat(190, 10, "TOTAL NETO", "0", 0, "L",
		false, 0, "")

	pdf.SetY(229)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVenta.Subtotal, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVenta.Descuento, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, FormatoFlotanteEntero(Flotante(miVenta.Subtotal)  -  Flotante(miVenta.Descuento)), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVenta.Subtotaliva19, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVenta.Subtotaliva5, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miVenta.Total, "0", 0, "R",
		false, 0, "")

	pdf.Image(imageFile("QR.jpg"), 20, 229, 25, 0, false,
		"", 0, "")
	pdf.SetFont("Arial", "", 8)
	pdf.SetY(249)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "Cufexxxxxxxxxx", "", 0,
		"L", false, 0, "")

}
// TERMINA PEI DE PAGINA

// INICIA FOOTER
func VentaFooter(pdf *gofpdf.Fpdf, miVenta venta){
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
}
// TERMINA FOOTER