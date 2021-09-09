package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// CUOTA TABLA
type Cuota struct {
	Filas           string  `json:"Filas"`
	Fecha          string `json:"Fecha"`
	Inicial   		float64  `json:"Inicial"`
	Intereses       float64 `json:"Intereses"`
	Capital         float64  `json:"Capital"`
	Final         	float64  `json:"Final"`
	Cuota 			float64   `json:"Cuota"`
}

// CENTRO KARDEX
func CuotaLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cuota/cuotaLista.html")
	varmap := map[string]interface{}{
		"hosting":  ruta,
	}
	tmp.Execute(w, varmap)
}

func Calcularcuota(Monto string, Plazo string, Intereses string, FechaInicial string)[]Cuota{
	log.Println("Monto : " + Monto)
	log.Println("Plazo : " + Plazo)
	log.Println("Intereses : " + Intereses)
	log.Println("fecha Inicial : " + FechaInicial)

	var monto float64
	var plazo float64
	var intereses float64
	var cuota float64
	var fechainicial time.Time
	monto= Flotante(Monto)

	plazo,err1 := strconv.ParseFloat(Plazo,8)
	if err1 != nil {
		log.Fatalln(err1)
	}
	intereses,err2 := strconv.ParseFloat(Intereses,8)
	if err2 != nil {
		log.Fatalln(err2)
	}
	fechainicial,err := time.Parse("2006-01-02", FechaInicial)
	if err != nil {
		log.Fatalln(err)
	}
	var arriba float64
	arriba=monto * ((intereses/100)*((math.Pow(1+(intereses/100),plazo))))
	fmt.Println("Arriba2")
	fmt.Println(arriba)
	var abajo float64

	abajo=(math.Pow((1+(intereses/100)),plazo))-1

	cuota = math.Round(arriba/abajo)
	fmt.Println("abajo")
	fmt.Println(abajo)

	fmt.Println("cuota")
	fmt.Println(cuota)

	var montolinea float64
	var intereslinea float64
	var capitallinea float64
	//var fechalinea time.Time
	//var abonoslinea float64
	var finallinea float64
	montolinea=monto
	listadocuota := []Cuota{}

	for i:=1 ;i<int(plazo)+1;i++{
		fmt.Println("linea"+strconv.Itoa(i))
		intereslinea=math.Round(montolinea*(intereses/100))
		capitallinea=cuota-intereslinea
		finallinea=montolinea-capitallinea

		fmt.Println("interees")
		fmt.Println(intereslinea)
		fmt.Println("capital")
		fmt.Println(capitallinea)
		fmt.Println("final")
		fmt.Println(finallinea)

		listadocuota=append(listadocuota,Cuota{strconv.Itoa(i),fechainicial.AddDate(0,i-1,0).Format("02/01/2006"),montolinea, intereslinea, capitallinea,finallinea ,cuota})
		montolinea=finallinea
	}
	return listadocuota
}


func CuotaDatos(w http.ResponseWriter, r *http.Request) {
	//db := dbConn()
	Monto := mux.Vars(r)["monto"]
	Plazo := mux.Vars(r)["plazo"]
	Intereses := mux.Vars(r)["intereses"]
	FechaInicial := mux.Vars(r)["fechainicial"]
	listadocuota := []Cuota{}
	listadocuota= Calcularcuota(Monto, Plazo, Intereses, FechaInicial)

		data, _ := json.Marshal(listadocuota)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
}


// CENTRO KARDEX
func cuotaLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/kardex/kardexLista.html")
	//	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	//res := []inventario{}
	//listadokardex := []kardex{}

	if Codigo == "False" {

	} else {

		//	FechaInicial := mux.Vars(r)["fechainicial"]

	}

	varmap := map[string]interface{}{
		//"res":     listadokardex,
		"hosting":  ruta,
		"bodega":   ListaBodega(),
		"producto": ListaProducto(),
	}
	tmp.Execute(w, varmap)
}


