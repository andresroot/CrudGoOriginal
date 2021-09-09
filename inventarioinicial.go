package main

// INICIA INVENTARIO INICIAL IMPORTAR PAQUETES
import (
	"bytes"
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


// INICIA INVENTARIO INICIAL ESTRUCTURA JSON
type inventarioinicialJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}


// INICIA INVENTARIO INICIAL ESTRUCTURA
type inventarioinicialLista struct {
	Codigo				string
	Fecha        		time.Time
	Total             	string
	Almacenista       	string
	AlmacenistaNombre	string
}

// INICIA INVENTARIO INICIAL ESTRUCTURA
type inventarioinicial struct {
	Items			string
	Codigo          string
	Fecha           time.Time
	Almacenista     string
	Subtotalbase19  string
	Subtotalbase5   string
	Subtotalbase0   string
	Total           string
	Accion          string
	Detalle         []inventarioinicialdetalle `json:"Detalle"`
	DetalleEditar   []inventarioinicialdetalleeditar `json:"DetalleEditar"`
	Tipo			string
	Centro	        string
}

// TERMINA INVENTARIO INICIAL ESTRUCTURA

// INICIA INVENTARIO INICIALDETALLE ESTRUCTURA
type inventarioinicialdetalle struct {
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
	Tipo		      string
	Fecha			  time.Time
}
// TERMINA INVENTARIO INICIAL ESTRUCTURA

// estructura para editar
type inventarioinicialdetalleeditar struct {
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
func InventarioinicialConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "inventarioinicialdetalle.Id as id ,"
	consulta += "inventarioinicialdetalle.Codigo as codigo,"
	consulta += "inventarioinicialdetalle.Fila as fila,"
	consulta += "inventarioinicialdetalle.Cantidad as cantidad,"
	consulta += "inventarioinicialdetalle.Precio as precio,"
	consulta += "inventarioinicialdetalle.Descuento as descuento,"
	consulta += "inventarioinicialdetalle.Montodescuento as montodescuento,"
	consulta += "inventarioinicialdetalle.Sigratis as sigratis,"
	consulta += "inventarioinicialdetalle.Subtotal as subtotal,"
	consulta += "inventarioinicialdetalle.Subtotaldescuento  as subtotaldescuento,"
	consulta += "inventarioinicialdetalle.Pagina as pagina ,"
	consulta += "inventarioinicialdetalle.Bodega as bodega,"
	consulta += "inventarioinicialdetalle.Producto as producto,"
	consulta += "inventarioinicialdetalle.Fecha as fecha,"
	consulta += "producto.nombre as ProductoNombre, "
	consulta += "producto.iva as ProductoIva, "
	consulta += "producto.unidad as ProductoUnidad "
	consulta += "from inventarioinicialdetalle "
	consulta += "inner join producto on producto.codigo=inventarioinicialdetalle.producto "
	consulta += " where inventarioinicialdetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA INVENTARIO INICIAL LISTA
func InventarioinicialLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/inventarioinicial/inventarioinicialLista.html")
	log.Println("Error inventarioinicial 0")
	var consulta string

	consulta = "  SELECT inventarioinicial.almacenista,inventarioinicial.total,inventarioinicial.codigo,fecha,almacenista.nombre as almacenistanombre"
	consulta += " FROM inventarioinicial "
	consulta += " inner join almacenista on almacenista.codigo=inventarioinicial.almacenista "
	consulta += " inner join centro on centro.codigo=inventarioinicial.centro "
	consulta += " ORDER BY inventarioinicial.codigo ASC"

	db := dbConn()
	res := []inventarioinicialLista{}
	//db.Select(&res, consulta)

	//error1 = db.Select(&res, consulta)
	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error inventarioinicial888")
	tmp.Execute(w, varmap)
}

// INICIA INVENTARIO INICIAL NUEVO
func InventarioinicialNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error inventarioinicial nuevo 1")
	log.Println("Error inventarioinicial nuevo 2")
	parametros := map[string]interface{}{
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"centro":      ListaCentro(),
	}

	t, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/inventarioinicial/inventarioinicialNuevo.html", "vista/inventarioinicial/inventarioinicialScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error inventarioinicial nuevo 3")
	t.Execute(w, parametros)
}

