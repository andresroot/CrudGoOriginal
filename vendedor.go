package main
// INICIA VENDEDOR IMPORTAR PAQUETES

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
"log"
"net/http"
"strings"
)

// TERMINA VENDEDOR IMPORTAR PAQUETES

// INICIA VENDEDOR ESTRUCTURA JSON
type vendedorJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// TERMINA VENDEDOR ESTRUCTURA JSON

// INICIA VENDEDOR ESTRUCTURA
type vendedor struct {
	Codigo          string
	Nit				string
	Dv              string
	Nombre          string
	Comision        string
}

// TERMINA VENDEDOR ESTRUCTURA

// INICIA VENDEDOR LISTA
func VendedorLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/vendedor/vendedorLista.html")
	log.Println("Error vendedor 0")
	db := dbConn()
	res := []vendedor{}
	db.Select(&res, "SELECT * FROM vendedor ORDER BY cast(codigo as integer ) ASC")
	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error vendedor888")
	tmp.Execute(w, varmap)
}

// TERMINA VENDEDOR LISTA

// INICIA VENDEDOR NUEVO
func VendedorNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error vendedor nuevo 1")
	Codigo := mux.Vars(r)["codigo"]
	Panel := mux.Vars(r)["panel"]
	Elemento := mux.Vars(r)["elemento"]
	log.Println("Error vendedor nuevo 2")
	parametros := map[string]interface{}{
		// INICIA VENDEDOR NUEVO AUTOCOMPLETADO
		"Codigo":                  Codigo,
		"Panel":                   Panel,
		"Elemento":                Elemento,
		"hosting":                 ruta,
		// TERMINA VENDEDOR NUEVO AUTOCOMPLETADO
	}
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html", "vista/vendedor/vendedorNuevo.html", "vista/autocompleta/autocompletaTercero.html")
	log.Println("Error vendedor nuevo 3")
	tmp.Execute(w, parametros)
}


// TERMINA VENDEDOR NUEVO

// INICIA VENDEDOR INSERTAR
func VendedorInsertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	var t vendedor
	err = decoder.Decode(&t, r.PostForm)
	if err != nil {
		// Handle error
		panic(err.Error())
	}
	var q string
	q = "insert into vendedor ("
	q += "Codigo,"
	q += "Nit,"
	q += "Dv,"
	q += "Nombre,"
	q += "Comision"
	q += " ) values("
	q += parametros(5)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	// INICIA GRABAR VENDEDOR INSERTAR
	t.Nit = strings.Replace(t.Nit, ".", "", -1)
	t.Nombre = Titulo(t.Nombre)
	// TERMINA VENDEDOR GRABAR INSERTAR
	_, err = insForm.Exec(
		t.Codigo,
		t.Nit,
		t.Dv,
		t.Nombre,
		t.Comision)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/VendedorLista", 301)
}

// TERMINA VENDEDOR INSERTAR

// INICIA VENDEDOR BUSCAR
func VendedorBuscar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	selDB, err := db.Query("SELECT codigo,"+
		"nombre FROM vendedor where codigo LIKE '%' || $1 || '%'  or  nombre LIKE '%' || $1 || '%' ORDER BY"+
		" codigo DESC", Codigo)
	if err != nil {
		panic(err.Error())
	}
	var resJson []vendedorJson
	var contar int
	contar = 0
	for selDB.Next() {
		contar++
		var id string
		var label string
		var value string
		var nombre string
		err = selDB.Scan(&id, &nombre)
		if err != nil {
			panic(err.Error())
		}
		value = id
		label = id + " " + nombre
		resJson = append(resJson, vendedorJson{id, label, value, nombre})
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

// TERMINA VENDEDOR BUSCAR

// INICIA VENDEDOR EXISTE
func VendedorExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM vendedor  WHERE codigo=$1", Codigo)
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

// TERMINA VENDEDOR EXISTE

// INICIA VENDEDOR ACTUAL
func VendedorActual(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := vendedor{}
	var res []vendedor
	err := db.Get(&t, "SELECT * FROM vendedor where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
	res = append(res, t)
	data, err := json.Marshal(res)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// TERMINA VENDEDOR ACTUAL

// INICIA VENDEDOR EDITAR
func VendedorEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio vendedor editar" + Codigo)
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/vendedor/vendedorEditar.html",
		"vista/autocompleta/autocompletaTercero.html")
	db := dbConn()
	t := vendedor{}
	err := db.Get(&t, "SELECT * FROM vendedor where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
	varmap := map[string]interface{}{
		// INICIA VENDEDOR EDITAR AUTOCOMPLETADO
		"emp":                     t,
		"hosting":                 ruta,
		"ciudad":                  ListaCiudad(),
		// TERMINA VENDEDOR EDITAR AUTOCOMPLETADO
	}
	tmp.Execute(w, varmap)
}

// INICIA VENDEDOR ACTUALIZAR
func VendedorActualizar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
		// Handle error
	}
	var t vendedor
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(&t, r.PostForm)
	if err != nil {
		// Handle error
		panic(err.Error())
	}
	var q string
	q = "UPDATE vendedor set "
	q += "Nit=$2,"
	q += "Dv=$3,"
	q += "Nombre=$4,"
	q += "Comision=$5"
	q += " where "
	q += "Codigo=$1"

	log.Println("cadena" + q)

	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	// INICIA GRABAR VENDEDOR ACTUALIZAR
	t.Nit = strings.Replace(t.Nit, ".", "", -1)
	t.Nombre = Titulo(t.Nombre)
	// TERMINA GRABAR VENDEDOR ACTUALIZAR
	_, err = insForm.Exec(
		t.Codigo,
		t.Nit,
		t.Dv,
		t.Nombre,
		t.Comision)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/VendedorLista", 301)

}

// INICIA VENDEDOR BORRAR
func VendedorBorrar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/vendedor/vendedorBorrar.html")
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio vendedor borrar" + Codigo)
	db := dbConn()
	t := vendedor{}
	err := db.Get(&t, "SELECT * FROM vendedor where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codigo nombre99 borrar" + t.Codigo + t.Nombre)
	varmap := map[string]interface{}{
		// INICIA VENDEDOR BORRAR AUTOCOMPLETADO
		"emp":                     t,
		"hosting":                 ruta,
		// TERMINA VENDEDOR BORRAR AUTOCOMPLETADO
	}
	tmp.Execute(w, varmap)
}

// INICIA VENDEDOR ELIMINAR
func VendedorEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	emp := mux.Vars(r)["codigo"]
	delForm, err := db.Prepare("DELETE from vendedor WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("Registro Eliminado" + emp)
	http.Redirect(w, r, "/VendedorLista", 301)
}
// TERMINA VENDEDOR ELIMINAR

// INICIA VENDEDOR PDF
func VendedorPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := vendedor{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Get(&t, "SELECT * FROM vendedor where codigo=$1", Codigo)
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
		log.Println("vendedor 3")
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+ " - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")
		log.Println("vendedor 4")
		pdf.Ln(10)
		pdf.CellFormat(190, 10, "Datos Vendedor", "0", 0,
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
	pdf.CellFormat(40, 4, "Nit No.:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, Coma(t.Nit)+ " - "+t.Dv, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Nombre:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Nombre, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Comision %:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Comision, "", 0,
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
// TERMINA VENDEDOR PDF

