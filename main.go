package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"
)

// CONEXION A BASE DE DATOS POSTGRES
const (
	host     = "192.168.1.3"
	port     = 5432
	user     = "postgres"
	password = "Murc4505"
	dbname   = "Base2020"
)

var mapaRuta = map[string]interface{}{
	"hosting": ruta,
}
var (
	once sync.Once
	db   *sqlx.DB
)

func dbConn() *sqlx.DB {
	var err error
	once.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		//db, err = sql.Open("postgres", psqlInfo)
		db, err = sqlx.Open("postgres", psqlInfo)

		log.Println("conectando...")
		if err != nil {
			panic(err.Error())
			log.Println(err.Error())
		}
		log.Println("Base de Datos Conectada")
	})
	return db
}

// LISTA DE REGISTROS//
func Index(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Index.html","vista/inicio/appInicio.html",
		)
	//db := dbConn()
	//Empresa := empresa{}
	//err := db.Get(&Empresa, "SELECT * FROM empresa limit 1")
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//Empresa.Nombre=Titulo(Empresa.Nombre)
	//Empresa.Codigo=Coma(Empresa.Codigo)
	log.Println("inicio listado")

	varmap := map[string]interface{}{
		"hosting": ruta,
		//"empresa": Empresa,
	}

	log.Println("Error empresa888")
	tmp.Execute(w, varmap)
}

