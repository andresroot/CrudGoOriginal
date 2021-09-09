package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	//"strings"
)

// CONFIGURACIONINVENTARIO TABLA
type configuracioncontabilidad struct {
	// INICIA ESTRUCTURA
	Pagocuentaefectivo     			string
	Pagonombreefectivo				string
	Pagocuentasaldoafavor     		string
	Pagonombresaldoafavor			string
	Pagocuentatdebito     			string
	Pagonombretdebito				string
	Pagocuentatcredito     			string
	Pagonombretcredito				string
	Pagocuentaconsignacion     		string
	Pagonombreconsignacion			string
	Pagocuentamayorvalor     		string
	Pagonombremayorvalor			string
	Pagocuentamenorvalor     		string
	Pagonombremenorvalor			string
	Phinicial						string
	Textodescuento1 				string
	Textodescuento2 				string
	Textodescuento3 				string
	Textoaviso1 	 				string
	Textoaviso2 	 				string
	Textoaviso3 	 				string
	Textoaviso4 	 				string
}

// CONFIGURACIONINVENTARIO NUEVO
func ConfiguracioncontabilidadNuevo(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/configuracioncontabilidad/configuracioncontabilidadNuevo.html",
		"vista/autocompleta/autocompletaPlandecuentaempresa."+
			"html")
	db := dbConn()
	//Compracuenta19 := mux.Vars(r)["compracuenta19"]
	panel := mux.Vars(r)["panel"]

	selDB, err := db.Query("SELECT * FROM configuracioncontabilidad ")
	if err != nil {
		panic(err.Error())
	}
	emp := configuracioncontabilidad{}
	for selDB.Next() {
		//  INICIA NUEVA COMPRA
		var	Pagocuentaefectivo             	string
		var	Pagonombreefectivo            	string
		var	Pagocuentasaldoafavor           string
		var	Pagonombresaldoafavor          	string
		var	Pagocuentatdebito              	string
		var	Pagonombretdebito              	string
		var	Pagocuentatcredito             	string
		var	Pagonombretcredito            	string
		var	Pagocuentaconsignacion         	string
		var	Pagonombreconsignacion         	string
		var	Pagocuentamayorvalor            string
		var	Pagonombremayorvalor          	string
		var	Pagocuentamenorvalor            string
		var	Pagonombremenorvalor          	string
		var	Phinicial			          	string
		var	Textodescuento1		          	string
		var	Textodescuento2		          	string
		var	Textodescuento3		          	string
		var	Textoaviso1 		          	string
		var	Textoaviso2	    	          	string
		var	Textoaviso3		            	string
		var	Textoaviso4		            	string
		//  TERMINA NUEVO SERVICIO
		err = selDB.Scan(
			//  INICIA SCAN COMPRA
			&Pagocuentaefectivo,
			&Pagonombreefectivo,
			&Pagocuentasaldoafavor,
			&Pagonombresaldoafavor,
			&Pagocuentatdebito,
			&Pagonombretdebito,
			&Pagocuentatcredito,
			&Pagonombretcredito,
			&Pagocuentaconsignacion,
			&Pagonombreconsignacion,
			&Pagocuentamayorvalor,
			&Pagonombremayorvalor,
			&Pagocuentamenorvalor,
			&Pagonombremenorvalor,
			&Phinicial,
			&Textodescuento1,
			&Textodescuento2,
			&Textodescuento3,
			&Textoaviso1,
			&Textoaviso2,
			&Textoaviso3,
			&Textoaviso4)

			emp.Pagocuentaefectivo=Pagocuentaefectivo
			emp.Pagonombreefectivo=Pagonombreefectivo
			emp.Pagocuentasaldoafavor=Pagocuentasaldoafavor
			emp.Pagonombresaldoafavor=Pagonombresaldoafavor
			emp.Pagocuentatdebito=Pagocuentatdebito
			emp.Pagonombretdebito=Pagonombretdebito
			emp.Pagocuentatcredito=Pagocuentatcredito
			emp.Pagonombretcredito=Pagonombretcredito
			emp.Pagocuentaconsignacion=Pagocuentaconsignacion
			emp.Pagonombreconsignacion=Pagonombreconsignacion
			emp.Pagocuentamayorvalor=Pagocuentamayorvalor
			emp.Pagonombremayorvalor=Pagonombremayorvalor
			emp.Pagocuentamenorvalor=Pagocuentamenorvalor
		    emp.Pagonombremenorvalor=Pagonombremenorvalor
			emp.Phinicial=Phinicial
			emp.Textodescuento1=Textodescuento1
			emp.Textodescuento2=Textodescuento2
			emp.Textodescuento3=Textodescuento3
			emp.Textoaviso1=Textoaviso1
			emp.Textoaviso2=Textoaviso2
			emp.Textoaviso3=Textoaviso3
			emp.Textoaviso4=Textoaviso4
	}

	varmap := map[string]interface{}{
		"parametro":     emp,
		"hosting": ruta,
		"panel":panel,
	}

	tmp.Execute(w, varmap)

}

