package main

// INICIA PRODUCTO IMPORTAR PAQUETES
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
	"log"
	"net/http"
	"strings"
)

// TERMINA PRODUCTO IMPORTAR PAQUETES

// INICIA PRODUCTO ESTRUCTURA JSON
type productoJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
	Iva		        string `json:"Iva"`
	Unidad		    string `json:"Unidad"`
	Precio			string `json:"Precio"`
}

// TERMINA PRODUCTO ESTRUCTURA JSON

// INICIA PRODUCTO ESTRUCTURA
type producto struct {
	Codigo          string
	Nombre          string
	Iva		        string
	Unidad		    string
	Subgrupo		string
	Tipo			string
	Precio			string
	Costo           string
	Cantidad        string
	Total           string
}

// TERMINA PRODUCTO ESTRUCTURA

// INICIA PRODUCTO LISTA
func ProductoLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/producto/productoLista.html")
	log.Println("Error producto 0")
	db := dbConn()
	res := []producto{}
	db.Select(&res, "SELECT * FROM producto ORDER BY cast(codigo as integer) ASC")
	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error producto888")
	tmp.Execute(w, varmap)
}

// TERMINA PRODUCTO LISTA

// INICIA PRODUCTO NUEVO
func ProductoNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error producto nuevo 1")
	Codigo := mux.Vars(r)["codigo"]
	Panel := mux.Vars(r)["panel"]
	Elemento := mux.Vars(r)["elemento"]
	log.Println("Error producto nuevo 2")
	parametros := map[string]interface{}{
		// INICIA PRODUCTO NUEVO AUTOCOMPLETADO
		"Codigo":                  Codigo,
		"Panel":                   Panel,
		"Elemento":                Elemento,
		"hosting":                 ruta,
		"subgrupo":                ListaSubgrupo(),
		"unidaddemedida":          ListaUnidaddemedida(),
		// TERMINA PRODUCTO NUEVO AUTOCOMPLETADO
	}
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html", "vista/producto/productoNuevo.html")
	log.Println("Error producto nuevo 3")
	tmp.Execute(w, parametros)
}

// TERMINA PRODUCTO NUEVO

// INICIA PRODUCTO INSERTAR
func ProductoInsertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	var t producto
	err = decoder.Decode(&t, r.PostForm)
	if err != nil {
		// Handle error
		panic(err.Error())
	}
	var q string
	q = "insert into producto ("
	q += "Codigo,"
	q += "Nombre,"
	q += "Iva,"
	q += "Unidad,"
	q += "Subgrupo,"
	q += "Tipo,"
	q += "Precio,"
	q += "Costo,"
	q += "Cantidad,"
	q += "Total"
	q += " ) values("
	q += parametros(10)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	// INICIA GRABAR PRODUCTO INSERTAR
	t.Codigo = strings.Replace(t.Codigo, ".", "", -1)
	t.Nombre = Titulo(t.Nombre)
	//t.Unidad = Titulo(t.Unidad)
	t.Precio = Puntos(t.Precio)
	t.Costo = Puntos(t.Costo)
	t.Cantidad = Puntos(t.Cantidad)
	t.Total = Puntos(t.Total)
	// TERMINA PRODUCTO GRABAR INSERTAR
	_, err = insForm.Exec(
		t.Codigo,
		t.Nombre,
		t.Iva,
		t.Unidad,
		t.Subgrupo,
		t.Tipo,
		t.Precio,
		t.Costo,
		t.Cantidad,
		t.Total)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/ProductoLista", 301)
}

// TERMINA PRODUCTO INSERTAR