// INICIA COMPROBANTE TODOS PDF
func CuotaTodosCabecera(pdf *gofpdf.Fpdf){
	pdf.SetFont("Arial", "", 10)
	// RELLENO TITULO
	pdf.SetY(60)
	pdf.SetFillColor(224,231,239)
	pdf.SetTextColor(0,0,0)
	pdf.Ln(6)
	pdf.SetX(20)
	pdf.CellFormat(184, 6, "No", "0", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(190, 6, "Fecha", "0", 0,
		"L", false, 0, "")
	pdf.SetX(64)
	pdf.CellFormat(190, 6, "Saldo Inicial", "0", 0,
		"L", false, 0, "")
	pdf.SetX(98)
	pdf.CellFormat(190, 6, "Intereses", "0", 0,
		"L", false, 0, "")
	pdf.SetX(132)
	pdf.CellFormat(190, 6, "Capital", "0", 0,
		"L", false, 0, "")
	pdf.SetX(156)
	pdf.CellFormat(190, 6, "Saldo Final", "0", 0,
		"L", false, 0, "")
	pdf.Ln(8)
}
func CuotaTodosDetalle(pdf *gofpdf.Fpdf,miFila Cuota, a int ){
	pdf.SetFont("Arial", "", 9)
	pdf.SetX(21)
	pdf.CellFormat(181, 4, strconv.Itoa(a), "", 0,
		"L", false, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 4, miFila.Fecha, "", 0,
		"L", false, 0, "")
	pdf.SetX(45)
	pdf.CellFormat(40, 4, FormatoFlotante(miFila.Inicial), "", 0,"R", false, 0, "")
	pdf.SetX(75)
	pdf.CellFormat(40, 4, FormatoFlotante(miFila.Intereses), "", 0,
		"R", false, 0, "")
	pdf.SetX(105)
	pdf.CellFormat(40, 4, FormatoFlotante(miFila.Capital), "", 0,
		"R", false, 0, "")
	pdf.SetX(135)
	pdf.CellFormat(40, 4, FormatoFlotante(miFila.Final), "", 0,
		"R", false, 0, "")
	pdf.Ln(4)
}

func CuotaTodosPdf(w http.ResponseWriter, r *http.Request) {
	//db := dbConn()
	Monto := mux.Vars(r)["monto"]
	Plazo := mux.Vars(r)["plazo"]
	Intereses := mux.Vars(r)["intereses"]
	FechaInicial := mux.Vars(r)["fechainicial"]
	listadocuota := []Cuota{}
	listadocuota= Calcularcuota(Monto, Plazo, Intereses, FechaInicial)
	var Cuota float64
	for _, miFila := range listadocuota {
		Cuota= miFila.Cuota
	}

	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
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
		pdf.Ln(10)

		// RELLENO TITULO
		pdf.SetX(20)
		pdf.SetFillColor(224,231,239)
		pdf.SetTextColor(0,0,0)
		pdf.CellFormat(184, 6, "DATOS PRESTAMO", "0", 0,
			"C", true, 0, "")
		pdf.Ln(6)
		pdf.SetX(20)
		pdf.CellFormat(20, 10, "Monto:", "0", 0, "L", false, 0,
			"")
		pdf.SetX(35)
		pdf.CellFormat(20, 10, FormatoFlotante(Flotante(Monto)), "0", 0, "L", false, 0,
			"")
		pdf.SetX(65)
		pdf.CellFormat(20, 10, "Plazo:", "0", 0, "L", false, 0,
			"")
		pdf.SetX(78)
		pdf.CellFormat(20, 10, Plazo, "0", 0, "L", false, 0,
			"")
		pdf.SetX(92)
		pdf.CellFormat(20, 10, "Intereses:", "0", 0, "L", false, 0,
			"")
		pdf.SetX(111)
		pdf.CellFormat(20, 10, Intereses, "0", 0, "L", false, 0,
			"")
		pdf.SetX(125)
		pdf.CellFormat(20, 10, "Fecha:", "0", 0, "L", false, 0,
			"")
		pdf.SetX(141)
		pdf.CellFormat(20, 10, FechaInicial, "0", 0, "L", false, 0,
			"")
		pdf.SetX(165)
		pdf.CellFormat(20, 10, "Cuota:", "0", 0, "L", false, 0,
			"")
		pdf.SetX(180)
		pdf.CellFormat(20, 10, FormatoFlotante(Cuota), "0", 0, "L", false, 0,
			"")
		pdf.Ln(6)
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

	CuotaTodosCabecera(pdf)
	// tercera pagina

	for i, miFila := range listadocuota {
		CuotaTodosDetalle(pdf,miFila,i+1)
		if math.Mod(float64(i+1),46)==0 {
			pdf.AliasNbPages("")
			pdf.AddPage()
			pdf.SetFont("Arial", "", 10)
			pdf.SetX(30)
			CuotaTodosCabecera(pdf)
		}

	}

	err1 = pdf.Output(&buf)
	if err1 != nil {
		panic(err1.Error())
	}
	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	w.Write(buf.Bytes())
}
// TERMINA COMPROBANTE TODOS PDF

// COMPROBANTE EXCEL
func CuotaExcel(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	t := []documento{}
	var e  empresa=ListaEmpresa()
	var c  ciudad=TraerCiudad(e.Ciudad)
	err := db.Select(&t, "SELECT * FROM cuota ORDER BY cast(codigo as integer) ")
	if err != nil {
		log.Fatalln(err)
	}

	f := excelize.NewFile()

	// FUNCION ANCHO DE LA COLUMNA
	if err =f.SetColWidth("Sheet1", "A", "A", 13); err != nil {
		fmt.Println(err)
		return
	}
	if err =f.SetColWidth("Sheet1", "B", "B", 50); err != nil {
		fmt.Println(err)
		return
	}
	if err =f.SetColWidth("Sheet1", "C", "C", 13); err != nil {
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
	f.SetCellValue("Sheet1", "A8","LISTADO DE DOCUMENTOS")

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
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(filaExcel+i), miFila.Consecutivo)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(filaExcel+i), Entero(miFila.Inicial))

		f.SetCellStyle("Sheet1","A"+strconv.Itoa(filaExcel+i),"A"+strconv.Itoa(filaExcel+i),estiloNumeroDetalle)
		f.SetCellStyle("Sheet1","B"+strconv.Itoa(filaExcel+i),"B"+strconv.Itoa(filaExcel+i),estiloTexto)
		f.SetCellStyle("Sheet1","C"+strconv.Itoa(filaExcel+i),"C"+strconv.Itoa(filaExcel+i),estiloTexto)
		f.SetCellStyle("Sheet1","D"+strconv.Itoa(filaExcel+i),"D"+strconv.Itoa(filaExcel+i),estiloNumeroDetalle)

		//van=i
	}

	// LINEA FINAL
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


