package main

// INICIA COTIZACIONSERVICIO IMPORTAR PAQUETES
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

// INICIA COTIZACIONSERVICIO ESTRUCTURA JSON
type cotizacionservicioJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA COTIZACIONSERVICIO ESTRUCTURA
type cotizacionservicioLista struct {
	Codigo        string
	Fecha         time.Time
	Total         string
	Tercero       string
	TerceroNombre string
	CentroNombre  string
	VendedorNombre string
}

// INICIA COTIZACIONSERVICIO ESTRUCTURA
type cotizacionservicio struct {
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
	Detalle                   []cotizacionserviciodetalle `json:"Detalle"`
	DetalleEditar			  []cotizacionserviciodetalleeditar `json:"DetalleEditar"`
	TerceroDetalle 			  tercero
	Tipo 					  string
	Ret2201					  string
	Centro					  string
}

// INICIA COTIZACIONSERVICIODETALLE ESTRUCTURA
type cotizacionserviciodetalle struct {
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
	Tipo 			  string
	Fecha             time.Time
}

// estructura para editar
type cotizacionserviciodetalleeditar struct {
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

// INICIA COMPRA CONSULTA DETALLE
func CotizacionservicioConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "cotizacionserviciodetalle.Id as id ,"
	consulta += "cotizacionserviciodetalle.Codigo as codigo,"
	consulta += "cotizacionserviciodetalle.Fila as fila,"
	consulta += "cotizacionserviciodetalle.Cantidad as cantidad,"
	consulta += "cotizacionserviciodetalle.Precio as precio,"
	consulta += "cotizacionserviciodetalle.Descuento as descuento,"
	consulta += "cotizacionserviciodetalle.Montodescuento as montodescuento,"
	consulta += "cotizacionserviciodetalle.Sigratis as sigratis,"
	consulta += "cotizacionserviciodetalle.Subtotal as subtotal,"
	consulta += "cotizacionserviciodetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "cotizacionserviciodetalle.Pagina as pagina ,"
	consulta += "cotizacionserviciodetalle.Fecha as fecha,"
	consulta += "cotizacionserviciodetalle.Nombreservicio as Nombreservicio,"
	consulta += "cotizacionserviciodetalle.Unidadservicio as Unidadservicio, "
	consulta += "cotizacionserviciodetalle.Codigoservicio as Codigoservicio, "
	consulta += "cotizacionserviciodetalle.Ivaservicio as Ivaservicio "
	consulta += "from cotizacionserviciodetalle "
	consulta += " where cotizacionserviciodetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA COTIZACIONSERVICIO LISTA
func CotizacionservicioLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cotizacionservicio/cotizacionservicioLista.html")
	log.Println("Error cotizacionservicio 0")
	var consulta string

	consulta = "  SELECT vendedor.nombre as VendedorNombre,centro.nombre as CentroNombre,total,cotizacionservicio.codigo,fecha,tercero,tercero.nombre as TerceroNombre "
	consulta += " FROM cotizacionservicio "
	consulta += " inner join tercero on tercero.codigo=cotizacionservicio.tercero "
	consulta += " inner join centro on centro.codigo=cotizacionservicio.centro "
	consulta += " inner join vendedor on vendedor.codigo=cotizacionservicio.vendedor "
	consulta += " ORDER BY cotizacionservicio.codigo ASC"

	db := dbConn()
	res := []cotizacionservicioLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error cotizacionservicio888")
	tmp.Execute(w, varmap)
}

// INICIA COTIZACIONSERVICIO NUEVO
func CotizacionservicioNuevo(w http.ResponseWriter, r *http.Request) {

	// TRAE COPIA DE EDITAR
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacionservicio editar" + Codigo)

	db := dbConn()
	v := cotizacionservicio{}
	tc := tercero{}
	det := []cotizacionserviciodetalleeditar{}
	if Codigo == "False" {

	} else {

		err := db.Get(&v, "SELECT * FROM cotizacionservicio where codigo=$1", Codigo)
		if err != nil {
			log.Fatalln(err)
		}

		err2 := db.Select(&det, CotizacionservicioConsultaDetalle(), Codigo)
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
		"cotizacionservicio":       v,
		"detalle":     det,
		"tercero":     tc,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" :ListaCentro(),
		"ventaserviciotipoiva":TraerParametrosInventario().Ventaserviciotipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
		"codigo": Codigo,
	}
	//TERMINA TRAE COPIA DE EDITAR

	t, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/cotizacionservicio/cotizacionservicioNuevo.html", "vista/cotizacionservicio/cotizacionservicioScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error cotizacionservicio nuevo 3")
	t.Execute(w, parametros)
}

