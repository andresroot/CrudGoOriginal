package main

// INICIA PEDIDO SOPORTE  SERVICIOI MPORTAR PAQUETES
import (
	"bytes"
	"math"

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
	"net/http"
	"strconv"
	"time"
)

// INICIA PEDIDO SOPORTE  SERVICIO ESTRUCTURA JSON
type pedidosoporteservicioJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA PEDIDO SOPORTE  SERVICIOESTRUCTURA
type pedidosoporteservicioLista struct {
	Codigo        string
	Fecha         time.Time
	Neto          string
	Tercero       string
	TerceroNombre string
	CentroNombre  string
	AlmacenistaNombre string
}

// INICIA PEDIDO SOPORTE  SERVICIOESTRUCTURA
type pedidosoporteservicio struct {
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
	Detalle                   []pedidosoporteserviciodetalle `json:"Detalle"`
	DetalleEditar			  []pedidosoporteserviciodetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Tipo      			  	  string
	Centro					  string
	Soporteserviciocuenta0	  string
	Soporteservicionombre0	  string
	Soporteserviciocuentaporcentajeretfte		      string
	Soporteserviciocuentaretfte	  string
	Soporteservicionombreretfte   string
}

// INICIA PEDIDO SOPORTEDETALLE ESTRUCTURA
type pedidosoporteserviciodetalle struct {
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
type pedidosoporteserviciodetalleeditar struct {
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
func PedidosoporteservicioConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "pedidosoporteserviciodetalle.Id as id ,"
	consulta += "pedidosoporteserviciodetalle.Codigo as codigo,"
	consulta += "pedidosoporteserviciodetalle.Fila as fila,"
	consulta += "pedidosoporteserviciodetalle.Cantidad as cantidad,"
	consulta += "pedidosoporteserviciodetalle.Precio as precio,"
	consulta += "pedidosoporteserviciodetalle.Descuento as descuento,"
	consulta += "pedidosoporteserviciodetalle.Montodescuento as montodescuento,"
	consulta += "pedidosoporteserviciodetalle.Sigratis as sigratis,"
	consulta += "pedidosoporteserviciodetalle.Subtotal as subtotal,"
	consulta += "pedidosoporteserviciodetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "pedidosoporteserviciodetalle.Pagina as pagina ,"
	consulta += "pedidosoporteserviciodetalle.Bodega as bodega,"
	consulta += "pedidosoporteserviciodetalle.Fecha as fecha,"
	consulta += "pedidosoporteserviciodetalle.Nombreservicio,"
	consulta += "pedidosoporteserviciodetalle.Unidadservicio,"
	consulta += "pedidosoporteserviciodetalle.Codigoservicio "
	consulta += "from pedidosoporteserviciodetalle "
	consulta += " where pedidosoporteserviciodetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA PEDIDO SOPORTE  SERVICIOLISTA
func PedidosoporteservicioLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/pedidosoporteservicio/pedidosoporteservicioLista.html",)
	log.Println("Error pedidosoporteservicio 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,centro.nombre as CentroNombre,pedidosoporteservicio.neto,pedidosoporteservicio.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM pedidosoporteservicio "
	consulta += " inner join tercero on tercero.codigo=pedidosoporteservicio.tercero "
	consulta += " inner join centro on centro.codigo=pedidosoporteservicio.centro "
	consulta += " inner join almacenista on almacenista.codigo=pedidosoporteservicio.almacenista "
	consulta += " ORDER BY pedidosoporteservicio.codigo ASC"

	db := dbConn()
	res := []pedidosoporteservicioLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error pedidosoporteservicio888")
	tmp.Execute(w, varmap)
}

// INICIA PEDIDO SOPORTE  SERVICION UEVO
func PedidosoporteservicioNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidosoporteservicio editar" + Codigo)

	db := dbConn()
	v := pedidosoporteservicio{}
	tc := tercero{}
	det := []pedidosoporteserviciodetalleeditar{}

	if Codigo == "False" {

	} else {

		err := db.Get(&v, "SELECT * FROM pedidosoporteservicio where codigo=$1", Codigo)
		if err != nil {
			log.Fatalln(err)
		}

		err2 := db.Select(&det, PedidosoporteservicioConsultaDetalle(), Codigo)
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
		"pedidosoporteservicio":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"codigo":     Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	t, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/pedidosoporteservicio/pedidosoporteservicioNuevo.html",
		"vista/pedidosoporteservicio/pedidosoporteservicioScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	log.Println("Error pedidosoporteservicio nuevo 3")
	t.Execute(w, parametros)
}

// INICIA PEDIDO SOPORTE  SERVICIOINSERTAR AJAX
func PedidosoporteservicioAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempPedidosoporteservicio pedidosoporteservicio

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la PEDIDO SOPORTE
	err = json.Unmarshal(b, &tempPedidosoporteservicio)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	if tempPedidosoporteservicio.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from pedidosoporteserviciodetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempPedidosoporteservicio.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from pedidosoporteservicio WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempPedidosoporteservicio.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempPedidosoporteservicio.Detalle {
		var a = i
		var q string
		q = "insert into pedidosoporteserviciodetalle ("
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

		// TERMINA PEDIDO SOPORTE  SERVICIOGRABAR INSERTAR
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

	// INICIA INSERTAR PEDIDO SOPORTE
	log.Println("Got %s age %s club %s\n", tempPedidosoporteservicio.Codigo, tempPedidosoporteservicio.Tercero, tempPedidosoporteservicio.Neto)
	var q string
	q = "insert into pedidosoporteservicio ("
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
	q += "Soporteserviciocuenta0,"
	q += "Soporteservicionombre0,"
	q += "Soporteserviciocuentaporcentajeretfte,"
	q += "Soporteserviciocuentaretfte,"
	q += "Soporteservicionombreretfte"
	q += " ) values("
	q+=parametros(37)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempPedidosoporteservicio.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempPedidosoporteservicio.Codigo,
		tempPedidosoporteservicio.Fecha.Format(layout),
		tempPedidosoporteservicio.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempPedidosoporteservicio.Descuento,
		tempPedidosoporteservicio.Subtotaldescuento19,
		tempPedidosoporteservicio.Subtotaldescuento5,
		tempPedidosoporteservicio.Subtotaldescuento0,
		tempPedidosoporteservicio.Subtotal,
		tempPedidosoporteservicio.Subtotal19,
		tempPedidosoporteservicio.Subtotal5,
		tempPedidosoporteservicio.Subtotal0,
		tempPedidosoporteservicio.Subtotaliva19,
		tempPedidosoporteservicio.Subtotaliva5,
		tempPedidosoporteservicio.Subtotaliva0,
		tempPedidosoporteservicio.Subtotalbase19,
		tempPedidosoporteservicio.Subtotalbase5,
		tempPedidosoporteservicio.Subtotalbase0,
		tempPedidosoporteservicio.TotalIva,
		tempPedidosoporteservicio.Total,
		tempPedidosoporteservicio.PorcentajeRetencionFuente,
		tempPedidosoporteservicio.TotalRetencionFuente,
		tempPedidosoporteservicio.PorcentajeRetencionIca,
		tempPedidosoporteservicio.TotalRetencionIca,
		tempPedidosoporteservicio.Neto,
		tempPedidosoporteservicio.Items,
		tempPedidosoporteservicio.Formadepago,
		tempPedidosoporteservicio.Mediodepago,
		tempPedidosoporteservicio.Tercero,
		tempPedidosoporteservicio.Almacenista,
		tempPedidosoporteservicio.Centro,
		tempPedidosoporteservicio.Tipo,
		tempPedidosoporteservicio.Soporteserviciocuenta0,
		tempPedidosoporteservicio.Soporteservicionombre0,
		tempPedidosoporteservicio.Soporteserviciocuentaporcentajeretfte,
		tempPedidosoporteservicio.Soporteserviciocuentaretfte,
		tempPedidosoporteservicio.Soporteservicionombreretfte)
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

// INICIA PEDIDO SOPORTE  SERVICIOEXISTE
func PedidosoporteservicioExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM pedidosoporteservicio  WHERE codigo=$1", Codigo)
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

// INICIA PEDIDO SOPORTE  SERVICIOEDITAR
func PedidosoporteservicioEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidosoporteservicio editar" + Codigo)

	db := dbConn()

	// traer PEDIDO SOPORTE
	v := pedidosoporteservicio{}
	err := db.Get(&v, "SELECT * FROM pedidosoporteservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []pedidosoporteserviciodetalleeditar{}

	err2 := db.Select(&det, PedidosoporteservicioConsultaDetalle(), Codigo)
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
		"pedidosoporteservicio":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),

	}
	log.Println("Error pedidosoporteservicio nuevo 555")
	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/pedidosoporteservicio/pedidosoporteservicioEditar.html",
		"vista/pedidosoporteservicio/pedidosoporteservicioScript.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa.html",)
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error pedidosoporteservicio nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA PEDIDO SOPORTE  SERVICIOBORRAR
func PedidosoporteservicioBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidosoporteservicio editar" + Codigo)

	db := dbConn()

	// traer PEDIDO SOPORTE
	v := pedidosoporteservicio{}
	err := db.Get(&v, "SELECT * FROM pedidosoporteservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []pedidosoporteserviciodetalleeditar{}
	err2 := db.Select(&det, PedidosoporteservicioConsultaDetalle(), Codigo)
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
		"pedidosoporteservicio":       v,
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
		"vista/pedidosoporteservicio/pedidosoporteservicioBorrar.html",
		"vista/pedidosoporteservicio/pedidosoporteservicioScript.html",)
	fmt.Printf("%v, %v", t, err)
	log.Println("Error pedidosoporteservicio nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA PEDIDO SOPORTE  SERVICIOELIMINAR
func PedidosoporteservicioEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar PEDIDO SOPORTE
	delForm, err := db.Prepare("DELETE from pedidosoporteservicio WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from pedidosoporteserviciodetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/PedidosoporteservicioLista", 301)
}

