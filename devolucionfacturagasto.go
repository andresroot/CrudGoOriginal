package main

// INICIA DEVOLUCION FACTURA GASTO IMPORTAR PAQUETES
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

// INICIA DEVOLUCION FACTURA GASTO ESTRUCTURA JSON
type devolucionfacturagastoJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA DEVOLUCION FACTURA GASTO ESTRUCTURA
type devolucionfacturagastoLista struct {
	Codigo			  string
	Fecha        	  time.Time
	Neto          	  string
	Tercero       	  string
	TerceroNombre	  string
	CentroNombre 	  string
	AlmacenistaNombre string
}

// INICIA DEVOLUCION FACTURA GASTO ESTRUCTURA
type devolucionfacturagasto struct {
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
	Detalle                   []devolucionfacturagastodetalle `json:"Detalle"`
	DetalleEditar			  []devolucionfacturagastodetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Facturagasto	          string
	Tipo					  string
	Centro					  string
	Facturagastocuenta0	      string
	Facturagastonombre0	      string
	Facturagastocuentaporcentajeretfte	      string
	Facturagastocuentaretfte	  string
	Facturagastonombreretfte	  string
	Facturagastocuentaiva	      string
	Facturagastonombreiva	      string
	Facturagastoporcentajeiva	      string
}

// INICIA SOPORTE SERVICIODETALLE ESTRUCTURA
type devolucionfacturagastodetalle struct {
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
	Tipo              string
	Fecha			  time.Time
	Nombreservicio    string
	Unidadservicio    string
	Codigoservicio    string

}

// INICIA DEVOLUCION FACTURA GASTO DETALLE EDITAR
type devolucionfacturagastodetalleeditar struct {
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
	Tipo              string
	Fecha			  time.Time
	Nombreservicio    string
	Unidadservicio    string
	Codigoservicio    string

}

// INICIA DEVOLUCION FACTURA GASTO CONSULTA DETALLE
func DevolucionfacturagastoConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "facturagastodetalle.Id as id ,"
	consulta += "facturagastodetalle.Codigo as codigo,"
	consulta += "facturagastodetalle.Fila as fila,"
	consulta += "facturagastodetalle.Cantidad as cantidad,"
	consulta += "facturagastodetalle.Precio as precio,"
	consulta += "facturagastodetalle.Descuento as descuento,"
	consulta += "facturagastodetalle.Montodescuento as montodescuento,"
	consulta += "facturagastodetalle.Sigratis as sigratis,"
	consulta += "facturagastodetalle.Subtotal as subtotal,"
	consulta += "facturagastodetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "facturagastodetalle.Pagina as pagina ,"
	consulta += "facturagastodetalle.Bodega as bodega,"
	consulta += "facturagastodetalle.Fecha as fecha,"
	consulta += "facturagastodetalle.Nombreservicio,"
	consulta += "facturagastodetalle.Unidadservicio,"
	consulta += "facturagastodetalle.Codigoservicio"
	consulta += " from facturagastodetalle "
	consulta += " where facturagastodetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA DEVOLUCION FACTURA GASTO LISTA
func DevolucionfacturagastoLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionfacturagasto/devolucionfacturagastoLista.html")
	log.Println("Error devolucionfacturagasto 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,centro.nombre as CentroNombre,devolucionfacturagasto.neto,devolucionfacturagasto.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM devolucionfacturagasto "
	consulta += " inner join tercero on tercero.codigo=devolucionfacturagasto.tercero "
	consulta += " inner join centro on centro.codigo=devolucionfacturagasto.centro "
	consulta += " inner join almacenista on almacenista.codigo=devolucionfacturagasto.almacenista "
	consulta += " ORDER BY devolucionfacturagasto.codigo ASC"

	db := dbConn()
	res := []devolucionfacturagastoLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error facturagasto888")
	tmp.Execute(w, varmap)
}

// INICIA DEVOLUCION FACTURA GASTO NUEVO
func DevolucionfacturagastoNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error devolucionfacturagasto nuevo 1")
	log.Println("Error devolucionfacturagasto nuevo 2")
	parametros := map[string]interface{}{
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro":      ListaCentro(),
	}

	t, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionfacturagasto/devolucionfacturagastoNuevo.html",
		"vista/devolucionfacturagasto/devolucionfacturagastoScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	log.Println("Error devolucionfacturagasto nuevo 3")
	t.Execute(w, parametros)
}

