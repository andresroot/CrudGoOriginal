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
	"strings"
	"time"
)



// RESOLUCION TABLA
type Resolucionventa struct {
	Codigo        string
	Numero        string
	Prefijo       string
	Tipo          string
	FechaInicial  time.Time
	FechaFinal    time.Time
	NumeroInicial string
	NumeroFinal   string
	NumeroActual  string
	Local         string
	Direccion     string
	Ciudad        string
	Telefono      string
	Informe       string
	Clavetecnica  string
	Idesoftware	  string
	Testid        string
	Pin           string
	Ambiente      string
}

// INICIA RESOLUCION LISTA
func ResolucionventaLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/resolucionventa/resolucionventaLista.html")
	log.Println("Error resolucionventa 0")
	db := dbConn()
	res := []Resolucionventa{}
	db.Select(&res, "SELECT * FROM resolucionventa ORDER BY cast(codigo as integer ) ASC")
	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error resolucionventa888")
	tmp.Execute(w, varmap)
}

// RESOLUCION NUEVO
func ResolucionventaNuevo(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/resolucionventa/resolucionventaNuevo.html")
	parametros := map[string]interface{}{
		// INICIA TERCERO NUEVO AUTOCOMPLETADO
		"hosting":                 ruta,
		"ciudad":                  ListaCiudad(),
		// TERMINA TERCERO NUEVO AUTOCOMPLETADO
	}
	tmp.Execute(w, parametros)
}

// RESOLUCION INSERTAR
func ResolucionventaInsertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Codigo := r.FormValue("Codigo")
		Numero := r.FormValue("Numero")
		Prefijo := r.FormValue("Prefijo")
		Tipo := r.FormValue("Tipo")
		FechaInicial := r.FormValue("FechaInicial")
		FechaFinal := r.FormValue("FechaFinal")
		NumeroInicial := r.FormValue("NumeroInicial")
		NumeroFinal := r.FormValue("NumeroFinal")
		NumeroActual := r.FormValue("NumeroActual")
		Local := r.FormValue("Local")
		Direccion := r.FormValue("Direccion")
		Ciudad := r.FormValue("Ciudad")
		Telefono := r.FormValue("Telefono")
		Informe := r.FormValue("Informe")
		Clavetecnica := r.FormValue("Clavetecnica")
		Idesoftware := r.FormValue("Idesoftware")
		Testid := r.FormValue("Testid")
		Pin := r.FormValue("Pin")
		Ambiente := r.FormValue("Ambiente")

		Prefijo = Mayuscula(Prefijo)
		NumeroInicial = strings.Replace(NumeroInicial, ".", "", -1)
		NumeroFinal = strings.Replace(NumeroFinal, ".", "", -1)
		NumeroActual = strings.Replace(NumeroActual, ".", "", -1)
		Local = Titulo(Local)
		Direccion = Titulo(Direccion)

		var q = "INSERT INTO resolucionventa(" +
			"codigo," +
			"numero," +
			"prefijo," +
			"tipo," +
			"fechainicial," +
			"fechafinal," +
			"numeroinicial," +
			"numerofinal, " +
			"numeroactual, " +
			"local, " +
			"direccion, " +
			"ciudad, " +
			"telefono, " +
			"informe, " +
			"clavetecnica, " +
			"idesoftware, " +
			"testid, " +
			"pin, " +
			"ambiente " +
			")" +
			"VALUES("
		q += parametros(19)
		q += ")"

		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())

		}

		log.Println("fechafinal:" + FechaFinal)

		_, err =insForm.Exec(Codigo,
			Numero,
			Prefijo,
			Tipo,
			FechaInicial,
			FechaFinal,
			NumeroInicial,
			NumeroFinal,
			NumeroActual,
			Local,
			Direccion,
			Ciudad,
			Telefono,
			Informe,
			Clavetecnica,
			Idesoftware,
			Testid,
			Pin,
			Ambiente)

		if err != nil {
			panic(err)
		}
		log.Println("Nuevo Registro:" + Codigo + "," + Numero)
	}
	http.Redirect(w, r, "/ResolucionventaLista", 301)
}

// RESOLUCION EXISTE
func ResolucionventaExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM resolucionventa  WHERE codigo=$1", Codigo)
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

// INICIA RESOLUCION EDITAR
func ResolucionventaEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio resolucionventa editar" + Codigo)
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/resolucionventa/resolucionventaEditar.html")
	db := dbConn()
	t := Resolucionventa{}
	err := db.Get(&t, "SELECT * FROM resolucionventa where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codigo Numero99" + t.Codigo + t.Numero)
	varmap := map[string]interface{}{
		// INICIA TERCERO EDITAR AUTOCOMPLETADO
		"emp":                     t,
		"hosting":                 ruta,
		"ciudad":                  ListaCiudad(),

		// TERMINA TERCERO EDITAR AUTOCOMPLETADO
	}
	tmp.Execute(w, varmap)
}

