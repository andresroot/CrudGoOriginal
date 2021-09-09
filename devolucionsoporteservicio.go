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
type devolucionsoporteservicioJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA SOPORTE ESTRUCTURA
type devolucionsoporteservicioLista struct {
	Codigo			  string
	Fecha        	  time.Time
	Neto          	  string
	Tercero       	  string
	TerceroNombre	  string
	CentroNombre 	  string
	AlmacenistaNombre string
}

// INICIA SOPORTE ESTRUCTURA
type devolucionsoporteservicio struct {
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
	Detalle                   []devolucionsoporteserviciodetalle `json:"Detalle"`
	DetalleEditar			  []devolucionsoporteserviciodetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Soporteservicio			  string
	Tipo					  string
	Centro					  string
	Soporteserviciocuenta0	  string
	Soporteservicionombre0	  string
	Soporteserviciocuentaporcentajeretfte	      string
	Soporteserviciocuentaretfte	  string
	Soporteservicionombreretfte   string
}

// INICIA SOPORTEDETALLE ESTRUCTURA
type devolucionsoporteserviciodetalle struct {
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

// INICIA SOPORTE DETALLE EDITARr
type devolucionsoporteserviciodetalleeditar struct {
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

// INICIA SOPORTE CONSULTA DETALLE
func DevolucionsoporteservicioConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "devolucionsoporteserviciodetalle.Id as id ,"
	consulta += "devolucionsoporteserviciodetalle.Codigo as codigo,"
	consulta += "devolucionsoporteserviciodetalle.Fila as fila,"
	consulta += "devolucionsoporteserviciodetalle.Cantidad as cantidad,"
	consulta += "devolucionsoporteserviciodetalle.Precio as precio,"
	consulta += "devolucionsoporteserviciodetalle.Descuento as descuento,"
	consulta += "devolucionsoporteserviciodetalle.Montodescuento as montodescuento,"
	consulta += "devolucionsoporteserviciodetalle.Sigratis as sigratis,"
	consulta += "devolucionsoporteserviciodetalle.Subtotal as subtotal,"
	consulta += "devolucionsoporteserviciodetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "devolucionsoporteserviciodetalle.Pagina as pagina ,"
	consulta += "devolucionsoporteserviciodetalle.Bodega as bodega,"
	consulta += "devolucionsoporteserviciodetalle.Fecha as fecha,"
	consulta += "devolucionsoporteserviciodetalle.Nombreservicio,"
	consulta += "devolucionsoporteserviciodetalle.Unidadservicio,"
	consulta += "devolucionsoporteserviciodetalle.Codigoservicio "
	consulta += "from devolucionsoporteserviciodetalle"
	consulta += " where devolucionsoporteserviciodetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA SOPORTE LISTA
func DevolucionsoporteservicioLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionsoporteservicio/devolucionsoporteservicioLista.html")
	log.Println("Error devolucionsoporteservicio 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,centro.nombre as CentroNombre,devolucionsoporteservicio.neto,devolucionsoporteservicio.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM devolucionsoporteservicio "
	consulta += " inner join tercero on tercero.codigo=devolucionsoporteservicio.tercero "
	consulta += " inner join centro on centro.codigo=devolucionsoporteservicio.centro "
	consulta += " inner join almacenista on almacenista.codigo=devolucionsoporteservicio.almacenista "
	consulta += " ORDER BY devolucionsoporteservicio.codigo ASC"

	db := dbConn()
	res := []devolucionsoporteservicioLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error devolucionsoporteservicio888")
	tmp.Execute(w, varmap)
}

// INICIA SOPORTE NUEVO
func DevolucionsoporteservicioNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionsoporteservicio editar" + Codigo)

	db := dbConn()
	v := devolucionsoporteservicio{}
	tc := tercero{}
	det := []devolucionsoporteserviciodetalleeditar{}
	if Codigo == "False" {

	} else {

	err := db.Get(&v, "SELECT * FROM devolucionsoporteservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	err2 := db.Select(&det, DevolucionsoporteservicioConsultaDetalle(), Codigo)
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
		"devolucionsoporteservicio":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"resolucionsoporte" : ListaResolucionsoporte(),
		"codigo": Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	t, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionsoporteservicio/devolucionsoporteservicioNuevo.html",
		"vista/devolucionsoporteservicio/devolucionsoporteservicioScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	log.Println("Error devolucionsoporteservicio nuevo 3")
	t.Execute(w, parametros)
}

// INICIA INSERTAR COMPROBANTE DE SOPORTE
func InsertaDetalleComprobanteDevolucionsoporteservicio(miFilaComprobante comprobantedetalle, miComprobante comprobante, miDevolucionsoporteservicio devolucionsoporteservicio){
	db := dbConn()

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miDevolucionsoporteservicio.Tercero)
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
	miDevolucionsoporteservicio.Tercero,
	miDevolucionsoporteservicio.Centro,
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
func DevolucionsoporteservicioAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempDevolucionsoporteservicio devolucionsoporteservicio

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la SOPORTE
	err = json.Unmarshal(b, &tempDevolucionsoporteservicio)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	if tempDevolucionsoporteservicio.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from devolucionsoporteserviciodetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempDevolucionsoporteservicio.Codigo)

		// borra detalle inventario
		delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo='Devolucionsoporteservicio'")
		if err != nil {
			panic(err.Error())
		}
		delForm2.Exec(tempDevolucionsoporteservicio.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from devolucionsoporteservicio WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempDevolucionsoporteservicio.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempDevolucionsoporteservicio.Detalle {
		var a = i
		var q string
		q = "insert into devolucionsoporteserviciodetalle ("
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

	// INICIA INSERTAR SOPORTES
	log.Println("Got %s age %s club %s\n", tempDevolucionsoporteservicio.Codigo, tempDevolucionsoporteservicio.Tercero, tempDevolucionsoporteservicio.Subtotal)
	var q string
	q = "insert into devolucionsoporteservicio ("
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
	q += "Soporteservicio,"
	q += "Centro,"
	q += "Tipo,"
	q += "Soporteserviciocuenta0,"
	q += "Soporteservicionombre0,"
	q += "Soporteserviciocuentaporcentajeretfte,"
	q += "Soporteserviciocuentaretfte,"
	q += "Soporteservicionombreretfte"
	q += " ) values("
	q+=parametros(39)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempDevolucionsoporteservicio.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempDevolucionsoporteservicio.Resolucionsoporte,
		tempDevolucionsoporteservicio.Codigo,
		tempDevolucionsoporteservicio.Fecha.Format(layout),
		tempDevolucionsoporteservicio.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempDevolucionsoporteservicio.Descuento,
		tempDevolucionsoporteservicio.Subtotaldescuento19,
		tempDevolucionsoporteservicio.Subtotaldescuento5,
		tempDevolucionsoporteservicio.Subtotaldescuento0,
		tempDevolucionsoporteservicio.Subtotal,
		tempDevolucionsoporteservicio.Subtotal19,
		tempDevolucionsoporteservicio.Subtotal5,
		tempDevolucionsoporteservicio.Subtotal0,
		tempDevolucionsoporteservicio.Subtotaliva19,
		tempDevolucionsoporteservicio.Subtotaliva5,
		tempDevolucionsoporteservicio.Subtotaliva0,
		tempDevolucionsoporteservicio.Subtotalbase19,
		tempDevolucionsoporteservicio.Subtotalbase5,
		tempDevolucionsoporteservicio.Subtotalbase0,
		tempDevolucionsoporteservicio.TotalIva,
		tempDevolucionsoporteservicio.Total,
		tempDevolucionsoporteservicio.PorcentajeRetencionFuente,
		tempDevolucionsoporteservicio.TotalRetencionFuente,
		tempDevolucionsoporteservicio.PorcentajeRetencionIca,
		tempDevolucionsoporteservicio.TotalRetencionIca,
		tempDevolucionsoporteservicio.Neto,
		tempDevolucionsoporteservicio.Items,
		tempDevolucionsoporteservicio.Formadepago,
		tempDevolucionsoporteservicio.Mediodepago,
		tempDevolucionsoporteservicio.Tercero,
		tempDevolucionsoporteservicio.Almacenista,
		tempDevolucionsoporteservicio.Soporteservicio,
		tempDevolucionsoporteservicio.Centro,
		tempDevolucionsoporteservicio.Tipo,
		tempDevolucionsoporteservicio.Soporteserviciocuenta0,
		tempDevolucionsoporteservicio.Soporteservicionombre0,
		tempDevolucionsoporteservicio.Soporteserviciocuentaporcentajeretfte,
		tempDevolucionsoporteservicio.Soporteserviciocuentaretfte,
		tempDevolucionsoporteservicio.Soporteservicionombreretfte)
	if err != nil {
		panic(err)
	}

	// INSERTAR COMPROBANTE CONTABILIDAD
	var tempComprobante comprobante
	var tempComprobanteDetalle comprobantedetalle
	tempComprobante.Documento="24"
	tempComprobante.Numero=tempDevolucionsoporteservicio.Codigo
	tempComprobante.Fecha =tempDevolucionsoporteservicio.Fecha
	tempComprobante.Fechaconsignacion =tempDevolucionsoporteservicio.Fecha
	tempComprobante.Debito = tempDevolucionsoporteservicio.Neto + ".00"
	tempComprobante.Credito	= tempDevolucionsoporteservicio.Neto + ".00"
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
	if (tempDevolucionsoporteservicio.Neto!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionsoporteservicio.Neto)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soporteserviciocuentaproveedor
		tempComprobanteDetalle.Debito = tempDevolucionsoporteservicio.Neto
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionsoporteservicio(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporteservicio)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO DESCUENTO
	//if (tempDevolucionsoporteservicio.Descuento!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalDebito+=Flotante(tempDevolucionsoporteservicio.Descuento)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Soporteserviciodevolucioncuentadescuento
	//	tempComprobanteDetalle.Debito = tempDevolucionsoporteservicio.Descuento
	//	tempComprobanteDetalle.Credito = ""
	//	InsertaDetalleComprobanteDevolucionsoporteservicio(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporteservicio)
	//	log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	//}

	// INSERTAR CUENTA CREDITO RET. FTE.
	if (tempDevolucionsoporteservicio.TotalRetencionFuente!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionsoporteservicio.TotalRetencionFuente)
		tempComprobanteDetalle.Cuenta  = tempDevolucionsoporteservicio.Soporteserviciocuentaretfte
		tempComprobanteDetalle.Debito = tempDevolucionsoporteservicio.TotalRetencionFuente
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionsoporteservicio(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporteservicio)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO RET. ICA.
	if (tempDevolucionsoporteservicio.TotalRetencionIca!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempDevolucionsoporteservicio.TotalRetencionIca)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Soporteserviciodevolucioncuentaretica
		tempComprobanteDetalle.Debito = tempDevolucionsoporteservicio.TotalRetencionIca
		tempComprobanteDetalle.Credito = ""
		InsertaDetalleComprobanteDevolucionsoporteservicio(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporteservicio)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA DEBITO SOPORTE 0%
	if (tempDevolucionsoporteservicio.Subtotal0!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempDevolucionsoporteservicio.Subtotal0)-Flotante(tempDevolucionsoporteservicio.Subtotaldescuento0)
		tempComprobanteDetalle.Cuenta  = tempDevolucionsoporteservicio.Soporteserviciocuenta0
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = FormatoFlotanteComprobante(Flotante(tempDevolucionsoporteservicio.Subtotal0)-Flotante(tempDevolucionsoporteservicio.Subtotaldescuento0))
		InsertaDetalleComprobanteDevolucionsoporteservicio(tempComprobanteDetalle,tempComprobante,tempDevolucionsoporteservicio)
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

// INICIA SOPORTE EXISTE
func DevolucionsoporteservicioExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM devolucionsoporteservicio  WHERE codigo=$1", Codigo)
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
func DevolucionsoporteservicioEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionsoporteservicio editar" + Codigo)

	db := dbConn()

	// traer devolucion soporte
	v := devolucionsoporteservicio{}
	err := db.Get(&v, "SELECT * FROM devolucionsoporteservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []devolucionsoporteserviciodetalleeditar{}

	err2 := db.Select(&det, DevolucionsoporteservicioConsultaDetalle(), Codigo)
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
		"devolucionsoporteservicio":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"resolucionsoporte" : ListaResolucionsoporte(),

	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionsoporteservicio/devolucionsoporteservicioEditar.html",
		"vista/devolucionsoporteservicio/devolucionsoporteservicioScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error devolucionsoporteservicio nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA SOPORTE BORRAR
func DevolucionsoporteservicioBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio devolucionsoporteservicio editar" + Codigo)

	db := dbConn()

	// traer SOPORTE
	v := devolucionsoporteservicio{}
	err := db.Get(&v, "SELECT * FROM devolucionsoporteservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []devolucionsoporteserviciodetalleeditar{}
	err2 := db.Select(&det, DevolucionsoporteservicioConsultaDetalle(), Codigo)
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
		"devolucionsoporteservicio":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro":ListaCentro(),
		"resolucionsoporte" : ListaResolucionsoporte(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/devolucionsoporteservicio/devolucionsoporteservicioBorrar.html", "vista/devolucionsoporteservicio/devolucionsoporteservicioScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error devolucionsoporteservicio nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA SOPORTE ELIMINAR
func DevolucionsoporteservicioEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar SOPORTE
	delForm, err := db.Prepare("DELETE from devolucionsoporteservicio WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from devolucionsoporteserviciodetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borra detalle anterior
	delForm, err = db.Prepare("DELETE from comprobantedetalle WHERE documento=$1 and numero=$2")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec("24", codigo)

	// borra cabecera anterior

	delForm1, err = db.Prepare("DELETE from comprobante WHERE documento=$1 and numero=$2 ")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec("24", codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/DevolucionsoporteservicioLista", 301)
}

// TRAER SOPORTE SERVICIO
func Datossoporteservicio(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedido editar" + Codigo)
	db := dbConn()
	var res []soporteservicio

	// traer SOPORTE
	v := soporteservicio{}
	err := db.Get(&v, "SELECT * FROM soporteservicio where codigo=$1", Codigo)
	var valida bool
	valida=true

	switch err {
	case nil:
		log.Printf("soporte existe: %+v\n", v)
	case sql.ErrNoRows:
		log.Println("soporteservicio NO Existe")
		valida=false
	default:
		log.Printf("soporteservicio error: %s\n", err)
	}
	det := []soporteserviciodetalleeditar{}
	t := tercero{}

	// trae datos si existe pedido
	if valida==true {
		err2 := db.Select(&det, SoporteservicioConsultaDetalle(), Codigo)
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
func DevolucionsoporteservicioPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer SOPORTE
	miDevolucionsoporteservicio := devolucionsoporteservicio{}
	err := db.Get(&miDevolucionsoporteservicio, "SELECT * FROM devolucionsoporteservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []devolucionsoporteserviciodetalleeditar{}
	err2 := db.Select(&miDetalle, DevolucionsoporteservicioConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miDevolucionsoporteservicio.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miDevolucionsoporteservicio.Almacenista)
	if err4 != nil {
		log.Fatalln(err4)
	}

	var e empresa = ListaEmpresa()
	var c ciudad = TraerCiudad(e.Ciudad)
	var re Resolucionsoporte = TraerResolucionsoporte(miDevolucionsoporteservicio.Resolucionsoporte)
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

		// DEVOLUCION NUMERO
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
	DevolucionsoporteservicioCabecera(pdf,miTercero,miDevolucionsoporteservicio,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			DevolucionsoporteservicioFilaDetalle(pdf,miFila,a)
		}
		DevolucionsoporteservicioPieDePagina(pdf,miTercero,miDevolucionsoporteservicio)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					DevolucionsoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			DevolucionsoporteservicioCabecera(pdf,miTercero,miDevolucionsoporteservicio,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					DevolucionsoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			DevolucionsoporteservicioPieDePagina(pdf,miTercero,miDevolucionsoporteservicio)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					DevolucionsoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			DevolucionsoporteservicioCabecera(pdf,miTercero,miDevolucionsoporteservicio,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					DevolucionsoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			DevolucionsoporteservicioCabecera(pdf,miTercero,miDevolucionsoporteservicio,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					DevolucionsoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			DevolucionsoporteservicioPieDePagina(pdf,miTercero,miDevolucionsoporteservicio)
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
func DevolucionsoporteservicioCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miDevolucionsoporteservicio devolucionsoporteservicio, miAlmacenista almacenista ){

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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miDevolucionsoporteservicio.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miDevolucionsoporteservicio.Fecha.Format("02/01/2006")+" "+Titulo(miDevolucionsoporteservicio.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miDevolucionsoporteservicio.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miDevolucionsoporteservicio.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Soporte No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, miDevolucionsoporteservicio.Soporteservicio, "", 0,
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
func DevolucionsoporteservicioFilaDetalle(pdf *gofpdf.Fpdf,miFila devolucionsoporteserviciodetalleeditar, a int ){
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
	pdf.SetX(165)
	pdf.CellFormat(40, 4, miFila.Subtotal, "", 0,
		"R", false, 0, "")
	pdf.Ln(yfinal-yinicial+3)

}

func DevolucionsoporteservicioPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miDevolucionsoporteservicio devolucionsoporteservicio ){

	Totalletras,err := IntLetra(Cadenaentero(miDevolucionsoporteservicio.Neto))
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
	pdf.CellFormat(190, 10, miDevolucionsoporteservicio.Subtotal, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miDevolucionsoporteservicio.TotalRetencionFuente, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miDevolucionsoporteservicio.TotalRetencionIca, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miDevolucionsoporteservicio.Neto, "0", 0, "R",
		false, 0, "")

	pdf.Image(imageFile("QR.jpg"), 20, 229, 25, 0, false,
		"", 0, "")
	pdf.SetFont("Arial", "", 8)

	pdf.SetY(249)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "Cufexxxxxxxxxx", "", 0,
		"L", false, 0, "")
}
