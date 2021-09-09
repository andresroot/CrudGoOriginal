package main

// INICIA TERCERO IMPORTAR PAQUETES
import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// TERMINA  a IMPORTAR PAQUETES

// INICIA TERCERO ESTRUCTURA JSON
type terceroJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// TERMINA TERCERO ESTRUCTURA JSON
type tercerolista struct {
	Codigo          string
	Dv              string
	Nombre          string

}

// INICIA TERCERO ESTRUCTURA
type tercero struct {
	Codigo          string
	Dv              string
	Nombre          string
	Juridica        string
	PrimerNombre    string
	SegundoNombre   string
	PrimerApellido  string
	SegundoApellido string
	Direccion       string
	Barrio          string
	Telefono1       string
	Telefono2       string
	Email1          string
	Email2          string
	Contacto        string
	Rut             string
	Descuento1      string
	Descuento2      string
	Cuotap          string
	Cuota1          string
	Cuota2          string
	Cuota3          string
	Area            string
	Factor          string
	Matricula       string
	Catastral       string
	Banco           string
	PhCodigo        string
	PhDv            string
	PhNombre        string
	Ciudad          string
	Documento       string
	Fiscal          string
	Regimen         string
	Tipo            string
	Ica				string
}

// TERMINA TERCERO ESTRUCTURA

// INICIA TERCERO LISTA
func TerceroLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/tercero/terceroLista.html")
	log.Println("Error tercero 0")
	db := dbConn()
	res := []tercerolista{}
	db.Select(&res, "SELECT codigo,dv,nombre FROM tercero ORDER BY cast(codigo as float) ASC")
	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error tercero888")
	tmp.Execute(w, varmap)
}

// TERMINA TERCERO LISTA


//INICIA TERCERO NUEVO
func TerceroNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error tercero nuevo 1")
	Codigo := mux.Vars(r)["codigo"]
	Panel := mux.Vars(r)["panel"]
	Elemento := mux.Vars(r)["elemento"]
	t := tercero{}


	log.Println("Error tercero nuevo 2")
	parametros := map[string]interface{}{
		// INICIA TERCERO NUEVO AUTOCOMPLETADO
		"Codigo":                  Codigo,
		"Panel":                   Panel,
		"Elemento":                Elemento,
		"hosting":                 ruta,
		"ciudad":                  ListaCiudad(),
		"tipoorganizacion":        ListaTipoOrganizacion(),
		"regimenfiscal":           ListaRegimenFiscal(),
		"responsabilidadfiscal":   ListaResponsabilidadFiscal(),
		"documentoidentificacion": ListaDocumentoIdentificacion(),
		"copiar": "False",
		"emp": t,
		// TERMINA TERCERO NUEVO AUTOCOMPLETADO
	}
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html", "vista/tercero/terceroNuevo.html", "vista/autocompleta/autocompletaTercero.html")
	log.Println("Error tercero nuevo 3")
	tmp.Execute(w, parametros)
}

