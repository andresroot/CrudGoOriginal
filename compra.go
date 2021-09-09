package main

// INICIA COMPRA IMPORTAR PAQUETES
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

// INICIA COMPRA ESTRUCTURA JSON
type compraJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA COMPRA ESTRUCTURA
type compraLista struct {
	Codigo			  string
	Fecha        	  time.Time
	Neto          	  string
	Tercero       	  string
	TerceroNombre	  string
	CentroNombre 	  string
	AlmacenistaNombre string
}

// INICIA COMPRA ESTRUCTURA
type compra struct {
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
    Detalle                   []compradetalle `json:"Detalle"`
	DetalleEditar			  []compradetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Pedido					  string
	Tipo					  string
	Centro					  string
}

// INICIA COMPRADETALLE ESTRUCTURA
type compradetalle struct {
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

// INICIA COMPRA DETALLE EDITARr
type compradetalleeditar struct {
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

// INICIA COMPRA CONSULTA DETALLE
func CompraConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "compradetalle.Id as id ,"
	consulta += "compradetalle.Codigo as codigo,"
	consulta += "compradetalle.Fila as fila,"
	consulta += "compradetalle.Cantidad as cantidad,"
	consulta += "compradetalle.Precio as precio,"
	consulta += "compradetalle.Descuento as descuento,"
	consulta += "compradetalle.Montodescuento as montodescuento,"
	consulta += "compradetalle.Sigratis as sigratis,"
	consulta += "compradetalle.Subtotal as subtotal,"
	consulta += "compradetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "compradetalle.Pagina as pagina ,"
	consulta += "compradetalle.Bodega as bodega,"
	consulta += "compradetalle.Producto as producto,"
	consulta += "compradetalle.Fecha as fecha,"
	consulta += "producto.nombre as ProductoNombre, "
	consulta += "producto.iva as ProductoIva, "
	consulta += "producto.unidad as ProductoUnidad "
	consulta += "from compradetalle "
	consulta += "inner join producto on producto.codigo=compradetalle.producto "
	consulta += " where compradetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA COMPRA LISTA
func CompraLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/compra/compraLista.html")
	log.Println("Error compra 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,centro.nombre as CentroNombre,compra.neto,compra.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM compra "
	consulta += " inner join tercero on tercero.codigo=compra.tercero "
	consulta += " inner join centro on centro.codigo=compra.centro "
	consulta += " inner join almacenista on almacenista.codigo=compra.almacenista "
	consulta += " ORDER BY compra.codigo ASC"

	db := dbConn()
	res := []compraLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error compra888")
	tmp.Execute(w, varmap)
}

// INICIA COMPRA NUEVO
func CompraNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio compra editar" + Codigo)

	db := dbConn()
	v := compra{}
	tc := tercero{}
	det := []compradetalleeditar{}
	if Codigo == "False" {

	} else {

		err := db.Get(&v, "SELECT * FROM compra where codigo=$1", Codigo)
		if err != nil {
			log.Fatalln(err)
		}

		err2 := db.Select(&det, CompraConsultaDetalle(), Codigo)
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
		"compra":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"retfte":TraerParametrosInventario().Compracuentaporcentajeretfte,
		"codigo": Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	t, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/compra/compraNuevo.html", "vista/compra/compraScript.html")
	fmt.Printf("%v, %v", tc, err)
	log.Println("Error compra nuevo 3")
	t.Execute(w, parametros)
}

// INICIA INSERTAR COMPROBANTE DE COMPRA
func InsertaDetalleComprobanteCompra(miFilaComprobante comprobantedetalle, miComprobante comprobante, miCompra compra){
	db := dbConn()

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miCompra.Tercero)
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
	miCompra.Tercero,
	miCompra.Centro,
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

// INICIA COMPRA INSERTAR AJAX
func CompraAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempCompra compra

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la COMPRA
	err = json.Unmarshal(b, &tempCompra)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	if tempCompra.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from compradetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempCompra.Codigo)

		// borra detalle inventario
		delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo='Compra'")
		if err != nil {
			panic(err.Error())
		}
		delForm2.Exec(tempCompra.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from compra WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempCompra.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempCompra.Detalle {
		var a = i
		var q string
		q = "insert into compradetalle ("
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

		// TERMINA COMPRA GRABAR INSERTAR
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
	for i, x := range tempCompra.Detalle {
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

		// TERMINA COMPRA GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Fecha,
			x.Tipo,
			x.Codigo,
			x.Bodega,
			x.Producto,
			x.Cantidad,
			x.Precio,
			operacionCompra)
		if err != nil {
			panic(err)
		}
		log.Println("Insertar Producto \n", x.Producto, a)
	}

	// INICIA INSERTAR COMPRAS
	log.Println("Got %s age %s club %s\n", tempCompra.Codigo, tempCompra.Tercero, tempCompra.Total)
	var q string
	q = "insert into compra ("
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
	q += "Pedido,"
	q += "Centro,"
	q += "Tipo"
	q += " ) values("
	q+=parametros(33)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempCompra.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempCompra.Codigo,
		tempCompra.Fecha.Format(layout),
		tempCompra.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempCompra.Descuento,
		tempCompra.Subtotaldescuento19,
		tempCompra.Subtotaldescuento5,
		tempCompra.Subtotaldescuento0,
		tempCompra.Subtotal,
		tempCompra.Subtotal19,
		tempCompra.Subtotal5,
		tempCompra.Subtotal0,
		tempCompra.Subtotaliva19,
		tempCompra.Subtotaliva5,
		tempCompra.Subtotaliva0,
		tempCompra.Subtotalbase19,
		tempCompra.Subtotalbase5,
		tempCompra.Subtotalbase0,
		tempCompra.TotalIva,
		tempCompra.Total,
		tempCompra.PorcentajeRetencionFuente,
		tempCompra.TotalRetencionFuente,
		tempCompra.PorcentajeRetencionIca,
		tempCompra.TotalRetencionIca,
		tempCompra.Neto,
		tempCompra.Items,
		tempCompra.Formadepago,
		tempCompra.Mediodepago,
		tempCompra.Tercero,
		tempCompra.Almacenista,
		tempCompra.Pedido,
		tempCompra.Centro,
		tempCompra.Tipo)

	if err != nil {
		panic(err)
	}

	// INSERTAR COMPROBANTE CONTABILIDAD
	var tempComprobante comprobante
	var tempComprobanteDetalle comprobantedetalle
	tempComprobante.Documento="8"
	tempComprobante.Numero=tempCompra.Codigo
	tempComprobante.Fecha =tempCompra.Fecha
	tempComprobante.Fechaconsignacion =tempCompra.Fecha
	tempComprobante.Debito = tempCompra.Neto + ".00"
	tempComprobante.Credito	= tempCompra.Neto + ".00"
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

	// INSERTAR CUENTA DEBITO COMPRAS 19%
	if (tempCompra.Subtotal19!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempCompra.Subtotal19)-Flotante(tempCompra.Subtotaldescuento19)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compracuenta19
		tempComprobanteDetalle.Debito = FormatoFlotanteComprobante(Flotante(tempCompra.Subtotal19)-Flotante(tempCompra.Subtotaldescuento19))
		tempComprobanteDetalle.Credito =""
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA DEBITO COMPRAS 5%
	if (tempCompra.Subtotal5!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempCompra.Subtotal5)-Flotante(tempCompra.Subtotaldescuento5)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compracuenta5
		tempComprobanteDetalle.Debito = FormatoFlotanteComprobante(Flotante(tempCompra.Subtotal5)-Flotante(tempCompra.Subtotaldescuento5))
		tempComprobanteDetalle.Credito =""
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA DEBITO COMPRAS 0%
	if (tempCompra.Subtotal0!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempCompra.Subtotal0)-Flotante(tempCompra.Subtotaldescuento0)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compracuenta0
		tempComprobanteDetalle.Debito = FormatoFlotanteComprobante(Flotante(tempCompra.Subtotal0)-Flotante(tempCompra.Subtotaldescuento0))
		tempComprobanteDetalle.Credito =""
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA DEBITO IVA 19%
	if (tempCompra.Subtotaliva19!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempCompra.Subtotaliva19)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compraiva19
		tempComprobanteDetalle.Debito = tempCompra.Subtotaliva19
		tempComprobanteDetalle.Credito =""
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA DEBITO IVA 5%
	if (tempCompra.Subtotaliva5!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalDebito+=Flotante(tempCompra.Subtotaliva5)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compraiva5
		tempComprobanteDetalle.Debito = tempCompra.Subtotaliva5
		tempComprobanteDetalle.Credito =""
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
		log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	}

	// INSERTAR CUENTA CREDITO DESCUENTO
	//if (tempCompra.Descuento!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalCredito+=Flotante(tempCompra.Descuento)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Compracuentadescuento
	//	tempComprobanteDetalle.Debito = ""
	//	tempComprobanteDetalle.Credito = tempCompra.Descuento
	//	InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
	//	log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	//}

	// INSERTAR CUENTA CREDITO RET. FTE.
	if (tempCompra.TotalRetencionFuente!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempCompra.TotalRetencionFuente)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compracuentaretfte
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempCompra.TotalRetencionFuente
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO RET. ICA.
	if (tempCompra.TotalRetencionIca!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempCompra.TotalRetencionIca)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compracuentaretica
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempCompra.TotalRetencionIca
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
		log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	}

	// INSERTAR CUENTA CREDITO PROVEEDOR
	if (tempCompra.Neto!="0")	{
		fila=fila+1
		tempComprobanteDetalle.Fila=strconv.Itoa(fila)
		totalCredito+=Flotante(tempCompra.Neto)
		tempComprobanteDetalle.Cuenta  = parametrosinventario.Compracuentaproveedor
		tempComprobanteDetalle.Debito = ""
		tempComprobanteDetalle.Credito = tempCompra.Neto
		InsertaDetalleComprobanteCompra(tempComprobanteDetalle,tempComprobante,tempCompra)
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

// INICIA COMPRA EXISTE
func CompraExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM compra  WHERE codigo=$1", Codigo)
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

// INICIA COMPRA EDITAR
func CompraEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio compra editar" + Codigo)

	db := dbConn()

	// traer COMPRA
	v := compra{}
	err := db.Get(&v, "SELECT * FROM compra where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []compradetalleeditar{}

	err2 := db.Select(&det, CompraConsultaDetalle(), Codigo)
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
		"compra":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"retfte":TraerParametrosInventario().Compracuentaporcentajeretfte,

	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/compra/compraEditar.html", "vista/compra/compraScript.html")
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error compra nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA COMPRA BORRAR
func CompraBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio compra editar" + Codigo)

	db := dbConn()

	// traer COMPRA
	v := compra{}
	err := db.Get(&v, "SELECT * FROM compra where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []compradetalleeditar{}
	err2 := db.Select(&det, CompraConsultaDetalle(), Codigo)
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
		"compra":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"retfte":TraerParametrosInventario().Compracuentaporcentajeretfte,
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/compra/compraBorrar.html", "vista/compra/compraScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error compra nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA COMPRA ELIMINAR
func CompraEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar COMPRA
	delForm, err := db.Prepare("DELETE from compra WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from compradetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borar detalle invenario
	Borrarinventario(codigo,"Compra")

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
	http.Redirect(w, r, "/CompraLista", 301)
}

// TRAER PEDIDO
func DatosPedido(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedido editar" + Codigo)
	db := dbConn()
	var res []pedido

	// traer PEDIDO
	v := pedido{}
	err := db.Get(&v, "SELECT * FROM pedido where codigo=$1", Codigo)
	var valida bool
	valida=true

	switch err {
	case nil:
		log.Printf("pedido existe: %+v\n", v)
	case sql.ErrNoRows:
		log.Println("pedido NO Existe")
		valida=false
	default:
		log.Printf("pedido error: %s\n", err)
	}
	det := []pedidodetalleeditar{}
	t := tercero{}

	// trae datos si existe pedido
	if valida==true {
		err2 := db.Select(&det, PedidoConsultaDetalle(), Codigo)
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
func CompraPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer COMPRA
	miCompra := compra{}
	err := db.Get(&miCompra, "SELECT * FROM compra where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []compradetalleeditar{}
	err2 := db.Select(&miDetalle, CompraConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miCompra.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miCompra.Almacenista)
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
		pdf.SetFont("Arial", "", 11)
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

		pdf.SetY(20)
		pdf.SetX(80)
		pdf.Ln(8)
		pdf.SetX(80)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(190, 10, "FACTURA DE COMPRA", "0", 0, "C",
			false, 0, "")
		pdf.Ln(6)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, " No.  "+Codigo, "0", 0, "C",
			false, 0, "")
		pdf.Ln(10)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(256)
		pdf.SetX(20)
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(80, 10, "www.Sadconf.com.co", "",
			0, "L", false, 0, "")
		pdf.SetX(130)
		pdf.CellFormat(78, 10, fmt.Sprintf(" %d de {nb}", pdf.PageNo()), "",
			0, "R", false, 0, "")

	})

	// inicia suma de paginas
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.Ln(1)
	CompraCabecera(pdf,miTercero,miCompra,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			CompraFilaDetalle(pdf,miFila,a)
		}
		CompraPieDePagina(pdf,miTercero,miCompra)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CompraFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			CompraCabecera(pdf,miTercero,miCompra,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					CompraFilaDetalle(pdf,miFila,a)
				}
			}

			CompraPieDePagina(pdf,miTercero,miCompra)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CompraFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CompraCabecera(pdf,miTercero,miCompra,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					CompraFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CompraCabecera(pdf,miTercero,miCompra,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					CompraFilaDetalle(pdf,miFila,a)
				}
			}

			CompraPieDePagina(pdf,miTercero,miCompra)
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
func CompraCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miCompra compra, miAlmacenista almacenista ){

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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miCompra.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miCompra.Fecha.Format("02/01/2006")+" "+Titulo(miCompra.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miCompra.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miCompra.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Pedido No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, miCompra.Pedido, "", 0,
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
func CompraFilaDetalle(pdf *gofpdf.Fpdf,miFila compradetalleeditar, a int ){
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

func CompraPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miCompra compra ){

	Totalletras,err := IntLetra(Cadenaentero(miCompra.Neto))
	if err!= nil{
		fmt.Println(err)
	}

	pdf.SetFont("Arial", "", 8)
	pdf.SetY(222)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "SON: " +Mayuscula(Totalletras)+" PESOS MDA. CTE.", "0", 0,
		"L", false, 0, "")

	pdf.Ln(20)
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
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miCompra.Subtotal, "0", 0, "R",
		false, 0, "")

	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miCompra.TotalIva, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miCompra.TotalRetencionFuente, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miCompra.TotalRetencionIca, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCompra.Neto, "0", 0, "R",
		false, 0, "")
}