// INICIA INSERTAR COMPROBANTE DE DEVOLUCION FACTURA GASTO
func InsertaDetalleComprobanteDevolucionfacturagasto(miFilaComprobante comprobantedetalle, miComprobante comprobante, miDevolucionfacturagasto devolucionfacturagasto){
	db := dbConn()

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miDevolucionfacturagasto.Tercero)
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
	q += " )"

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	if len(miFilaComprobante.Debito)>0 {
		miFilaComprobante.Debito=	miFilaComprobante.Debito
		//+ ".00"
	}

	if len(miFilaComprobante.Credito)>0 {
		miFilaComprobante.Credito=	miFilaComprobante.Credito
		//+ ".00"
	}

	// TERMINA COMPROBANTE GRABAR INSERTAR
	_, err = insForm.Exec(
		miFilaComprobante.Fila,
	miFilaComprobante.Cuenta  ,
	miDevolucionfacturagasto.Tercero,
	miDevolucionfacturagasto.Centro,
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

// INICIA DEVOLUCION FACTURA GASTO INSERTAR AJAX
func DevolucionfacturagastoAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempDevolucionfacturagasto devolucionfacturagasto

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la DEVOLUCION FACTURA GASTO
	err = json.Unmarshal(b, &tempDevolucionfacturagasto)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}
	var DocumentoContable string
	DocumentoContable="31"

	if tempDevolucionfacturagasto.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from devolucionfacturagastodetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempDevolucionfacturagasto.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from devolucionfacturagasto WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempDevolucionfacturagasto.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempDevolucionfacturagasto.Detalle {
		var a = i
		var q string
		q = "insert into devolucionfacturagastodetalle ("
		q += "Id,"
		q += "Codigo,"
		q += "Fila,"
		q += "Cantidad,"
		q += "Precio,"
		q += "Subtotal,"
		q += "Pagina,"
		q += "Bodega,"
		q += "Descuento,"
		q += "Montodescuento,"
		q += "Sigratis,"
		q += "Subtotaldescuento,"
		q += "Tipo,"
		q += "Fecha,"
		q += "Nombreservicio,"
		q += "Unidadservicio,"
		q += "CodigoServicio"
		q += " ) values("
		q += parametros(17)
		q += " )"

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA DEVOLUCION FACTURA GASTO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Id,
			x.Codigo,
			x.Fila,
			x.Cantidad,
			x.Precio,
			x.Subtotal,
			x.Pagina,
			x.Bodega,
			x.Descuento,
			x.Montodescuento,
			x.Sigratis,
			x.Subtotaldescuento,
			x.Tipo,
			x.Fecha,
			x.Nombreservicio,
			x.Unidadservicio,
			x.Codigoservicio)
		if err != nil {
			panic(err)
		}

		log.Println("Insertar Producto \n", x.Nombreservicio, a)
	}

	// INICIA INSERTAR DEVOLUCION FACTURA GASTO
	log.Println("Got %s age %s club %s\n", tempDevolucionfacturagasto.Codigo, tempDevolucionfacturagasto.Tercero, tempDevolucionfacturagasto.Subtotal)
	var q string
	q = "insert into devolucionfacturagasto ("
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
	q += "Facturagasto,"
	q += "Centro,"
	q += "Tipo,"
	q += "Facturagastocuenta0,"
	q += "Facturagastonombre0,"
	q += "Facturagastocuentaporcentajeretfte,"
	q += "Facturagastocuentaretfte,"
	q += "Facturagastonombreretfte,"
	q += "Facturagastocuentaiva,"
	q += "Facturagastonombreiva,"
	q += "Facturagastoporcentajeiva"
	q += " ) values("
	q+=parametros(41)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempDevolucionfacturagasto.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempDevolucionfacturagasto.Codigo,
		tempDevolucionfacturagasto.Fecha.Format(layout),
		tempDevolucionfacturagasto.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempDevolucionfacturagasto.Descuento,
		tempDevolucionfacturagasto.Subtotaldescuento19,
		tempDevolucionfacturagasto.Subtotaldescuento5,
		tempDevolucionfacturagasto.Subtotaldescuento0,
		tempDevolucionfacturagasto.Subtotal,
		tempDevolucionfacturagasto.Subtotal19,
		tempDevolucionfacturagasto.Subtotal5,
		tempDevolucionfacturagasto.Subtotal0,
		tempDevolucionfacturagasto.Subtotaliva19,
		tempDevolucionfacturagasto.Subtotaliva5,
		tempDevolucionfacturagasto.Subtotaliva0,
		tempDevolucionfacturagasto.Subtotalbase19,
		tempDevolucionfacturagasto.Subtotalbase5,
		tempDevolucionfacturagasto.Subtotalbase0,
		tempDevolucionfacturagasto.TotalIva,
		tempDevolucionfacturagasto.Total,
		tempDevolucionfacturagasto.PorcentajeRetencionFuente,
		tempDevolucionfacturagasto.TotalRetencionFuente,
		tempDevolucionfacturagasto.PorcentajeRetencionIca,
		tempDevolucionfacturagasto.TotalRetencionIca,
		tempDevolucionfacturagasto.Neto,
		tempDevolucionfacturagasto.Items,
		tempDevolucionfacturagasto.Formadepago,
		tempDevolucionfacturagasto.Mediodepago,
		tempDevolucionfacturagasto.Tercero,
		tempDevolucionfacturagasto.Almacenista,
		tempDevolucionfacturagasto.Facturagasto,
		tempDevolucionfacturagasto.Centro,
		tempDevolucionfacturagasto.Tipo,
		tempDevolucionfacturagasto.Facturagastocuenta0,
		tempDevolucionfacturagasto.Facturagastonombre0,
		tempDevolucionfacturagasto.Facturagastocuentaporcentajeretfte,
		tempDevolucionfacturagasto.Facturagastocuentaretfte,
		tempDevolucionfacturagasto.Facturagastonombreretfte,
		tempDevolucionfacturagasto.Facturagastocuentaiva,
		tempDevolucionfacturagasto.Facturagastonombreiva,
		tempDevolucionfacturagasto.Facturagastoporcentajeiva)

	if err != nil {
		panic(err)
	}

	// INSERTAR COMPROBANTE CONTABILIDAD
	var tempComprobante comprobante
	var tempComprobanteDetalle comprobantedetalle
	tempComprobante.Documento=DocumentoContable
	tempComprobante.Numero=tempDevolucionfacturagasto.Codigo
	tempComprobante.Fecha =tempDevolucionfacturagasto.Fecha
	tempComprobante.Fechaconsignacion =tempDevolucionfacturagasto.Fecha
	tempComprobante.Debito = tempDevolucionfacturagasto.Neto + ".00"
	tempComprobante.Credito	= tempDevolucionfacturagasto.Neto + ".00"
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

	// INSERTAR CUENTA DEBITO PROVEEDOR
	if (tempDevolucionfacturagasto.Neto!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionfacturagasto.Neto)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Facturagastocuentaproveedor
		tempComprobanteDetalle.Debito = tempDevolucionfacturagasto.Neto
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionfacturagasto(tempComprobanteDetalle,tempComprobante,tempDevolucionfacturagasto)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO RET. FTE.
	if (tempDevolucionfacturagasto.TotalRetencionFuente!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionfacturagasto.TotalRetencionFuente)
		tempComprobanteDetalle.Cuenta  = tempDevolucionfacturagasto.Facturagastocuentaretfte
		tempComprobanteDetalle.Debito = tempDevolucionfacturagasto.TotalRetencionFuente
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionfacturagasto(tempComprobanteDetalle,tempComprobante,tempDevolucionfacturagasto)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO RET. ICA.
	if (tempDevolucionfacturagasto.TotalRetencionIca!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionfacturagasto.TotalRetencionIca)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Facturagastocuentaretica
		tempComprobanteDetalle.Debito = tempDevolucionfacturagasto.TotalRetencionIca
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionfacturagasto(tempComprobanteDetalle,tempComprobante,tempDevolucionfacturagasto)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR TOTAL IVA
	if (tempDevolucionfacturagasto.TotalIva!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempDevolucionfacturagasto.TotalIva)
		tempComprobanteDetalle.Cuenta  = tempDevolucionfacturagasto.Facturagastocuentaiva
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempDevolucionfacturagasto.TotalIva
		InsertaDetalleComprobanteDevolucionfacturagasto(tempComprobanteDetalle,tempComprobante,tempDevolucionfacturagasto)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA DEBITO DEVOLUCION FACTURA GASTO 0%
	if (tempDevolucionfacturagasto.Subtotal0!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempDevolucionfacturagasto.Subtotal0)-Flotante(tempDevolucionfacturagasto.Subtotaldescuento0)
		tempComprobanteDetalle.Cuenta  = tempDevolucionfacturagasto.Facturagastocuenta0
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = FormatoFlotanteComprobante(Flotante(tempDevolucionfacturagasto.Subtotal0)-Flotante(tempDevolucionfacturagasto.Subtotaldescuento0))
		InsertaDetalleComprobanteDevolucionfacturagasto(tempComprobanteDetalle,tempComprobante,tempDevolucionfacturagasto)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO DESCUENTO
	//if (tempDevolucionfacturagasto.Descuento!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalCredito+=Flotante(tempDevolucionfacturagasto.Descuento)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Facturagastocuentadescuento
	//	tempComprobanteDetalle.Debito = ""
	//	tempComprobanteDetalle.Credito = tempDevolucionfacturagasto.Descuento
	//	InsertaDetalleComprobanteDevolucionfacturagasto(tempComprobanteDetalle,tempComprobante,tempDevolucionfacturagasto)
	//	log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	//}

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
	q += " ) values("
	q += parametros(10)
	q += " )"

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

