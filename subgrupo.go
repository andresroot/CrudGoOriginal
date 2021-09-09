package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

// SUBGRUPO TABLA
type Subgrupo struct {
	Codigo      string
	Nombre      string
	Grupo		string
}


// SUBGRUPO LISTA
func SubgrupoLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/subgrupo/subgrupoLista.html")
	db := dbConn()

	res := []Subgrupo{}
	db.Select(&res, "SELECT * FROM subgrupo ORDER BY cast(codigo as integer) ASC")
	log.Println("Lista Grupo 1" )

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	tmp.Execute(w, varmap)
}

// SUBGRUPO NUEVO
func SubgrupoNuevo(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/subgrupo/subgrupoNuevo.html")
	//tmp.Execute(w, mapaRuta)

	varmap := map[string]interface{}{
		"grupo": ListaGrupo(),
		"hosting": ruta,
	}
	tmp.Execute(w, varmap)
}

// SUBGRUPO INSERTAR
func SubgrupoInsertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Codigo := r.FormValue("Codigo")
		Nombre := r.FormValue("Nombre")
		Grupo := r.FormValue("Grupo")

		Nombre = Titulo(Nombre)
		insForm, err := db.Prepare("INSERT INTO subgrupo(codigo, nombre, grupo)VALUES($1, $2, $3)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Codigo, Nombre, Grupo)
		log.Println("Nuevo Registro:" + Codigo + "," + Nombre)
	}
	http.Redirect(w, r, "/SubgrupoLista", 301)
}

// SUBGRUPO EXISTE
func SubgrupoExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM subgrupo  WHERE codigo=$1", Codigo)
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

// SUBGRUPO EDITAR
func SubgrupoEditar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/subgrupo/subgrupoEditar.html")
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	selDB, err := db.Query("SELECT * FROM subgrupo WHERE codigo=$1", Codigo)
	if err != nil {
		panic(err.Error())
	}
	emp := Subgrupo{}
	for selDB.Next() {
		var codigo string
		var nombre string
		var grupo string
		err = selDB.Scan(&codigo, &nombre, &grupo)
		if err != nil {
			panic(err.Error())
		}
		emp.Codigo = codigo
		emp.Nombre = nombre
		emp.Grupo = grupo
	}
	varmap := map[string]interface{}{
		"emp":     emp,
		"hosting": ruta,
		"grupo": ListaGrupo(),
	}
	//vistaSubgrupo.ExecuteTemplate(w, "SubgrupoEditar", varmap)
	tmp.Execute(w, varmap)
}

// SUBGRUPO ACTUALIZAR
func SubgrupoActualizar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		codigo := r.FormValue("Codigo")
		nombre := r.FormValue("Nombre")

		nombre = Titulo(nombre)
		insForm, err := db.Prepare("UPDATE subgrupo set	nombre=$2  " + " WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(codigo, nombre)
		log.Println("Registro Actualizado:" + codigo + "," +
			"" + nombre)
	}
	http.Redirect(w, r, "/SubgrupoLista", 301)
}

// SUBGRUPO BORRAR
func SubgrupoBorrar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/subgrupo/subgrupoBorrar.html")
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := Subgrupo{}
	err := db.Get(&t, "SELECT * FROM subgrupo where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	varmap := map[string]interface{}{
		"emp":     t,
		"hosting": ruta,
		"grupo": ListaGrupo(),
	}
	tmp.Execute(w, varmap)
}

// SUBGRUPO ELIMINAR
func SubgrupoEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	emp := mux.Vars(r)["codigo"]
	//Codigo, _ := strconv.ParseInt(emp, 10, 0)
	delForm, err := db.Prepare("DELETE from subgrupo WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("Registro Eliminado" + emp)
	http.Redirect(w, r, "/SubgrupoLista", 301)
}

// SUBGRUPO PDF
func SubgrupoPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := Subgrupo{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Get(&t, "SELECT * FROM subgrupo where codigo=$1", Codigo)
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
		log.Println("tercero 3")
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+ " - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")
		log.Println("tercero 4")
		pdf.Ln(10)
		pdf.CellFormat(190, 10, "Datos Subgrupo", "0", 0,
			"C", false, 0, "")
		pdf.Ln(10)
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Grupo", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, TraerGrupo(t.Grupo), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
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

