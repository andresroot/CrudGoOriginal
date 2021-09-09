package main

// INICIA TRASLADO IMPORTAR PAQUETES
import (
	"bytes"
	"math"

	//"database/sql"
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

// INICIA TRASLADO ESTRUCTURA JSON
type trasladoJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// INICIA TRASLADO ESTRUCTURA
type trasladoLista struct {
	Codigo			  string
	Fecha        	  time.Time
	AlmacenistaNombre string
}

// INICIA TRASLADO ESTRUCTURA
type traslado struct {
	Codigo                    string
	Fecha                     time.Time
	Items                     string
	Almacenista               string
	Accion                    string
	Detalle                   []trasladodetalle `json:"Detalle"`
	DetalleEditar			  []trasladodetalleeditar `json:"DetalleEditar"`
	Tipo					  string

}

// INICIA TRASLADODETALLE ESTRUCTURA
type trasladodetalle struct {
	Id                string
	Codigo            string
	Fila              string
	Bodega            string
	Producto          string
	Tipo              string
	Entra			  string
	Sale			  string
	Fecha			  time.Time
}

// INICIA TRASLADO DETALLE EDITARr
type trasladodetalleeditar struct {
	Id                string
	Codigo            string
	Fila              string
	Bodega            string
	Producto          string
	Fecha			  time.Time
	BodegaNombre      string
	ProductoNombre    string
	ProductoIva       string
	ProductoUnidad    string
	Tipo              string
	Entra			  string
	Sale			  string

}

// INICIA TRASLADO CONSULTA DETALLE
func TrasladoConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "trasladodetalle.Id as id ,"
	consulta += "trasladodetalle.Codigo as codigo,"
	consulta += "trasladodetalle.Entra as entra,"
	consulta += "trasladodetalle.Sale as sale,"
	consulta += "trasladodetalle.Fila as fila,"
	consulta += "trasladodetalle.Bodega as bodega,"
	consulta += "trasladodetalle.Producto as producto,"
	consulta += "trasladodetalle.Fecha as fecha,"
	consulta += "bodega.nombre as BodegaNombre, "
	consulta += "producto.nombre as ProductoNombre, "
	consulta += "producto.iva as ProductoIva, "
	consulta += "producto.unidad as ProductoUnidad "
	consulta += "from trasladodetalle "
	consulta += "inner join producto on producto.codigo=trasladodetalle.producto "
	consulta += "inner join bodega on bodega.codigo=trasladodetalle.bodega "
	consulta += " where trasladodetalle.codigo=$1"
	log.Println(consulta)
	return consulta
}

// INICIA TRASLADO LISTA
func TrasladoLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/traslado/trasladoLista.html")
	log.Println("Error traslado 0")
	var consulta string

	consulta = "  SELECT almacenista.nombre as AlmacenistaNombre,traslado.codigo,fecha"
	consulta += " FROM traslado "
	consulta += " inner join almacenista on almacenista.codigo=traslado.almacenista "
	consulta += " ORDER BY traslado.codigo ASC"

	db := dbConn()
	res := []trasladoLista{}

	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error traslado888")
	tmp.Execute(w, varmap)
}

// INICIA TRASLADO NUEVO
func TrasladoNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error traslado nuevo 1")
	log.Println("Error traslado nuevo 2")
	parametros := map[string]interface{}{
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),

	}

	t, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/traslado/trasladoNuevo.html", "vista/traslado/trasladoScript.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error traslado nuevo 3")
	t.Execute(w, parametros)
}


