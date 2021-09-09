package main

// INICIA SOPORTE IMPORTAR PAQUETES
import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "github.com/bitly/go-simplejson"
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

// INICIA SOPORTE ESTRUCTURA JSON
type devolucionsoporteJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA SOPORTE ESTRUCTURA
type devolucionsoporteLista struct {
	Codigo			  string
	Fecha        	  time.Time
	Neto          	  string
	Tercero       	  string
	TerceroNombre	  string
	CentroNombre 	  string
	AlmacenistaNombre string
}

// INICIA SOPORTE ESTRUCTURA
type devolucionsoporte struct {
	Resolucionsoporte      	  string
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
	PorcentajeRetencionFuente string
	TotalRetencionFuente      string
	PorcentajeRetencionIca    string
	TotalRetencionIca         string
	Neto                      string
	Items                     string
	Formadepago               string
	Mediodepago               string
	Tercero                   string
	Almacenista               string
	Accion                    string
	Detalle                   []devolucionsoportedetalle `json:"Detalle"`
	DetalleEditar			  []devolucionsoportedetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Soporte					  string
	Tipo					  string
	Centro					  string
}

// INICIA SOPORTEDETALLE ESTRUCTURA
type devolucionsoportedetalle struct {
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
	Fecha			  time.Time
}

// INICIA SOPORTE DETALLE EDITARr
type devolucionsoportedetalleeditar struct {
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
	Fecha			  time.Time
}

// INICIA SOPORTE CONSULTA DETALLE
func DevolucionsoporteConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "devolucionsoportedetalle.Id as id ,"
	consulta += "devolucionsoportedetalle.Codigo as codigo,"
	consulta += "devolucionsoportedetalle.Fila as fila,"
	consulta += "devolucionsoportedetalle.Cantidad as cantidad,"
	consulta += "devolucionsoportedetalle.Precio as precio,"
	consulta += "devolucionsoportedetalle.Descuento as descuento,"
	consulta += "devolucionsoportedetalle.Montodescuento as montodescuento,"
	consulta += "devolucionsoportedetalle.Sigratis as sigratis,"
	consulta += "devolucionsoportedetalle.Subtotal as subtotal,"
	consulta += "devolucionsoportedetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "devolucionsoportedetalle.Pagina as pagina ,"
	consulta += "devolucionsoportedetalle.Bodega as bodega,"
	consulta += "devolucionsoportedetalle.Producto as producto,"
	consulta += "devolucionsoportedetalle.Fecha as fecha,"
	consulta += "producto.nombre as ProductoNombre, "
	consulta += "producto.iva as ProductoIva, "
	consulta += "producto.unidad as ProductoUnidad "
	consulta += "from devolucionsoportedetalle "
	consulta += "inner join producto on producto.codigo=devolucionsoportedetalle.producto "
	consulta += " where devolucionsoportedetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA SOPORTE LISTA
func DevolucionsoporteLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionsoporte/devolucionsoporteLista.html")
	log.Println("Error devolucionsoporte 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,centro.nombre as CentroNombre,devolucionsoporte.neto,devolucionsoporte.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM devolucionsoporte "
	consulta += " inner join tercero on tercero.codigo=devolucionsoporte.tercero "
	consulta += " inner join centro on centro.codigo=devolucionsoporte.centro "
	consulta += " inner join almacenista on almacenista.codigo=devolucionsoporte.almacenista "
	consulta += " ORDER BY devolucionsoporte.codigo ASC"

	db := dbConn()
	res := []devolucionsoporteLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error devolucionsoporte888")
	tmp.Execute(w, varmap)
}

