package main

// INICIA PEDIDO SOPORTE IMPORTAR PAQUETES
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

// INICIA PEDIDO SOPORTE ESTRUCTURA JSON
type pedidosoporteJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA PEDIDO SOPORTE ESTRUCTURA
type pedidosoporteLista struct {
	Codigo        string
	Fecha         time.Time
	Neto          string
	Tercero       string
	TerceroNombre string
	CentroNombre  string
	AlmacenistaNombre string
}

// INICIA PEDIDO SOPORTE ESTRUCTURA
type pedidosoporte struct {
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
	Detalle                   []pedidosoportedetalle `json:"Detalle"`
	DetalleEditar			  []pedidosoportedetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Tipo      			  	  string
	Centro					  string
}

// INICIA PEDIDO SOPORTEDETALLE ESTRUCTURA
type pedidosoportedetalle struct {
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

// estructura para editar
type pedidosoportedetalleeditar struct {
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
func PedidosoporteConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "pedidosoportedetalle.Id as id ,"
	consulta += "pedidosoportedetalle.Codigo as codigo,"
	consulta += "pedidosoportedetalle.Fila as fila,"
	consulta += "pedidosoportedetalle.Cantidad as cantidad,"
	consulta += "pedidosoportedetalle.Precio as precio,"
	consulta += "pedidosoportedetalle.Descuento as descuento,"
	consulta += "pedidosoportedetalle.Montodescuento as montodescuento,"
	consulta += "pedidosoportedetalle.Sigratis as sigratis,"
	consulta += "pedidosoportedetalle.Subtotal as subtotal,"
	consulta += "pedidosoportedetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "pedidosoportedetalle.Pagina as pagina ,"
	consulta += "pedidosoportedetalle.Bodega as bodega,"
	consulta += "pedidosoportedetalle.Producto as producto,"
	consulta += "pedidosoportedetalle.Fecha as fecha,"
	consulta += "producto.nombre as ProductoNombre, "
	consulta += "producto.iva as ProductoIva, "
	consulta += "producto.unidad as ProductoUnidad "
	consulta += "from pedidosoportedetalle "
	consulta += "inner join producto on producto.codigo=pedidosoportedetalle.producto "
	consulta += " where pedidosoportedetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA PEDIDO SOPORTE LISTA
func PedidosoporteLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/pedidosoporte/pedidosoporteLista.html")
	log.Println("Error pedidosoporte 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,centro.nombre as CentroNombre,pedidosoporte.neto,pedidosoporte.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM pedidosoporte "
	consulta += " inner join tercero on tercero.codigo=pedidosoporte.tercero "
	consulta += " inner join centro on centro.codigo=pedidosoporte.centro "
	consulta += " inner join almacenista on almacenista.codigo=pedidosoporte.almacenista "
	consulta += " ORDER BY pedidosoporte.codigo ASC"

	db := dbConn()
	res := []pedidosoporteLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error pedidosoporte888")
	tmp.Execute(w, varmap)
}

// INICIA PEDIDO SOPORTE NUEVO
func PedidosoporteNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidosoporte editar" + Codigo)

	db := dbConn()
	v := pedidosoporte{}
	tc := tercero{}
	det := []pedidosoportedetalleeditar{}
	if Codigo == "False" {

	} else {

	err := db.Get(&v, "SELECT * FROM pedidosoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	err2 := db.Select(&det, PedidosoporteConsultaDetalle(), Codigo)
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
		"pedidosoporte":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" : ListaCentro(),
		"retfte":TraerParametrosInventario().Compracuentaporcentajeretfte,
		"codigo":     Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	t, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/pedidosoporte/pedidosoporteNuevo.html", "vista/pedidosoporte/pedidosoporteScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error pedidosoporte nuevo 3")
	t.Execute(w, parametros)
}

// INICIA PEDIDO SOPORTE INSERTAR AJAX
func PedidosoporteAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempPedidosoporte pedidosoporte

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la PEDIDO SOPORTE
	err = json.Unmarshal(b, &tempPedidosoporte)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	if tempPedidosoporte.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from pedidosoportedetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempPedidosoporte.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from pedidosoporte WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempPedidosoporte.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempPedidosoporte.Detalle {
		var a = i
		var q string
		q = "insert into pedidosoportedetalle ("
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

		// TERMINA PEDIDO SOPORTE GRABAR INSERTAR
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

	// INICIA INSERTAR PEDIDO SOPORTE
	log.Println("Got %s age %s club %s\n", tempPedidosoporte.Codigo, tempPedidosoporte.Tercero, tempPedidosoporte.Neto)
	var q string
	q = "insert into pedidosoporte ("
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
	q += "Tipo"
	q += " ) values("
	q+=parametros(32)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempPedidosoporte.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempPedidosoporte.Codigo,
		tempPedidosoporte.Fecha.Format(layout),
		tempPedidosoporte.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempPedidosoporte.Descuento,
		tempPedidosoporte.Subtotaldescuento19,
		tempPedidosoporte.Subtotaldescuento5,
		tempPedidosoporte.Subtotaldescuento0,
		tempPedidosoporte.Subtotal,
		tempPedidosoporte.Subtotal19,
		tempPedidosoporte.Subtotal5,
		tempPedidosoporte.Subtotal0,
		tempPedidosoporte.Subtotaliva19,
		tempPedidosoporte.Subtotaliva5,
		tempPedidosoporte.Subtotaliva0,
		tempPedidosoporte.Subtotalbase19,
		tempPedidosoporte.Subtotalbase5,
		tempPedidosoporte.Subtotalbase0,
		tempPedidosoporte.TotalIva,
		tempPedidosoporte.Total,
		tempPedidosoporte.PorcentajeRetencionFuente,
		tempPedidosoporte.TotalRetencionFuente,
		tempPedidosoporte.PorcentajeRetencionIca,
		tempPedidosoporte.TotalRetencionIca,
		tempPedidosoporte.Neto,
		tempPedidosoporte.Items,
		tempPedidosoporte.Formadepago,
		tempPedidosoporte.Mediodepago,
		tempPedidosoporte.Tercero,
		tempPedidosoporte.Almacenista,
		tempPedidosoporte.Centro,
		tempPedidosoporte.Tipo)

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

// INICIA PEDIDO SOPORTE EXISTE
func PedidosoporteExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM pedidosoporte  WHERE codigo=$1", Codigo)
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

// INICIA PEDIDO SOPORTE EDITAR
func PedidosoporteEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidosoporte editar" + Codigo)

	db := dbConn()

	// traer PEDIDO SOPORTE
	v := pedidosoporte{}
	err := db.Get(&v, "SELECT * FROM pedidosoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []pedidosoportedetalleeditar{}

	err2 := db.Select(&det, PedidosoporteConsultaDetalle(), Codigo)
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
		"pedidosoporte":       v,
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

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/pedidosoporte/pedidosoporteEditar.html", "vista/pedidosoporte/pedidosoporteScript.html")
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error pedidosoporte nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA PEDIDO SOPORTE BORRAR
func PedidosoporteBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio pedidosoporte editar" + Codigo)

	db := dbConn()

	// traer PEDIDO SOPORTE
	v := pedidosoporte{}
	err := db.Get(&v, "SELECT * FROM pedidosoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []pedidosoportedetalleeditar{}
	err2 := db.Select(&det, PedidosoporteConsultaDetalle(), Codigo)
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
		"pedidosoporte":       v,
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
		"vista/pedidosoporte/pedidosoporteBorrar.html", "vista/pedidosoporte/pedidosoporteScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error pedidosoporte nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA PEDIDO SOPORTE ELIMINAR
func PedidosoporteEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar PEDIDO SOPORTE
	delForm, err := db.Prepare("DELETE from pedidosoporte WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from pedidosoportedetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/PedidosoporteLista", 301)
}

// INICIA PDF
func PedidosoportePdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer PEDIDO SOPORTE
	miPedidosoporte := pedidosoporte{}
	err := db.Get(&miPedidosoporte, "SELECT * FROM pedidosoporte where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []pedidosoportedetalleeditar{}
	err2 := db.Select(&miDetalle, PedidosoporteConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miPedidosoporte.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miPedidosoporte.Almacenista)
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
		pdf.CellFormat(190, 10, "PEDIDO SOPORTE", "0", 0, "C",
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
	PedidosoporteCabecera(pdf,miTercero,miPedidosoporte,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			PedidosoporteFilaDetalle(pdf,miFila,a)
		}
		PedidosoportePieDePagina(pdf,miTercero,miPedidosoporte)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					PedidosoporteFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			PedidosoporteCabecera(pdf,miTercero,miPedidosoporte,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					PedidosoporteFilaDetalle(pdf,miFila,a)
				}
			}

			PedidosoportePieDePagina(pdf,miTercero,miPedidosoporte)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					PedidosoporteFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			PedidosoporteCabecera(pdf,miTercero,miPedidosoporte,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					PedidosoporteFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			PedidosoporteCabecera(pdf,miTercero,miPedidosoporte,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					PedidosoporteFilaDetalle(pdf,miFila,a)
				}
			}

			PedidosoportePieDePagina(pdf,miTercero,miPedidosoporte)
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
func PedidosoporteCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miPedidosoporte pedidosoporte, miAlmacenista almacenista ){
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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miPedidosoporte.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miPedidosoporte.Fecha.Format("02/01/2006")+" "+Titulo(miPedidosoporte.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miPedidosoporte.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miPedidosoporte.Vence.Format("02/01/2006"), "", 0,
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
func PedidosoporteFilaDetalle(pdf *gofpdf.Fpdf,miFila pedidosoportedetalleeditar, a int ){
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

func PedidosoportePieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miPedidosoporte pedidosoporte ){

	Totalletras,err := IntLetra(Cadenaentero(miPedidosoporte.Neto))
	if err!= nil{
		fmt.Println(err)
	}

	pdf.SetFont("Arial", "", 8)
	pdf.SetY(225)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "SON: " +Mayuscula(Totalletras)+" PESOS MDA. CTE.", "0", 0,
		"L", false, 0, "")

	pdf.Ln(20)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "__________________________________________", "0", 0,
		"C", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "A C E P T A D A ", "0", 0, "C",
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
	pdf.CellFormat(190, 10, miPedidosoporte.Subtotal, "0", 0, "R",
		false, 0, "")

	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miPedidosoporte.TotalRetencionFuente, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miPedidosoporte.TotalRetencionIca, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(14)
	pdf.CellFormat(190, 10, miPedidosoporte.Neto, "0", 0, "R",
		false, 0, "")
}