// INICIA TERCERO DUPLICAR
func TerceroNuevoCopia(w http.ResponseWriter, r *http.Request) {
	log.Println("Error tercero nuevo 1")
	Codigo :="False"
	Panel := "False"
	Elemento := "False"

	copiarCodigo := Quitacoma(mux.Vars(r)["copiacodigo"])
	log.Println("inicio tercero editar" + Codigo)

	db := dbConn()
	t := tercero{}

	if copiarCodigo == "False" {

	} else {
		// traer comprobante

		err := db.Get(&t, "SELECT * FROM tercero WHERE codigo=$1", copiarCodigo)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("Error tercero nuevo 2")
	parametros := map[string]interface{}{
		// INICIA TERCERO NUEVO AUTOCOMPLETADO
		"Codigo":                  Codigo,
		"Panel":                   Panel,
		"Elemento":                Elemento,
		"hosting":                 ruta,
		"ciudad":                  ListaCiudad(),
		"tipoorganizacion":        ListaTipoOrganizacion(),
		"regimenfiscal":           ListaRegimenFiscal(),
		"responsabilidadfiscal":   ListaResponsabilidadFiscal(),
		"documentoidentificacion": ListaDocumentoIdentificacion(),
		"emp": t,
		"copiar": "True",

		// TERMINA TERCERO NUEVO AUTOCOMPLETADO
	}
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html", "vista/tercero/terceroNuevo.html", "vista/autocompleta/autocompletaTercero.html")
	log.Println("Error tercero nuevo 3")
	tmp.Execute(w, parametros)
}


// TERMINA TERCERO NUEVO

// INICIA TERCERO INSERTAR
func TerceroInsertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	var t tercero
	err = decoder.Decode(&t, r.PostForm)
	if err != nil {
		// Handle error
		panic(err.Error())
	}
	var q string
	q = "insert into tercero ("
	q += "Codigo,"
	q += "Dv,"
	q += "Nombre,"
	q += "Juridica,"
	q += "PrimerNombre,"
	q += "SegundoNombre,"
	q += "PrimerApellido,"
	q += "SegundoApellido,"
	q += "Direccion,"
	q += "Barrio,"
	q += "Telefono1,"
	q += "Telefono2,"
	q += "Email1,"
	q += "Email2,"
	q += "Contacto,"
	q += "Rut,"
	q += "Ciudad,"
	q += "Documento,"
	q += "Fiscal,"
	q += "Regimen,"
	q += "Tipo,"
	q += "Ica,"
	q += "Descuento1,"
	q += "Descuento2,"
	q += "Cuotap,"
	q += "Cuota1,"
	q += "Cuota2,"
	q += "Cuota3,"
	q += "Area,"
	q += "Factor,"
	q += "Matricula,"
	q += "Catastral,"
	q += "Banco,"
	q += "PhCodigo,"
	q += "PhDv,"
	q += "PhNombre"

	q += " ) values("
	q += parametros(36)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	// INICIA GRABAR TERCERO INSERTAR
	t.Codigo = Quitacoma(t.Codigo)
	t.PhCodigo = Quitacoma(t.PhCodigo)
	//t.Descuento1 = Quitacoma(t.Descuento1)
	//t.Descuento2 = Quitacoma(t.Descuento2)
	//t.Cuotap = Quitacoma(t.Cuotap)
	//t.Cuota1 = Quitacoma(t.Cuota1)
	//t.Cuota2 = Quitacoma(t.Cuota2)
	//t.Cuota3 = Quitacoma(t.Cuota3)
	t.Nombre = t.Juridica
	if t.Tipo == "2" {
		t.PrimerNombre = Titulo(t.PrimerNombre)
		t.SegundoNombre = Titulo(t.SegundoNombre)
		t.PrimerApellido = Titulo(t.PrimerApellido)
		t.SegundoApellido = Titulo(t.SegundoApellido)
		t.Nombre = t.PrimerNombre + " " + t.SegundoNombre + " " + t.PrimerApellido + " " + t.SegundoApellido
	} else {
		t.Juridica = Titulo(t.Juridica)
		t.Nombre = t.Juridica
	}
	t.Nombre = Titulo(t.Nombre)
	t.Direccion = Titulo(t.Direccion)
	t.Barrio = Titulo(t.Barrio)
	t.Contacto = Titulo(t.Contacto)
	t.Banco = Titulo(t.Banco)
	t.PhNombre = Titulo(t.PhNombre)
	t.Email1 = Minuscula(t.Email1)
	t.Email2 = Minuscula(t.Email2)
	// TERMINA TERCERO GRABAR INSERTAR
	_, err = insForm.Exec(
		t.Codigo,
		t.Dv,
		t.Nombre,
		t.Juridica,
		t.PrimerNombre,
		t.SegundoNombre,
		t.PrimerApellido,
		t.SegundoApellido,
		t.Direccion,
		t.Barrio,
		t.Telefono1,
		t.Telefono2,
		t.Email1,
		t.Email2,
		t.Contacto,
		t.Rut,
		t.Ciudad,
		t.Documento,
		t.Fiscal,
		t.Regimen,
		t.Tipo,
		t.Ica,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		t.PhCodigo,
		t.PhDv,
		t.PhNombre)

	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/TerceroLista", 301)
}