// INICIA SOPORTE NUEVO
func DevolucionsoporteNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionsoporte editar" + Codigo)

	db := dbConn()
	v := devolucionsoporte{}
	tc := tercero{}
	det := []devolucionsoportedetalleeditar{}
	if Codigo == "False" {

	} else {

	err := db.Get(&v, "SELECT * FROM devolucionsoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	err2 := db.Select(&det, DevolucionsoporteConsultaDetalle(), Codigo)
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
		"devolucionsoporte":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"resolucionsoporte" : ListaResolucionsoporte(),
		"retfte":TraerParametrosInventario().Soportecuentaporcentajeretfte,
		"codigo": Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	t, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/devolucionsoporte/devolucionsoporteNuevo.html", "vista/devolucionsoporte/devolucionsoporteScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error devolucionsoporte nuevo 3")
	t.Execute(w, parametros)
}

// INICIA INSERTAR COMPROBANTE DE SOPORTE
func InsertaDetalleComprobanteDevolucionsoporte(miFilaComprobante comprobantedetalle, miComprobante comprobante, miDevolucionsoporte devolucionsoporte){
	db := dbConn()

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miDevolucionsoporte.Tercero)
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
		miFilaComprobante.Debito=	miFilaComprobante.Debito
	}

	if len(miFilaComprobante.Credito)>0 {
		miFilaComprobante.Credito=	miFilaComprobante.Credito
	}

	// TERMINA COMPROBANTE GRABAR INSERTAR
	_, err = insForm.Exec(
		miFilaComprobante.Fila,
	miFilaComprobante.Cuenta  ,
	miDevolucionsoporte.Tercero,
	miDevolucionsoporte.Centro,
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

// INICIA SOPORTE INSERTAR AJAX
func DevolucionsoporteAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempDevolucionsoporte devolucionsoporte

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la SOPORTE
	err = json.Unmarshal(b, &tempDevolucionsoporte)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	if tempDevolucionsoporte.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from devolucionsoportedetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempDevolucionsoporte.Codigo)

		// borra detalle inventario
		delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo='Devolucionsoporte'")
		if err != nil {
			panic(err.Error())
		}
		delForm2.Exec(tempDevolucionsoporte.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from devolucionsoporte WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempDevolucionsoporte.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempDevolucionsoporte.Detalle {
		var a = i
		var q string
		q = "insert into devolucionsoportedetalle ("
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

		// TERMINA SOPORTE GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Id,
			x.Codigo,
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
	for i, x := range tempDevolucionsoporte.Detalle {
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

		// TERMINA SOPORTE GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Fecha,
			x.Tipo,
			x.Codigo,
			x.Bodega,
			x.Producto,
			x.Cantidad,
			x.Precio,
			operacionDevolucionSoporte)
		if err != nil {
			panic(err)
		}
		log.Println("Insertar Producto \n", x.Producto, a)
	}

	// INICIA INSERTAR SOPORTES
	log.Println("Got %s age %s club %s\n", tempDevolucionsoporte.Codigo, tempDevolucionsoporte.Tercero, tempDevolucionsoporte.Subtotal)
	var q string
	q = "insert into devolucionsoporte ("
	q += "Resolucionsoporte,"
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
	q += "Total,"
	q += "PorcentajeRetencionFuente,"
	q += "TotalRetencionFuente,"
	q += "PorcentajeRetencionIca,"
	q += "TotalRetencionIca,"
	q += "Neto,"
	q += "Items,"
	q += "Formadepago,"
	q += "Mediodepago,"
	q += "Tercero,"
	q += "Almacenista,"
	q += "Soporte,"
	q += "Centro,"
	q += "Tipo"
	q += " ) values("
	q+=parametros(34)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempDevolucionsoporte.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempDevolucionsoporte.Resolucionsoporte,
		tempDevolucionsoporte.Codigo,
		tempDevolucionsoporte.Fecha.Format(layout),
		tempDevolucionsoporte.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempDevolucionsoporte.Descuento,
		tempDevolucionsoporte.Subtotaldescuento19,
		tempDevolucionsoporte.Subtotaldescuento5,
		tempDevolucionsoporte.Subtotaldescuento0,
		tempDevolucionsoporte.Subtotal,
		tempDevolucionsoporte.Subtotal19,
		tempDevolucionsoporte.Subtotal5,
		tempDevolucionsoporte.Subtotal0,
		tempDevolucionsoporte.Subtotaliva19,
		tempDevolucionsoporte.Subtotaliva5,
		tempDevolucionsoporte.Subtotaliva0,
		tempDevolucionsoporte.Subtotalbase19,
		tempDevolucionsoporte.Subtotalbase5,
		tempDevolucionsoporte.Subtotalbase0,
		tempDevolucionsoporte.TotalIva,
		tempDevolucionsoporte.Total,
		tempDevolucionsoporte.PorcentajeRetencionFuente,
		tempDevolucionsoporte.TotalRetencionFuente,
		tempDevolucionsoporte.PorcentajeRetencionIca,
		tempDevolucionsoporte.TotalRetencionIca,
		tempDevolucionsoporte.Neto,
		tempDevolucionsoporte.Items,
		tempDevolucionsoporte.Formadepago,
		tempDevolucionsoporte.Mediodepago,
		tempDevolucionsoporte.Tercero,
		tempDevolucionsoporte.Almacenista,
		tempDevolucionsoporte.Soporte,
		tempDevolucionsoporte.Centro,
		tempDevolucionsoporte.Tipo)

	if err != nil {
		panic(err)
	}

	// INSERTAR COMPROBANTE CONTABILIDAD
	var tempComprobante comprobante
	var tempComprobanteDetalle comprobantedetalle
	tempComprobante.Documento="24"
	tempComprobante.Numero=tempDevolucionsoporte.Codigo
	tempComprobante.Fecha =tempDevolucionsoporte.Fecha
	tempComprobante.Fechaconsignacion =tempDevolucionsoporte.Fecha
	tempComprobante.Debito = tempDevolucionsoporte.Neto + ".00"
	tempComprobante.Credito	= tempDevolucionsoporte.Neto + ".00"
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

	// INSERTAR CUENTA CREDITO PROVEEDOR
	if (tempDevolucionsoporte.Neto!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionsoporte.Neto)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soportecuentaproveedor
		tempComprobanteDetalle.Debito = tempDevolucionsoporte.Neto
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionsoporte(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporte)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO DESCUENTO
	//if (tempDevolucionsoporte.Descuento!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalDebito+=Flotante(tempDevolucionsoporte.Descuento)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Soportedevolucioncuentadescuento
	//	tempComprobanteDetalle.Debito = tempDevolucionsoporte.Descuento
	//	tempComprobanteDetalle.Credito = ""
	//	InsertaDetalleComprobanteDevolucionsoporte(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporte)
	//	log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	//}

	// INSERTAR CUENTA CREDITO RET. FTE.
	if (tempDevolucionsoporte.TotalRetencionFuente!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionsoporte.TotalRetencionFuente)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soportedevolucioncuentaretfte
		tempComprobanteDetalle.Debito = tempDevolucionsoporte.TotalRetencionFuente
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionsoporte(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporte)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO RET. ICA.
	if (tempDevolucionsoporte.TotalRetencionIca!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionsoporte.TotalRetencionIca)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soportedevolucioncuentaretica
		tempComprobanteDetalle.Debito = tempDevolucionsoporte.TotalRetencionIca
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionsoporte(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporte)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}


	// INSERTAR CUENTA DEBITO SOPORTE 19%
	if (tempDevolucionsoporte.Subtotal19!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempDevolucionsoporte.Subtotal19)-Flotante(tempDevolucionsoporte.Subtotaldescuento19)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soportedevolucioncuenta19
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = FormatoFlotanteComprobante(Flotante(tempDevolucionsoporte.Subtotal19)-Flotante(tempDevolucionsoporte.Subtotaldescuento19))
		InsertaDetalleComprobanteDevolucionsoporte(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporte)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA DEBITO SOPORTE 5%
	if (tempDevolucionsoporte.Subtotal5!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempDevolucionsoporte.Subtotal5)-Flotante(tempDevolucionsoporte.Subtotaldescuento5)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soportedevolucioncuenta5
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = FormatoFlotanteComprobante(Flotante(tempDevolucionsoporte.Subtotal5)-Flotante(tempDevolucionsoporte.Subtotaldescuento5))
		InsertaDetalleComprobanteDevolucionsoporte(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporte)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA DEBITO SOPORTE 0%
	if (tempDevolucionsoporte.Subtotal0!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempDevolucionsoporte.Subtotal0)-Flotante(tempDevolucionsoporte.Subtotaldescuento0)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soportedevolucioncuenta0
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = FormatoFlotanteComprobante(Flotante(tempDevolucionsoporte.Subtotal0)-Flotante(tempDevolucionsoporte.Subtotaldescuento0))
		InsertaDetalleComprobanteDevolucionsoporte(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporte)
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