// INICIA RESOLUCION ACTUALIZAR
func ResolucionventaActualizar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
		// Handle error
	}
	var t Resolucionventa
	err = decoder.Decode(&t, r.PostForm)
	decoder.RegisterConverter(time.Time{}, timeConverter)

	if err != nil {
		// Handle error
		panic(err.Error())
	}
	var q string
	q = "UPDATE resolucionventa set "
	q += "Numero=$2,"
	q += "Prefijo=$3,"
	q += "Tipo=$4,"
	q += "FechaInicial=$5,"
	q += "FechaFinal=$6,"
	q += "NumeroInicial=$7,"
	q += "NumeroFinal=$8,"
	q += "NumeroActual=$9,"
	q += "Local=$10,"
	q += "Direccion=$11,"
	q += "Ciudad=$12,"
	q += "Telefono=$13,"
	q += "Informe=$14,"
	q += "Clavetecnica=$15,"
	q += "Idesoftware=$16,"
	q += "Testid=$17,"
	q += "Pin=$18,"
	q += "Ambiente=$19"
	q += " where "
	q += "Codigo=$1"

	log.Println("cadena" + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	// INICIA GRABAR RESOLUCION ACTUALIZAR
	t.Prefijo = Mayuscula(t.Prefijo)
	t.NumeroInicial = strings.Replace(t.NumeroInicial, ".", "", -1)
	t.NumeroFinal = strings.Replace(t.NumeroFinal, ".", "", -1)
	t.NumeroActual = strings.Replace(t.NumeroActual, ".", "", -1)
	t.Local = Titulo(t.Local)
	t.Direccion = Titulo(t.Direccion)
	// TERMINA GRABAR RESOLUCION ACTUALIZAR
	_, err = insForm.Exec(
		t.Codigo,
		t.Numero,
		t.Prefijo,
		t.Tipo,
		t.FechaInicial,
		t.FechaFinal,
		t.NumeroInicial,
		t.NumeroFinal,
		t.NumeroActual,
		t.Local,
		t.Direccion,
		t.Ciudad,
		t.Telefono,
		t.Informe,
		t.Clavetecnica,
		t.Idesoftware,
		t.Testid,
		t.Pin,
		t.Ambiente)

	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/ResolucionventaLista", 301)

}


// RESOLUCION BORRAR
func ResolucionventaBorrar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/resolucionventa/resolucionventaBorrar.html")
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := Resolucionventa{}
	err := db.Get(&t, "SELECT * FROM resolucionventa where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	varmap := map[string]interface{}{
		"emp":    t,
		"hosting": ruta,
		"ciudad":  ListaCiudad(),
	}
	tmp.Execute(w, varmap)
}

// RESOLUCION ELIMINAR
func ResolucionventaEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	emp := mux.Vars(r)["codigo"]
	delForm, err := db.Prepare("DELETE from resolucionventa WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("Registro Eliminado" + emp)
	http.Redirect(w, r, "/ResolucionventaLista", 301)
}

// INICIA RESOLUCION PDF
func ResolucionventaPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	t := Resolucionventa{}
	var e empresa = ListaEmpresa()
	var c ciudad = TraerCiudad(e.Ciudad)

	err := db.Get(&t, "SELECT * FROM resolucionventa where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	var ciudadresolucionventa ciudad = TraerCiudad(t.Ciudad)
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
		log.Println("tercero 3")
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+" - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")
		log.Println("tercero 4")
		pdf.Ln(10)
		pdf.CellFormat(190, 10, "Datos Resolucion Dian", "0", 0,
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
	pdf.CellFormat(40, 4, "Resolucion No.:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Numero, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Prefijo", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Prefijo, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Tipo", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Tipo, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Fecha Inicial", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.FechaInicial.Format("02/01/2006"), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Fecha Final", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.FechaFinal.Format("02/01/2006"), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Numero Inicial", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, Coma(t.NumeroInicial), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Numero Final", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, Coma(t.NumeroFinal), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Numero Actual", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, Coma(t.NumeroActual), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Nombre Local", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Local, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Direccion", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Direccion, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Ciudad", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, ciudadresolucionventa.NombreCiudad, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Telefono", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Telefono, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Informe", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Informe, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Clave Tecnica", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Clavetecnica, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Ide Software", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Idesoftware, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Test Id", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Testid, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Pin", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Pin, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Ambiente", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Ambiente, "", 0,
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

// TERMINA RESOLUCION PDF