// INICIA TRASLADO INSERTAR AJAX
func TrasladoAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempTraslado traslado

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}


	// carga informacion de la TRASLADO
	err = json.Unmarshal(b, &tempTraslado)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parametrosinventario := configuracioninventario{}
	err = db.Get(&parametrosinventario, "SELECT * FROM configuracioninventario")
	if err != nil {
		panic(err.Error())
	}

	var Codigoactual string
	//if tempTraslado.Accion == "Nuevo" {
	//	log.Println("Resolucion " + tempTraslado.Resoluciontraslado)
	//	Codigoactual=Numerotraslado(tempTraslado.Resoluciontraslado)
	//	tempTraslado.Codigo=Codigoactual
	//}else{
		Codigoactual=tempTraslado.Codigo
	//}



	if tempTraslado.Accion == "Actualizar" {

		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from trasladodetalle WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempTraslado.Codigo)

		// borra detalle inventario
		delForm2, err := db.Prepare("DELETE from inventario WHERE codigo=$1 and tipo='Traslado'")
		if err != nil {
			panic(err.Error())
		}
		delForm2.Exec(tempTraslado.Codigo)

		// borra cabecera anterior
		delForm1, err := db.Prepare("DELETE from traslado WHERE codigo=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempTraslado.Codigo)
	}

	// INSERTA DETALLE
	for i, x := range tempTraslado.Detalle {
		var a = i
		var q string
		q = "insert into trasladodetalle ("
		q += "Id,"
		q += "Codigo,"
		q += "Fila,"
		q += "Bodega,"
		q += "Producto,"
		q += "Tipo,"
		q += "Entra,"
		q += "Sale,"
		q += "Fecha"
		q += " ) values("
		q += parametros(9)
		q += " ) "

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA TRASLADO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Id,
			Codigoactual,
			x.Fila,
			x.Bodega,
			x.Producto,
			x.Tipo,
			x.Entra,
			x.Sale,
			x.Fecha)
		if err != nil {
			panic(err)
		}

		log.Println("Insertar Producto \n", x.Producto, a)
	}

	// INSERTA DETALLE INVENTARIO
	for i, x := range tempTraslado.Detalle {
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

		var tipooperacion string
		var cantidadtraslado string

		if x.Sale == "" {
			cantidadtraslado = x.Entra

			tipooperacion = operacionTrasladoEntrada
		} else {
			cantidadtraslado = x.Sale
			tipooperacion = operacionTrasladoSalida

		}


		// TERMINA TRASLADO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Fecha,
			x.Tipo,
			Codigoactual,
			x.Bodega,
			x.Producto,
			cantidadtraslado,
			"",
			tipooperacion)
		if err != nil {
			panic(err)
		}
		log.Println("Insertar Producto \n", x.Producto, a)
	}

	// INICIA INSERTAR TRASLADOS
	log.Println("Got %s age %s club %s\n", tempTraslado.Codigo,)
	var q string
	q = "insert into traslado ("
	q += "Codigo,"
	q += "Fecha,"
	q += "Items,"
	q += "Almacenista,"
	q += "Tipo"
	q += " ) values("
	q+=parametros(5)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}

	layout := "2006-01-02"
	log.Println("Hora", tempTraslado.Fecha.Format("02/01/2006"))

	_, err = insForm.Exec(
		tempTraslado.Codigo,
		tempTraslado.Fecha.Format(layout),
		tempTraslado.Items,
		tempTraslado.Almacenista,
		tempTraslado.Tipo)

	if err != nil {
		panic(err)
	}

	//// INSERTAR COMPROBANTE CONTABILIDAD
	//var tempComprobante comprobante
	//var tempComprobanteDetalle comprobantedetalle
	//tempComprobante.Documento="10"
	//tempComprobante.Numero=tempTraslado.Codigo
	//tempComprobante.Fecha =tempTraslado.Fecha
	//tempComprobante.Fechaconsignacion =tempTraslado.Fecha
	//tempComprobante.Debito = tempTraslado.Neto + ".00"
	//tempComprobante.Credito	= tempTraslado.Neto + ".00"
	//tempComprobante.Periodo	= ""
	//tempComprobante.Licencia = ""
	//tempComprobante.Usuario	= ""
	//tempComprobante.Estado	= ""
	//
	//// borra detalle anterior
	//delForm, err := db.Prepare("DELETE from comprobantedetalle WHERE documento=$1 and numero=$2")
	//if err != nil {
	//	panic(err.Error())
	//}
	//delForm.Exec(tempComprobante.Documento, tempComprobante.Numero)
	//
	//// borra cabecera anterior
	//
	//delForm1, err := db.Prepare("DELETE from comprobante WHERE documento=$1 and numero=$2 ")
	//if err != nil {
	//	panic(err.Error())
	//}
	//delForm1.Exec(tempComprobante.Documento, tempComprobante.Numero)
	//
	//// INSERTAR CABECERA COMPROBANTE
	//
	//log.Println("Got %s age %s club %s\n", tempComprobante.Documento, tempComprobante.Numero)
	//
	//var totalDebito float64
	//var totalCredito float64
	//var fila int
	//fila=0
	//totalDebito=0
	//totalCredito=0
	//
	//// INSERTAR CUENTA DEBITO TRASLADO 19%
	//if (tempTraslado.Subtotal19!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalDebito+=Flotante(tempTraslado.Subtotal19)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Trasladocuenta19
	//	tempComprobanteDetalle.Debito = tempTraslado.Subtotal19
	//	tempComprobanteDetalle.Credito =""
	//	InsertaDetalleComprobanteTraslado(tempComprobanteDetalle,tempComprobante,tempTraslado)
	//	log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	//}
	//
	//// INSERTAR CUENTA DEBITO TRASLADO 5%
	//if (tempTraslado.Subtotal5!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalDebito+=Flotante(tempTraslado.Subtotal5)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Trasladocuenta5
	//	tempComprobanteDetalle.Debito = tempTraslado.Subtotal5
	//	tempComprobanteDetalle.Credito =""
	//	InsertaDetalleComprobanteTraslado(tempComprobanteDetalle,tempComprobante,tempTraslado)
	//	log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	//}
	//
	//// INSERTAR CUENTA DEBITO TRASLADO 0%
	//if (tempTraslado.Subtotal0!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalDebito+=Flotante(tempTraslado.Subtotal0)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Trasladocuenta0
	//	tempComprobanteDetalle.Debito = tempTraslado.Subtotal0
	//	tempComprobanteDetalle.Credito =""
	//	InsertaDetalleComprobanteTraslado(tempComprobanteDetalle,tempComprobante,tempTraslado)
	//	log.Println("debito linea" + fmt.Sprintf("%.2f",totalDebito))
	//}
	//
	//// INSERTAR CUENTA CREDITO DESCUENTO
	//if (tempTraslado.Descuento!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalCredito+=Flotante(tempTraslado.Descuento)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Trasladocuentadescuento
	//	tempComprobanteDetalle.Debito = ""
	//	tempComprobanteDetalle.Credito = tempTraslado.Descuento
	//	InsertaDetalleComprobanteTraslado(tempComprobanteDetalle,tempComprobante,tempTraslado)
	//	log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	//}
	//
	//// INSERTAR CUENTA CREDITO RET. FTE.
	//if (tempTraslado.TotalRetencionFuente!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalCredito+=Flotante(tempTraslado.TotalRetencionFuente)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Trasladocuentaretfte
	//	tempComprobanteDetalle.Debito = ""
	//	tempComprobanteDetalle.Credito = tempTraslado.TotalRetencionFuente
	//	InsertaDetalleComprobanteTraslado(tempComprobanteDetalle,tempComprobante,tempTraslado)
	//	log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	//}
	//
	//// INSERTAR CUENTA CREDITO RET. ICA.
	//if (tempTraslado.TotalRetencionIca!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalCredito+=Flotante(tempTraslado.TotalRetencionIca)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Trasladocuentaretica
	//	tempComprobanteDetalle.Debito = ""
	//	tempComprobanteDetalle.Credito = tempTraslado.TotalRetencionIca
	//	InsertaDetalleComprobanteTraslado(tempComprobanteDetalle,tempComprobante,tempTraslado)
	//	log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	//}
	//
	//// INSERTAR CUENTA CREDITO PROVEEDOR
	//if (tempTraslado.Neto!="0")	{
	//	fila=fila+1
	//	tempComprobanteDetalle.Fila=strconv.Itoa(fila)
	//	totalCredito+=Flotante(tempTraslado.Neto)
	//	tempComprobanteDetalle.Cuenta  = parametrosinventario.Trasladocuentaproveedor
	//	tempComprobanteDetalle.Debito = ""
	//	tempComprobanteDetalle.Credito = tempTraslado.Neto
	//	InsertaDetalleComprobanteTraslado(tempComprobanteDetalle,tempComprobante,tempTraslado)
	//	log.Println("credito linea" + fmt.Sprintf("%.2f",totalCredito))
	//}
	//
	//var cadenaDebito=FormatoFlotante(totalDebito)
	//var cadenaCredito=FormatoFlotante(totalCredito)
	//
	//q = "insert into comprobante ("
	//q += "Documento,"
	//q += "Numero,"
	//q += "Fecha,"
	//q += "Fechaconsignacion,"
	//q += "Debito,"
	//q += "Credito,"
	//q += "Periodo,"
	//q += "Licencia,"
	//q += "Usuario,"
	//q += "Estado"
	//q += " ) values("
	//q += parametros(10)
	//q += " ) "
	//
	//log.Println("Cadena SQL " + q)
	//insForm, err = db.Prepare(q)
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//_, err = insForm.Exec(
	//	tempComprobante.Documento,
	//	tempComprobante.Numero,
	//	tempComprobante.Fecha.Format(layout),
	//	tempComprobante.Fechaconsignacion.Format(layout),
	//	cadenaDebito,
	//	cadenaCredito,
	//	tempComprobante.Periodo,
	//	tempComprobante.Licencia,
	//	tempComprobante.Usuario,
	//	tempComprobante.Estado)
	//if err != nil {
	//	panic(err)
	//}

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