// INICIA DEVOLUCION FACTURA GASTO EXISTE
func DevolucionfacturagastoExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM devolucionfacturagasto  WHERE codigo=$1", Codigo)
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

// INICIA DEVOLUCION FACTURA GASTO EDITAR
func DevolucionfacturagastoEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionfacturagasto editar" + Codigo)

	db := dbConn()

	// TRAE DEVOLUCION FACTURA GASTO
	v := devolucionfacturagasto{}
	err := db.Get(&v, "SELECT * FROM devolucionfacturagasto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []devolucionfacturagastodetalleeditar{}

	err2 := db.Select(&det, DevolucionfacturagastoConsultaDetalle(), Codigo)
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
		"devolucionfacturagasto":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),

	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionfacturagasto/devolucionfacturagastoEditar.html",
		"vista/devolucionfacturagasto/devolucionfacturagastoScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error devolucionfacturagasto nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA DEVOLUCION FACTURA GASTO BORRAR
func DevolucionfacturagastoBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionfacturagasto editar" + Codigo)

	db := dbConn()

	// TRAE DEVOLUCION FACTURA GASTO
	v := devolucionfacturagasto{}
	err := db.Get(&v, "SELECT * FROM devolucionfacturagasto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []devolucionfacturagastodetalleeditar{}
	err2 := db.Select(&det, DevolucionfacturagastoConsultaDetalle(), Codigo)
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
		"devolucionfacturagasto":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionfacturagasto/devolucionfacturagastoBorrar.html", "vista/devolucionfacturagasto/devolucionfacturagastoScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error devolucionfacturagasto nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA DEVOLUCION FACTURA GASTO ELIMINAR
func DevolucionfacturagastoEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar DEVOLUCION FACTURA GASTO
	delForm, err := db.Prepare("DELETE from devolucionfacturagasto WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from devolucionfacturagastodetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borra detalle anterior
	delForm, err = db.Prepare("DELETE from comprobantedetalle WHERE documento=$1 and numero=$2")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec("10", codigo)

	// borra cabecera anterior

	delForm1, err = db.Prepare("DELETE from comprobante WHERE documento=$1 and numero=$2 ")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec("10", codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/DevolucionfacturagastoLista", 301)
}

// TRAER PEDIDO DEVOLUCION FACTURA GASTO
func Datosfacturagasto(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio factura gasto editar" + Codigo)
	db := dbConn()
	var res []facturagasto

	// TRAE FACTURA GASTO
	v := facturagasto{}
	err := db.Get(&v, "SELECT * FROM facturagasto where codigo=$1", Codigo)
	var valida bool
	valida=true

	switch err {
	case nil:
		log.Printf(" facturagasto existe: %+v\n", v)
	case sql.ErrNoRows:
		log.Println(" facturagasto NO Existe")
		valida=false
	default:
		log.Printf(" facturagasto error: %s\n", err)
	}
	det := []facturagastodetalleeditar{}
	t := tercero{}

	// trae datos si existe factura de gastos
	if valida==true {
		err2 := db.Select(&det, FacturagastoConsultaDetalle(), Codigo)
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
func DevolucionfacturagastoPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer DEVOLUCION FACTURA GASTO
	miDevolucionfacturagasto := devolucionfacturagasto{}
	err := db.Get(&miDevolucionfacturagasto, "SELECT * FROM devolucionfacturagasto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []devolucionfacturagastodetalleeditar{}
	err2 := db.Select(&miDetalle, DevolucionfacturagastoConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miDevolucionfacturagasto.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miDevolucionfacturagasto.Almacenista)
	if err4 != nil {
		log.Fatalln(err4)
	}

	var e empresa = ListaEmpresa()
	var c ciudad = TraerCiudad(e.Ciudad)
	var buf bytes.Buffer
	var err1 error

	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)
	ene := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetHeaderFunc(func() {
		// LOGO
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

		// DEVOLUCION FACTURA GASTO NUMERO
		pdf.SetY(20)
		pdf.SetX(80)
		pdf.Ln(8)
		pdf.SetX(75)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(190, 10, "DEVOLUCION FACTURA GASTO", "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
		pdf.SetX(75)
		pdf.CellFormat(190, 10, "No. "+Codigo, "0", 0, "C",
			false, 0, "")
		pdf.Ln(6)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetFont("Arial", "", 8)
		pdf.SetY(256)
		pdf.SetX(20)
		pdf.CellFormat(80, 10, "www.Sadconf.com.co", "",
			0, "L", false, 0, "")
		pdf.SetX(130)
		pdf.CellFormat(75, 10, fmt.Sprintf(" %d de {nb}", pdf.PageNo()), "",
			0, "R", false, 0, "")

	})

	// inicia suma de paginas
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.Ln(1)
	DevolucionfacturagastoCabecera(pdf,miTercero,miDevolucionfacturagasto,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			DevolucionfacturagastoFilaDetalle(pdf,miFila,a)
		}
		DevolucionfacturagastoPieDePagina(pdf,miTercero,miDevolucionfacturagasto)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					DevolucionfacturagastoFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			DevolucionfacturagastoCabecera(pdf,miTercero,miDevolucionfacturagasto,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					DevolucionfacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			DevolucionfacturagastoPieDePagina(pdf,miTercero,miDevolucionfacturagasto)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					DevolucionfacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			DevolucionfacturagastoCabecera(pdf,miTercero,miDevolucionfacturagasto,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					DevolucionfacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			DevolucionfacturagastoCabecera(pdf,miTercero,miDevolucionfacturagasto,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					DevolucionfacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			DevolucionfacturagastoPieDePagina(pdf,miTercero,miDevolucionfacturagasto)
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
func DevolucionfacturagastoCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miDevolucionfacturagasto devolucionfacturagasto, miAlmacenista almacenista ){

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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miDevolucionfacturagasto.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miDevolucionfacturagasto.Fecha.Format("02/01/2006")+" "+Titulo(miDevolucionfacturagasto.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miDevolucionfacturagasto.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miDevolucionfacturagasto.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Factura Gasto No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, miDevolucionfacturagasto.Facturagasto, "", 0,
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
	pdf.CellFormat(190, 5, "DESCRIPCION", "0", 0,
		"L", false, 0, "")
	pdf.SetX(116)
	pdf.CellFormat(190, 5, "UNIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(141)
	pdf.CellFormat(190, 5, "CANTIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(162)
	pdf.CellFormat(190, 5, "P. UNITARIO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(190)
	pdf.CellFormat(190, 5, "TOTAL", "0", 0,
		"L", false, 0, "")
	pdf.Ln(8)
}
func DevolucionfacturagastoFilaDetalle(pdf *gofpdf.Fpdf,miFila devolucionfacturagastodetalleeditar, a int ){
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
	var yinicial  float64
	yinicial=pdf.GetY()

	pdf.CellFormat(183, 4, strconv.Itoa(a), "", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 4, miFila.Codigoservicio, "", 0,
		"L", false, 0, "")
	var y float64
	y=pdf.GetY()

	pdf.SetX(42)
	pdf.MultiCell(80,4, Mayuscula(miFila.Nombreservicio), "","L", false)
	var yfinal float64
	yfinal=pdf.GetY()
	pdf.SetY(y)
	pdf.SetX(92)
	pdf.CellFormat(40, 4, Titulo(Subcadena(miFila.Unidadservicio, 0,6)), "", 0,
		"R", false, 0, "")
	pdf.SetX(121)
	pdf.CellFormat(40, 4, miFila.Cantidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(143)
	pdf.CellFormat(40, 4, miFila.Precio, "", 0,
		"R", false, 0, "")
	pdf.SetX(163)
	pdf.CellFormat(40, 4, miFila.Subtotal, "", 0,
		"R", false, 0, "")
	pdf.Ln(yfinal-yinicial+3)

}

func DevolucionfacturagastoPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miDevolucionfacturagasto devolucionfacturagasto ){

	Totalletras,err := IntLetra(Cadenaentero(miDevolucionfacturagasto.Neto))
	if err!= nil{
		fmt.Println(err)
	}

	pdf.SetFont("Arial", "", 8)
	pdf.SetY(226)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "SON: " +Mayuscula(Totalletras)+" PESOS MDA. CTE.", "0", 0,
		"L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetY(233)
	pdf.SetX(1)

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
	pdf.CellFormat(190, 10, "TOTAL IVA", "0", 0, "L",
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
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miDevolucionfacturagasto.Subtotal), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miDevolucionfacturagasto.TotalIva), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miDevolucionfacturagasto.TotalRetencionFuente), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miDevolucionfacturagasto.TotalRetencionIca), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miDevolucionfacturagasto.Neto), "0", 0, "R",
		false, 0, "")
}