// CONFIGURACIONINVENTARIO INSERTAR
func ConfiguracioncontabilidadInsertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	panel:=r.FormValue("panel")
	if r.Method == "POST" {
		// INICIA INSERTAR CONFIGURACIONCONTABILIDAD
		Pagocuentaefectivo:=r.FormValue("Pagocuentaefectivo")
		Pagonombreefectivo:=r.FormValue("Pagonombreefectivo")
		Pagocuentasaldoafavor:=r.FormValue("Pagocuentasaldoafavor")
		Pagonombresaldoafavor:=r.FormValue("Pagonombresaldoafavor")
		Pagocuentatdebito:=r.FormValue("Pagocuentatdebito")
		Pagonombretdebito:=r.FormValue("Pagonombretdebito")
		Pagocuentatcredito:=r.FormValue("Pagocuentatcredito")
		Pagonombretcredito:=r.FormValue("Pagonombretcredito")
		Pagocuentaconsignacion:=r.FormValue("Pagocuentaconsignacion")
		Pagonombreconsignacion:=r.FormValue("Pagonombreconsignacion")
		Pagocuentamayorvalor:=r.FormValue("Pagocuentamayorvalor")
		Pagonombremayorvalor:=r.FormValue("Pagonombremayorvalor")
		Pagocuentamenorvalor:=r.FormValue("Pagocuentamenorvalor")
		Pagonombremenorvalor:=r.FormValue("Pagonombremenorvalor")
		Phinicial:=r.FormValue("Phinicial")
		Textodescuento1:=r.FormValue("Textodescuento1")
		Textodescuento2:=r.FormValue("Textodescuento2")
		Textodescuento3:=r.FormValue("Textodescuento3")
		Textoaviso1:=r.FormValue("Textoaviso1")
		Textoaviso2:=r.FormValue("Textoaviso2")
		Textoaviso3:=r.FormValue("Textoaviso3")
		Textoaviso4:=r.FormValue("Textoaviso4")
		Pagonombreefectivo=Titulo(Pagonombreefectivo)
		Pagonombresaldoafavor=Titulo(Pagonombresaldoafavor)
		Pagonombretdebito=Titulo(Pagonombretdebito)
		Pagonombretcredito=Titulo(Pagonombretcredito)
		Pagonombreconsignacion=Titulo(Pagonombreconsignacion)
		Pagonombremayorvalor=Titulo(Pagonombremayorvalor)
		Pagonombremenorvalor=Titulo(Pagonombremenorvalor)
		Textodescuento1=Titulo(Textodescuento1)
		Textodescuento2=Titulo(Textodescuento2)
		Textodescuento3=Titulo(Textodescuento3)
		Textoaviso1=Mayuscula(Textoaviso1)
		Textoaviso2=Mayuscula(Textoaviso2)
		Textoaviso3=Mayuscula(Textoaviso3)
		Textoaviso4=Mayuscula(Textoaviso4)

		// TERMINA INSERTAR COMPRA

		var consulta="INSERT INTO configuracioncontabilidad("
		// INICIA CONSULTA COMPRA
		consulta+="Pagocuentaefectivo,"
		consulta+="Pagonombreefectivo,"
		consulta+="Pagocuentasaldoafavor,"
		consulta+="Pagonombresaldoafavor,"
		consulta+="Pagocuentatdebito,"
		consulta+="Pagonombretdebito,"
		consulta+="Pagocuentatcredito,"
		consulta+="Pagonombretcredito,"
		consulta+="Pagocuentaconsignacion,"
		consulta+="Pagonombreconsignacion,"
		consulta+="Pagocuentamayorvalor,"
		consulta+="Pagonombremayorvalor,"
		consulta+="Pagocuentamenorvalor,"
		consulta+="Pagonombremenorvalor,"
		consulta+="Phinicial,"
		consulta+="Textodescuento1,"
		consulta+="Textodescuento2,"
		consulta+="Textodescuento3,"
		consulta+="Textoaviso1,"
		consulta+="Textoaviso2,"
		consulta+="Textoaviso3,"
		consulta+="Textoaviso4"
		consulta+=")VALUES("
		consulta+=parametros(22)
		consulta+=")"

		delForm, err := db.Prepare("DELETE from configuracioncontabilidad")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec()
		insForm, err := db.Prepare(consulta)
		if err != nil {
			panic(err.Error())
		}
		_, err =insForm.Exec(
			// INICIA BORRAR COMPRA
			Pagocuentaefectivo,
			Pagonombreefectivo,
			Pagocuentasaldoafavor,
			Pagonombresaldoafavor,
			Pagocuentatdebito,
			Pagonombretdebito,
			Pagocuentatcredito,
			Pagonombretcredito,
			Pagocuentaconsignacion,
			Pagonombreconsignacion,
			Pagocuentamayorvalor,
			Pagonombremayorvalor,
			Pagocuentamenorvalor,
			Pagonombremenorvalor,
			Phinicial,
			Textodescuento1,
			Textodescuento2,
			Textodescuento3,
			Textoaviso1,
			Textoaviso2,
			Textoaviso3,
			Textoaviso4)

			// TERMINA BORRAR SERVICIO
		if err != nil {
			panic(err)
		}
		log.Println("Nuevo Registro:" + Pagocuentaefectivo + "," + Pagonombreefectivo)
	}
	http.Redirect(w, r, "/ConfiguracioncontabilidadNuevo/"+panel, 301)
}

// INICIA CONFIGURACIONINVENTARIO PDF
func ConfiguracioncontabilidadPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Pagocuentaefectivo := mux.Vars(r)["pagocuentaefectivo"]
	t := configuracioncontabilidad{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Get(&t, "SELECT * FROM configuracioncontabilidad where pagocuentaefectivo=$1",Pagocuentaefectivo)
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
		pdf.CellFormat(190, 10, "Datos Centro de Costos", "0", 0,
			"C", false, 0, "")
		pdf.Ln(10)
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Cuenta", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Pagocuentaefectivo, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Nombre:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Pagonombreefectivo, "", 0,
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
// TERMINA CONFIGURACIONINVENTARIO PDF