// INICIA PDF
func PedidosoporteservicioPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer PEDIDO SOPORTE
	miPedidosoporteservicio := pedidosoporteservicio{}
	err := db.Get(&miPedidosoporteservicio, "SELECT * FROM pedidosoporteservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []pedidosoporteserviciodetalleeditar{}
	err2 := db.Select(&miDetalle, PedidosoporteservicioConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miPedidosoporteservicio.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miPedidosoporteservicio.Almacenista)
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

		// EMPRESA
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

		pdf.SetFont("Arial", "", 10)
		pdf.SetY(20)
		pdf.SetX(80)
		pdf.Ln(8)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "PEDIDO SOPORTE SERVICIO", "0", 0, "C",
			false, 0, "")
		pdf.Ln(6)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, " No.  "+Codigo, "0", 0, "C",
			false, 0, "")
		pdf.Ln(10)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetFont("Arial", "", 8)
		pdf.SetY(256)
		pdf.SetX(20)
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
	PedidosoporteservicioCabecera(pdf,miTercero,miPedidosoporteservicio,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			PedidosoporteservicioFilaDetalle(pdf,miFila,a)
		}
		PedidosoporteservicioPieDePagina(pdf,miTercero,miPedidosoporteservicio)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					PedidosoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			PedidosoporteservicioCabecera(pdf,miTercero,miPedidosoporteservicio,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					PedidosoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			PedidosoporteservicioPieDePagina(pdf,miTercero,miPedidosoporteservicio)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					PedidosoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			PedidosoporteservicioCabecera(pdf,miTercero,miPedidosoporteservicio,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					PedidosoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			PedidosoporteservicioCabecera(pdf,miTercero,miPedidosoporteservicio,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					PedidosoporteservicioFilaDetalle(pdf,miFila,a)
				}
			}

			PedidosoporteservicioPieDePagina(pdf,miTercero,miPedidosoporteservicio)
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
func PedidosoporteservicioCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miPedidosoporteservicio pedidosoporteservicio, miAlmacenista almacenista ){
	pdf.SetFont("Arial", "", 10)
	// RELLENO TITULO
	pdf.SetY(46)
	pdf.SetX(20)
	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)

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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miPedidosoporteservicio.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miPedidosoporteservicio.Fecha.Format("02/01/2006")+" "+Titulo(miPedidosoporteservicio.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miPedidosoporteservicio.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miPedidosoporteservicio.Vence.Format("02/01/2006"), "", 0,
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
	pdf.SetX(135)
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
func PedidosoporteservicioFilaDetalle(pdf *gofpdf.Fpdf,miFila pedidosoporteserviciodetalleeditar, a int ){
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
	pdf.CellFormat(40, 4, Subcadena(miFila.Codigoservicio,0,12), "", 0,
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
	pdf.SetX(100)
	pdf.CellFormat(40, 4, miFila.Descuento, "", 0,
		"R", false, 0, "")
	pdf.SetX(115)
	pdf.CellFormat(40, 4, miFila.Cantidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(143)
	pdf.CellFormat(40, 4, miFila.Precio, "", 0,
		"R", false, 0, "")
	pdf.SetX(164)
	pdf.CellFormat(40, 4, miFila.Subtotal, "", 0,
		"R", false, 0, "")
	pdf.Ln(yfinal-yinicial+3)

}

func PedidosoporteservicioPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miPedidosoporteservicio pedidosoporteservicio ){

	Totalletras,err := IntLetra(Cadenaentero(miPedidosoporteservicio.Neto))
	if err!= nil{
		fmt.Println(err)
	}

	pdf.SetFont("Arial", "", 8)
	pdf.SetY(226)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "SON: " +Mayuscula(Totalletras)+" PESOS MDA. CTE.", "0", 0,
		"L", false, 0, "")

	pdf.Ln(14)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "__________________________________________", "0", 0,
		"C", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "FIRMA RESPONSABLE ", "0", 0, "C",
		false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetY(233)
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

	pdf.SetY(233)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miPedidosoporteservicio.Subtotal, "0", 0, "R",
		false, 0, "")

	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miPedidosoporteservicio.TotalRetencionFuente, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miPedidosoporteservicio.TotalRetencionIca, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miPedidosoporteservicio.Neto, "0", 0, "R",
		false, 0, "")
}