// TRAER SOPORTE DE LA DEVOLUCION
func DatosSoporte(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio soporte editar" + Codigo)
	db := dbConn()
	var res []soporte

	// traer SOPORTE
	v := soporte{}
	err := db.Get(&v, "SELECT * FROM soporte where codigo=$1", Codigo)
	var valida bool
	valida=true

	switch err {
	case nil:
		log.Printf("soporte existe: %+v\n", v)
	case sql.ErrNoRows:
		log.Println("soporte NO Existe")
		valida=false
	default:
		log.Printf("soporte error: %s\n", err)
	}
	det := []soportedetalleeditar{}
	t := tercero{}

	// trae datos si existe soporte
	if valida==true {
		err2 := db.Select(&det, SoporteConsultaDetalle(), Codigo)
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

// INICIA SOPORTE EXISTE
func DevolucionsoporteExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM devolucionsoporte  WHERE codigo=$1", Codigo)
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

// INICIA SOPORTE EDITAR
func DevolucionsoporteEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionsoporte editar" + Codigo)

	db := dbConn()

	// traer devolucion soporte
	v := devolucionsoporte{}
	err := db.Get(&v, "SELECT * FROM devolucionsoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []devolucionsoportedetalleeditar{}

	err2 := db.Select(&det, DevolucionsoporteConsultaDetalle(), Codigo)
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
		"devolucionsoporte":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"resolucionsoporte" : ListaResolucionsoporte(),
		"retfte":TraerParametrosInventario().Soportecuentaporcentajeretfte,

	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionsoporte/devolucionsoporteEditar.html",
		"vista/devolucionsoporte/devolucionsoporteScript.html")
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error devolucionsoporte nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA SOPORTE BORRAR
func DevolucionsoporteBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionsoporte editar" + Codigo)

	db := dbConn()

	// traer SOPORTE
	v := devolucionsoporte{}
	err := db.Get(&v, "SELECT * FROM devolucionsoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []devolucionsoportedetalleeditar{}
	err2 := db.Select(&det, DevolucionsoporteConsultaDetalle(), Codigo)
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
		"devolucionsoporte":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"retfte":TraerParametrosInventario().Soportecuentaporcentajeretfte,
		"centro":ListaCentro(),
		"resolucionsoporte" : ListaResolucionsoporte(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionsoporte/devolucionsoporteBorrar.html", "vista/devolucionsoporte/devolucionsoporteScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error devolucionsoporte nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA SOPORTE ELIMINAR
func DevolucionsoporteEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar SOPORTE
	delForm, err := db.Prepare("DELETE from devolucionsoporte WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from devolucionsoportedetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borar detalle invenario
	Borrarinventario(codigo,"Devolucionsoporte")

	// borra detalle anterior
	delForm, err = db.Prepare("DELETE from comprobantedetalle WHERE documento=$1 and numero=$2")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec("8", codigo)

	// borra cabecera anterior

	delForm1, err = db.Prepare("DELETE from comprobante WHERE documento=$1 and numero=$2 ")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec("8", codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/DevolucionsoporteLista", 301)
}

// INICIA PDF
func DevolucionsoportePdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer SOPORTE
	miDevolucionsoporte := devolucionsoporte{}
	err := db.Get(&miDevolucionsoporte, "SELECT * FROM devolucionsoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []devolucionsoportedetalleeditar{}
	err2 := db.Select(&miDetalle, DevolucionsoporteConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miDevolucionsoporte.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miDevolucionsoporte.Almacenista)
	if err4 != nil {
		log.Fatalln(err4)
	}

	var e empresa = ListaEmpresa()
	var c ciudad = TraerCiudad(e.Ciudad)
	var re Resolucionsoporte = TraerResolucionsoporte(miDevolucionsoporte.Resolucionsoporte)
	var buf bytes.Buffer
	var err1 error

	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)
	ene := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetHeaderFunc(func() {
		// logo
		pdf.Image(imageFile("logo.png"), 20, 30, 40, 0, false,
			"", 0, "")
		pdf.SetY(5)

		// DATOS EMPRESA
		pdf.SetFont("Arial", "", 9)
		// CUADRO EMPRESA
		pdf.SetY(20)
		pdf.SetX(20)
		pdf.SetDrawColor(119,134,153)
		pdf.CellFormat(184, 69, "", "1", 0, "C",
			false, 0, "")

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
		pdf.CellFormat(190, 10, e.Telefono1+" "+e.Telefono2, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+" - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")

		// RESOLUCION
		pdf.SetFont("Arial", "", 9)
		pdf.SetY(20)
		pdf.SetX(75)
		pdf.CellFormat(190, 10, "Resolucion No. "+re.Numero, "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.SetX(75)
		pdf.CellFormat(190, 10, "Del Numero "+re.Prefijo+" "+re.NumeroInicial+" al "+re.Prefijo+" "+Coma(re.NumeroFinal), "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.SetX(75)

		pdf.CellFormat(190, 10, "Vigencia del "+re.FechaInicial.Format("02/01/2006")+" al "+re.FechaFinal.Format("02/01/2006"), "0", 0, "C",
			false, 0, "")

		// DEVOLUCION SOPORTE NUMERO
		pdf.Ln(5)
		pdf.SetX(75)
		pdf.CellFormat(190, 10, "DEVOLUCION SOPORTE ADQUISICIONES", "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
		pdf.SetX(75)
		pdf.CellFormat(190, 10, "EFECTUADAS A NO OBLIGADOS A ", "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
		pdf.SetX(75)
		pdf.CellFormat(190, 10, "FACTURAR No. "+re.Prefijo+" "+Codigo, "0", 0, "C",
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
	DevolucionsoporteCabecera(pdf,miTercero,miDevolucionsoporte,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			DevolucionsoporteFilaDetalle(pdf,miFila,a)
		}
		DevolucionsoportePieDePagina(pdf,miTercero,miDevolucionsoporte)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					DevolucionsoporteFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			DevolucionsoporteCabecera(pdf,miTercero,miDevolucionsoporte,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					DevolucionsoporteFilaDetalle(pdf,miFila,a)
				}
			}

			DevolucionsoportePieDePagina(pdf,miTercero,miDevolucionsoporte)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					DevolucionsoporteFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			DevolucionsoporteCabecera(pdf,miTercero,miDevolucionsoporte,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					DevolucionsoporteFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			DevolucionsoporteCabecera(pdf,miTercero,miDevolucionsoporte,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					DevolucionsoporteFilaDetalle(pdf,miFila,a)
				}
			}

			DevolucionsoportePieDePagina(pdf,miTercero,miDevolucionsoporte)
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
func DevolucionsoporteCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miDevolucionsoporte devolucionsoporte, miAlmacenista almacenista ){

	// RELLENO TITULO
	pdf.SetY(46)
	pdf.SetX(20)
	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(8)
	pdf.SetX(20)
	pdf.CellFormat(184, 5, "DATOS DEL PROVEEDOR", "0", 0,
		"L", true, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(94, 5, "LUGAR DE ENTREGA O SERVICIO", "1", 0,
		"L", false, 0, "")
	pdf.Ln(8)
	pdf.SetX(20)
	pdf.SetTextColor(0,0,0)
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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miDevolucionsoporte.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miDevolucionsoporte.Fecha.Format("02/01/2006")+" "+Titulo(miDevolucionsoporte.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miDevolucionsoporte.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miDevolucionsoporte.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Soporte No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, miDevolucionsoporte.Soporte, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Almacenista", "", 0,
		"L", false, 0, "")
	pdf.SetX(135)
	pdf.CellFormat(40, 4, miAlmacenista.Nombre, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.SetY(88)
	pdf.SetX(20)

	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)

	pdf.CellFormat(184, 5, "ITEM", "1", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(190, 5, "CODIGO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(190, 5, "PRODUCTO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(106)
	pdf.CellFormat(40, 5, "UNIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(122)
	pdf.CellFormat(40, 5, "IVA", "0", 0,
		"L", false, 0, "")
	pdf.SetX(135)
	pdf.CellFormat(40, 5, "CANTIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(160)
	pdf.CellFormat(40, 5, "P. UNITARIO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(190)
	pdf.CellFormat(40, 5, "TOTAL", "0", 0,
		"L", false, 0, "")
	pdf.Ln(8)
}
func DevolucionsoporteFilaDetalle(pdf *gofpdf.Fpdf,miFila devolucionsoportedetalleeditar, a int ){
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
	pdf.CellFormat(183, 4, strconv.Itoa(a), "", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 4, Subcadena(miFila.Producto,0,12), "", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(40, 4, Subcadena(miFila.ProductoNombre,0,32),  "", 0,
		"L", false, 0, "")
	pdf.SetX(82)
	pdf.CellFormat(40, 4, miFila.ProductoUnidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(89)
	pdf.CellFormat(40, 4, miFila.ProductoIva, "", 0,
		"R", false, 0, "")
	pdf.SetX(115)
	pdf.CellFormat(40, 4, miFila.Cantidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(140)
	pdf.CellFormat(40, 4, miFila.Precio, "", 0,
		"R", false, 0, "")
	pdf.SetX(164)
	pdf.CellFormat(40, 4, miFila.Subtotal, "", 0,
		"R", false, 0, "")
	pdf.Ln(4)
}

func DevolucionsoportePieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miDevolucionsoporte devolucionsoporte ){

	Totalletras,err := IntLetra(Cadenaentero(miDevolucionsoporte.Neto))
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
	pdf.CellFormat(190, 10, "FIRMA RESPONSABLE ", "0", 0, "C",
		false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetY(229)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "SUBTOTAL", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "RET. FTE.", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "RET. ICA.", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "TOTAL", "0", 0, "L",
		false, 0, "")

	pdf.SetY(229)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miDevolucionsoporte.Subtotal, "0", 0, "R",
		false, 0, "")

	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miDevolucionsoporte.TotalRetencionFuente, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miDevolucionsoporte.TotalRetencionIca, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miDevolucionsoporte.Neto, "0", 0, "R",
		false, 0, "")

	pdf.Image(imageFile("QR.jpg"), 20, 229, 25, 0, false,
		"", 0, "")
	pdf.SetFont("Arial", "", 8)

	pdf.SetY(249)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "Cufexxxxxxxxxx", "", 0,
		"L", false, 0, "")
}