// TERMINA TERCERO INSERTAR

// INICIA TERCERO BUSCAR
func TerceroBuscar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	Codigo = Quitacoma(Codigo)
	selDB, err := db.Query("SELECT codigo,"+
		"nombre FROM tercero where codigo LIKE '%' || $1 || '%'  or  nombre LIKE '%' || $1 || '%' ORDER BY"+
		" codigo DESC", Codigo)
	if err != nil {
		panic(err.Error())
	}
	var resJson []terceroJson
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
		resJson = append(resJson, terceroJson{id, label, value, nombre})
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

// TERMINA TERCERO BUSCAR

// INICIA TERCERO EXISTE
func TerceroExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM tercero  WHERE codigo=$1", Codigo)
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

// TERMINA TERCERO EXISTE

// INICIA TERCERO ACTUAL
func TerceroActual(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	Codigo = Quitacoma(Codigo)

	t := tercero{}
	var res []tercero
	err := db.Get(&t, "SELECT * FROM tercero where codigo=$1", Codigo)

	switch err {
	case nil:
		log.Printf("tercero found: %+v\n", t)
	case sql.ErrNoRows:
		log.Println("tercero NOT found, no error")
	default:
		log.Printf("tercero error: %s\n", err)
	}

	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
	res = append(res, t)
	data, err := json.Marshal(res)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// TERMINA TERCERO ACTUAL

// INICIA TERCERO EDITAR
func TerceroEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio tercero editar" + Codigo)
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/tercero/terceroEditar.html",
		"vista/autocompleta/autocompletaTercero.html")
	db := dbConn()
	t := tercero{}
	err := db.Get(&t, "SELECT * FROM tercero where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
	varmap := map[string]interface{}{
		// INICIA TERCERO EDITAR AUTOCOMPLETADO
		"emp":                     t,
		"hosting":                 ruta,
		"ciudad":                  ListaCiudad(),
		"tipoorganizacion":        ListaTipoOrganizacion(),
		"regimenfiscal":           ListaRegimenFiscal(),
		"responsabilidadfiscal":   ListaResponsabilidadFiscal(),
		"documentoidentificacion": ListaDocumentoIdentificacion(),
		// TERMINA TERCERO EDITAR AUTOCOMPLETADO
	}
	tmp.Execute(w, varmap)
}