// RUTAS GENERALES
var router = mux.NewRouter()
var ruta="http://localhost:9002/"
//var ruta="http://192.168.1.3:9002/"
//var ruta="/"

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
func favicon(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=\n")
}
func CacheControlWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		h.ServeHTTP(w, r)
	})
}
func main() {
	//router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
	//http.FileServer(http.Dir("static"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		CacheControlWrapper(http.FileServer(http.Dir("static")))))

	router.Path("/").HandlerFunc(Index).Name("Index")
	// favicon
	http.HandleFunc("/favicon.ico", favicon)

	// ARCHIVO CENTRO
	router.Path("/CentroNuevo/{codigo}").HandlerFunc(CentroNuevo).Name("CentroNuevo")
	router.Path("/CentroBuscar/{codigo}").HandlerFunc(CentroBuscar).
		Name("CentroBuscar")
	router.Path("/CentroActual/{codigo}").HandlerFunc(CentroActual).
		Name("CentroActual")
	router.Path("/CentroLista").HandlerFunc(CentroLista).Name("CentroLista")
	router.Path("/CentroExiste/{codigo:[0-9]+}").HandlerFunc(CentroExiste).
		Name("CentroExiste")
	router.Path("/CentroInsertar").HandlerFunc(CentroInsertar).Name(
		"CentroInsertar")
	router.Path("/CentroActualizar/{codigo:[0-9]+}").HandlerFunc(CentroActualizar).Name(
		"CentroActualizar")
	router.Path("/CentroBorrar/{codigo:[0-9]+}").HandlerFunc(CentroBorrar).Name(
		"CentroBorrar")
	router.Path("/CentroEliminar/{codigo:[0-9]+}").HandlerFunc(
		CentroEliminar).Name("CentroEliminar")
	router.Path("/CentroEditar/{codigo:[0-9]+}").HandlerFunc(CentroEditar).Name(
		"CentroEditar")
	router.Path("/CentroPdf/{codigo:[0-9]+}").HandlerFunc(CentroPdf).Name(
		"CentroPdf")

	router.Path("/CentroTodosPdf").HandlerFunc(CentroTodosPdf).
		Name("CentroTodosPdf")
	router.Path("/CentroExcel").HandlerFunc(CentroExcel).
		Name("CentroExcel")

	// DOCUMENTO
	router.Path("/DocumentoNuevo/{codigo}").HandlerFunc(DocumentoNuevo).Name("DocumentoNuevo")
	router.Path("/DocumentoLista").HandlerFunc(DocumentoLista).Name("DocumentoLista")
	router.Path("/DocumentoExiste/{codigo:[0-9]+}").HandlerFunc(DocumentoExiste).
		Name("DocumentoExiste")
	router.Path("/DocumentoActual/{codigo}").HandlerFunc(DocumentoActual).
		Name("DocumentoActual")
	router.Path("/DocumentoBuscar/{codigo}").HandlerFunc(DocumentoBuscar).
		Name("DocumentoBuscar")
	router.Path("/DocumentoInsertar").HandlerFunc(DocumentoInsertar).Name(
		"DocumentoInsertar")
	router.Path("/DocumentoActualizar/{codigo:[0-9]+}").HandlerFunc(DocumentoActualizar).Name(
		"DocumentoActualizar")
	router.Path("/DocumentoBorrar/{codigo:[0-9]+}").HandlerFunc(DocumentoBorrar).Name(
		"DocumentoBorrar")
	router.Path("/DocumentoEliminar/{codigo:[0-9]+}").HandlerFunc(
		DocumentoEliminar).Name("DocumentoEliminar")
	router.Path("/DocumentoEditar/{codigo:[0-9]+}").HandlerFunc(DocumentoEditar).Name(
		"DocumentoEditar")
	router.Path("/DocumentoPdf/{codigo:[0-9]+}").HandlerFunc(DocumentoPdf).
			Name("DocumentoPdf")

	router.Path("/DocumentoTodosPdf").HandlerFunc(DocumentoTodosPdf).
		Name("DocumentoTodosPdf")
	router.Path("/DocumentoExcel").HandlerFunc(DocumentoExcel).
		Name("DocumentoExcel")

	// ARCHIVO BODEGA
	router.Path("/BodegaNuevo").HandlerFunc(BodegaNuevo).Name("BodegaNuevo")
	router.Path("/BodegaLista").HandlerFunc(BodegaLista).Name("BodegaLista")
	router.Path("/BodegaExiste/{codigo:[0-9]+}").HandlerFunc(BodegaExiste).
		Name("BodegaExiste")
	router.Path("/BodegaInsertar").HandlerFunc(BodegaInsertar).Name(
		"BodegaInsertar")
	router.Path("/BodegaActualizar/{codigo:[0-9]+}").HandlerFunc(BodegaActualizar).Name(
		"BodegaActualizar")
	router.Path("/BodegaBorrar/{codigo:[0-9]+}").HandlerFunc(BodegaBorrar).Name(
		"BodegaBorrar")
	router.Path("/BodegaEliminar/{codigo:[0-9]+}").HandlerFunc(
		BodegaEliminar).Name("BodegaEliminar")
	router.Path("/BodegaEditar/{codigo:[0-9]+}").HandlerFunc(BodegaEditar).Name(
		"BodegaEditar")
	router.Path("/BodegaPdf/{codigo:[0-9]+}").HandlerFunc(BodegaPdf).Name(
		"BodegaPdf")

	// ARCHIVO GRUPO
	router.Path("/GrupoNuevo").HandlerFunc(GrupoNuevo).Name("GrupoNuevo")
	router.Path("/GrupoLista").HandlerFunc(GrupoLista).Name("GrupoLista")
	router.Path("/GrupoExiste/{codigo:[0-9]+}").HandlerFunc(GrupoExiste).
		Name("GrupoExiste")
	router.Path("/GrupoInsertar").HandlerFunc(GrupoInsertar).Name(
		"GrupoInsertar")
	router.Path("/GrupoActualizar/{codigo:[0-9]+}").HandlerFunc(GrupoActualizar).Name(
		"GrupoActualizar")
	router.Path("/GrupoBorrar/{codigo:[0-9]+}").HandlerFunc(GrupoBorrar).Name(
		"GrupoBorrar")
	router.Path("/GrupoEliminar/{codigo:[0-9]+}").HandlerFunc(
		GrupoEliminar).Name("GrupoEliminar")
	router.Path("/GrupoEditar/{codigo:[0-9]+}").HandlerFunc(GrupoEditar).Name(
		"GrupoEditar")
	router.Path("/GrupoPdf/{codigo:[0-9]+}").HandlerFunc(GrupoPdf).Name(
		"GrupoPdf")

	// TERCERO
	router.Path("/TerceroLista").HandlerFunc(TerceroLista).Name("TerceroLista")
	router.Path("/TerceroNuevo/{panel}/{codigo}/{elemento}").HandlerFunc(TerceroNuevo).Name("TerceroNuevo")
	router.Path("/TerceroNuevoCopia/{copiacodigo}").HandlerFunc(TerceroNuevoCopia).Name("TerceroNuevoCopia")

	router.Path("/TerceroBuscar/{codigo}").HandlerFunc(TerceroBuscar).
		Name("TerceroBuscar")
	router.Path("/TerceroActual/{codigo}").HandlerFunc(TerceroActual).
		Name("TerceroActual")
	router.Path("/TerceroActualBanco/{codigo}").HandlerFunc(TerceroActualBanco).
		Name("TerceroActualBanco")


	router.Path("/TerceroExiste/{codigo:[0-9]+}").HandlerFunc(TerceroExiste).
		Name("TerceroExiste")
	router.Path("/TerceroEditar/{codigo:[0-9]+}").HandlerFunc(TerceroEditar).
		Name("TerceroEditar")
	router.Path("/TerceroActualizar/{codigo:[0-9]+}").HandlerFunc(
		TerceroActualizar).Name("TerceroActualizar")
	router.Path("/TerceroInsertar").HandlerFunc(TerceroInsertar).Name(
		"TerceroInsertar")
	router.Path("/TerceroBorrar/{codigo:[0-9]+}").HandlerFunc(TerceroBorrar).
			Name("TerceroBorrar")
	router.Path("/TerceroEliminar/{codigo:[0-9]+}").HandlerFunc(
		TerceroEliminar).Name("TerceroEliminar")

	router.Path("/TerceroPdf/{codigo}").HandlerFunc(TerceroPdf).Name(
		"TerceroPdf")

	router.Path("/TerceroTodosPdf").HandlerFunc(TerceroTodosPdf).
		Name("TerceroTodosPdf")
	router.Path("/TerceroExcel").HandlerFunc(TerceroExcel).
		Name("TerceroExcel")


	//// TERCERO HORIZONTAL
	//router.Path("/TercerohorizontalLista").HandlerFunc(TercerohorizontalLista).Name("TercerohorizontalLista")
	//router.Path("/TercerohorizontalNuevo/{panel}/{codigo}/{elemento}").HandlerFunc(TercerohorizontalNuevo).Name("TercerohorizontalNuevo")
	//router.Path("/TercerohorizontalBuscar/{codigo}").HandlerFunc(TercerohorizontalBuscar).
	//	Name("TercerohorizontalBuscar")
	//router.Path("/TercerohorizontalActual/{codigo}").HandlerFunc(TercerohorizontalActual).
	//	Name("TercerohorizontalActual")
	//router.Path("/TercerohorizontalExiste/{codigo:[0-9]+}").HandlerFunc(TercerohorizontalExiste).
	//	Name("TercerohorizontalExiste")
	//router.Path("/TercerohorizontalEditar/{codigo:[0-9]+}").HandlerFunc(TercerohorizontalEditar).
	//	Name("TercerohorizontalEditar")
	//router.Path("/TercerohorizontalActualizar/{codigo:[0-9]+}").HandlerFunc(
	//	TercerohorizontalActualizar).Name("TercerohorizontalActualizar")
	//router.Path("/TercerohorizontalInsertar").HandlerFunc(TercerohorizontalInsertar).Name(
	//	"TercerohorizontalInsertar")
	//router.Path("/TercerohorizontalBorrar/{codigo:[0-9]+}").HandlerFunc(TercerohorizontalBorrar).
	//	Name("TercerohorizontalBorrar")
	//router.Path("/TercerohorizontalEliminar/{codigo:[0-9]+}").HandlerFunc(
	//	TercerohorizontalEliminar).Name("TercerohorizontalEliminar")
	//router.Path("/TercerohorizontalPdf/{codigo}").HandlerFunc(TercerohorizontalPdf).Name(
	//	"TercerohorizontalPdf")

	// EMPRESA
	router.Path("/EmpresaLista").HandlerFunc(EmpresaLista).Name("EmpresaLista")
	router.Path("/EmpresaNuevo/{codigo}").HandlerFunc(EmpresaNuevo).Name("EmpresaNuevo")
	router.Path("/EmpresaBuscar/{codigo}").HandlerFunc(EmpresaBuscar).
		Name("EmpresaBuscar")
	router.Path("/EmpresaActual/{codigo}").HandlerFunc(EmpresaActual).
		Name("EmpresaActual")
	router.Path("/EmpresaExiste/{codigo:[0-9]+}").HandlerFunc(EmpresaExiste).
		Name("EmpresaExiste")
	router.Path("/EmpresaEditar/{codigo:[0-9]+}").HandlerFunc(EmpresaEditar).
		Name("EmpresaEditar")
	router.Path("/EmpresaActualizar/{codigo:[0-9]+}").HandlerFunc(
		EmpresaActualizar).Name("EmpresaActualizar")
	router.Path("/EmpresaInsertar").HandlerFunc(EmpresaInsertar).Name(
		"EmpresaInsertar")
	router.Path("/EmpresaBorrar/{codigo:[0-9]+}").HandlerFunc(EmpresaBorrar).
		Name("EmpresaBorrar")
	router.Path("/EmpresaEliminar/{codigo:[0-9]+}").HandlerFunc(
		EmpresaEliminar).Name("EmpresaEliminar")
	router.Path("/EmpresaPdf/{codigo:[0-9]+}").HandlerFunc(EmpresaPdf).Name(
		"EmpresaPdf")

	router.Path("/EmpresaTodosPdf").HandlerFunc(EmpresaTodosPdf).
		Name("EmpresaTodosPdf")
	router.Path("/EmpresaExcel").HandlerFunc(EmpresaExcel).
		Name("EmpresaExcel")

	// USUARIO
	router.Path("/UsuarioLista").HandlerFunc(UsuarioLista).Name("UsuarioLista")
	router.Path("/UsuarioNuevo").HandlerFunc(UsuarioNuevo).Name("UsuarioNuevo")
	router.Path("/UsuarioBuscar/{codigo}").HandlerFunc(UsuarioBuscar).
		Name("UsuarioBuscar")
	router.Path("/UsuarioActual/{codigo}").HandlerFunc(UsuarioActual).
		Name("UsuarioActual")
	router.Path("/UsuarioExiste/{codigo:[0-9]+}").HandlerFunc(UsuarioExiste).
		Name("UsuarioExiste")
	router.Path("/UsuarioEditar/{codigo:[0-9]+}").HandlerFunc(UsuarioEditar).
		Name("UsuarioEditar")
	router.Path("/UsuarioActualizar/{codigo:[0-9]+}").HandlerFunc(
		UsuarioActualizar).Name("UsuarioActualizar")
	router.Path("/UsuarioInsertar").HandlerFunc(UsuarioInsertar).Name(
		"UsuarioInsertar")
	router.Path("/UsuarioBorrar/{codigo:[0-9]+}").HandlerFunc(UsuarioBorrar).
		Name("UsuarioBorrar")
	router.Path("/UsuarioEliminar/{codigo:[0-9]+}").HandlerFunc(
		UsuarioEliminar).Name("UsuarioEliminar")
	router.Path("/UsuarioPdf/{codigo}").HandlerFunc(UsuarioPdf).Name(
		"UsuarioPdf")

	router.Path("/UsuarioTodosPdf").HandlerFunc(UsuarioTodosPdf).
		Name("UsuarioTodosPdf")
	router.Path("/UsuarioExcel").HandlerFunc(UsuarioExcel).
		Name("UsuarioaExcel")


	// ARCHIVO SUBGRUPO
	router.Path("/SubgrupoNuevo").HandlerFunc(SubgrupoNuevo).Name("SubgrupoNuevo")
	router.Path("/SubgrupoLista").HandlerFunc(SubgrupoLista).Name("SubgrupoLista")
	router.Path("/SubgrupoExiste/{codigo:[0-9]+}").HandlerFunc(SubgrupoExiste).
		Name("SubgrupoExiste")
	router.Path("/SubgrupoInsertar").HandlerFunc(SubgrupoInsertar).Name(
		"SubgrupoInsertar")
	router.Path("/SubgrupoActualizar/{codigo:[0-9]+}").HandlerFunc(SubgrupoActualizar).Name(
		"SubgrupoActualizar")
	router.Path("/SubgrupoBorrar/{codigo:[0-9]+}").HandlerFunc(SubgrupoBorrar).Name(
		"SubgrupoBorrar")
	router.Path("/SubgrupoEliminar/{codigo:[0-9]+}").HandlerFunc(
		SubgrupoEliminar).Name("SubgrupoEliminar")
	router.Path("/SubgrupoEditar/{codigo:[0-9]+}").HandlerFunc(SubgrupoEditar).Name(
		"SubgrupoEditar")
	router.Path("/SubgrupoPdf/{codigo:[0-9]+}").HandlerFunc(SubgrupoPdf).Name(
		"SubgrupoPdf")

	// VENDEDOR
	router.Path("/VendedorLista").HandlerFunc(VendedorLista).Name("VendedorLista")
	router.Path("/VendedorNuevo").HandlerFunc(VendedorNuevo).Name("VendedorNuevo")
	router.Path("/VendedorBuscar/{codigo}").HandlerFunc(VendedorBuscar).
		Name("VendedorBuscar")
	router.Path("/VendedorActual/{codigo}").HandlerFunc(VendedorActual).
		Name("VendedorActual")
	router.Path("/VendedorExiste/{codigo:[0-9]+}").HandlerFunc(VendedorExiste).
		Name("VendedorExiste")
	router.Path("/VendedorEditar/{codigo:[0-9]+}").HandlerFunc(VendedorEditar).
		Name("VendedorEditar")
	router.Path("/VendedorActualizar/{codigo:[0-9]+}").HandlerFunc(
		VendedorActualizar).Name("VendedorActualizar")
	router.Path("/VendedorInsertar").HandlerFunc(VendedorInsertar).Name(
		"VendedorInsertar")
	router.Path("/VendedorBorrar/{codigo:[0-9]+}").HandlerFunc(VendedorBorrar).
		Name("VendedorBorrar")
	router.Path("/VendedorEliminar/{codigo:[0-9]+}").HandlerFunc(
		VendedorEliminar).Name("VendedorEliminar")
	router.Path("/VendedorPdf/{codigo}").HandlerFunc(VendedorPdf).Name(
		"VendedorPdf")

	// ALMACENISTA
	router.Path("/AlmacenistaLista").HandlerFunc(AlmacenistaLista).Name("AlmacenistaLista")
	router.Path("/AlmacenistaNuevo").HandlerFunc(AlmacenistaNuevo).Name("AlmacenistaNuevo")
	router.Path("/AlmacenistaBuscar/{codigo}").HandlerFunc(AlmacenistaBuscar).
		Name("AlmacenistaBuscar")
	router.Path("/AlmacenistaActual/{codigo}").HandlerFunc(AlmacenistaActual).
		Name("AlmacenistaActual")
	router.Path("/AlmacenistaExiste/{codigo:[0-9]+}").HandlerFunc(AlmacenistaExiste).
		Name("AlmacenistaExiste")
	router.Path("/AlmacenistaEditar/{codigo:[0-9]+}").HandlerFunc(AlmacenistaEditar).
		Name("AlmacenistaEditar")
	router.Path("/AlmacenistaActualizar/{codigo:[0-9]+}").HandlerFunc(
		AlmacenistaActualizar).Name("AlmacenistaActualizar")
	router.Path("/AlmacenistaInsertar").HandlerFunc(AlmacenistaInsertar).Name(
		"AlmacenistaInsertar")
	router.Path("/AlmacenistaBorrar/{codigo:[0-9]+}").HandlerFunc(AlmacenistaBorrar).
		Name("AlmacenistaBorrar")
	router.Path("/AlmacenistaEliminar/{codigo:[0-9]+}").HandlerFunc(
		AlmacenistaEliminar).Name("AlmacenistaEliminar")
	router.Path("/AlmacenistaPdf/{codigo}").HandlerFunc(AlmacenistaPdf).Name(
		"AlmacenistaPdf")

	// PRODUCTO
	router.Path("/ProductoLista").HandlerFunc(ProductoLista).Name("ProductoLista")
	router.Path("/ProductoNuevo/{panel}/{codigo}/{elemento}").HandlerFunc(ProductoNuevo).Name("ProductoNuevo")
	router.Path("/ProductoBuscar/{codigo}").HandlerFunc(ProductoBuscar).
		Name("ProductoBuscar")
	router.Path("/ProductoActual/{codigo}").HandlerFunc(ProductoActual).
		Name("ProductoActual")
	router.Path("/ProductoExiste/{codigo:[0-9]+}").HandlerFunc(ProductoExiste).
		Name("ProductoExiste")
	router.Path("/ProductoEditar/{codigo:[0-9]+}").HandlerFunc(ProductoEditar).
		Name("ProductoEditar")
	router.Path("/ProductoActualizar/{codigo:[0-9]+}").HandlerFunc(
		ProductoActualizar).Name("ProductoActualizar")
	router.Path("/ProductoInsertar").HandlerFunc(ProductoInsertar).Name(
		"ProductoInsertar")
	router.Path("/ProductoBorrar/{codigo:[0-9]+}").HandlerFunc(ProductoBorrar).
		Name("ProductoBorrar")
	router.Path("/ProductoEliminar/{codigo:[0-9]+}").HandlerFunc(
		ProductoEliminar).Name("ProductoEliminar")
	router.Path("/ProductoPdf/{codigo}").HandlerFunc(ProductoPdf).Name(
		"ProductoPdf")

	// ARCHIVO RESOLUCION VENTA
	router.Path("/ResolucionventaNuevo").HandlerFunc(ResolucionventaNuevo).Name("ResolucionventaNuevo")
	router.Path("/ResolucionventaLista").HandlerFunc(ResolucionventaLista).Name("ResolucionventaLista")
	router.Path("/ResolucionventaExiste/{codigo:[0-9]+}").HandlerFunc(ResolucionventaExiste).
		Name("ResolucionventaExiste")
	router.Path("/ResolucionventaInsertar").HandlerFunc(ResolucionventaInsertar).Name(
		"ResolucionventaInsertar")
	router.Path("/ResolucionventaActualizar/{codigo:[0-9]+}").HandlerFunc(ResolucionventaActualizar).Name(
		"ResolucionventaActualizar")
	router.Path("/ResolucionventaBorrar/{codigo:[0-9]+}").HandlerFunc(ResolucionventaBorrar).Name(
		"ResolucionventaBorrar")
	router.Path("/ResolucionventaEliminar/{codigo:[0-9]+}").HandlerFunc(
		ResolucionventaEliminar).Name("ResolucionventaEliminar")
	router.Path("/ResolucionventaEditar/{codigo:[0-9]+}").HandlerFunc(ResolucionventaEditar).Name(
		"ResolucionEditar")
	router.Path("/ResolucionventaPdf/{codigo:[0-9]+}").HandlerFunc(ResolucionventaPdf).Name(
		"ResolucionventaPdf")

	// ARCHIVO RESOLUCION SOPORTE
	router.Path("/ResolucionsoporteNuevo").HandlerFunc(ResolucionsoporteNuevo).Name("ResolucionsoporteNuevo")
	router.Path("/ResolucionsoporteLista").HandlerFunc(ResolucionsoporteLista).Name("ResolucionsoporteLista")
	router.Path("/ResolucionsoporteExiste/{codigo:[0-9]+}").HandlerFunc(ResolucionsoporteExiste).
		Name("ResolucionsoporteExiste")
	router.Path("/ResolucionsoporteInsertar").HandlerFunc(ResolucionsoporteInsertar).Name(
		"ResolucionsoporteInsertar")
	router.Path("/ResolucionsoporteActualizar/{codigo:[0-9]+}").HandlerFunc(ResolucionsoporteActualizar).Name(
		"ResolucionsoporteActualizar")
	router.Path("/ResolucionsoporteBorrar/{codigo:[0-9]+}").HandlerFunc(ResolucionsoporteBorrar).Name(
		"ResolucionsoporteBorrar")
	router.Path("/ResolucionsoporteEliminar/{codigo:[0-9]+}").HandlerFunc(
		ResolucionsoporteEliminar).Name("ResolucionsoporteEliminar")
	router.Path("/ResolucionsoporteEditar/{codigo:[0-9]+}").HandlerFunc(ResolucionsoporteEditar).Name(
		"ResolucionEditar")
	router.Path("/ResolucionsoportePdf/{codigo:[0-9]+}").HandlerFunc(ResolucionsoportePdf).Name(
		"ResolucionsoportePdf")

	// COTIZACION
	router.Path("/CotizacionLista").HandlerFunc(CotizacionLista).Name("CotizacionLista")
	router.Path("/CotizacionNuevo/{codigo}").HandlerFunc(CotizacionNuevo).Name("CotizacionNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/CotizacionExiste/{codigo}").HandlerFunc(CotizacionExiste).
		Name("CotizacionExiste")
	router.Path("/CotizacionEditar/{codigo}").HandlerFunc(CotizacionEditar).
		Name("CotizacionEditar")
	router.Path("/CotizacionAgregar").HandlerFunc(CotizacionAgregar).Name(
		"CotizacionAgregar")
	router.Path("/CotizacionBorrar/{codigo}").HandlerFunc(CotizacionBorrar).
		Name("CotizacionBorrar")
	router.Path("/CotizacionEliminar/{codigo}").HandlerFunc(
		CotizacionEliminar).Name("CotizacionEliminar")
	router.Path("/CotizacionPdf/{codigo}").HandlerFunc(CotizacionPdf).Name(
		"CotizacionPdf")

	// COTIZACIONSERVICIO
	router.Path("/CotizacionservicioLista").HandlerFunc(CotizacionservicioLista).Name("CotizacionservicioLista")
	router.Path("/CotizacionservicioNuevo/{codigo}").HandlerFunc(CotizacionservicioNuevo).Name("CotizacionservicioNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/CotizacionservicioExiste/{codigo}").HandlerFunc(CotizacionservicioExiste).
		Name("CotizacionservicioExiste")
	router.Path("/CotizacionservicioEditar/{codigo}").HandlerFunc(CotizacionservicioEditar).
		Name("CotizacionservicioEditar")
	router.Path("/CotizacionservicioAgregar").HandlerFunc(CotizacionservicioAgregar).Name(
		"CotizacionservicioAgregar")
	router.Path("/CotizacionservicioBorrar/{codigo}").HandlerFunc(CotizacionservicioBorrar).
		Name("CotizacionservicioBorrar")
	router.Path("/CotizacionservicioEliminar/{codigo}").HandlerFunc(
		CotizacionservicioEliminar).Name("CotizacionservicioEliminar")
	router.Path("/CotizacionservicioPdf/{codigo}").HandlerFunc(CotizacionservicioPdf).Name(
		"CotizacionservicioPdf")

	// VENTA
	router.Path("/VentaLista").HandlerFunc(VentaLista).Name("VentaLista")
	router.Path("/VentaNuevo/{codigo}").HandlerFunc(VentaNuevo).Name("VentaNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/VentaExiste/{codigo}").HandlerFunc(VentaExiste).
		Name("VentaExiste")
	router.Path("/VentaEditar/{codigo}").HandlerFunc(VentaEditar).
		Name("VentaEditar")
	router.Path("/VentaAgregar").HandlerFunc(VentaAgregar).Name(
		"VentaAgregar")
	router.Path("/VentaBorrar/{codigo}").HandlerFunc(VentaBorrar).
		Name("VentaBorrar")
	router.Path("/VentaEliminar/{codigo}").HandlerFunc(
		VentaEliminar).Name("VentaEliminar")
	router.Path("/VentaPdf/{codigo}").HandlerFunc(VentaPdf).Name(
		"VentaPdf")

	// VENTA SERVICIO
	router.Path("/VentaservicioLista").HandlerFunc(VentaservicioLista).Name("VentaservicioLista")
	router.Path("/VentaservicioNuevo/{codigo}").HandlerFunc(VentaservicioNuevo).Name("VentaservicioNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/VentaservicioExiste/{codigo}").HandlerFunc(VentaservicioExiste).
		Name("VentaservicioExiste")
	router.Path("/VentaservicioEditar/{codigo}").HandlerFunc(VentaservicioEditar).
		Name("VentaservicioEditar")
	router.Path("/VentaservicioAgregar").HandlerFunc(VentaservicioAgregar).Name(
		"VentaservicioAgregar")
	router.Path("/VentaservicioBorrar/{codigo}").HandlerFunc(VentaservicioBorrar).
		Name("VentaservicioBorrar")
	router.Path("/VentaservicioEliminar/{codigo}").HandlerFunc(
		VentaservicioEliminar).Name("VentaservicioEliminar")
	router.Path("/VentaservicioPdf/{codigo}").HandlerFunc(VentaservicioPdf).Name(
		"VentaservicioPdf")

	// COMPRA
	router.Path("/CompraLista").HandlerFunc(CompraLista).Name("CompraLista")
	router.Path("/CompraNuevo/{codigo}").HandlerFunc(CompraNuevo).Name("CompraNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/CompraExiste/{codigo}").HandlerFunc(CompraExiste).
		Name("CompraExiste")
	router.Path("/CompraEditar/{codigo}").HandlerFunc(CompraEditar).
		Name("CompraEditar")
	router.Path("/CompraAgregar").HandlerFunc(CompraAgregar).Name(
		"CompraAgregar")
	router.Path("/CompraBorrar/{codigo}").HandlerFunc(CompraBorrar).
		Name("CompraBorrar")
	router.Path("/CompraEliminar/{codigo}").HandlerFunc(
		CompraEliminar).Name("CompraEliminar")
	router.Path("/CompraPdf/{codigo}").HandlerFunc(CompraPdf).Name(
		"CompraPdf")

	// TRAE EL PEDIDO EN LA COMPRA
	router.Path("/DatosPedido/{codigo}").HandlerFunc(DatosPedido).
		Name("DatosPedido")

	// TRAE EL PEDIDO SOPORTE EN EL SOPORTE
	router.Path("/DatosPedidosoporte/{codigo}").HandlerFunc(Datospedidosoporte).
		Name("DatosPedidosoporte")

	// TRAE EL PEDIDO SOPORTE SERVICIO EN EL SOPORTE SERVICIO
	router.Path("/DatosPedidosoporteservicio/{codigo}").HandlerFunc(Datospedidosoporteservicio).
		Name("DatosPedidosoporteservicio")

	// TRAE EL PEDIDO FACTURA GASTO LA FACTURA GASTO
	router.Path("/DatosPedidofacturagasto/{codigo}").HandlerFunc(Datospedidofacturagasto).
		Name("DatosPedidofacturagasto")

	// TRAE LA COMPRA EN DEVOLUCION
		router.Path("/DatosCompra/{codigo}").HandlerFunc(DatosCompra).
		Name("DatosCompra")

	// TRAE LA COTIZACION EN LA VENTA
	router.Path("/DatosCotizacion/{codigo}").HandlerFunc(DatosCotizacion).
		Name("DatosCotizacion")

	// TRAE LA COTIZACIONSERVICIO EN LA VENTA
	router.Path("/DatosCotizacionservicio/{codigo}").HandlerFunc(DatosCotizacionservicio).
		Name("DatosCotizacionservicio")

	// TRAE LA VENTA EN DEVOLUCION
	router.Path("/DatosVenta/{codigo}").HandlerFunc(DatosVenta).
		Name("DatosVenta")

	// TRAE LA VENTA SERVICIO EN DEVOLUCION
	router.Path("/DatosVentaservicio/{codigo}").HandlerFunc(DatosVentaservicio).
		Name("DatosVentaservicio")

	// TRAE EL SOPORTE EN LA DEVOLUCION
	router.Path("/DatosSoporte/{codigo}").HandlerFunc(DatosSoporte).
		Name("DatosSoporte")

	// TRAE EL SOPORTE SERVICIO  EN LA DEVOLUCION SOPORTE SERVICIO
	router.Path("/Datossoporteservicio/{codigo}").HandlerFunc(Datossoporteservicio).
		Name("Datossoporteservicio")

	// TRAE LA FACTURA GASTO  EN LA DEVOLUCION FACTURA GASTO
	router.Path("/Datosfacturagasto/{codigo}").HandlerFunc(Datosfacturagasto).
		Name("Datosfacturagasto")

	// PEDIDO
	router.Path("/PedidoLista").HandlerFunc(PedidoLista).Name("PedidoLista")
	router.Path("/PedidoNuevo/{codigo}").HandlerFunc(PedidoNuevo).Name("PedidoNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/PedidoExiste/{codigo}").HandlerFunc(PedidoExiste).
		Name("PedidoExiste")
	router.Path("/PedidoEditar/{codigo}").HandlerFunc(PedidoEditar).
		Name("PedidoEditar")
	router.Path("/PedidoAgregar").HandlerFunc(PedidoAgregar).Name(
		"PedidoAgregar")
	router.Path("/PedidoBorrar/{codigo}").HandlerFunc(PedidoBorrar).
		Name("PedidoBorrar")
	router.Path("/PedidoEliminar/{codigo}").HandlerFunc(
		PedidoEliminar).Name("PedidoEliminar")
	router.Path("/PedidoPdf/{codigo}").HandlerFunc(PedidoPdf).Name(
		"PedidoPdf")

	// PEDIDO SOPORTE
	router.Path("/PedidosoporteLista").HandlerFunc(PedidosoporteLista).Name("PedidosoporteLista")
	router.Path("/PedidosoporteNuevo/{codigo}").HandlerFunc(PedidosoporteNuevo).Name("PedidosoporteNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/PedidosoporteExiste/{codigo}").HandlerFunc(PedidosoporteExiste).
		Name("PedidosoporteExiste")
	router.Path("/PedidosoporteEditar/{codigo}").HandlerFunc(PedidosoporteEditar).
		Name("PedidosoporteEditar")
	router.Path("/PedidosoporteAgregar").HandlerFunc(PedidosoporteAgregar).Name(
		"PedidosoporteAgregar")
	router.Path("/PedidosoporteBorrar/{codigo}").HandlerFunc(PedidosoporteBorrar).
		Name("PedidosoporteBorrar")
	router.Path("/PedidosoporteEliminar/{codigo}").HandlerFunc(
		PedidosoporteEliminar).Name("PedidosoporteEliminar")
	router.Path("/PedidosoportePdf/{codigo}").HandlerFunc(PedidosoportePdf).Name(
		"PedidosoportePdf")

	// PEDIDO SOPORTE SERVICIO
	router.Path("/PedidosoporteservicioLista").HandlerFunc(PedidosoporteservicioLista).Name("PedidosoporteservicioLista")
	router.Path("/PedidosoporteservicioNuevo/{codigo}").HandlerFunc(PedidosoporteservicioNuevo).Name("PedidosoporteservicioNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/PedidosoporteservicioExiste/{codigo}").HandlerFunc(PedidosoporteservicioExiste).
		Name("PedidosoporteservicioExiste")
	router.Path("/PedidosoporteservicioEditar/{codigo}").HandlerFunc(PedidosoporteservicioEditar).
		Name("PedidosoporteservicioEditar")
	router.Path("/PedidosoporteservicioAgregar").HandlerFunc(PedidosoporteservicioAgregar).Name(
		"PedidosoporteservicioAgregar")
	router.Path("/PedidosoporteservicioBorrar/{codigo}").HandlerFunc(PedidosoporteservicioBorrar).
		Name("PedidosoporteservicioBorrar")
	router.Path("/PedidosoporteservicioEliminar/{codigo}").HandlerFunc(
		PedidosoporteservicioEliminar).Name("PedidosoporteservicioEliminar")
	router.Path("/PedidosoporteservicioPdf/{codigo}").HandlerFunc(PedidosoporteservicioPdf).Name(
		"PedidosoporteservicioPdf")

	// DEVOLUCION COMPRA
	router.Path("/DevolucioncompraLista").HandlerFunc(DevolucioncompraLista).Name("DevolucioncompraLista")
	router.Path("/DevolucioncompraNuevo/{codigo}").HandlerFunc(DevolucioncompraNuevo).Name("DevolucioncompraNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/DevolucioncompraExiste/{codigo}").HandlerFunc(DevolucioncompraExiste).
		Name("DevolucioncompraExiste")
	router.Path("/DevolucioncompraEditar/{codigo}").HandlerFunc(DevolucioncompraEditar).
		Name("DevolucioncompraEditar")
	router.Path("/DevolucioncompraAgregar").HandlerFunc(DevolucioncompraAgregar).Name(
		"DevolucioncompraAgregar")
	router.Path("/DevolucioncompraBorrar/{codigo}").HandlerFunc(DevolucioncompraBorrar).
		Name("DevolucioncompraBorrar")
	router.Path("/DevolucioncompraEliminar/{codigo}").HandlerFunc(
		DevolucioncompraEliminar).Name("DevolucioncompraEliminar")
	router.Path("/DevolucioncompraPdf/{codigo}").HandlerFunc(DevolucioncompraPdf).Name(
		"DevolucioncompraPdf")

	// DEVOLUCION VENTA
	router.Path("/DevolucionventaLista").HandlerFunc(DevolucionventaLista).Name("DevolucionventaLista")
	router.Path("/DevolucionventaNuevo/{codigo}").HandlerFunc(DevolucionventaNuevo).Name("DevolucionventaNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/DevolucionventaExiste/{codigo}").HandlerFunc(DevolucionventaExiste).
		Name("DevolucionventaExiste")
	router.Path("/DevolucionventaEditar/{codigo}").HandlerFunc(DevolucionventaEditar).
		Name("DevolucionventaEditar")
	router.Path("/DevolucionventaAgregar").HandlerFunc(DevolucionventaAgregar).Name(
		"DevolucionventaAgregar")
	router.Path("/DevolucionventaBorrar/{codigo}").HandlerFunc(DevolucionventaBorrar).
		Name("DevolucionventaBorrar")
	router.Path("/DevolucionventaEliminar/{codigo}").HandlerFunc(
		DevolucionventaEliminar).Name("DevolucionventaEliminar")
	router.Path("/DevolucionventaPdf/{codigo}").HandlerFunc(DevolucionventaPdf).Name(
		"DevolucionventaPdf")

	// DEVOLUCION VENTA SERVICIO
	router.Path("/DevolucionventaservicioLista").HandlerFunc(DevolucionventaservicioLista).Name("DevolucionventaservicioLista")
	router.Path("/DevolucionventaservicioNuevo").HandlerFunc(DevolucionventaservicioNuevo).Name("DevolucionventaservicioNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/DevolucionventaservicioExiste/{codigo}").HandlerFunc(DevolucionventaservicioExiste).
		Name("DevolucionventaservicioExiste")
	router.Path("/DevolucionventaservicioEditar/{codigo}").HandlerFunc(DevolucionventaservicioEditar).
		Name("DevolucionventaservicioEditar")
	router.Path("/DevolucionventaservicioAgregar").HandlerFunc(DevolucionventaservicioAgregar).Name(
		"DevolucionventaservicioAgregar")
	router.Path("/DevolucionventaservicioBorrar/{codigo}").HandlerFunc(DevolucionventaservicioBorrar).
		Name("DevolucionventaservicioBorrar")
	router.Path("/DevolucionventaservicioEliminar/{codigo}").HandlerFunc(
		DevolucionventaservicioEliminar).Name("DevolucionventaservicioEliminar")
	router.Path("/DevolucionventaservicioPdf/{codigo}").HandlerFunc(DevolucionventaservicioPdf).Name(
		"DevolucionventaservicioPdf")

	// INVENTARIO INICIAL
	router.Path("/InventarioinicialLista").HandlerFunc(InventarioinicialLista).Name("InventarioinicialLista")
	router.Path("/InventarioinicialNuevo").HandlerFunc(InventarioinicialNuevo).Name("InventarioinicialNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/InventarioinicialExiste/{codigo}").HandlerFunc(InventarioinicialExiste).
		Name("InventarioinicialExiste")
	router.Path("/InventarioinicialEditar/{codigo}").HandlerFunc(InventarioinicialEditar).
		Name("InventarioinicialEditar")
	router.Path("/InventarioinicialAgregar").HandlerFunc(InventarioinicialAgregar).Name(
		"InventarioinicialAgregar")
	router.Path("/InventarioinicialBorrar/{codigo}").HandlerFunc(InventarioinicialBorrar).
		Name("InventarioinicialBorrar")
	router.Path("/InventarioinicialEliminar/{codigo}").HandlerFunc(
		InventarioinicialEliminar).Name("InventarioinicialEliminar")
	router.Path("/InventarioinicialPdf/{codigo}").HandlerFunc(InventarioinicialPdf).Name(
		"InventarioinicialPdf")

	// ARCHIVO CONFIGURACION
	router.Path("/ConfiguracioninventarioNuevo/{panel}").HandlerFunc(ConfiguracioninventarioNuevo).Name("ConfiguracioninventarioNuevo")
	router.Path("/ConfiguracioninventarioInsertar").HandlerFunc(ConfiguracioninventarioInsertar).Name(
		"ConfiguracioninventarioInsertar")
	router.Path("/ConfiguracioninventarioPdf/{codigo:[0-9]+}").HandlerFunc(ConfiguracioninventarioPdf).Name(
		"ConfiguracioninventarioPdf")

	// ARCHIVO CONFIGURACION

	router.Path("/ConfiguracioncontabilidadNuevo/{panel}").HandlerFunc(ConfiguracioncontabilidadNuevo).Name("ConfiguracioninventarioNuevo")

	router.Path("/ConfiguracioncontabilidadInsertar").HandlerFunc(ConfiguracioncontabilidadInsertar).Name(
		"ConfiguracioncontabilidadInsertar")
	router.Path("/ConfiguracioncontabilidadPdf/{codigo:[0-9]+}").HandlerFunc(ConfiguracioncontabilidadPdf).Name(
		"ConfiguracioncontabilidadPdf")


	// ARCHIVO PLAN DE CUENTAS NIIF
	router.Path("/PlandecuentaniifLista").HandlerFunc(PlandecuentaniifLista).Name("PlandecuentaniifLista")
	router.Path("/PlandecuentaniifLista/{panel}/{codigo}/{elemento}").HandlerFunc(PlandecuentaniifLista).Name("PlandecuentaniifLista")

	router.Path("/PlandecuentaniifPdf/{codigo:[0-9]+}").HandlerFunc(PlandecuentaniifPdf).Name(
		"PlandecuentaniifPdf")

	// ARCHIVO PLAN DE CUENTAS PUC
	router.Path("/PlandecuentapucLista").HandlerFunc(PlandecuentapucLista).Name("PlandecuentapucLista")
	router.Path("/PlandecuentapucPdf/{codigo:[0-9]+}").HandlerFunc(PlandecuentapucPdf).Name(
		"PlandecuentapucPdf")

	// ARCHIVO PLAN DE CUENTAS EMPRESA
	router.Path("/PlandecuentaempresaBuscar/{codigo}").HandlerFunc(PlandecuentaempresaBuscar).
		Name("PlandecuentaempresaBuscar")

	router.Path("/PlandecuentaempresaBuscarAuxiliar/{codigo}").HandlerFunc(PlandecuentaempresaBuscarAuxiliar).
		Name("PlandecuentaempresaBuscarAuxiliar")

	router.Path("/PlandecuentaempresaNuevo/{panel}/{codigo}/{elemento}").HandlerFunc(
		PlandecuentaempresaNuevo).	Name("PlandecuentaempresaNuevo")

	router.Path("/PlandecuentaempresaNuevoCopia/{copiacodigo}").HandlerFunc(
		PlandecuentaempresaNuevoCopia).	Name("PlandecuentaempresaNuevoCopia")

	router.Path("/PlandecuentaempresaInsertar").HandlerFunc(PlandecuentaempresaInsertar).Name(
		"PlandecuentaempresaInsertar")
	router.Path("/PlandecuentaempresaAgregar").HandlerFunc(PlandecuentaempresaAgregar).Name(
		"PlandecuentaempresaAgregar")
	router.Path("/PlandecuentaempresaActual/{codigo}").HandlerFunc(PlandecuentaempresaActual).
		Name("PlandecuentaempresaActual")
	router.Path("/PlandecuentaempresaLista/{panel}/{codigo}/{elemento}").HandlerFunc(PlandecuentaempresaLista).Name("PlandecuentaempresaLista")
	router.Path("/PlandecuentaempresaExiste/{codigo:[0-9]+}").HandlerFunc(PlandecuentaempresaExiste).
		Name("PlandecuentaempresaExiste")
	router.Path("/PlandecuentaempresaEditar/{codigo:[0-9]+}").HandlerFunc(PlandecuentaempresaEditar).Name(
		"PlandecuentaempresaEditar")
	router.Path("/PlandecuentaempresaActualizar/{codigo}").HandlerFunc(PlandecuentaempresaActualizar).Name(
		"PlandecuentaempresaActualizar")
	router.Path("/PlandecuentaempresaBorrar/{codigo}").HandlerFunc(PlandecuentaempresaBorrar).Name(
		"PlandecuentaempresaBorrar")
	router.Path("/PlandecuentaempresaEliminar/{codigo}").HandlerFunc(PlandecuentaempresaEliminar).Name("PlandecuentaempresaEliminar")
	router.Path("/PlandecuentaempresaPdf/{codigo:[0-9]+}").HandlerFunc(PlandecuentaempresaPdf).Name(
		"PlandecuentaempresaPdf")

	router.Path("/PlandecuentaempresaTodosPdf").HandlerFunc(PlandecuentaempresaTodosPdf).
		Name("PlandecuentaempresaTodosPdf")
	router.Path("/PlandecuentaempresaExcel").HandlerFunc(PlandecuentaempresaExcel).
		Name("PlandecuentaempresaExcel")


	// ARCHIVO RETENCION EN LA FUENTE
	router.Path("/RetencionenlafuenteLista").HandlerFunc(RetencionenlafuenteLista).Name("RetencionenlafuenteLista")
	router.Path("/RetencionenlafuentePdf/{codigo:[0-9]+}").HandlerFunc(RetencionenlafuentePdf).Name(
		"RetencionenlafuentePdf")

	// ARCHIVO DEPRECIACION
	router.Path("/DepreciacionLista").HandlerFunc(DepreciacionLista).Name("DepreciacionLista")
	router.Path("/DepreciacionPdf/{codigo:[0-9]+}").HandlerFunc(DepreciacionPdf).Name(
		"DepreciacionPdf")

	// COMPROBANTE
	router.Path("/ComprobanteLista").HandlerFunc(ComprobanteLista).Name("ComprobanteLista")
	router.Path("/ComprobanteNuevo/{documento}/{numero}").HandlerFunc(ComprobanteNuevo).Name("ComprobanteNuevo")
	router.Path("/ComprobanteExiste/{documento}/{numero}").HandlerFunc(ComprobanteExiste).
		Name("ComprobanteExiste")
	router.Path("/ComprobanteAgregar").HandlerFunc(ComprobanteAgregar).Name(
		"ComprobanteAgregar")
	router.Path("/ComprobanteEditar/{documento}/{numero}").HandlerFunc(ComprobanteEditar).
		Name("ComprobanteEditar")
	router.Path("/ComprobanteBorrar/{documento}/{numero}").HandlerFunc(ComprobanteBorrar).
		Name("ComprobanteBorrar")
	router.Path("/ComprobanteEliminar/{documento}/{numero}").HandlerFunc(
		ComprobanteEliminar).Name("ComprobanteEliminar")
	router.Path("/ComprobantePdf/{documento}/{numero}").HandlerFunc(ComprobantePdf).Name(
		"ComprobantePdf")

	router.Path("/ComprobanteTodosPdf").HandlerFunc(ComprobanteTodosPdf).
		Name("ComprobanteTodosPdf")
	router.Path("/ComprobanteExcel").HandlerFunc(ComprobanteExcel).
		Name("ComprobanteExcel")

	// CUENTADECOBRO
	router.Path("/CuentadecobroGenerar").HandlerFunc(CuentadecobroGenerar).Name("CuentadecobroGenerar")
	router.Path("/CuentadecobroGenerarMes/{mes}/{centro}/{porcentaje}").HandlerFunc(CuentadecobroGenerarMes).Name("CuentadecobroGenerarMes")

	router.Path("/CuentadecobroLista").HandlerFunc(CuentadecobroLista).Name("CuentadecobroLista")
	router.Path("/CuentadecobroNuevo").HandlerFunc(CuentadecobroNuevo).Name("CuentadecobroNuevo")
	router.Path("/CuentadecobroExiste/{numero}").HandlerFunc(CuentadecobroExiste).
		Name("CuentadecobroExiste")
	router.Path("/CuentadecobroAgregar").HandlerFunc(CuentadecobroAgregar).Name(
		"CuentadecobroAgregar")
	router.Path("/CuentadecobroEditar/{numero}").HandlerFunc(CuentadecobroEditar).
		Name("CuentadecobroEditar")
	router.Path("/CuentadecobroBorrar/{numero}").HandlerFunc(CuentadecobroBorrar).
		Name("CuentadecobroBorrar")
	router.Path("/CuentadecobroEliminar/{numero}").HandlerFunc(
		CuentadecobroEliminar).Name("CuentadecobroEliminar")
	router.Path("/CuentadecobroPdf/{numero}").HandlerFunc(CuentadecobroPdf).Name(
		"CuentadecobroPdf")

	// CUENTA DE COBRO DATO
	router.Path("/CuentadecobroDato").HandlerFunc(CuentadecobroDato).Name("CuentadecobroGenerar")
	router.Path("/CuentadecobroDatoAgregar").HandlerFunc(CuentadecobroDatoAgregar).Name("CuentadecobroGenerarAgregar")

	// ARCHIVO INVENTARIO
	router.Path("/InventarioLista").HandlerFunc(InventarioLista).Name("InventarioLista")

	// CONCILIACION

	router.Path("/ConciliacionLista").HandlerFunc(ConciliacionLista).Name("ConciliacionLista")
	router.Path("/ConciliacionDatos/{cuenta}/{mes}").HandlerFunc(ConciliacionDato).Name("ConciliacionDato")
	router.Path("/ConciliacionInsertar").HandlerFunc(ConciliacionInsertar).Name("ConciliacionInsertar")
	router.Path("/ConciliacionTodosPdf/{cuenta}/{mes}").HandlerFunc(ConciliacionTodosPdf).Name("ConciliacionCuentaTodosPdf")

	// BANCO
	router.Path("/BancoLista").HandlerFunc(BancoLista).Name("BancoLista")
	router.Path("/BancoDatos/{tercero}/{documento}").HandlerFunc(BancoDato).Name("BancoDato")
	router.Path("/BancoDatoAgregar").HandlerFunc(BancoDatoAgregar).Name("BancoDatoAgregar")



	// ARCHIVO KARDEX
	router.Path("/KardexLista").HandlerFunc(KardexLista).Name("KardexLista")
	router.Path("/KardexDatos/{codigo}/{fechainicial}/{fechafinal}/{bodega}/{tipo}/{discriminar}").HandlerFunc(KardexDatos).Name("KardexDatos")
	router.Path("/KardexDatosTodos/{fechainicial}/{fechafinal}/{bodega}/{tipo}/{discriminar}").HandlerFunc(KardexDatosTodos).Name("KardexDatosTodos")

	// KARDEX COSTO
	router.Path("/KardexCostoLista").HandlerFunc(KardexCostoLista).Name("KardexCostoLista")
	router.Path("/KardexCostoDatos/{mes}/{centro}").HandlerFunc(KardexCostoDatos).Name("KardexCostoGenerar")


	// BALANCE DE PRUEBA
	router.Path("/BalancedepruebaLista").HandlerFunc(BalancedepruebaLista).Name("BalancedepruebaLista")
	router.Path("/BalancedepruebaDatos").HandlerFunc(BalancedepruebaDatos).Name("BalancedepruebaDatos")
	//router.Path("/BalancedepruebaPdf").HandlerFunc(BalancedepruebaPdf).Name("BalancedepruebaPdf")
	router.Path("/BalancedepruebaPdf/{FechaInicial}/{FechaFinal}/{CuentaInicial}/{CuentaFinal}/{TerceroInicial}/{TerceroFinal}/{CentroInicial}/{CentroFinal}/{DocumentoInicial}/{DocumentoFinal}/{NumeroInicial}/{NumeroFinal}/{Detalle}/{Nivel}/{Activa}/{Subtotal}").HandlerFunc(BalancedepruebaPdf).Name("BalancedepruebaPdf")
	router.Path("/BalancedepruebaExcel/{FechaInicial}/{FechaFinal}/{CuentaInicial}/{CuentaFinal}/{TerceroInicial}/{TerceroFinal}/{CentroInicial}/{CentroFinal}/{DocumentoInicial}/{DocumentoFinal}/{NumeroInicial}/{NumeroFinal}/{Detalle}/{Nivel}/{Activa}/{Subtotal}").HandlerFunc(BalancedepruebaExcel).Name("BalancedepruebaExcel")

	//{FechaInicial}/{FechaFinal}/{CuentaInicial}/{CuentaFinal}/{TerceroInicial}/{TerceroFinal}/{CentroInicial}/{CentroFinal}/{DocumentoInicial}/{DocumentoFinal}/{NumeroInicial}/{NumeroFinal}/{Detalle}/{Nivel}/{Activa}/{Subtotal}

	// ARCHIVO CONCEPTO
	router.Path("/ConceptoNuevo/{codigo}").HandlerFunc(ConceptoNuevo).Name("ConceptoNuevo")
	router.Path("/ConceptoBuscar/{codigo}").HandlerFunc(ConceptoBuscar).
		Name("ConceptoBuscar")
	router.Path("/ConceptoActual/{codigo}").HandlerFunc(ConceptoActual).
		Name("ConceptoActual")
	router.Path("/ConceptoLista").HandlerFunc(ConceptoLista).Name("ConceptoLista")
	router.Path("/ConceptoExiste/{codigo:[0-9]+}").HandlerFunc(ConceptoExiste).
		Name("ConceptoExiste")
	router.Path("/ConceptoInsertar").HandlerFunc(ConceptoInsertar).Name(
		"ConceptoInsertar")
	router.Path("/ConceptoActualizar/{codigo:[0-9]+}").HandlerFunc(ConceptoActualizar).Name(
		"ConceptoActualizar")
	router.Path("/ConceptoBorrar/{codigo:[0-9]+}").HandlerFunc(ConceptoBorrar).Name(
		"ConceptoBorrar")
	router.Path("/ConceptoEliminar/{codigo:[0-9]+}").HandlerFunc(
		ConceptoEliminar).Name("ConceptoEliminar")
	router.Path("/ConceptoEditar/{codigo:[0-9]+}").HandlerFunc(ConceptoEditar).Name(
		"ConceptoEditar")
	router.Path("/ConceptoPdf/{codigo:[0-9]+}").HandlerFunc(ConceptoPdf).Name(
		"ConceptoPdf")

	router.Path("/ConceptoTodosPdf").HandlerFunc(ConceptoTodosPdf).
		Name("ConceptoTodosPdf")
	router.Path("/ConceptoExcel").HandlerFunc(ConceptoExcel).
		Name("ConceptoExcel")



	// SOPORTE
	router.Path("/SoporteLista").HandlerFunc(SoporteLista).Name("SoporteLista")
	router.Path("/SoporteNuevo/{codigo}").HandlerFunc(SoporteNuevo).Name("SoporteNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/SoporteExiste/{codigo}").HandlerFunc(SoporteExiste).
		Name("SoporteExiste")
	router.Path("/SoporteEditar/{codigo}").HandlerFunc(SoporteEditar).
		Name("SoporteEditar")
	router.Path("/SoporteAgregar").HandlerFunc(SoporteAgregar).Name(
		"SoporteAgregar")
	router.Path("/SoporteBorrar/{codigo}").HandlerFunc(SoporteBorrar).
		Name("SoporteBorrar")
	router.Path("/SoporteEliminar/{codigo}").HandlerFunc(
		SoporteEliminar).Name("SoporteEliminar")
	router.Path("/SoportePdf/{codigo}").HandlerFunc(SoportePdf).Name(
		"SoportePdf")

	// SOPORTE SERVICIO
	router.Path("/SoporteservicioLista").HandlerFunc(SoporteservicioLista).Name("SoporteservicioLista")
	router.Path("/SoporteservicioNuevo/{codigo}").HandlerFunc(SoporteservicioNuevo).Name("SoporteservicioNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/SoporteservicioExiste/{codigo}").HandlerFunc(SoporteservicioExiste).
		Name("SoporteservicioExiste")
	router.Path("/SoporteservicioEditar/{codigo}").HandlerFunc(SoporteservicioEditar).
		Name("SoporteservicioEditar")
	router.Path("/SoporteservicioAgregar").HandlerFunc(SoporteservicioAgregar).Name(
		"SoporteservicioAgregar")
	router.Path("/SoporteservicioBorrar/{codigo}").HandlerFunc(SoporteservicioBorrar).
		Name("SoporteservicioBorrar")
	router.Path("/SoporteservicioEliminar/{codigo}").HandlerFunc(
		SoporteservicioEliminar).Name("SoporteservicioEliminar")
	router.Path("/SoporteservicioPdf/{codigo}").HandlerFunc(SoporteservicioPdf).Name(
		"SoporteservicioPdf")

	// DEVOLUCION SOPORTE
	router.Path("/DevolucionsoporteLista").HandlerFunc(DevolucionsoporteLista).Name("DevolucionsoporteLista")
	router.Path("/DevolucionsoporteNuevo/{codigo}").HandlerFunc(DevolucionsoporteNuevo).Name("DevolucionsoporteNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/DevolucionsoporteExiste/{codigo}").HandlerFunc(DevolucionsoporteExiste).
		Name("DevolucionsoporteExiste")
	router.Path("/DevolucionsoporteEditar/{codigo}").HandlerFunc(DevolucionsoporteEditar).
		Name("DevolucionsoporteEditar")
	router.Path("/DevolucionsoporteAgregar").HandlerFunc(DevolucionsoporteAgregar).Name(
		"DevolucionsoporteAgregar")
	router.Path("/DevolucionsoporteBorrar/{codigo}").HandlerFunc(DevolucionsoporteBorrar).
		Name("DevolucionsoporteBorrar")
	router.Path("/DevolucionsoporteEliminar/{codigo}").HandlerFunc(
		DevolucionsoporteEliminar).Name("DevolucionsoporteEliminar")
	router.Path("/DevolucionsoportePdf/{codigo}").HandlerFunc(DevolucionsoportePdf).Name(
		"DevolucionsoportePdf")

	// DEVOLUCION SOPORTE SERVICIO
	router.Path("/DevolucionsoporteservicioLista").HandlerFunc(DevolucionsoporteservicioLista).Name("DevolucionsoporteservicioLista")
	router.Path("/DevolucionsoporteservicioNuevo/{codigo}").HandlerFunc(DevolucionsoporteservicioNuevo).Name("DevolucionsoporteservicioNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/DevolucionsoporteservicioExiste/{codigo}").HandlerFunc(DevolucionsoporteservicioExiste).
		Name("DevolucionsoporteservicioExiste")
	router.Path("/DevolucionsoporteservicioEditar/{codigo}").HandlerFunc(DevolucionsoporteservicioEditar).
		Name("DevolucionsoporteservicioEditar")
	router.Path("/DevolucionsoporteservicioAgregar").HandlerFunc(DevolucionsoporteservicioAgregar).Name(
		"DevolucionsoporteservicioAgregar")
	router.Path("/DevolucionsoporteservicioBorrar/{codigo}").HandlerFunc(DevolucionsoporteservicioBorrar).
		Name("DevolucionsoporteservicioBorrar")
	router.Path("/DevolucionsoporteservicioEliminar/{codigo}").HandlerFunc(
		DevolucionsoporteservicioEliminar).Name("DevolucionsoporteservicioEliminar")
	router.Path("/DevolucionsoporteservicioPdf/{codigo}").HandlerFunc(DevolucionsoporteservicioPdf).Name(
		"DevolucionsoporteservicioPdf")

	// PEDIDO FACTURA GASTO
	router.Path("/PedidofacturagastoLista").HandlerFunc(PedidofacturagastoLista).Name("PedidofacturagastoLista")
	router.Path("/PedidofacturagastoNuevo").HandlerFunc(PedidofacturagastoNuevo).Name("PedidofacturagastoNuevo")
	router.Path("/PedidofacturagastoExiste/{codigo}").HandlerFunc(PedidofacturagastoExiste).
		Name("PedidofacturagastoExiste")
	router.Path("/PedidofacturagastoEditar/{codigo}").HandlerFunc(PedidofacturagastoEditar).
		Name("PedidofacturagastoEditar")
	router.Path("/PedidofacturagastoAgregar").HandlerFunc(PedidofacturagastoAgregar).Name(
		"PedidofacturagastoAgregar")
	router.Path("/PedidofacturagastoBorrar/{codigo}").HandlerFunc(PedidofacturagastoBorrar).
		Name("PedidofacturagastoBorrar")
	router.Path("/PedidofacturagastoEliminar/{codigo}").HandlerFunc(
		PedidofacturagastoEliminar).Name("PedidofacturagastoEliminar")
	router.Path("/PedidofacturagastoPdf/{codigo}").HandlerFunc(PedidofacturagastoPdf).Name(
		"PedidofacturagastoPdf")

	//router.Path("/PedidofacturagastoTodosPdf").HandlerFunc(PedidofacturagastoTodosPdf).
	//	Name("PedidofacturagastoTodosPdf")
	//router.Path("/PedidofacturagastoExcel").HandlerFunc(PedidofacturagastoExcel).
	//	Name("PedidofacturagastoExcel")

	// FACTURA GASTO
	router.Path("/FacturagastoLista").HandlerFunc(FacturagastoLista).Name("FacturagastoLista")
	router.Path("/FacturagastoNuevo").HandlerFunc(FacturagastoNuevo).Name("FacturagastoNuevo")
	router.Path("/FacturagastoExiste/{codigo}").HandlerFunc(FacturagastoExiste).
		Name("FacturagastoExiste")
	router.Path("/FacturagastoEditar/{codigo}").HandlerFunc(FacturagastoEditar).
		Name("FacturagastoEditar")
	router.Path("/FacturagastoAgregar").HandlerFunc(FacturagastoAgregar).Name(
		"FacturagastoAgregar")
	router.Path("/FacturagastoBorrar/{codigo}").HandlerFunc(FacturagastoBorrar).
		Name("FacturagastoBorrar")
	router.Path("/FacturagastoEliminar/{codigo}").HandlerFunc(
		FacturagastoEliminar).Name("FacturagastoEliminar")
	router.Path("/FacturagastoPdf/{codigo}").HandlerFunc(FacturagastoPdf).Name(
		"FacturagastoPdf")

	// FACTURA GASTO DEVOLUCION
	router.Path("/DevolucionfacturagastoLista").HandlerFunc(DevolucionfacturagastoLista).Name("DevolucionfacturagastoLista")
	router.Path("/DevolucionfacturagastoNuevo").HandlerFunc(DevolucionfacturagastoNuevo).Name("DevolucionfacturagastoNuevo")
	router.Path("/DevolucionfacturagastoExiste/{codigo}").HandlerFunc(DevolucionfacturagastoExiste).
		Name("DevolucionfacturagastoExiste")
	router.Path("/DevolucionfacturagastoEditar/{codigo}").HandlerFunc(DevolucionfacturagastoEditar).
		Name("DevolucionfacturagastoEditar")
	router.Path("/DevolucionfacturagastoAgregar").HandlerFunc(DevolucionfacturagastoAgregar).Name(
		"DevolucionfacturagastoAgregar")
	router.Path("/DevolucionfacturagastoBorrar/{codigo}").HandlerFunc(DevolucionfacturagastoBorrar).
		Name("DevolucionfacturagastoBorrar")
	router.Path("/DevolucionfacturagastoEliminar/{codigo}").HandlerFunc(
		DevolucionfacturagastoEliminar).Name("DevolucionfacturagastoEliminar")
	router.Path("/DevolucionfacturagastoPdf/{codigo}").HandlerFunc(DevolucionfacturagastoPdf).Name(
		"DevolucionfacturagastoPdf")


	// TRASLADO AJUSTES
	router.Path("/TrasladoLista").HandlerFunc(TrasladoLista).Name("TrasladoLista")
	router.Path("/TrasladoNuevo").HandlerFunc(TrasladoNuevo).Name("TrasladoNuevo")
	router.Path("/BodegaLlenar").HandlerFunc(BodegaLlenar).Name("BodegaLlenar")
	router.Path("/TrasladoExiste/{codigo}").HandlerFunc(TrasladoExiste).
		Name("TrasladoExiste")
	router.Path("/TrasladoEditar/{codigo}").HandlerFunc(TrasladoEditar).
		Name("TrasladoEditar")
	router.Path("/TrasladoAgregar").HandlerFunc(TrasladoAgregar).Name(
		"TrasladoAgregar")
	router.Path("/TrasladoBorrar/{codigo}").HandlerFunc(TrasladoBorrar).
		Name("TrasladoBorrar")
	router.Path("/TrasladoEliminar/{codigo}").HandlerFunc(
		TrasladoEliminar).Name("TrasladoEliminar")
	router.Path("/TrasladoPdf/{codigo}").HandlerFunc(TrasladoPdf).Name(
		"TrasladoPdf")

	// ARCHIVO FINANCIERO
	router.Path("/FinancieroNuevo/{codigo}").HandlerFunc(FinancieroNuevo).Name("FinancieroNuevo")
	router.Path("/FinancieroBuscar/{codigo}").HandlerFunc(FinancieroBuscar).
		Name("FinancieroBuscar")
	router.Path("/FinancieroActual/{codigo}").HandlerFunc(FinancieroActual).
		Name("FinancieroActual")
	router.Path("/FinancieroLista").HandlerFunc(FinancieroLista).Name("FinancieroLista")
	router.Path("/FinancieroExiste/{codigo:[0-9]+}").HandlerFunc(FinancieroExiste).
		Name("FinancieroExiste")
	router.Path("/FinancieroInsertar").HandlerFunc(FinancieroInsertar).Name(
		"FinancieroInsertar")
	router.Path("/FinancieroActualizar/{codigo:[0-9]+}").HandlerFunc(FinancieroActualizar).Name(
		"FinancieroActualizar")
	router.Path("/FinancieroBorrar/{codigo:[0-9]+}").HandlerFunc(FinancieroBorrar).Name(
		"FinancieroBorrar")
	router.Path("/FinancieroEliminar/{codigo:[0-9]+}").HandlerFunc(
		FinancieroEliminar).Name("FinancieroEliminar")
	router.Path("/FinancieroEditar/{codigo:[0-9]+}").HandlerFunc(FinancieroEditar).Name(
		"FinancieroEditar")
	router.Path("/FinancieroPdf/{codigo:[0-9]+}").HandlerFunc(FinancieroPdf).Name(
		"FinancieroPdf")

	router.Path("/FinancieroTodosPdf").HandlerFunc(FinancieroTodosPdf).
		Name("FinancieroTodosPdf")
	router.Path("/FinancieroExcel").HandlerFunc(FinancieroExcel).
		Name("FinancieroExcel")

	// ARCHIVO PROPIEDAD
	router.Path("/PropiedadNuevo").HandlerFunc(PropiedadNuevo).Name("PropiedadNuevo")
	router.Path("/PropiedadBuscar/{codigo}").HandlerFunc(PropiedadBuscar).
		Name("PropiedadBuscar")
	router.Path("/PropiedadActual/{codigo}").HandlerFunc(PropiedadActual).
		Name("PropiedadActual")
	router.Path("/PropiedadLista").HandlerFunc(PropiedadLista).Name("PropiedadLista")
	router.Path("/PropiedadExiste/{codigo:[0-9]+}").HandlerFunc(PropiedadExiste).
		Name("PropiedadExiste")
	router.Path("/PropiedadInsertar").HandlerFunc(PropiedadInsertar).Name(
		"PropiedadInsertar")
	router.Path("/PropiedadActualizar/{codigo:[0-9]+}").HandlerFunc(PropiedadActualizar).Name(
		"PropiedadActualizar")
	router.Path("/PropiedadBorrar/{codigo:[0-9]+}").HandlerFunc(PropiedadBorrar).Name(
		"PropiedadBorrar")
	router.Path("/PropiedadEliminar/{codigo:[0-9]+}").HandlerFunc(
		PropiedadEliminar).Name("PropiedadEliminar")
	router.Path("/PropiedadEditar/{codigo:[0-9]+}").HandlerFunc(PropiedadEditar).Name(
		"PropiedadEditar")
	router.Path("/PropiedadPdf/{codigo:[0-9]+}").HandlerFunc(PropiedadPdf).Name(
		"PropiedadPdf")

	router.Path("/PropiedadTodosPdf").HandlerFunc(PropiedadTodosPdf).
		Name("PropiedadTodosPdf")
	router.Path("/PropiedadExcel").HandlerFunc(PropiedadExcel).
		Name("PropiedadExcel")

	router.Path("/PropiedadGenerar").HandlerFunc(PropiedadGenerar).
		Name("PropiedadGenerar")

	router.Path("/PropiedadMes/{mes}").HandlerFunc(PropiedadMes).
		Name("PropiedadMes")

	// ARCHIVO DIFERIDO
	router.Path("/DiferidoNuevo").HandlerFunc(DiferidoNuevo).Name("DiferidoNuevo")
	router.Path("/DiferidoBuscar/{codigo}").HandlerFunc(DiferidoBuscar).
		Name("DiferidoBuscar")
	router.Path("/DiferidoActual/{codigo}").HandlerFunc(DiferidoActual).
		Name("DiferidoActual")
	router.Path("/DiferidoLista").HandlerFunc(DiferidoLista).Name("DiferidoLista")
	router.Path("/DiferidoExiste/{codigo:[0-9]+}").HandlerFunc(DiferidoExiste).
		Name("DiferidoExiste")
	router.Path("/DiferidoInsertar").HandlerFunc(DiferidoInsertar).Name(
		"DiferidoInsertar")
	router.Path("/DiferidoActualizar/{codigo:[0-9]+}").HandlerFunc(DiferidoActualizar).Name(
		"DiferidoActualizar")
	router.Path("/DiferidoBorrar/{codigo:[0-9]+}").HandlerFunc(DiferidoBorrar).Name(
		"DiferidoBorrar")
	router.Path("/DiferidoEliminar/{codigo:[0-9]+}").HandlerFunc(
		DiferidoEliminar).Name("DiferidoEliminar")
	router.Path("/DiferidoEditar/{codigo:[0-9]+}").HandlerFunc(DiferidoEditar).Name(
		"DiferidoEditar")
	router.Path("/DiferidoPdf/{codigo:[0-9]+}").HandlerFunc(DiferidoPdf).Name(
		"DiferidoPdf")

	router.Path("/DiferidoTodosPdf").HandlerFunc(DiferidoTodosPdf).
		Name("DiferidoTodosPdf")
	router.Path("/DiferidoExcel").HandlerFunc(DiferidoExcel).
		Name("DiferidoExcel")

	router.Path("/DiferidoGenerar").HandlerFunc(DiferidoGenerar).
		Name("DiferidoGenerar")

	router.Path("/DiferidoMes/{mes}").HandlerFunc(DiferidoMes).
		Name("DiferidoMes")


	// ARCHIVO CUOTA
	router.Path("/CuotaLista").HandlerFunc(CuotaLista).Name("CuotaLista")
	router.Path("/CuotaDatos/{monto}/{plazo}/{intereses}/{fechainicial}").HandlerFunc(CuotaDatos).Name("CuotaDatos")

	router.Path("/CuotaTodosPdf/{monto}/{plazo}/{intereses}/{fechainicial}").HandlerFunc(CuotaTodosPdf).
		Name("CuotaTodosPdf")
	router.Path("/CuotaExcel/{monto}/{plazo}/{intereses}/{fechainicial}").HandlerFunc(CuotaExcel).
		Name("CuotaExcel")


	// LOCAL HOST 9002
	log.Println("Servidor Corriendo en "+ruta)
	if err := http.ListenAndServe(":9002", router); err != nil {
		log.Fatal(err)
	}
}

// FORMATO DE FECHA
var decoder = schema.NewDecoder()
var timeConverter = func(value string) reflect.Value {
	if v, err := time.Parse("2006-01-02", value); err == nil {

		return reflect.ValueOf(v)
	}
	return reflect.Value{} // this is the same as the private const invalidType
}