// INICIA COTIZACIONSERVICIO INSERTAR AJAX
func CotizacionservicioAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempCotizacion cotizacionservicio

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la cotizacionservicio
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
		delForm, err := db.Prepare("DELETE from cotizacionserviciodetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempCotizacion.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from cotizacionservicio WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempCotizacion.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempCotizacion.Detalle {
		var a = i
		var q string
		q = "insert into cotizacionserviciodetalle ("
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

		// TERMINA COTIZACIONSERVICIO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Id,
			x.Codigo,
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

	// INICIA INSERTAR COTIZACIONSERVICIO
	log.Println("Got %s age %s club %s\n", tempCotizacion.Codigo, tempCotizacion.Tercero, tempCotizacion.Total)
	var q string
	q = "insert into cotizacionservicio ("
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

// INICIA COTIZACIONSERVICIO EXISTE
func CotizacionservicioExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM cotizacionservicio  WHERE codigo=$1", Codigo)
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

// INICIA COTIZACIONSERVICIO EDITAR
func CotizacionservicioEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacionservicio editar" + Codigo)

	db := dbConn()

	// traer cotizacionservicio
	v := cotizacionservicio{}
	err := db.Get(&v, "SELECT * FROM cotizacionservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []cotizacionserviciodetalleeditar{}

	err2 := db.Select(&det, CotizacionservicioConsultaDetalle(), Codigo)
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
		"cotizacionservicio":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"centro" :ListaCentro(),
		"ventaserviciotipoiva":TraerParametrosInventario().Ventaserviciotipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/cotizacionservicio/cotizacionservicioEditar.html", "vista/cotizacionservicio/cotizacionservicioScript.html")
	fmt.Printf("%v, %v", t, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error cotizacionservicio nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA COTIZACIONSERVICIO BORRAR
func CotizacionservicioBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio cotizacionservicio editar" + Codigo)

	db := dbConn()

	// traer cotizacionservicio
	v := cotizacionservicio{}
	err := db.Get(&v, "SELECT * FROM cotizacionservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []cotizacionserviciodetalleeditar{}
	err2 := db.Select(&det, CotizacionservicioConsultaDetalle(), Codigo)
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
		"cotizacionservicio":       v,
		"detalle":     det,
		"tercero":     t,
		"hosting":     ruta,
		"vendedor":    ListaVendedor(),
		"bodega":      ListaBodega(),
		"mediodepago": ListaMedioDePago(),
		"formadepago": ListaFormaDePago(),
		"ventaserviciotipoiva":TraerParametrosInventario().Ventaserviciotipoiva,
		"autoret2201":TraerParametrosInventario().Ventacuentaporcentajeret2201,
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cotizacionservicio/cotizacionservicioBorrar.html", "vista/cotizacionservicio/cotizacionservicioScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error cotizacionservicio nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA COTIZACIONSERVICIO ELIMINAR
func CotizacionservicioEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar cotizacionservicio
	delForm, err := db.Prepare("DELETE from cotizacionservicio WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from cotizacionserviciodetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/CotizacionservicioLista", 301)
}

// INICIA PDF
func CotizacionservicioPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer cotizacionservicio
	miCotizacionservicio := cotizacionservicio{}
	err := db.Get(&miCotizacionservicio, "SELECT * FROM cotizacionservicio where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []cotizacionserviciodetalleeditar{}
	err2 := db.Select(&miDetalle, CotizacionservicioConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miCotizacionservicio.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	// traer Vendedor
	miVendedor := vendedor{}
	err4 := db.Get(&miVendedor, "SELECT * FROM vendedor where codigo=$1", miCotizacionservicio.Vendedor)
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

		// FACTURA NUMERO
		pdf.SetFont("Arial", "", 10)
		pdf.Ln(5)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "COTIZACION SERVICIO", "0", 0, "C",
			false, 0, "")
		pdf.Ln(5)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, " No. "+Codigo, "0", 0, "C",
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
	CotizacionservicioCabecera(pdf,miTercero,miCotizacionservicio, miVendedor)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			CotizacionservicioFilaDetalle(pdf,miFila,a)
		}
		CotizacionservicioPieDePagina(pdf,miTercero,miCotizacionservicio)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CotizacionservicioFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			CotizacionservicioCabecera(pdf,miTercero,miCotizacionservicio, miVendedor)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					CotizacionservicioFilaDetalle(pdf,miFila,a)
				}
			}

			CotizacionservicioPieDePagina(pdf,miTercero,miCotizacionservicio)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CotizacionservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CotizacionservicioCabecera(pdf,miTercero,miCotizacionservicio, miVendedor)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					CotizacionservicioFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CotizacionservicioCabecera(pdf,miTercero,miCotizacionservicio, miVendedor)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					CotizacionservicioFilaDetalle(pdf,miFila,a)
				}
			}

			CotizacionservicioPieDePagina(pdf,miTercero,miCotizacionservicio)
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
func CotizacionservicioCabecera(pdf *gofpdf.Fpdf,miTercero tercero,miCotizacionservicio cotizacionservicio, miVendedor vendedor ){
	pdf.SetFont("Arial", "", 10)
	pdf.Ln(4)
	pdf.SetX(20)
	pdf.CellFormat(90, 5, "DATOS DEL COMPRADOR", "1", 0,
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
	pdf.CellFormat(40, 4, Titulo(TraerFormadepago(miCotizacionservicio.Formadepago)), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miCotizacionservicio.Fecha.Format("02/01/2006")+" "+Titulo(miCotizacionservicio.Hora), "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Medio de Pago", "", 0,
		"L", false, 0, "")
	pdf.SetX(50)
	pdf.CellFormat(40, 4, TraerMediodepago(miCotizacionservicio.Mediodepago), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Fecha de Vencimiento", "", 0,
		"L", false, 0, "")
	pdf.SetX(150)
	pdf.CellFormat(40, 4, miCotizacionservicio.Vence.Format("02/01/2006"), "", 0,
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
func CotizacionservicioFilaDetalle(pdf *gofpdf.Fpdf,miFila cotizacionserviciodetalleeditar, a int ){
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

func CotizacionservicioPieDePagina(pdf *gofpdf.Fpdf,miTercero tercero,miCotizacionservicio cotizacionservicio ){

	Totalletras,err := IntLetra(Cadenaentero(miCotizacionservicio.Total))
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
	pdf.CellFormat(190, 10, miCotizacionservicio.Subtotal, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacionservicio.Subtotaliva19, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacionservicio.Subtotaliva5, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacionservicio.Descuento, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miCotizacionservicio.Total, "0", 0, "R",
		false, 0, "")

}