// INICIA TERCERO ACTUALIZAR
func TerceroActualizar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
		// Handle error
	}
	var t tercero
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(&t, r.PostForm)
	if err != nil {
		// Handle error
		panic(err.Error())
	}
	var q string
	q = "UPDATE tercero set "
	q += "Dv=$2,"
	q += "Nombre=$3,"
	q += "Juridica=$4,"
	q += "PrimerNombre=$5,"
	q += "SegundoNombre=$6,"
	q += "PrimerApellido=$7,"
	q += "SegundoApellido=$8,"
	q += "Direccion=$9,"
	q += "Barrio=$10,"
	q += "Telefono1=$11,"
	q += "Telefono2=$12,"
	q += "Email1=$13,"
	q += "Email2=$14,"
	q += "Contacto=$15,"
	q += "Rut=$16,"
	q += "Ciudad=$17,"
	q += "Documento=$18,"
	q += "Fiscal=$19,"
	q += "Regimen=$20,"
	q += "Ica=$21,"
	q += "Tipo=$22,"
	q += "PhCodigo=$23,"
	q += "PhDv=$24,"
	q += "PhNombre=$25"
	q += " where "
	q += "Codigo=$1"

	log.Println("cadena" + q)

	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	// INICIA GRABAR TERCERO ACTUALIZAR
	t.Codigo = Quitacoma(t.Codigo)
	t.PhCodigo = Quitacoma(t.PhCodigo)
	//t.Descuento1 = Quitacoma(t.Descuento1)
	//t.Descuento2 = Quitacoma(t.Descuento2)
	//t.Cuotap = Quitacoma(t.Cuotap)
	//t.Cuota1 = Quitacoma(t.Cuota1)
	//t.Cuota2 = Quitacoma(t.Cuota2)
	//t.Cuota3 = Quitacoma(t.Cuota3)
	t.Nombre = t.Juridica
	if t.Tipo == "2" {
		t.PrimerNombre = Titulo(t.PrimerNombre)
		t.SegundoNombre = Titulo(t.SegundoNombre)
		t.PrimerApellido = Titulo(t.PrimerApellido)
		t.SegundoApellido = Titulo(t.SegundoApellido)
		t.Nombre = t.PrimerNombre + " " + t.SegundoNombre + " " + t.PrimerApellido + " " + t.SegundoApellido
	} else {
		t.Juridica = Titulo(t.Juridica)
		t.Nombre = t.Juridica
	}
	t.Nombre = Titulo(t.Nombre)
	t.Direccion = Titulo(t.Direccion)
	t.Barrio = Titulo(t.Barrio)
	t.Contacto = Titulo(t.Contacto)
	t.PhNombre = Titulo(t.PhNombre)
	t.Email1 = Minuscula(t.Email1)
	t.Email2 = Minuscula(t.Email2)
	// TERMINA GRABAR TERCERO ACTUALIZAR

	_, err = insForm.Exec(
		t.Codigo,
		t.Dv,
		t.Nombre,
		t.Juridica,
		t.PrimerNombre,
		t.SegundoNombre,
		t.PrimerApellido,
		t.SegundoApellido,
		t.Direccion,
		t.Barrio,
		t.Telefono1,
		t.Telefono2,
		t.Email1,
		t.Email2,
		t.Contacto,
		t.Rut,
		t.Ciudad,
		t.Documento,
		t.Fiscal,
		t.Regimen,
		t.Ica,
		t.Tipo,
		t.PhCodigo,
		t.PhDv,
		t.PhNombre)

	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/TerceroLista", 301)

}

// INICIA TERCERO BORRAR
func TerceroBorrar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/tercero/terceroBorrar.html")
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio tercero borrar" + Codigo)
	db := dbConn()
	t := tercero{}
	err := db.Get(&t, "SELECT * FROM tercero where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codigo nombre99 borrar" + t.Codigo + t.Nombre)
	varmap := map[string]interface{}{
		// INICIA TERCERO BORRAR AUTOCOMPLETADO
		"emp":                     t,
		"hosting":                 ruta,
		"ciudad":                  ListaCiudad(),
		"tipoorganizacion":        ListaTipoOrganizacion(),
		"regimenfiscal":           ListaRegimenFiscal(),
		"responsabilidadfiscal":   ListaResponsabilidadFiscal(),
		"documentoidentificacion": ListaDocumentoIdentificacion(),
		// TERMINA TERCERO BORRAR AUTOCOMPLETADO
	}
	tmp.Execute(w, varmap)
}

// INICIA TERCERO ELIMINAR
func TerceroEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	emp := mux.Vars(r)["codigo"]
	delForm, err := db.Prepare("DELETE from tercero WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("Registro Eliminado" + emp)
	http.Redirect(w, r, "/TerceroLista", 301)
}

