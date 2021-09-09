package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"html/template"
	"log"
	"net/http"
)

// INICIA USUARIO ESTRUCTURA JSON
type bodegaJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// BODEGA TABLA
type bodega struct {
	Codigo      string
	Nombre      string
}

// BODEGA LISTA
func BodegaLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/bodega/bodegaLista.html")
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM bodega ORDER BY cast(codigo as integer) ASC")
	if err != nil {
		panic(err.Error())
	}
	res := []bodega{}
	for selDB.Next() {
		var Codigo string
		var Nombre string
		err = selDB.Scan(&Codigo, &Nombre)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, bodega{Codigo, Nombre })
	}
	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	tmp.Execute(w, varmap)
}

// BODEGA NUEVO
func BodegaNuevo(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/bodega/bodegaNuevo.html")
	tmp.Execute(w, mapaRuta)
}

// BODEGA INSERTAR
func BodegaInsertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Codigo := r.FormValue("Codigo")
		Nombre := r.FormValue("Nombre")
		Nombre = Titulo(Nombre)
		insForm, err := db.Prepare("INSERT INTO bodega(codigo, nombre)VALUES($1, $2)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Codigo, Nombre)
		log.Println("Nuevo Registro:" + Codigo + "," + Nombre)
	}
	http.Redirect(w, r, "/BodegaLista", 301)
}

// BODEGA EXISTE
func BodegaExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM bodega  WHERE codigo=$1", Codigo)
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

// BODEGA EDITAR
func BodegaEditar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/bodega/bodegaEditar.html")
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	selDB, err := db.Query("SELECT * FROM bodega WHERE codigo=$1", Codigo)
	if err != nil {
		panic(err.Error())
	}
	emp := bodega{}
	for selDB.Next() {
		var codigo string
		var nombre string
		err = selDB.Scan(&codigo, &nombre)
		if err != nil {
			panic(err.Error())
		}
		emp.Codigo = codigo
		emp.Nombre = nombre
	}
	varmap := map[string]interface{}{
		"emp":     emp,
		"hosting": ruta,
	}
	//vistaBodega.ExecuteTemplate(w, "BodegaEditar", varmap)
	tmp.Execute(w, varmap)
}

// BODEGA ACTUALIZAR
func BodegaActualizar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		codigo := r.FormValue("Codigo")
		nombre := r.FormValue("Nombre")
		nombre = Titulo(nombre)
		insForm, err := db.Prepare("UPDATE bodega set	nombre=$2  " + " WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(codigo, nombre)
		log.Println("Registro Actualizado:" + codigo + "," +
			"" + nombre)
	}
	http.Redirect(w, r, "/BodegaLista", 301)
}

// BODEGA BORRAR
func BodegaBorrar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/bodega/bodegaBorrar.html")
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	selDB, err := db.Query("SELECT * FROM bodega WHERE codigo=$1", Codigo)
	if err != nil {
		panic(err.Error())
	}
	emp := bodega{}
	for selDB.Next() {
		var codigo string
		var nombre string
		err = selDB.Scan(&codigo, &nombre)
		if err != nil {
			panic(err.Error())
		}
		emp.Codigo = codigo
		emp.Nombre = nombre
	}
	varmap := map[string]interface{}{
		"emp":     emp,
		"hosting": ruta,
	}
	tmp.Execute(w, varmap)
}

// BODEGA ELIMINAR
func BodegaEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	emp := mux.Vars(r)["codigo"]
	//Codigo, _ := strconv.ParseInt(emp, 10, 0)
	delForm, err := db.Prepare("DELETE from bodega WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("Registro Eliminado" + emp)
	http.Redirect(w, r, "/BodegaLista", 301)
}

// INICIA BODEGA PDF
func BodegaPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := bodega{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Get(&t, "SELECT * FROM bodega where codigo=$1", Codigo)
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
		pdf.CellFormat(190, 10, "Datos Bodega", "0", 0,
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

	pdf.SetFooterFunc(func() {
		pdf.SetY(-20)
		pdf.SetFont("Arial", "", 9)
		pdf.SetX(30)
		pdf.CellFormat(90, 10, "Sadconf.com", "", 0,
			"L", false, 0, "")
		pdf.CellFormat(80, 10, fmt.Sprintf(" %d de {nb}", pdf.PageNo()), "",
			0, "R", false, 0, "")


		// ROTACION DE TEXTO
		margin := 25.0
		markFontHt := 50.0
		markLineHt := pdf.PointToUnitConvert(markFontHt)
		markY := (297.0 - markLineHt) / 2.0
		ctrX := 210.0 / 2.0
		ctrY := 297.0 / 2.0
		pdf.SetFont("Arial", "B", markFontHt)
		// COLOR ROJO
		//pdf.SetTextColor(255, 0, 0)
		// COLOR VERDE
		//pdf.SetTextColor(0, 255, 0)
		// COLOR AMARILLO
		//pdf.SetTextColor(255, 255, 0)
		// COLOR AZUL
		//pdf.SetTextColor(0, 0, 255)
		// COLOR CIAN
		//pdf.SetTextColor(0, 255, 255)
		// COLOR MAGENTA
		//pdf.SetTextColor(255, 0, 255)
		// COLOR NEGRO
		//pdf.SetTextColor(0, 0, 0)
		// COLOR BLANCO
		//pdf.SetTextColor(255, 255, 255)
		pdf.SetXY(margin, markY)
		pdf.TransformBegin()
		pdf.TransformRotate(45, ctrX, ctrY)
		//pdf.CellFormat(0, markLineHt, "ANULADA", "", 0, "C", false, 0, "")
		pdf.TransformEnd()



	})
	err1 = pdf.Output(&buf)
	if err1 != nil {
		panic(err1.Error())
	}
	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	w.Write(buf.Bytes())


}
// TERMINA BODEGA PDF