// INICIA INVENTARIO INICIAL INSERTAR AJAX
func InventarioinicialAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempInventarioinicial inventarioinicial

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la Inventarioinicial
	err = json.Unmarshal(b, &tempInventarioinicial)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if tempInventarioinicial.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from inventarioinicialdetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempInventarioinicial.Codigo)

		// borra detalle inventario
		delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo='Inventarioinicial'")
		if err != nil {
			panic(err.Error())
		}
		delForm2.Exec(tempInventarioinicial.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from inventarioinicial WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempInventarioinicial.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempInventarioinicial.Detalle {
		var a = i
		var q string
		q = "insert into inventarioinicialdetalle ("
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

		// TERMINA INVENTARIO INICIAL GRABAR INSERTAR
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
	for i, x := range tempInventarioinicial.Detalle {
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
			operacionInventarioInicial)
		if err != nil {
			panic(err)
		}
		log.Println("Insertar Producto \n", x.Producto, a)
	}

	// INICIA INSERTAR INVENTARIO INICIAL
	log.Println("Got %s age %s club %s\n", tempInventarioinicial.Codigo, tempInventarioinicial.Total)
	var q string
	q = "insert into inventarioinicial ("
	q += "Codigo,"
	q += "Fecha,"
	q += "Subtotalbase19,"
	q += "Subtotalbase5,"
	q += "Subtotalbase0,"
	q += "Total,"
	q += "Items,"
	q += "Almacenista,"
	q += "Centro"
	q += " ) values("
	q+=parametros(9)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"

		_, err = insForm.Exec(
		tempInventarioinicial.Codigo,
		tempInventarioinicial.Fecha.Format(layout),
		tempInventarioinicial.Subtotalbase19,
		tempInventarioinicial.Subtotalbase5,
		tempInventarioinicial.Subtotalbase0,
		tempInventarioinicial.Total,
		tempInventarioinicial.Items,
		tempInventarioinicial.Almacenista,
		tempInventarioinicial.Centro)

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

// INICIA INVENTARIO INICIAL EXISTE
func InventarioinicialExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM inventarioinicial  WHERE codigo=$1", Codigo)
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

// INICIA INVENTARIO INICIAL EDITAR
func InventarioinicialEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio inventarioinicial editar" + Codigo)

	db := dbConn()

	// traer INVENTARIO INICIAL
	v := inventarioinicial{}
	err := db.Get(&v, "SELECT * FROM inventarioinicial where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []inventarioinicialdetalleeditar{}

	err2 := db.Select(&det, InventarioinicialConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"inventarioinicial":       v,
		"detalle":     det,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"centro" : ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/inventarioinicial/inventarioinicialEditar.html", "vista/inventarioinicial/inventarioinicialScript.html")
	fmt.Printf("%v, %v", miTemplate, err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error inventarioinicial nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA INVENTARIO INICIAL BORRAR
func InventarioinicialBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio inventarioinicial editar" + Codigo)

	db := dbConn()

	// traer INVENTARIO INICIAL
	v := inventarioinicial{}
	err := db.Get(&v, "SELECT * FROM inventarioinicial where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []inventarioinicialdetalleeditar{}
	err2 := db.Select(&det, InventarioinicialConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}


	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"inventarioinicial":       v,
		"detalle":     det,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),
		"centro":      ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/inventarioinicial/inventarioinicialBorrar.html", "vista/inventarioinicial/inventarioinicialScript.html")
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error inventarioinicial nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA INVENTARIO INICIAL ELIMINAR
func InventarioinicialEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar INVENTARIO INICIAL
	delForm, err := db.Prepare("DELETE from inventarioinicial WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from inventarioinicialdetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borar detalle invenario
	Borrarinventario(codigo,"Inventarioinicial")


	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/InventarioinicialLista", 301)
}

// INICIA PDF
func InventarioinicialPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer INVENTARIO INICIAL
	miInventarioinicial := inventarioinicial{}
	err := db.Get(&miInventarioinicial, "SELECT * FROM inventarioinicial where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []inventarioinicialdetalleeditar{}
	err2 := db.Select(&miDetalle, InventarioinicialConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}


	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miInventarioinicial.Almacenista)
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
		pdf.CellFormat(190, 10, "INVENTARIO INICIAL", "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
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
	InventarioinicialCabecera(pdf,miInventarioinicial,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			InventarioinicialFilaDetalle(pdf,miFila,a)
		}
		InventarioinicialPieDePagina(pdf,miInventarioinicial)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					InventarioinicialFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			InventarioinicialCabecera(pdf,miInventarioinicial,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					InventarioinicialFilaDetalle(pdf,miFila,a)
				}
			}

			InventarioinicialPieDePagina(pdf,miInventarioinicial)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					InventarioinicialFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			InventarioinicialCabecera(pdf,miInventarioinicial,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					InventarioinicialFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			InventarioinicialCabecera(pdf,miInventarioinicial,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					InventarioinicialFilaDetalle(pdf,miFila,a)
				}
			}

			InventarioinicialPieDePagina(pdf,miInventarioinicial)
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

func InventarioinicialCabecera(pdf *gofpdf.Fpdf,miInventarioinicial inventarioinicial, miAlmacenista almacenista ){
	pdf.SetFont("Arial", "", 10)
	pdf.Ln(8)
	pdf.SetX(20)
	pdf.CellFormat(184, 5, "DATOS", "1", 0,
		"C", false, 0, "")
	pdf.Ln(8)

	pdf.SetX(20)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(60)
	pdf.CellFormat(40, 4, miInventarioinicial.Fecha.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Almacenista", "", 0,
		"L", false, 0, "")
	pdf.SetX(140)
	pdf.CellFormat(40, 4, miAlmacenista.Nombre, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.SetY(62)
	pdf.Ln(3)
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
func InventarioinicialFilaDetalle(pdf *gofpdf.Fpdf,miFila inventarioinicialdetalleeditar, a int ){
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
func InventarioinicialPieDePagina(pdf *gofpdf.Fpdf,miInventarioinicial inventarioinicial ){

	Totalletras,err := IntLetra(Cadenaentero(miInventarioinicial.Total))
	if err!= nil{
		fmt.Println(err)
	}

	ene := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Arial", "", 8)
	pdf.SetY(222)
	pdf.SetX(20)
	pdf.CellFormat(190, 10, "SON: " +ene(Mayuscula(Totalletras))+" PESOS MDA. CTE.", "0", 0,
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
	pdf.SetY(233)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "SUBTOTAL IVA 19%", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "SUBTOTAL IVA 5%", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "SUBTOTAL IVA 0%", "0", 0, "L",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(150)
	pdf.CellFormat(190, 10, "TOTAL", "0", 0, "L",
		false, 0, "")

	pdf.SetY(229)
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miInventarioinicial.Subtotalbase19, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miInventarioinicial.Subtotalbase5, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miInventarioinicial.Subtotalbase0, "0", 0, "R",
		false, 0, "")
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.CellFormat(190, 10, miInventarioinicial.Total, "0", 0, "R",
		false, 0, "")
}
