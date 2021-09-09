package main

// INICIA PEDIDO FACTURA GASTO
import (
	"bytes"
	//"database/sql"
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

// INICIA PEDIDO FACTURA GASTO JSON
type pedidofacturagastoJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA PEDIDO FACTURA GASTO ESTRUCTURA
type pedidofacturagastoLista struct {
	Codigo            string
	Fecha             time.Time
	Neto              string
	Tercero           string
	TerceroNombre     string
	CentroNombre	  string
	AlmacenistaNombre string
}

// INICIA PEDIDO FACTURA GASTO
type pedidofacturagasto struct {
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
	Detalle                   []pedidofacturagastodetalle `json:"Detalle"`
	DetalleEditar			  []pedidofacturagastodetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Tipo      			  	  string
	Centro					  string
	Facturagastocuenta0	  string
	Facturagastonombre0	  string
	Facturagastocuentaporcentajeretfte		      string
	Facturagastocuentaretfte	  string
	Facturagastonombreretfte      string
	Facturagastocuentaiva	      string
	Facturagastonombreiva	      string
	Facturagastoporcentajeiva	  string
}

// INICIA PEDIDO FACTURA GASTODETALLE ESTRUCTURA
type pedidofacturagastodetalle struct {
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

// estructura para editar
type pedidofacturagastodetalleeditar struct {
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

// INICIA COMPRA CONSULTA DETALLE
func PedidofacturagastoConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "pedidofacturagastodetalle.Id as id ,"
	consulta += "pedidofacturagastodetalle.Codigo as codigo,"
	consulta += "pedidofacturagastodetalle.Fila as fila,"
	consulta += "pedidofacturagastodetalle.Cantidad as cantidad,"
	consulta += "pedidofacturagastodetalle.Precio as precio,"
	consulta += "pedidofacturagastodetalle.Descuento as descuento,"
	consulta += "pedidofacturagastodetalle.Montodescuento as montodescuento,"
	consulta += "pedidofacturagastodetalle.Sigratis as sigratis,"
	consulta += "pedidofacturagastodetalle.Subtotal as subtotal,"
	consulta += "pedidofacturagastodetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "pedidofacturagastodetalle.Pagina as pagina ,"
	consulta += "pedidofacturagastodetalle.Bodega as bodega,"
	consulta += "pedidofacturagastodetalle.Fecha as fecha,"
	consulta += "pedidofacturagastodetalle.Nombreservicio,"
	consulta += "pedidofacturagastodetalle.Unidadservicio,"
	consulta += "pedidofacturagastodetalle.Codigoservicio "
	consulta += "from pedidofacturagastodetalle "
	consulta += " where pedidofacturagastodetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA PEDIDO FACTURA GASTO  SERVICIOLISTA
func PedidofacturagastoLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/pedidofacturagasto/pedidofacturagastoLista.html",)
	log.Println("Error pedidofacturagasto 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,centro.nombre as CentroNombre,pedidofacturagasto.neto,pedidofacturagasto.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM pedidofacturagasto "
	consulta += " inner join tercero on tercero.codigo=pedidofacturagasto.tercero "
	consulta += " inner join centro on centro.codigo=pedidofacturagasto.centro "
	consulta += " inner join almacenista on almacenista.codigo=pedidofacturagasto.almacenista "
	consulta += " ORDER BY pedidofacturagasto.codigo ASC"

	db := dbConn()
	res := []pedidofacturagastoLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error pedidofacturagasto888")
	tmp.Execute(w, varmap)
}

// INICIA PEDIDO FACTURA GASTO  SERVICION UEVO
func PedidofacturagastoNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error pedidofacturagasto nuevo 1")
	log.Println("Error pedidofacturagasto nuevo 2")
	parametros := map[string]interface{}{
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro":      ListaCentro(),
	}

	t, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/pedidofacturagasto/pedidofacturagastoNuevo.html",
		"vista/pedidofacturagasto/pedidofacturagastoScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	log.Println("Error pedidofacturagasto nuevo 3")
	t.Execute(w, parametros)
}

// INICIA PEDIDO FACTURA GASTO  SERVICIOINSERTAR AJAX
func PedidofacturagastoAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempPedidofacturagasto pedidofacturagasto

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la PEDIDO FACTURA GASTO
	err = json.Unmarshal(b, &tempPedidofacturagasto)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	if tempPedidofacturagasto.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from pedidofacturagastodetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempPedidofacturagasto.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from pedidofacturagasto WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempPedidofacturagasto.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempPedidofacturagasto.Detalle {
		var a = i
		var q string
		q = "insert into pedidofacturagastodetalle ("
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

		// TERMINA PEDIDO FACTURA GASTO  SERVICIOGRABAR INSERTAR
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

	// INICIA INSERTAR PEDIDO FACTURA GASTO
	log.Println("Got %s age %s club %s\n", tempPedidofacturagasto.Codigo, tempPedidofacturagasto.Tercero, tempPedidofacturagasto.Neto)
	var q string
	q = "insert into pedidofacturagasto ("
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
	q+=parametros(40)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempPedidofacturagasto.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempPedidofacturagasto.Codigo,
		tempPedidofacturagasto.Fecha.Format(layout),
		tempPedidofacturagasto.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempPedidofacturagasto.Descuento,
		tempPedidofacturagasto.Subtotaldescuento19,
		tempPedidofacturagasto.Subtotaldescuento5,
		tempPedidofacturagasto.Subtotaldescuento0,
		tempPedidofacturagasto.Subtotal,
		tempPedidofacturagasto.Subtotal19,
		tempPedidofacturagasto.Subtotal5,
		tempPedidofacturagasto.Subtotal0,
		tempPedidofacturagasto.Subtotaliva19,
		tempPedidofacturagasto.Subtotaliva5,
		tempPedidofacturagasto.Subtotaliva0,
		tempPedidofacturagasto.Subtotalbase19,
		tempPedidofacturagasto.Subtotalbase5,
		tempPedidofacturagasto.Subtotalbase0,
		tempPedidofacturagasto.TotalIva,
		tempPedidofacturagasto.Total,
		tempPedidofacturagasto.PorcentajeRetencionFuente,
		tempPedidofacturagasto.TotalRetencionFuente,
		tempPedidofacturagasto.PorcentajeRetencionIca,
		tempPedidofacturagasto.TotalRetencionIca,
		tempPedidofacturagasto.Neto,
		tempPedidofacturagasto.Items,
		tempPedidofacturagasto.Formadepago,
		tempPedidofacturagasto.Mediodepago,
		tempPedidofacturagasto.Tercero,
		tempPedidofacturagasto.Almacenista,
		tempPedidofacturagasto.Centro,
		tempPedidofacturagasto.Tipo,
		tempPedidofacturagasto.Facturagastocuenta0,
		tempPedidofacturagasto.Facturagastonombre0,
		tempPedidofacturagasto.Facturagastocuentaporcentajeretfte,
		tempPedidofacturagasto.Facturagastocuentaretfte,
		tempPedidofacturagasto.Facturagastonombreretfte,
		tempPedidofacturagasto.Facturagastocuentaiva,
		tempPedidofacturagasto.Facturagastonombreiva,
		tempPedidofacturagasto.Facturagastoporcentajeiva)

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

// INICIA PEDIDO FACTURA GASTO EXISTE
func PedidofacturagastoExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM pedidofacturagasto  WHERE codigo=$1", Codigo)
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

// INICIA PEDIDO FACTURA GASTO  SERVICIOEDITAR
func PedidofacturagastoEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidofacturagasto editar" + Codigo)

	db := dbConn()

	// traer PEDIDO FACTURA GASTO
	v := pedidofacturagasto{}
	err := db.Get(&v, "SELECT * FROM pedidofacturagasto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []pedidofacturagastodetalleeditar{}

	err2 := db.Select(&det, PedidofacturagastoConsultaDetalle(), Codigo)
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
		"pedidofacturagasto":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),

	}
	log.Println("Error pedidofacturagasto nuevo 555")
	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/pedidofacturagasto/pedidofacturagastoEditar.html",
		"vista/pedidofacturagasto/pedidofacturagastoScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error pedidofacturagasto nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA PEDIDO FACTURA GASTO BORRAR
func PedidofacturagastoBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidofacturagasto editar" + Codigo)

	db := dbConn()

	// traer PEDIDO FACTURA GASTO
	v := pedidofacturagasto{}
	err := db.Get(&v, "SELECT * FROM pedidofacturagasto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []pedidofacturagastodetalleeditar{}
	err2 := db.Select(&det, PedidofacturagastoConsultaDetalle(), Codigo)
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
		"pedidofacturagasto":       v,
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
		"vista/pedidofacturagasto/pedidofacturagastoBorrar.html",
		"vista/pedidofacturagasto/pedidofacturagastoScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	log.Println("Error pedidofacturagasto nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA PEDIDO FACTURA GASTO  SERVICIOELIMINAR
func PedidofacturagastoEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar PEDIDO FACTURA GASTO
	delForm, err := db.Prepare("DELETE from pedidofacturagasto WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from pedidofacturagastodetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/PedidofacturagastoLista", 301)
}

// INICIA PDF
func PedidofacturagastoPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer PEDIDO FACTURA GASTO
	miPedidofacturagasto := pedidofacturagasto{}
	err := db.Get(&miPedidofacturagasto, "SELECT * FROM pedidofacturagasto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []pedidofacturagastodetalleeditar{}
	err2 := db.Select(&miDetalle, PedidofacturagastoConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miPedidofacturagasto.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miPedidofacturagasto.Almacenista)
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

		// FACTURA GASTO NUMERO
		pdf.SetY(20)
		pdf.SetX(80)
		pdf.Ln(8)
		pdf.SetX(75)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(190, 10, "PEDIDO FACTURA GASTO", "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
		pdf.SetX(75)
		pdf.CellFormat(190, 10, " No.  "+Codigo, "0", 0, "C",
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
	PedidofacturagastoCabecera(pdf,miTercero,miPedidofacturagasto,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			PedidofacturagastoFilaDetalle(pdf,miFila,a)
		}
		PedidofacturagastoPieDePagina(pdf,miTercero,miPedidofacturagasto)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					PedidofacturagastoFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			PedidofacturagastoCabecera(pdf,miTercero,miPedidofacturagasto,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					PedidofacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			PedidofacturagastoPieDePagina(pdf,miTercero,miPedidofacturagasto)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					PedidofacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			PedidofacturagastoCabecera(pdf,miTercero,miPedidofacturagasto,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					PedidofacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			PedidofacturagastoCabecera(pdf,miTercero,miPedidofacturagasto,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					PedidofacturagastoFilaDetalle(pdf,miFila,a)
				}
			}

			PedidofacturagastoPieDePagina(pdf,miTercero,miPedidofacturagasto)
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
func PedidofacturagastoCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miPedidofacturagasto pedidofacturagasto, miAlmacenista almacenista ){

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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miPedidofacturagasto.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miPedidofacturagasto.Fecha.Format("02/01/2006")+" "+Titulo(miPedidofacturagasto.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miPedidofacturagasto.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miPedidofacturagasto.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Condiciones", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, "", "", 0,
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
func PedidofacturagastoFilaDetalle(pdf *gofpdf.Fpdf,miFila pedidofacturagastodetalleeditar, a int ){
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

func PedidofacturagastoPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miPedidofacturagasto pedidofacturagasto ){

	Totalletras,err := IntLetra(Cadenaentero(miPedidofacturagasto.Neto))
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
	pdf.CellFormat(190, 10, Coma(miPedidofacturagasto.Subtotal), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miPedidofacturagasto.TotalIva), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miPedidofacturagasto.TotalRetencionFuente), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miPedidofacturagasto.TotalRetencionIca), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(12)
	pdf.CellFormat(190, 10, Coma(miPedidofacturagasto.Neto), "0", 0, "R",
		false, 0, "")
}