// INICIA TRASLADO EXISTE
func TrasladoExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM traslado  WHERE codigo=$1", Codigo)
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

// INICIA TRASLADO EDITAR
func TrasladoEditar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio traslado editar" + Codigo)

	db := dbConn()

	// traer TRASLADO
	v := traslado{}
	err := db.Get(&v, "SELECT * FROM traslado where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []trasladodetalleeditar{}

	err2 := db.Select(&det, TrasladoConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"traslado":       v,
		"detalle":     det,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),
		"bodega":      ListaBodega(),

	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html", "vista/traslado/trasladoEditar.html", "vista/traslado/trasladoScript.html")
	fmt.Printf("%v, %v", err)
	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error traslado nuevo 3")
	miTemplate.Execute(w, parametros)
}

// INICIA TRASLADO BORRAR
func TrasladoBorrar(w http.ResponseWriter, r *http.Request) {
	Codigo := mux.Vars(r)["codigo"]
	log.Println("inicio traslado editar" + Codigo)

	db := dbConn()

	// traer TRASLADO
	v := traslado{}
	err := db.Get(&v, "SELECT * FROM traslado where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}
	// traer detalle

	det := []trasladodetalleeditar{}
	err2 := db.Select(&det, TrasladoConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"traslado":       v,
		"detalle":     det,
		"hosting":     ruta,
		"almacenista": ListaAlmacenista(),

	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/traslado/trasladoBorrar.html", "vista/traslado/trasladoScript.html")
	fmt.Printf("%v, %v", err)
	log.Println("Error traslado nuevo 3")
	miTemplate.Execute(w, parametros)

}

// INICIA TRASLADO ELIMINAR
func TrasladoEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	codigo := mux.Vars(r)["codigo"]

	// borrar TRASLADO
	delForm, err := db.Prepare("DELETE from traslado WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(codigo)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from trasladodetalle WHERE codigo=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec(codigo)

	// borar detalle invenario
	Borrarinventario(codigo,"Traslado")

	log.Println("Registro Eliminado" + codigo)
	http.Redirect(w, r, "/TrasladoLista", 301)
}

// TRAER PEDIDO
//func Datospedidotraslado(w http.ResponseWriter, r *http.Request) {
//	Codigo := mux.Vars(r)["codigo"]
//	log.Println("inicio pedido editar" + Codigo)
//	db := dbConn()
//	var res []pedidotraslado

	// traer PEDIDO
	//v := pedidotraslado{}
	//err := db.Get(&v, "SELECT * FROM pedidotraslado where codigo=$1", Codigo)
	//var valida bool
	//valida=true
	//
	//switch err {
	//case nil:
	//	log.Printf("pedido existe: %+v\n", v)
	//case sql.ErrNoRows:
	//	log.Println("pedidotraslado NO Existe")
	//	valida=false
	//default:
	//	log.Printf("pedidotraslado error: %s\n", err)
	//}
	//det := []pedidotrasladodetalleeditar{}
	//t := tercero{}
	//
	//// trae datos si existe pedido
	//if valida==true {
	//	err2 := db.Select(&det, PedidotrasladoConsultaDetalle(), Codigo)
	//	if err2 != nil {
	//		fmt.Println(err2)
	//	}
	//	// traer tercero
	//	err1 := db.Get(&t, "SELECT * FROM tercero where codigo=$1", v.Tercero)
	//	if err1 != nil {
	//		log.Fatalln(err1)
	//	}
	//	v.TerceroDetalle=t;
	//	v.DetalleEditar=det;
	//	res = append(res, v)
	//}

//	log.Println("codigo nombre99" + t.Codigo + t.Nombre)
//	data, err := json.Marshal(res)
//	w.WriteHeader(200)
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(data)
//
//}

// INICIA PDF
func TrasladoPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Codigo := mux.Vars(r)["codigo"]

	// traer TRASLADO
	miTraslado := traslado{}
	err := db.Get(&miTraslado, "SELECT * FROM traslado where codigo=$1", Codigo)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	miDetalle := []trasladodetalleeditar{}
	err2 := db.Select(&miDetalle, TrasladoConsultaDetalle(), Codigo)
	if err2 != nil {
		fmt.Println(err2)
	}

	// traer tercero


	// traer almacenista
	miAlmacenista := almacenista{}
	err4 := db.Get(&miAlmacenista, "SELECT * FROM almacenista where codigo=$1", miTraslado.Almacenista)
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
		pdf.Image(imageFile("logo.png"), 20, 30, 40, 0, false,
			"", 0, "")
		pdf.SetY(5)

		// DATOS EMPRESA
		pdf.SetFont("Arial", "", 11)
		// CUADRO EMPRESA
		pdf.SetY(20)
		pdf.SetX(20)
		pdf.SetDrawColor(119,134,153)
		pdf.CellFormat(184, 50, "", "1", 0, "C",
			false, 0, "")

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
		pdf.CellFormat(190, 10, e.Telefono1+" "+e.Telefono2, "0", 0, "C", false, 0,
			"")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, ene(c.NombreCiudad+" - "+c.NombreDepartamento), "0", 0, "C", false, 0,
			"")

		pdf.SetY(20)
		pdf.SetX(80)
		pdf.Ln(8)
		pdf.SetX(80)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(190, 10, "TRASLADOS Y AJUSTES", "0", 0, "C",
			false, 0, "")
		pdf.Ln(6)
		pdf.SetX(80)
		pdf.CellFormat(190, 10, "No. " +miTraslado.Codigo, "0", 0, "C",
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
	TrasladoCabecera(pdf,miTraslado,miAlmacenista)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			TrasladoFilaDetalle(pdf,miFila,a)
		}
		TrasladoPieDePagina(pdf,miTraslado)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					TrasladoFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			TrasladoCabecera(pdf,miTraslado,miAlmacenista)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					TrasladoFilaDetalle(pdf,miFila,a)
				}
			}

			TrasladoPieDePagina(pdf,miTraslado)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					TrasladoFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			TrasladoCabecera(pdf,miTraslado,miAlmacenista)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					TrasladoFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			TrasladoCabecera(pdf,miTraslado,miAlmacenista)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					TrasladoFilaDetalle(pdf,miFila,a)
				}
			}

			TrasladoPieDePagina(pdf,miTraslado)
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
func TrasladoCabecera(pdf *gofpdf.Fpdf,miTraslado traslado, miAlmacenista almacenista ){

	// RELLENO TITULO
	pdf.SetY(46)
	pdf.SetX(20)
	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(8)
	pdf.SetX(20)
	pdf.CellFormat(184, 5, "TRASLADOS Y AJUSTES", "0", 0,
		"C", true, 0, "")
	pdf.Ln(8)
	pdf.SetX(20)
	pdf.SetTextColor(0,0,0)
	pdf.CellFormat(40, 4, "Fecha de Expedicion", "", 0,
		"L", false, 0, "")
	pdf.SetX(60)
	pdf.CellFormat(40, 4, miTraslado.Fecha.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Almacenista", "", 0,
		"L", false, 0, "")
	pdf.SetX(135)
	pdf.CellFormat(40, 4, miAlmacenista.Nombre, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.SetY(68)
	pdf.SetX(20)

	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)

	pdf.CellFormat(184, 5, "ITEM", "1", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(190, 5, "CODIGO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(190, 5, "PRODUCTO", "0", 0,
		"L", false, 0, "")
	pdf.SetX(96)
	pdf.CellFormat(190, 5, "UNIDAD", "0", 0,
		"L", false, 0, "")
	pdf.SetX(112)
	pdf.CellFormat(190, 5, "IVA", "0", 0,
		"L", false, 0, "")
	pdf.SetX(122)
	pdf.CellFormat(190, 5, "BODEGA", "0", 0,
		"L", false, 0, "")
	pdf.SetX(152)
	pdf.CellFormat(190, 5, "ENTRADA", "0", 0,
		"L", false, 0, "")
	pdf.SetX(180)
	pdf.CellFormat(190, 5, "SALIDA", "0", 0,
		"L", false, 0, "")
	pdf.Ln(8)
}
func TrasladoFilaDetalle(pdf *gofpdf.Fpdf,miFila trasladodetalleeditar, a int ){
	pdf.SetFont("Arial", "", 9)
	if math.Mod(float64(a),2)==0 {
		pdf.SetFillColor(224,231,239)
		pdf.SetTextColor(0,0,0)
	} else{
		pdf.SetFillColor(255,255,255)
		pdf.SetTextColor(0,0,0)
	}

	pdf.SetTextColor(0,0,0)
	pdf.SetX(21)
	pdf.CellFormat(183, 4, strconv.Itoa(a), "", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(40, 4, Subcadena(miFila.Producto,0,12), "", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(40, 4, Subcadena(miFila.ProductoNombre,0,32),  "", 0,
		"L", false, 0, "")
	pdf.SetX(72)
	pdf.CellFormat(40, 4, miFila.ProductoUnidad, "", 0,
		"R", false, 0, "")
	pdf.SetX(79)
	pdf.CellFormat(40, 4, miFila.ProductoIva, "", 0,
		"R", false, 0, "")
	pdf.SetX(95)
	pdf.CellFormat(40, 4, miFila.BodegaNombre, "", 0,
		"R", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, miFila.Entra,"", 0,
		"R", false, 0, "")
	pdf.SetX(155)
	pdf.CellFormat(40, 4, miFila.Sale, "", 0,
		"R", false, 0, "")
	pdf.Ln(4)
}

func TrasladoPieDePagina(pdf *gofpdf.Fpdf,miTraslado traslado ){
	pdf.SetFont("Arial", "", 9)
	pdf.SetY(232)
	pdf.Ln(10)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "__________________________________________", "0", 0,
		"C", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(1)
	pdf.CellFormat(190, 10, "FIRMA RESPONSABLE ", "0", 0, "C",
		false, 0, "")

}