// INICIA TERCERO PDF
func TerceroPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := tercero{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Get(&t, "SELECT * FROM tercero where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	var buf bytes.Buffer
	var err1 error
	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)
	ene := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetHeaderFunc(func() {
		pdf.Image(imageFile("logo.png"), 20, 30, 40, 0, false,
			"", 0, "")
		pdf.SetY(17)
		pdf.SetFont("Arial", "", 10)
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
		pdf.CellFormat(190, 10, e.Telefono1+" - "+e.Telefono2, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+ " - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")
		pdf.Ln(10)

		// RELLENO TITULO
		pdf.SetX(20)
		pdf.SetFillColor(224,231,239)
		pdf.SetTextColor(0,0,0)

		pdf.SetX(20)
		pdf.CellFormat(184, 5, "DATOS TERCERO", "0", 0,
			"C", true, 0, "")
		pdf.Ln(8)
	})

	pdf.SetTextColor(0,0,0)
	pdf.SetX(21)
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetX(21)

	pdf.CellFormat(40, 4, "Nit. No.", "", 0,
		"", false, 0, "")
	pdf.CellFormat(142, 4, Coma(t.Codigo)+ " - "+t.Dv, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Tipo:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, TraerTipo(t.Tipo), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Empresa:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Juridica, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Primer Nombre:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.PrimerNombre, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Segundo Nombre:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.SegundoNombre, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Primer Apellido:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.PrimerApellido, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Segundo Apellido:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.SegundoApellido, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Ret. Ica:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, (t.Ica), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Rut:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Rut, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Direccion:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Direccion, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Barrio:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Barrio, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Ciudad:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Ciudad, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "E-mail 1:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Email1, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "E-mail 2:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Email2, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Contacto:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Contacto, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Documento:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, TraerDocumentoIdentificacion(t.Documento), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Regimen:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, TraerRegimen(t.Regimen), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Fiscal:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, TraerFiscal(t.Fiscal), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Codigo No.", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, Coma(t.PhCodigo) + " - "+t.PhDv, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Nombre:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.PhNombre, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Descuento 1:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Descuento1, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Descuento 2:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Descuento2, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Cuota P:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Cuotap, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Cuota 1:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Cuota1, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Cuota 2:", "0", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Cuota2, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Cuota 3:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Cuota3, "", 0,
		"", false, 0, "")

	pdf.SetFooterFunc(func() {
		pdf.SetTextColor(0, 0, 0)
		pdf.SetY(252)
		pdf.SetFont("Arial", "", 9)
		pdf.SetX(20)

		// LINEA
		pdf.Line(20,259,204,259)
		pdf.Ln(6)
		pdf.SetX(20)
		pdf.CellFormat(90, 10, "Sadconf.com", "", 0,
			"L", false, 0, "")
		pdf.SetX(129)
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

// INICIA TERCEROS TODOS PDF
func TerceroTodosCabecera(pdf *gofpdf.Fpdf){
	pdf.SetFont("Arial", "", 10)
	// RELLENO TITULO
	pdf.SetY(50)
	pdf.SetFillColor(224,231,239)
	pdf.SetTextColor(0,0,0)
	pdf.Ln(7)
	pdf.SetX(20)
	pdf.CellFormat(181, 6, "No", "0", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 6, "Codigo", "0", 0,
		"L", false, 0, "")
	pdf.SetX(60)
	pdf.CellFormat(40, 6, "Nombre", "0", 0,
		"L", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(40, 6, "Direccion", "0", 0,
		"L", false, 0, "")
	pdf.SetX(171)
	pdf.CellFormat(40, 6, "Telefono", "0", 0,
		"L", false, 0, "")
	pdf.Ln(8)
}
func TerceroTodosDetalle(pdf *gofpdf.Fpdf,t tercero, a int ){
	pdf.SetFont("Arial", "", 9)

	pdf.SetX(21)
	pdf.CellFormat(181, 4, strconv.Itoa(a), "", 0,
		"L", false, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 4, Coma(t.Codigo)+" - "+t.Dv,"", 0,
		"L", false, 0, "")
	pdf.SetX(60)
	pdf.CellFormat(40, 4, t.Nombre, "", 0,"L", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(40, 4, t.Direccion, "", 0,
		"L", false, 0, "")
	pdf.SetX(155)
	pdf.CellFormat(40, 4, t.Telefono1, "", 0,
		"R", false, 0, "")
	pdf.Ln(4)
}

func TerceroTodosPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	//	Codigo := mux.Vars(r)["codigo"]

	t := []tercero{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Select(&t, "SELECT * FROM tercero ORDER BY cast(codigo as integer) ")
	if err != nil {
		log.Fatalln(err)
	}
	var buf bytes.Buffer
	var err1 error
	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)
	ene := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetHeaderFunc(func() {
		pdf.Image(imageFile("logo.png"), 20, 30, 40, 0, false,
			"", 0, "")
		pdf.SetY(17)
		pdf.SetFont("Arial", "", 10)
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
		pdf.CellFormat(190, 10, e.Telefono1+" "+e.Telefono2, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+ " - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")
		pdf.Ln(6)
		pdf.CellFormat(190, 10, "DATOS TERCERO", "0", 0,
			"C", false, 0, "")
		pdf.Ln(10)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetTextColor(0,0,0)
		pdf.SetY(252)
		pdf.SetFont("Arial", "", 9)
		pdf.SetX(20)

		// LINEA
		pdf.Line(20,259,204,259)
		pdf.Ln(6)
		pdf.SetX(20)
		pdf.CellFormat(90, 10, "Sadconf.com", "", 0,
			"L", false, 0, "")
		pdf.SetX(129)
		pdf.CellFormat(80, 10, fmt.Sprintf(" %d de {nb}", pdf.PageNo()), "",
			0, "R", false, 0, "")
	})

	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(30)

	TerceroTodosCabecera(pdf)
	// tercera pagina

	for i, miFila := range t {
		var a = i + 1
		if math.Mod(float64(a),49)==0 {
			pdf.AliasNbPages("")
			pdf.AddPage()
			pdf.SetFont("Arial", "", 10)
			pdf.SetX(30)
			TerceroTodosCabecera(pdf)
		}
		TerceroTodosDetalle(pdf,miFila,a)
	}
	//BalancePieDePagina(pdf)

	err1 = pdf.Output(&buf)
	if err1 != nil {
		panic(err1.Error())
	}
	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	w.Write(buf.Bytes())
}
// TERMINA TERCERO TODOS PDF

// TERCERO EXCEL
func TerceroExcel(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	t := []tercero{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Select(&t, "SELECT * FROM tercero ORDER BY cast(codigo as integer) ")
	if err != nil {
		log.Fatalln(err)
	}

	f := excelize.NewFile()

	// FUNCION ANCHO DE LA COLUMNA
	if err =f.SetColWidth("Sheet1", "A", "A", 13); err != nil {
		fmt.Println(err)
		return
	}
	if err =f.SetColWidth("Sheet1", "B", "B", 30); err != nil {
		fmt.Println(err)
		return
	}
	if err =f.SetColWidth("Sheet1", "C", "C", 30); err != nil {
		fmt.Println(err)
		return
	}

	if err =f.SetColWidth("Sheet1", "D", "D", 13); err != nil {
		fmt.Println(err)
		return
	}

	// FUNCION PARA UNIR DOS CELDAS
	if err = f.MergeCell("Sheet1", "A1", "D1"); err != nil {
		fmt.Println(err)
		return
	}
	if err = f.MergeCell("Sheet1", "A2", "D2"); err != nil {
		fmt.Println(err)
		return
	}
	if err = f.MergeCell("Sheet1", "A3", "D3"); err != nil {
		fmt.Println(err)
		return
	}
	if err = f.MergeCell("Sheet1", "A4", "D4"); err != nil {
		fmt.Println(err)
		return
	}

	if err = f.MergeCell("Sheet1", "A5", "D5"); err != nil {
		fmt.Println(err)
		return
	}

	if err = f.MergeCell("Sheet1", "A6", "D6"); err != nil {
		fmt.Println(err)
		return
	}

	if err = f.MergeCell("Sheet1", "A7", "D7"); err != nil {
		fmt.Println(err)
		return
	}

	if err = f.MergeCell("Sheet1", "A8", "D8"); err != nil {
		fmt.Println(err)
		return
	}

	estiloTitulo, err := f.NewStyle(`{  "alignment":{"horizontal": "center"},"font":{"bold":false,"italic":false,"family":"Arial","size":10,"color":"##000000"}}`)

	// titulo
	f.SetCellValue("Sheet1", "A1", e.Nombre)
	f.SetCellValue("Sheet1", "A2","Nit No. "+Coma(e.Codigo)+" - "+e.Dv)
	f.SetCellValue("Sheet1", "A3",e.Iva+" - "+e.ReteIva)
	f.SetCellValue("Sheet1", "A4","Actividad Ica - "+e.ActividadIca)
	f.SetCellValue("Sheet1", "A5",e.Direccion)
	f.SetCellValue("Sheet1", "A6",(e.Telefono1+" - "+e.Telefono2))
	f.SetCellValue("Sheet1", "A7",(c.NombreCiudad+" - "+c.NombreDepartamento))
	f.SetCellValue("Sheet1", "A8","LISTADO DE TERCEROS")

	f.SetCellStyle("Sheet1","A1","A1",estiloTitulo)
	f.SetCellStyle("Sheet1","A2","A2",estiloTitulo)
	f.SetCellStyle("Sheet1","A3","A3",estiloTitulo)
	f.SetCellStyle("Sheet1","A4","A4",estiloTitulo)
	f.SetCellStyle("Sheet1","A5","A5",estiloTitulo)
	f.SetCellStyle("Sheet1","A6","A6",estiloTitulo)
	f.SetCellStyle("Sheet1","A7","A7",estiloTitulo)
	f.SetCellStyle("Sheet1","A8","A8",estiloTitulo)

	var filaExcel=10

	estiloTexto, err := f.NewStyle(`{"font":{"bold":false,"italic":false,"family":"Arial","size":10,"color":"#000000"}}`)

	estiloCabecera, err := f.NewStyle(`{
"alignment":{"horizontal":"center"},
    "border": [
    {
        "type": "left",
        "color": "#000000",
        "style": 1
    },
    {
        "type": "top",
        "color": "#000000",
        "style": 1
    },
    {
        "type": "bottom",
        "color": "#000000",
        "style": 1
    },
    {
        "type": "right",
        "color": "#000000",
        "style": 1
    }]
}`)
	if err != nil {
		fmt.Println(err)
	}
	estiloNumeroDetalle, err := f.NewStyle(`{"number_format": 3,"font":{"bold":false,"italic":false,"family":"Arial","size":10,"color":"##000000"}}`)

	if err != nil {
		fmt.Println(err)
	}
	//cabecera
	f.SetCellValue("Sheet1", "A"+strconv.Itoa(filaExcel),"Codigo")
	f.SetCellValue("Sheet1", "B"+strconv.Itoa(filaExcel), "Nombre")
	f.SetCellValue("Sheet1", "C"+strconv.Itoa(filaExcel), "Consecutivo")
	f.SetCellValue("Sheet1", "D"+strconv.Itoa(filaExcel), "Inicial")


	f.SetCellStyle("Sheet1","A"+strconv.Itoa(filaExcel),"A"+strconv.Itoa(filaExcel),estiloCabecera)
	f.SetCellStyle("Sheet1","B"+strconv.Itoa(filaExcel),"B"+strconv.Itoa(filaExcel),estiloCabecera)
	f.SetCellStyle("Sheet1","C"+strconv.Itoa(filaExcel),"C"+strconv.Itoa(filaExcel),estiloCabecera)
	f.SetCellStyle("Sheet1","D"+strconv.Itoa(filaExcel),"D"+strconv.Itoa(filaExcel),estiloCabecera)
	filaExcel++

	for i, miFila := range t{
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(filaExcel+i), Entero(miFila.Codigo))
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(filaExcel+i), miFila.Nombre)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(filaExcel+i), miFila.Direccion)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(filaExcel+i), miFila.Telefono1)

		f.SetCellStyle("Sheet1","A"+strconv.Itoa(filaExcel+i),"A"+strconv.Itoa(filaExcel+i),estiloNumeroDetalle)
		f.SetCellStyle("Sheet1","B"+strconv.Itoa(filaExcel+i),"B"+strconv.Itoa(filaExcel+i),estiloTexto)
		f.SetCellStyle("Sheet1","C"+strconv.Itoa(filaExcel+i),"C"+strconv.Itoa(filaExcel+i),estiloTexto)
		f.SetCellStyle("Sheet1","D"+strconv.Itoa(filaExcel+i),"D"+strconv.Itoa(filaExcel+i),estiloNumeroDetalle)

		//van=i
	}

	// LIENA FINAL
	//a=strconv.Itoa(van+1+filaExcel)
	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=userInputData.xlsx")
	w.Header().Set("File-Name", "userInputData.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	err = f.Write(w)
	if err != nil {
		panic(err.Error())
	}
}

