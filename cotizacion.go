package main

// INICIA COTIZACION IMPORTAR PAQUETES
import (
	"bytes"
	//"database/sql"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "github.com/bitly/go-simplejson"
	//"github.com/dustin/go-humanize"
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

// INICIA COTIZACION ESTRUCTURA JSON
type cotizacionJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA COTIZACION ESTRUCTURA
type cotizacionLista struct {
	Codigo        string
	Fecha         time.Time
	Total         string
	Tercero       string
	TerceroNombre string
	CentroNombre  string
	VendedorNombre string
}

// INICIA COTIZACION ESTRUCTURA
type cotizacion struct {
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
	Tercero                   string
	Vendedor                  string
	Accion                    string
	Detalle                   []cotizaciondetalle `json:"Detalle"`
	DetalleEditar			  []cotizaciondetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Tipo 					  string
	Ret2201					  string
	Centro					  string
}

// INICIA COTIZACIONDETALLE ESTRUCTURA
type cotizaciondetalle struct {
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
	Tipo 			  string
	Fecha             time.Time
}

// estructura para editar
type cotizaciondetalleeditar struct {
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
func CotizacionConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "cotizaciondetalle.Id as id ,"
	consulta += "cotizaciondetalle.Codigo as codigo,"
	consulta += "cotizaciondetalle.Fila as fila,"
	consulta += "cotizaciondetalle.Cantidad as cantidad,"
	consulta += "cotizaciondetalle.Precio as precio,"
	consulta += "cotizaciondetalle.Descuento as descuento,"
	consulta += "cotizaciondetalle.Montodescuento as montodescuento,"
	consulta += "cotizaciondetalle.Sigratis as sigratis,"
	consulta += "cotizaciondetalle.Subtotal as subtotal,"
	consulta += "cotizaciondetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "cotizaciondetalle.Pagina as pagina ,"
	consulta += "cotizaciondetalle.Bodega as bodega,"
	consulta += "cotizaciondetalle.Producto as producto,"
	consulta += "cotizaciondetalle.Fecha as fecha,"
	consulta += "producto.nombre as ProductoNombre, "
	consulta += "producto.iva as ProductoIva, "
	consulta += "producto.unidad as ProductoUnidad "
	consulta += "from cotizaciondetalle "
	consulta += "inner join producto on producto.codigo=cotizaciondetalle.producto "
	consulta += " where cotizaciondetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA COTIZACION LISTA
func CotizacionLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cotizacion/cotizacionLista.html")
	log.Println("Error cotizacion 0")
	var consulta string

	consulta = "  SELECT vendedor.nombre as VendedorNombre,centro.nombre as CentroNombre,total,cotizacion.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM cotizacion "
	consulta += " inner join tercero on tercero.codigo=cotizacion.tercero "
	consulta += " inner join centro on centro.codigo=cotizacion.centro "
	consulta += " inner join vendedor on vendedor.codigo=cotizacion.vendedor "
	consulta += " ORDER BY cast(cotizacion.codigo as integer) ASC"

	db := dbConn()
	res := []cotizacionLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error cotizacion888")
	tmp.Execute(w, varmap)
}

// INICIA COTIZACION NUEVO
func CotizacionNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacion editar" + Codigo)

	db := dbConn()
	v := cotizacion{}
	tc := tercero{}
	det := []cotizaciondetalleeditar{}
	if Codigo == "False" {

	} else {

	err := db.Get(&v, "SELECT * FROM cotizacion where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	err2 := db.Select(&det, CotizacionConsultaDetalle(), Codigo)
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
		"codigo":     Codigo,
		"cotizacion":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" :ListaCentro(),
		"ventatipoiva":TraerParametrosInventario().Ventatipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
	}
	//TERMINA TRAE COPIA DE EDITAR

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/cotizacion/cotizacionEditar.html", "vista/cotizacion/cotizacionScript.html")
	fmt.Printf("%v, %v", tc, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error cotizacion nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA COTIZACION INSERTAR AJAX
func CotizacionAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempCotizacion cotizacion

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la cotizacion
	err = json.Unmarshal(b, &tempCotizacion)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	if tempCotizacion.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from cotizaciondetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempCotizacion.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from cotizacion WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempCotizacion.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempCotizacion.Detalle {
		var a = i
		var q string
		q = "insert into cotizaciondetalle ("
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

		// TERMINA COTIZACION GRABAR INSERTAR
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

	// INICIA INSERTAR COTIZACION
	log.Println("Got %s age %s club %s\n", tempCotizacion.Codigo, tempCotizacion.Tercero, tempCotizacion.Total)
	var q string
	q = "insert into cotizacion ("
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
	q += "Tercero,"
	q += "Vendedor,"
	q += "Centro,"
	q += "Tipo"
	q += " ) values("
	q+=parametros(29)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempCotizacion.Fecha.Format("02/01/2006"))
	currentTime := time.Now()

	_, err = insForm.Exec(
		tempCotizacion.Codigo,
		tempCotizacion.Fecha.Format(layout),
		tempCotizacion.Vence.Format(layout),
		currentTime.Format("3:4:5 PM"),
		tempCotizacion.Descuento,
		tempCotizacion.Subtotaldescuento19,
		tempCotizacion.Subtotaldescuento5,
		tempCotizacion.Subtotaldescuento0,
		tempCotizacion.Subtotal,
		tempCotizacion.Subtotal19,
		tempCotizacion.Subtotal5,
		tempCotizacion.Subtotal0,
		tempCotizacion.Subtotaliva19,
		tempCotizacion.Subtotaliva5,
		tempCotizacion.Subtotaliva0,
		tempCotizacion.Subtotalbase19,
		tempCotizacion.Subtotalbase5,
		tempCotizacion.Subtotalbase0,
		tempCotizacion.TotalIva,
		tempCotizacion.Ret2201,
		tempCotizacion.Total,
		tempCotizacion.Neto,
		tempCotizacion.Items,
		tempCotizacion.Formadepago,
		tempCotizacion.Mediodepago,
		tempCotizacion.Tercero,
		tempCotizacion.Vendedor,
		tempCotizacion.Centro,
		tempCotizacion.Tipo)

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

// INICIA COTIZACION EXISTE
func CotizacionExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM cotizacion  WHERE codigo=$1", Codigo)
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

// INICIA COTIZACION EDITAR
func CotizacionEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacion editar" + Codigo)

	db := dbConn()

	// traer cotizacion
	v := cotizacion{}
	err := db.Get(&v, "SELECT * FROM cotizacion where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []cotizaciondetalleeditar{}

	err2 := db.Select(&det, CotizacionConsultaDetalle(), Codigo)
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
		"cotizacion":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" :ListaCentro(),
		"ventatipoiva":TraerParametrosInventario().Ventatipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/cotizacion/cotizacionEditar.html", "vista/cotizacion/cotizacionScript.html")
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error cotizacion nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA COTIZACION BORRAR
func CotizacionBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacion editar" + Codigo)

	db := dbConn()

	// traer cotizacion
	v := cotizacion{}
	err := db.Get(&v, "SELECT * FROM cotizacion where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []cotizaciondetalleeditar{}
	err2 := db.Select(&det, CotizacionConsultaDetalle(), Codigo)
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
		"cotizacion":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"ventatipoiva":TraerParametrosInventario().Ventatipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cotizacion/cotizacionBorrar.html", "vista/cotizacion/cotizacionScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error cotizacion nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA COTIZACION ELIMINAR
func CotizacionEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar cotizacion
	delForm, err := db.Prepare("DELETE from cotizacion WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from cotizaciondetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/CotizacionLista", 301)
}

// INICIA PDF
func CotizacionPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer cotizacion
	miCotizacion := cotizacion{}
	err := db.Get(&miCotizacion, "SELECT * FROM cotizacion where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []cotizaciondetalleeditar{}
	err2 := db.Select(&miDetalle, CotizacionConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miCotizacion.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer Vendedor
	miVendedor := vendedor{}
	err4 := db.Get(&miVendedor, "SELECT * FROM vendedor where codigo=$1", miCotizacion.Vendedor)
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
		// logo
		pdf.Image(imageFile("logo.png"), 20, 30, 40, 0, false,
			"", 0, "")
		pdf.SetY(5)

		// nombre de la empresa centrado
		pdf.SetY(20)
		pdf.SetFont("Arial", "", 8)
		// linea 1
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

		pdf.SetY(20)
		pdf.SetX(80)
		pdf.Ln(8)
		pdf.SetX(80)
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(190, 10, "COTIZACION", "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, " No. "+Codigo, "0", 0, "C",
			false, 0, "")
		pdf.Ln(6)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(256)
		pdf.SetX(20)
		pdf.SetFont("Arial", "", 8)
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
	CotizacionCabecera(pdf,miTercero,miCotizacion, miVendedor)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			CotizacionFilaDetalle(pdf,miFila,a)
		}
		CotizacionPieDePagina(pdf,miTercero,miCotizacion)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CotizacionFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			CotizacionCabecera(pdf,miTercero,miCotizacion, miVendedor)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					CotizacionFilaDetalle(pdf,miFila,a)
				}
			}

			CotizacionPieDePagina(pdf,miTercero,miCotizacion)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CotizacionFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CotizacionCabecera(pdf,miTercero,miCotizacion, miVendedor)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					CotizacionFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CotizacionCabecera(pdf,miTercero,miCotizacion, miVendedor)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					CotizacionFilaDetalle(pdf,miFila,a)
				}
			}

			CotizacionPieDePagina(pdf,miTercero,miCotizacion)
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
func CotizacionCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miCotizacion cotizacion, miVendedor vendedor ){
	pdf.SetFont("Arial", "", 10)
	pdf.Ln(12)
	pdf.SetX(20)
	pdf.CellFormat(90, 5, "DATOS DEL COMPRADOR", "1", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(94, 5, "LUGAR DE ENTREGA O SERVICIO", "1", 0,
		"L", false, 0, "")
	pdf.Ln(8)

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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miCotizacion.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miCotizacion.Fecha.Format("02/01/2006")+" "+Titulo(miCotizacion.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miCotizacion.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miCotizacion.Vence.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Condiciones", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, "", "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Vendedor", "", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, miVendedor.Nombre, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetY(86)
	pdf.SetX(20)
	pdf.CellFormat(184, 5, "ITEM", "1", 0,
		"L", false, 0, "")
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
func CotizacionFilaDetalle(pdf *gofpdf.Fpdf,miFila cotizaciondetalleeditar, a int ){
	pdf.SetFont("Arial", "", 9)
	pdf.SetX(20)
	pdf.CellFormat(40, 4, strconv.Itoa(a), "", 0,
		"L", false, 0, "")
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
	pdf.SetX(100)
	pdf.CellFormat(40, 4, miFila.Descuento, "", 0,
		"R", false, 0, "")
	pdf.SetX(119)
	pdf.CellFormat(40, 4, miFila.Cantidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(141)
	pdf.CellFormat(40, 4, miFila.Precio, "", 0,
		"R", false, 0, "")
	pdf.SetX(165)
	pdf.CellFormat(40, 4, miFila.Subtotal, "", 0,
		"R", false, 0, "")
	pdf.Ln(4)
}
func CotizacionPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miCotizacion cotizacion ){

	Totalletras,err := IntLetra(Cadenaentero(miCotizacion.Total))
	if err!= nil{
		fmt.Println(err)
	}

	pdf.SetFont("Arial", "", 8)
	pdf.SetY(220)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "SON: " +Mayuscula(Totalletras)+" PESOS MDA. CTE.", "0", 0,
		"L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.Ln(20)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "__________________________________________", "0", 0,
		"C", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "A C E P T A D A ", "0", 0, "C",
		false, 0, "")

	pdf.SetFont("Arial", "", 9)

	pdf.SetY(225)
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

	pdf.SetY(225)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacion.Subtotal, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacion.Descuento, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, FormatoFlotanteEntero(Flotante(miCotizacion.Subtotal)  -  Flotante(miCotizacion.Descuento)), "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacion.Subtotaliva19, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacion.Subtotaliva5, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacion.Total, "0", 0, "R",
		false, 0, "")
}