// INICIA PRODUCTO BUSCAR
func ProductoBuscar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	selDB, err := db.Query("SELECT codigo,"+
		"nombre,unidad,iva,precio FROM producto where codigo LIKE '%' || $1 || '%'  or  nombre LIKE '%' || $1 || '%' ORDER BY"+
		" codigo DESC", Codigo)
	if err != nil {
		panic(err.Error())
	}
	var resJson []productoJson
	var contar int
	contar = 0
	for selDB.Next() {
		contar++
		var id string
		var label string
		var value string
		var nombre string
		var unidad string
		var iva string
		var precio1 string

		err = selDB.Scan(&id, &nombre, &iva, &unidad,&precio1)
		if err != nil {
			panic(err.Error())
		}
		value = id
		label = id + "  -  " + nombre
		resJson = append(resJson, productoJson{id, label, value, nombre,unidad,iva,precio1})
	}
	if err := selDB.Err(); err != nil { // make sure that there was no issue during the process
		log.Println(err)
		return
	}
	if contar == 0 {
		var slice []string
		slice = make([]string, 0)
		data, _ := json.Marshal(slice)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		data, _ := json.Marshal(resJson)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

// TERMINA PRODUCTO BUSCAR

// INICIA PRODUCTO EXISTE
func ProductoExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM producto  WHERE codigo=$1", Codigo)
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

// TERMINA PRODUCTO EXISTE

// INICIA PRODUCTO ACTUAL
func ProductoActual(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := producto{}
	var res []producto
	err := db.Get(&t, "SELECT * FROM producto where codigo=$1", Codigo)

	switch err {
	case nil:
		log.Printf("user found: %+v\n", t)
	case sql.ErrNoRows:
		log.Println("user NOT found, no error")
	default:
		log.Printf("error: %s\n", err)
	}

	//if err != nil {
	//	log.Fatalln(err)
	//}
	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
	res = append(res, t)
	data, err := json.Marshal(res)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// TERMINA PRODUCTO ACTUAL

// INICIA PRODUCTO EDITAR
func ProductoEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio producto editar" + Codigo)
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/producto/ProductoEditar.html")
	db := dbConn()
	t := producto{}
	err := db.Get(&t, "SELECT * FROM producto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
	varmap := map[string]interface{}{
		// INICIA PRODUCTO EDITAR AUTOCOMPLETADO
		"emp":                     t,
		"hosting":                 ruta,
		"subgrupo":                ListaSubgrupo(),
		"unidaddemedida":          ListaUnidaddemedida(),
		// TERMINA PRODUCTO EDITAR AUTOCOMPLETADO
	}
	tmp.Execute(w, varmap)
}

// INICIA PRODUCTO ACTUALIZAR
func ProductoActualizar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
		// Handle error
	}
	var t producto
	err = decoder.Decode(&t, r.PostForm)
	if err != nil {
		// Handle error
		panic(err.Error())
	}
	var q string
	q = "UPDATE producto set "
	q += "Nombre=$2,"
	q += "Iva=$3,"
	q += "Unidad=$4,"
	q += "Subgrupo=$5,"
	q += "Tipo=$6,"
	q += "Precio=$7,"
	q += "Costo=$8,"
	q += "Cantidad=$9,"
	q += "Total=$10"
	q += " where "
	q += "Codigo=$1"

	log.Println("cadena" + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	// INICIA GRABAR PRODUCTO ACTUALIZAR
	t.Codigo = strings.Replace(t.Codigo, ".", "", -1)
	t.Nombre = Titulo(t.Nombre)
	t.Precio = Puntos(t.Precio)
	t.Costo = Puntos(t.Costo)
	t.Cantidad = Puntos(t.Cantidad)
	t.Total = Puntos(t.Total)
	// TERMINA GRABAR PRODUCTO ACTUALIZAR
	_, err = insForm.Exec(
		t.Codigo,
		t.Nombre,
		t.Iva,
		t.Unidad,
		t.Subgrupo,
		t.Tipo,
		t.Precio,
		t.Costo,
		t.Cantidad,
		t.Total)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/ProductoLista", 301)

}

// INICIA PRODUCTO BORRAR
func ProductoBorrar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/Producto/ProductoBorrar.html")
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio Producto borrar" + Codigo)
	db := dbConn()
	t := producto{}
	err := db.Get(&t, "SELECT * FROM Producto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codigo nombre99 borrar" + t.Codigo + t.Nombre)
	varmap := map[string]interface{}{
		// INICIA PRODUCTO BORRAR AUTOCOMPLETADO
		"emp":                     t,
		"hosting":                 ruta,
		"subgrupo":                ListaSubgrupo(),
		"unidaddemedida":          ListaUnidaddemedida(),
		// TERMINA PRODUCTO BORRAR AUTOCOMPLETADO
	}
	tmp.Execute(w, varmap)
}

// INICIA PRODUCTO ELIMINAR
func ProductoEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	emp := mux.Vars(r)["codigo"]
	delForm, err := db.Prepare("DELETE from producto WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("Registro Eliminado" + emp)
	http.Redirect(w, r, "/ProductoLista", 301)
}
// TERMINA PRODUCTO ELIMINAR

// INICIA PRODUCTO PDF
func ProductoPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := producto{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Get(&t, "SELECT * FROM producto where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	var buf bytes.Buffer
	var err1 error
	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)
	ene := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetHeaderFunc(func() {
		pdf.Image(imageFile("logo.png"), 20, 20, 40, 0, false,
			"", 0, "")
		pdf.SetY(15)
		//pdf.AddFont("Helvetica", "", "cp1251.map")
		pdf.SetFont("Helvetica", "", 10)
		pdf.CellFormat(190, 10, e.Nombre, "0", 0,
			"C", false, 0, "")
		pdf.Ln(4)

		pdf.CellFormat(190, 10, "Nit No. " +Coma(e.Codigo)+ " - "+e.Dv, "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, e.Iva+ " - "+e.ReteIva, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, "Actividad Ica - "+e.ActividadIca, "0",
			0, "C", false, 0, "")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, e.Direccion, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		log.Println("PRODUCTO 3")
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+ " - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")
		log.Println("PRODUCTO 4")
		pdf.Ln(10)
		pdf.CellFormat(190, 10, "Datos Producto", "0", 0,
			"C", false, 0, "")
		pdf.Ln(10)
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Codigo", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Codigo, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Nombre:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Nombre, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Iva:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Iva, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Unidad:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, TraerUnidaddemedida(t.Unidad), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Subgrupo:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, TraerSubgrupo(t.Subgrupo), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Tipo:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Tipo, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Precio Venta:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, Coma(t.Precio), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Costo:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Costo, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Cantidad:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Cantidad, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Total:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Total, "", 0,
		"", false, 0, "")

	pdf.SetFooterFunc(func() {
		pdf.SetY(-20)
		pdf.SetFont("Arial", "", 9)
		pdf.SetX(30)
		pdf.CellFormat(90, 10, "Sadconf.com", "", 0,
			"L", false, 0, "")
		pdf.CellFormat(80, 10, fmt.Sprintf(" %d de {nb}", pdf.PageNo()), "",
			0, "R", false, 0, "")
	})
	err1 = pdf.Output(&buf)
	if err1 != nil {
		panic(err1.Error())
	}
	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	w.Write(buf.Bytes())
}
// TERMINA PRODUCTO PDF
