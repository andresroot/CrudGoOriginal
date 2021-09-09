package main

// INICIA CUENTADECOBRO IMPORTAR PAQUETES
import (
	"bytes"
	"math"
	"strings"

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

// TERMINA CUENTADECOBRO IMPORTAR PAQUETES

// INICIA CUENTADECOBRO ESTRUCTURA JSON
type cuentadecobroJson struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Nombre string `json:"nombre"`
}

// TERMINA CUENTADECOBRO ESTRUCTURA JSON

// INICIA CUENTADECOBRO ESTRUCTURA
type cuentadecobroLista struct {
	Numero     	        string
	Fecha         		time.Time
	Tercero				string
	Terceronombre		string
	Total				string
}

type cuentadecobroDato struct {
	Tercero	    	string
	Nombre      	string
	Descuento1      string
	Descuento2      string
	Cuotap          string
	Cuota1          string
	Cuota2          string
}

type cuentaTercero struct {
	Cuenta     	        string
}

type cuentaSaldo struct {
	Debito     	        string
	Credito             string
}

// INICIA CUENTADECOBRO ESTRUCTURA
type cuentadecobro struct {
	Numero				string
	Fecha               time.Time
	Centro				string
	Tercero             string
	Totalanterior		string
	Totalactual			string
	Total				string
	Accion              string
	Detalle             []cuentadecobrodetalle `json:"Detalle"`
	DetalleEditar		[]cuentadecobrodetalleeditar `json:"DetalleEditar"`
}

// TERMINA CUENTADECOBRO ESTRUCTURA

// INICIA CUENTADECOBRO DETALLE ESTRUCTURA
type cuentadecobrodetalle struct {
	Fila				string
	Numero				string
	Cuenta  			string
	Anterior			string
	Actual				string
	Total    			string
}

type cuentadecobrodetalleGenerar struct {
	Numero				string
	Tercero  			string
	TerceroNombre		string
	Totalanterior		string
	Totalactual			string
	Total    			string
}

// estructura para editar
type cuentadecobrodetalleeditar struct {
	Fila				string
	Numero				string
	Cuenta  			string
	Cuentanombre		string
	Anterior			string
	Actual				string
	Total    			string
}

// TERMINA COMPRA DETALLE EDITAR

// INICIA COMPRA CONSULTA DETALLE
func CuentadecobroConsultaDetalle() string {
	var consulta = ""
	consulta = "select "
	consulta += "cuentadecobrodetalle.Fila as fila, "
	consulta += "cuentadecobrodetalle.Cuenta as cuenta, "
	consulta += "plandecuentaempresa.nombre as cuentanombre, "
	consulta += "cuentadecobrodetalle.Anterior as anterior, "
	consulta += "cuentadecobrodetalle.Actual as actual, "
	consulta += "cuentadecobrodetalle.Total as total "
	consulta += "from cuentadecobrodetalle "
	consulta += " inner join plandecuentaempresa on Plandecuentaempresa.codigo = Cuentadecobrodetalle.cuenta "
	consulta += " where cuentadecobrodetalle.numero=$1 "
	consulta += " order by fila "
	log.Println(consulta)
	return consulta
}
func saldoAnterior(tercero string,fechaFinal time.Time,micuenta string) float64{
	var consultasaldoAnterior="SELECT debito,credito FROM comprobantedetalle  where"
	consultasaldoAnterior += " tercero=$1 and fecha<=$2 and cuenta=$3"

	listadoSaldoTercero := []cuentaSaldo{}
	db.Select(&listadoSaldoTercero, consultasaldoAnterior,tercero,fechaFinal,micuenta)
	log.Println("calculo cuenta anterior "+micuenta)
	var totalSaldo float64
	totalSaldo=0
	for _, miSaldo := range listadoSaldoTercero {
		//var c=k
		log.Println(miSaldo.Debito)
		log.Println(miSaldo.Credito)
		totalSaldo=Flotante(miSaldo.Debito)-Flotante(miSaldo.Credito)
	}
	return totalSaldo

}
// cuenta de cobro todos
func CuentadecobroGenerarMes(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	mes := mux.Vars(r)["mes"]
	miCentro := mux.Vars(r)["centro"]
	miPorcentaje:= mux.Vars(r)["porcentaje"]

	log.Println(miPorcentaje)

	var Documentocontable string
	Documentocontable = "30"
	var miPorcentajeNumero float64
	miPorcentajeNumero,err:=strconv.ParseFloat(miPorcentaje,64)
	if err != nil {
		fmt.Println(err)
		return
	}
	parametroscontabilidad := configuracioncontabilidad{}
	parametroscontabilidad=TraerParametrosContabilidad()



	//
	// borra datos anteriores
	listadoCuentaCobroBorrar := []cuentadecobro{}

	var consultaborra="select * from cuentadecobro where EXTRACT(MONTH FROM  fecha)>=$1"
	db.Select(&listadoCuentaCobroBorrar,consultaborra, mes)

	for _, miCuentaBorra := range listadoCuentaCobroBorrar {

		var consultaborradetalle="delete from cuentadecobrodetalle where numero=$1"
		db.Exec(consultaborradetalle,miCuentaBorra.Numero)

		consultaborradetalle="delete from cuentadecobro where numero=$1"
		db.Exec(consultaborradetalle,miCuentaBorra.Numero)

		// BORRA MOVIMIENTOS
		var consultaborracomprobante="delete from comprobante where documento=$2 and numero=$1"
		db.Exec(consultaborracomprobante,miCuentaBorra.Numero,Documentocontable)

		var consultaborracomprobantedetalle="delete from comprobantedetalle where documento=$2 and  numero=$1"
		db.Exec(consultaborracomprobantedetalle,miCuentaBorra.Numero,Documentocontable)

	}

	var ultimo int

	if parametroscontabilidad.Phinicial=="0"{

	}else{
		db := dbConn()
		Numero:= parametroscontabilidad.Phinicial
		var total int
		row := db.QueryRow("SELECT COUNT(*) FROM cuentadecobro  WHERE  Numero=$1",  Numero)
		err := row.Scan(&total)
		if err != nil {
			log.Fatal(err)
		}
		//var resultado bool
		if total > 0 {
			// BUSCAR ULTIMO
			row := db.QueryRow("SELECT MAX(CAST(NUMERO AS INTEGER)) AS NUMERO FROM cuentadecobro")
			err := row.Scan(&ultimo)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			ultimo,_ = strconv.Atoi(parametroscontabilidad.Phinicial)
			ultimo--
		}
	}
	log.Println("Ultimo numero")



	var fechaString string
	fechaString=fechaInicial("2021",mes)
	const (
		layoutISO = "2006-01-02"
	)
	fechaDate, _ := time.Parse(layoutISO, fechaString)

	log.Println(mes)
	log.Println(fechaString)
	cuentaInteresCodigo := plandecuentaempresa{}
	cuentaInteres := plandecuentaempresa{}

	db.Get(&cuentaInteresCodigo, "select distinct cuentaintereses as codigo from plandecuentaempresa where interes='SI'")
	db.Get(&cuentaInteres, "select * from plandecuentaempresa where codigo=$1",cuentaInteresCodigo.Codigo)

	var consultacuentaAnterior=" SELECT distinct cuenta FROM comprobantedetalle "
	consultacuentaAnterior+=   " inner join plandecuentaempresa on "
	consultacuentaAnterior+=   " plandecuentaempresa.codigo=comprobantedetalle.cuenta "
	consultacuentaAnterior+=   " where "
	consultacuentaAnterior +=  " plandecuentaempresa.cuota='NO' and "
	consultacuentaAnterior+=   " comprobantedetalle.tercero=$1 and "
	consultacuentaAnterior+=   " comprobantedetalle.fecha<=$2 and "
	consultacuentaAnterior +=  " substring(comprobantedetalle.cuenta,1,2)='13' and comprobantedetalle.cuenta<>'"+cuentaInteres.Codigo+"' group by cuenta,tercero"

	var consultasaldoAnterior="SELECT debito,credito FROM comprobantedetalle  where"
	consultasaldoAnterior += " tercero=$1 and fecha<=$2 and cuenta=$3"


	//listadoGenerar := []cuentadecobrodetalleGenerar{}
	listadoCuentaTercero := []cuentaTercero{}
	//res := []tercero{}
	listaTercero := []tercero{}

	db.Select(&listaTercero, "SELECT * FROM tercero where not (cuotap='0' and cuota1='0' and cuota2='0')")

	listadoSaldoTercero := []cuentaSaldo{}

	//FacturaInicial=""
	for i, miTercero := range listaTercero {
		//var miUltimo int
		//miUltimo=ultimo+1
		var numeroFactura string
		numeroFactura=strconv.Itoa(ultimo+i+1)
		var miFila int
		miFila=0
		var miFilaComprobante int
		miFilaComprobante=0

		miCuentaTercerodetalle :=[] cuentadecobrodetalle{}
		miComprobanteDetalleDebito :=[] comprobantedetalle{}
		miComprobanteDetalleCredito :=[] comprobantedetalle{}
		miComprobanteDetalle :=[] comprobantedetalle{}


		var totalCuotaActual float64
		var totalCuotaAnterior float64
		var totalCuota float64
		var totalBaseInteres float64

		totalCuotaActual=0
		totalCuotaAnterior=0
		totalCuota=0
		totalBaseInteres=0
		var totalDebito float64
		var totalCredito float64

		totalDebito=0
		totalCredito=0

		// sumar saldos anteriores
		log.Println("Tercero"+miTercero.Nombre)
		db.Select(&listadoCuentaTercero, consultacuentaAnterior,miTercero.Codigo,fechaString)
		for _, miCuentaSaldo := range listadoCuentaTercero {
			//var b=j
			db.Select(&listadoSaldoTercero, consultasaldoAnterior,miTercero.Codigo,fechaString,miCuentaSaldo.Cuenta)
			log.Println("Cuenta"+miCuentaSaldo.Cuenta)
			var totalSaldo float64
			totalSaldo=0
			for _, miSaldo := range listadoSaldoTercero {
				//var c=k
				log.Println(miSaldo.Debito)
				log.Println(miSaldo.Credito)
				totalSaldo=Flotante(miSaldo.Debito)-Flotante(miSaldo.Credito)
			}
			if totalSaldo==0{
			}else {
				// nueva fila anterior
				miFila++
				miCuentaTercerodetalle=append(miCuentaTercerodetalle,
					cuentadecobrodetalle{strconv.Itoa(miFila),
						numeroFactura,miCuentaSaldo.Cuenta,FormatoFlotanteEntero(totalSaldo),"0",FormatoFlotanteEntero(totalSaldo)})
				totalCuotaAnterior+=totalSaldo
				totalBaseInteres+=totalSaldo
				totalCuota+=totalSaldo
			}

			log.Println("Total Cuenta Anterior "+miCuentaSaldo.Cuenta)
			log.Println(FormatoFlotanteEntero(totalSaldo))
		}
		log.Println("Cuentap")

		// saldo actual CUENTAp
		CuentaP := plandecuentaempresa{}
		db.Get(&CuentaP, "SELECT * FROM plandecuentaempresa where cuota='SI' and tipo='CuotaP'")
		var totalSaldoAnterior float64
		totalSaldoAnterior=saldoAnterior(miTercero.Codigo,fechaDate,CuentaP.Codigo)
		if (miTercero.Cuotap=="" ||miTercero.Cuotap=="0") && totalSaldoAnterior==0{
		} else {
			miFila++
			var cuotapActual float64
			cuotapActual=Flotante(miTercero.Cuotap)
			miCuentaTercerodetalle=append(miCuentaTercerodetalle,
				cuentadecobrodetalle{strconv.Itoa(miFila),
					numeroFactura,CuentaP.Codigo,
					FormatoFlotanteEntero(totalSaldoAnterior),
					FormatoFlotanteEntero(Flotante(miTercero.Cuotap)),
					FormatoFlotanteEntero(totalSaldoAnterior+cuotapActual)})

			// sumatorias Interes
			if CuentaP.Interes=="SI"{
				totalBaseInteres+=totalSaldoAnterior
			}

			// inserta fila cuenta
			miFilaComprobante++;
			miComprobanteDetalleDebito=append(miComprobanteDetalleDebito,
			comprobantedetalle{strconv.Itoa(miFilaComprobante),
			CuentaP.Codigo,
			miTercero.Codigo,
			miCentro,
			strings.TrimSpace(CuentaP.Nombre)+" "+mesLetras(mes),
			"",
			FormatoFlotante(cuotapActual)	,
			"",
			Documentocontable,
				numeroFactura,
			fechaDate,
			fechaDate,"",""})
			// Inserta Fila contra
			miFilaComprobante++;
			miComprobanteDetalleCredito=append(miComprobanteDetalleCredito,
				comprobantedetalle{strconv.Itoa(miFilaComprobante),
					CuentaP.Contra,
					miTercero.Codigo,
					miCentro,
					strings.TrimSpace(CuentaP.Nombre)+" "+mesLetras(mes),
					"",
					"",
					FormatoFlotante(cuotapActual),
					Documentocontable,
					numeroFactura,
					fechaDate,
					fechaDate,"",""})

			totalDebito+=cuotapActual;
			totalCredito+=cuotapActual;


			totalCuotaAnterior+=totalSaldoAnterior
			totalCuotaActual+=cuotapActual
			totalCuota+=totalSaldoAnterior+cuotapActual
		}

		// cuenta 1
		Cuenta1 := plandecuentaempresa{}
		db.Get(&Cuenta1, "SELECT * FROM plandecuentaempresa where cuota='SI' and tipo='Cuota1'")
		totalSaldoAnterior=saldoAnterior(miTercero.Codigo,fechaDate,Cuenta1.Codigo)
		if (miTercero.Cuota1=="" ||miTercero.Cuota1=="0") && totalSaldoAnterior==0{
		} else {
			miFila++
			var cuota1Actual float64
			cuota1Actual=Flotante(miTercero.Cuota1)
			miCuentaTercerodetalle=append(miCuentaTercerodetalle,
				cuentadecobrodetalle{strconv.Itoa(miFila),
					numeroFactura,Cuenta1.Codigo,FormatoFlotanteEntero(totalSaldoAnterior),
					FormatoFlotanteEntero(Flotante(miTercero.Cuota1)),
					FormatoFlotanteEntero(totalSaldoAnterior+cuota1Actual)})

			// sumatorias Interes
			if Cuenta1.Interes=="SI"{
				totalBaseInteres+=totalSaldoAnterior
			}
			// inserta fila cuenta1
			miFilaComprobante++;
			miComprobanteDetalleDebito=append(miComprobanteDetalleDebito,
				comprobantedetalle{strconv.Itoa(miFilaComprobante),
					Cuenta1.Codigo,
					miTercero.Codigo,
					miCentro,
					strings.TrimSpace(Cuenta1.Nombre)+" "+mesLetras(mes),
					"",
					FormatoFlotante(cuota1Actual)	,
					"",
					Documentocontable,
					numeroFactura,
					fechaDate,
					fechaDate,"",""})
			// Inserta Fila contra
			miFilaComprobante++;
			miComprobanteDetalleCredito=append(miComprobanteDetalleCredito,
				comprobantedetalle{strconv.Itoa(miFilaComprobante),
					Cuenta1.Contra,
					miTercero.Codigo,
					miCentro,
					strings.TrimSpace(Cuenta1.Nombre)+" "+mesLetras(mes),
					"",
					"",
					FormatoFlotante(cuota1Actual),
					Documentocontable,
					numeroFactura,
					fechaDate,
					fechaDate,"",""})

			totalDebito+=cuota1Actual;
			totalCredito+=cuota1Actual;

			// sumatorias
			totalCuotaAnterior+=totalSaldoAnterior
			totalCuotaActual+=cuota1Actual
			totalCuota+=totalSaldoAnterior+cuota1Actual
		}


		//cuenta 2
		Cuenta2 := plandecuentaempresa{}
		db.Get(&Cuenta2, "SELECT * FROM plandecuentaempresa where cuota='SI' and tipo='Cuota2'")
		totalSaldoAnterior=saldoAnterior(miTercero.Codigo,fechaDate,Cuenta2.Codigo)
		if (miTercero.Cuota2=="" ||miTercero.Cuota2=="0") && totalSaldoAnterior==0{
		} else {
			miFila++
			var cuota2Actual float64
			cuota2Actual=Flotante(miTercero.Cuota2)
			miCuentaTercerodetalle=append(miCuentaTercerodetalle,
				cuentadecobrodetalle{strconv.Itoa(miFila),
					numeroFactura,Cuenta2.Codigo,FormatoFlotanteEntero(totalSaldoAnterior),FormatoFlotanteEntero(Flotante(miTercero.Cuota2)),FormatoFlotanteEntero(totalSaldoAnterior+cuota2Actual)})

			// sumatorias Interes
			if Cuenta2.Interes=="SI"{
				totalBaseInteres+=totalSaldoAnterior
			}


			// inserta fila cuenta1
			miFilaComprobante++;
			miComprobanteDetalleDebito=append(miComprobanteDetalleDebito,
				comprobantedetalle{strconv.Itoa(miFilaComprobante),
					Cuenta2.Codigo,
					miTercero.Codigo,
					miCentro,
					strings.TrimSpace(Cuenta2.Nombre)+" "+mesLetras(mes),
					"",
					FormatoFlotante(cuota2Actual)	,
					"",
					Documentocontable,
					numeroFactura,
					fechaDate,
					fechaDate,"",""})
			// Inserta Fila contra
			miFilaComprobante++;
			miComprobanteDetalleCredito=append(miComprobanteDetalleCredito,
				comprobantedetalle{strconv.Itoa(miFilaComprobante),
					Cuenta2.Contra,
					miTercero.Codigo,
					miCentro,
					strings.TrimSpace(Cuenta2.Nombre)+" "+mesLetras(mes),
					"",
					"",
					FormatoFlotante(cuota2Actual),
					Documentocontable,
					numeroFactura,
					fechaDate,
					fechaDate,"",""})

			totalDebito+=cuota2Actual;
			totalCredito+=cuota2Actual;
			// sumatorias
			totalCuotaAnterior+=totalSaldoAnterior
			totalCuotaActual+=cuota2Actual
			totalCuota+=totalSaldoAnterior+cuota2Actual
		}
		// cuenta 3

		//Cuenta3 := plandecuentaempresa{}
		//db.Get(&Cuenta3, "SELECT * FROM plandecuentaempresa where cuota='SI' and tipo='Cuota3'")
		//totalSaldoAnterior=saldoAnterior(miTercero.Codigo,fechaDate,Cuenta3.Codigo)
		//if (miTercero.Cuota3=="" ||miTercero.Cuota3=="0") && totalSaldoAnterior==0{
		//} else {
		//	miFila++
		//	var cuota3Actual float64
		//	cuota3Actual=Flotante(miTercero.Cuota3)
		//	miCuentaTercerodetalle=append(miCuentaTercerodetalle,
		//		cuentadecobrodetalle{strconv.Itoa(miFila),
		//			numeroFactura,Cuenta3.Codigo,FormatoFlotanteEntero(totalSaldoAnterior),FormatoFlotanteEntero(Flotante(miTercero.Cuota3)),FormatoFlotanteEntero(totalSaldoAnterior+cuota3Actual)})
		//	// sumatorias
		//
		//	// sumatorias Interes
		//	if Cuenta3.Interes=="SI"{
		//		totalBaseInteres+=totalSaldoAnterior
		//	}
		//
		//
		//	// inserta fila cuenta1
		//	miFilaComprobante++;
		//	miComprobanteDetalleDebito=append(miComprobanteDetalleDebito,
		//		comprobantedetalle{strconv.Itoa(miFilaComprobante),
		//			Cuenta3.Codigo,
		//			miTercero.Codigo,
		//			miCentro,
		//			strings.TrimSpace(Cuenta3.Nombre)+" "+mesLetras(mes),
		//			"",
		//			FormatoFlotante(cuota3Actual)	,
		//			"",
		//			Documentocontable,
		//			numeroFactura,
		//			fechaDate,
		//			fechaDate})
		//	// Inserta Fila contra
		//	miFilaComprobante++;
		//	miComprobanteDetalleCredito=append(miComprobanteDetalleCredito,
		//		comprobantedetalle{strconv.Itoa(miFilaComprobante),
		//			Cuenta3.Contra,
		//			miTercero.Codigo,
		//			miCentro,
		//			strings.TrimSpace(Cuenta3.Nombre)+" "+mesLetras(mes),
		//			"",
		//			"",
		//			FormatoFlotante(cuota3Actual),
		//			Documentocontable,
		//			numeroFactura,
		//			fechaDate,
		//			fechaDate})
		//
		//	totalDebito+=cuota3Actual;
		//	totalCredito+=cuota3Actual;
		//
		//	totalCuotaAnterior+=totalSaldoAnterior
		//	totalCuotaActual+=cuota3Actual
		//	totalCuota+=totalSaldoAnterior+cuota3Actual
		//}


		// cuenta interes

		totalSaldoAnterior=0
			//saldoAnterior(miTercero.Codigo,fechaFinal,cuentaInteres.Codigo)


		totalSaldoAnterior=0
		totalSaldoAnterior=saldoAnterior(miTercero.Codigo,fechaDate,cuentaInteres.Codigo)

		log.Println("saldo interes anterior ")
		log.Println(FormatoFlotante(totalSaldoAnterior))

		log.Println("base interes ")
		log.Println(FormatoFlotante(totalBaseInteres))

		if totalBaseInteres==0 && totalSaldoAnterior==0{
		} else {
			miFila++
			var cuotaInteresActual float64
			cuotaInteresActual=(totalBaseInteres)*(miPorcentajeNumero/100)

			miCuentaTercerodetalle=append(miCuentaTercerodetalle,
				cuentadecobrodetalle{strconv.Itoa(miFila),
					numeroFactura,cuentaInteres.Codigo,FormatoFlotanteEntero(totalSaldoAnterior),FormatoFlotanteEntero(cuotaInteresActual),FormatoFlotanteEntero(totalSaldoAnterior+cuotaInteresActual)})
			// sumatorias
			if cuotaInteresActual==0{

			}else {


			// inserta fila cuenta1
			miFilaComprobante++;
			miComprobanteDetalleDebito=append(miComprobanteDetalleDebito,
				comprobantedetalle{strconv.Itoa(miFilaComprobante),
					cuentaInteres.Codigo,
					miTercero.Codigo,
					miCentro,
					strings.TrimSpace(cuentaInteres.Nombre)+" "+mesLetras(mes),
					"",
					FormatoFlotante(cuotaInteresActual)	,
					"",
					Documentocontable,
					numeroFactura,
					fechaDate,
					fechaDate,"",""})
			// Inserta Fila contra
			miFilaComprobante++;
			miComprobanteDetalleCredito=append(miComprobanteDetalleCredito,
				comprobantedetalle{strconv.Itoa(miFilaComprobante),
					cuentaInteres.Contra,
					miTercero.Codigo,
					miCentro,
					strings.TrimSpace(cuentaInteres.Nombre)+" "+mesLetras(mes),
					"",
					"",
					FormatoFlotante(cuotaInteresActual),
					Documentocontable,
					numeroFactura,
					fechaDate,
					fechaDate,"",""})
			totalDebito+=cuotaInteresActual;
			totalCredito+=cuotaInteresActual;
			}



			totalCuotaAnterior+=totalSaldoAnterior
			totalCuotaActual+=cuotaInteresActual
			totalCuota+=totalSaldoAnterior+cuotaInteresActual
		}



		// genera ceunta de cobro
		CuentadecobroNuevaGenerar(cuentadecobro{numeroFactura,
			fechaDate,miCentro,
			miTercero.Codigo,
			FormatoFlotanteEntero(totalCuotaAnterior),
			FormatoFlotanteEntero(totalCuotaActual),
			FormatoFlotanteEntero(totalCuota),
			"Nueva",
			miCuentaTercerodetalle,
			nil})


		// agrega lineas debito
		var filavan=1
		for i, midetalle := range miComprobanteDetalleDebito {
			filavan=i+1
			midetalle.Fila=strconv.Itoa(filavan)
			miComprobanteDetalle= append(miComprobanteDetalle, midetalle)
			log.Println("fila"+midetalle.Fila)
			log.Println(midetalle.Cuenta)
			log.Println(midetalle.Debito)
		}

		filavan++
		for i, midetalle := range miComprobanteDetalleCredito {
			midetalle.Fila=strconv.Itoa(i+filavan)
			miComprobanteDetalle= append(miComprobanteDetalle, midetalle)
			log.Println("fila"+midetalle.Fila)
			log.Println(midetalle.Cuenta)
			log.Println(midetalle.Debito)
		}


		// crea comprobante
		ComprobanteAgregarGenerar(comprobante{Documentocontable,
			numeroFactura,fechaDate,
			fechaDate,
			"2021",
			"",
			"",
			"",
		FormatoFlotante(totalDebito),
		FormatoFlotante(totalCredito),
		"Actualizar",
	miComprobanteDetalle,nil})

		// fin tercero
	}

	var consulta string

	consulta="select Numero,Tercero,tercero.nombre as TerceroNombre,"
	consulta+=" Totalanterior,Totalactual,Total from cuentadecobro "
	consulta+=" inner join tercero on tercero.codigo=cuentadecobro.tercero"
	consulta+=" where EXTRACT(MONTH FROM  cuentadecobro.fecha)=$1"
	log.Println(consulta)
	log.Println(mes)
	listacuentadecobro := []cuentadecobrodetalleGenerar{}

	err1:=db.Select(&listacuentadecobro, consulta,mes)
	if err != nil {
		fmt.Println(err1)
		return
	}







	//if simueve == false {
	//	var slice []string
	//	slice = make([]string, 0)
	//	data, _ := json.Marshal(slice)
	//	w.WriteHeader(200)
	//	w.Header().Set("Content-Type", "application/json")
	//	w.Write(data)
	//} else {
		data, _ := json.Marshal(listacuentadecobro)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	//}
}
// TERMINA CUENTADECOBRODETALLE ESTRUCTURA
func CuentadecobroGenerar(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cuentadecobro/cuentadecobroGenerar.html")

	varmap := map[string]interface{}{
		"hosting":  ruta,
		"centro":ListaCentro(),
	}
	tmp.Execute(w, varmap)
}

// INICIA CUENTADECOBRO LISTA
func CuentadecobroLista(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cuentadecobro/cuentadecobroLista.html")
	log.Println("Error cuentadecobro 0")
	var consulta string

	consulta = "  SELECT cuentadecobro.numero,fecha,tercero,tercero.nombre as Terceronombre,cuentadecobro.total "
	consulta += " FROM cuentadecobro "
	consulta += " inner join tercero on tercero.codigo=cuentadecobro.tercero "
	consulta += " ORDER BY cuentadecobro.numero ASC"

	db := dbConn()
	res := []cuentadecobroLista{}
	//db.Select(&res, consulta)

	//error1 = db.Select(&res, consulta)
	err := db.Select(&res, consulta)

	if err != nil {
		fmt.Println(err)
		return
	}

	varmap := map[string]interface{}{
		"res":     res,
		"hosting": ruta,
	}
	log.Println("Error cuentadecobro888")
	tmp.Execute(w, varmap)
}

// TERMINA CUENTADECOBRO LISTA

// INICIA CUENTADECOBRO NUEVO
func CuentadecobroNuevo(w http.ResponseWriter, r *http.Request) {
	log.Println("Error cuentadecobro nuevo 1")
	log.Println("Error cuentadecobro nuevo 2")
	parametros := map[string]interface{}{
		"hosting":     ruta,
		"centro":ListaCentro(),
	}

	t, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cuentadecobro/cuentadecobroNuevo.html",
		"vista/cuentadecobro/cuentadecobroScript.html",
		"vista/cuentadecobro/autocompletaplandecuentaempresa.html",
		"vista/cuentadecobro/autocompletatercero.html")
	fmt.Printf("%v, %v", t, err)
	log.Println("Error cuentadecobro nuevo 3")
	t.Execute(w, parametros)
}


// proceso que crea la cuenta de cobro desde objeto
func CuentadecobroNuevaGenerar(tempCuentadecobro cuentadecobro) {
	db := dbConn()
	//var tempCuentadecobro cuentadecobro
	//
	//b, err := ioutil.ReadAll(r.Body)
	//defer r.Body.Close()
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//
	//// carga informacion de la CUENTADECOBRO
	//err = json.Unmarshal(b, &tempCuentadecobro)
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}

	if tempCuentadecobro.Accion == "Actualizar" {
		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from cuentadecobrodetalle WHERE numero=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempCuentadecobro.Numero)

		// borra cabecera anterior

		delForm1, err := db.Prepare("DELETE from cuentadecobro WHERE numero=$1 ")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempCuentadecobro.Numero)
	}

	// INSERTA DETALLE
	for i, x := range tempCuentadecobro.Detalle {
		var a = i
		var q string

		q = "insert into cuentadecobrodetalle ("
		q += "Fila,"
		q += "Numero,"
		q += "Cuenta,"
		q += "Anterior,"
		q += "Actual,"
		q += "Total"
		q += " ) values("
		q += parametros(6)
		q += ")"

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA CUENTADECOBRO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Fila,
			x.Numero,
			x.Cuenta,
			x.Anterior,
			x.Actual,
			x.Total)
		if err != nil {
			panic(err)
		}


		// crea detalle

		log.Println("Insertar Detalle \n", x.Cuenta, a)
	}

	log.Println("Got %s age %s club %s\n", tempCuentadecobro.Numero)
	var q string
	q += "insert into cuentadecobro ("
	q += "Numero,"
	q += "Fecha,"
	q += "Centro,"
	q += "Tercero,"
	q += "Totalanterior,"
	q += "Totalactual,"
	q += "Total"
	q += " ) values("
	q += parametros(7)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}
	layout := "2006-01-02"

	log.Println("Hora", tempCuentadecobro.Fecha.Format("02/01/2006"))

	// TERMINA CUENTADECOBRO GRABAR INSERTAR
	_, err = insForm.Exec(
		tempCuentadecobro.Numero,
		tempCuentadecobro.Fecha.Format(layout),
		tempCuentadecobro.Centro,
		tempCuentadecobro.Tercero,
		tempCuentadecobro.Totalanterior,
		tempCuentadecobro.Totalactual,
		tempCuentadecobro.Total)
	if err != nil {
		panic(err)
	}


	//http.Redirect(w, r, "/CUENTADECOBROLista", 301)
}
// INICIA CUENTADECOBRO INSERTAR AJAX
func CuentadecobroAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var tempCuentadecobro cuentadecobro

	var Documentocontable string
	Documentocontable = "30"

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la CUENTADECOBRO
	err = json.Unmarshal(b, &tempCuentadecobro)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if tempCuentadecobro.Accion == "Actualizar" {
		// borra detalle anterior
		delForm, err := db.Prepare("DELETE from cuentadecobrodetalle WHERE numero=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(tempCuentadecobro.Numero)

		// borra cabecera anterior

		delForm1, err := db.Prepare("DELETE from cuentadecobro WHERE numero=$1 ")
		if err != nil {
			panic(err.Error())
		}
		delForm1.Exec(tempCuentadecobro.Numero)

	}
	var miFilaComprobante=0;
	miComprobanteDetalleDebito :=[] comprobantedetalle{}
	miComprobanteDetalleCredito :=[] comprobantedetalle{}
	miComprobanteDetalle :=[] comprobantedetalle{}
	// INSERTA DETALLE
	var totalDebito float64;
	var totalCredito float64;
	totalDebito=0
	totalCredito=0

	for i, x := range tempCuentadecobro.Detalle {
		var a = i
		var q string

		q = "insert into cuentadecobrodetalle ("
		q += "Fila,"
		q += "Numero,"
		q += "Cuenta,"
		q += "Anterior,"
		q += "Actual,"
		q += "Total"
		q += " ) values("
		q += parametros(6)
		q += ")"

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA CUENTADECOBRO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Fila,
			x.Numero,
			x.Cuenta,
			x.Anterior,
			x.Actual,
			x.Total)
		if err != nil {
			panic(err)
		}

		CuentaP := plandecuentaempresa{}
		db.Get(&CuentaP, "SELECT * FROM plandecuentaempresa where codigo=$1",x.Cuenta)

		var cuotaActual=Flotante(x.Actual);

		log.Println("cuenta debito " + x.Cuenta)
		miFilaComprobante++;
		miComprobanteDetalleDebito=append(miComprobanteDetalleDebito,
			comprobantedetalle{strconv.Itoa(miFilaComprobante),
				CuentaP.Codigo,
				tempCuentadecobro.Tercero,
				tempCuentadecobro.Centro,
				strings.TrimSpace(CuentaP.Nombre)+" "+mesLetras(strconv.Itoa(int(tempCuentadecobro.Fecha.Month()))),
				"",
				FormatoFlotante(cuotaActual)	,
				"",
				Documentocontable,
				tempCuentadecobro.Numero,
				tempCuentadecobro.Fecha,
				tempCuentadecobro.Fecha,"",""})
		// Inserta Fila contra


		miFilaComprobante++;
		miComprobanteDetalleCredito=append(miComprobanteDetalleCredito,
			comprobantedetalle{strconv.Itoa(miFilaComprobante),
				CuentaP.Contra,
				tempCuentadecobro.Tercero,
				tempCuentadecobro.Centro,
				strings.TrimSpace(CuentaP.Nombre)+" "+mesLetras(strconv.Itoa(int(tempCuentadecobro.Fecha.Month()))),
				"",
				"",
				FormatoFlotante(cuotaActual),
				Documentocontable,
				tempCuentadecobro.Numero,
				tempCuentadecobro.Fecha,
				tempCuentadecobro.Fecha,"",""})

		totalDebito+=cuotaActual;
		totalCredito+=cuotaActual;

		log.Println("Insertar Detalle \n", x.Cuenta, a)
	}


	// agrega lineas debito
	var filavan=1
	for i, midetalle := range miComprobanteDetalleDebito {
		filavan=i+1
		midetalle.Fila=strconv.Itoa(filavan)
		miComprobanteDetalle= append(miComprobanteDetalle, midetalle)
		log.Println("fila"+midetalle.Fila)
		log.Println(midetalle.Cuenta)
		log.Println(midetalle.Debito)
	}

	filavan++
	for i, midetalle := range miComprobanteDetalleCredito {
		midetalle.Fila=strconv.Itoa(i+filavan)
		miComprobanteDetalle= append(miComprobanteDetalle, midetalle)
		log.Println("fila"+midetalle.Fila)
		log.Println(midetalle.Cuenta)
		log.Println(midetalle.Debito)
	}

	// crea comprobante
	ComprobanteAgregarGenerar(comprobante{Documentocontable,
		tempCuentadecobro.Numero,tempCuentadecobro.Fecha,
		tempCuentadecobro.Fecha,
		"2021",
		"",
		"",
		"",
		FormatoFlotante(totalDebito),
		FormatoFlotante(totalCredito),
		"Actualizar",
		miComprobanteDetalle,
		nil})

	log.Println("Got %s age %s club %s\n", tempCuentadecobro.Numero)
	var q string
	q += "insert into cuentadecobro ("
	q += "Numero,"
	q += "Fecha,"
	q += "Centro,"
	q += "Tercero,"
	q += "Totalanterior,"
	q += "Totalactual,"
	q += "Total"
	q += " ) values("
	q += parametros(7)
	q += " ) "

	log.Println("Cadena SQL " + q)
	insForm, err := db.Prepare(q)
	if err != nil {
		panic(err.Error())
	}
	layout := "2006-01-02"

	log.Println("Hora", tempCuentadecobro.Fecha.Format("02/01/2006"))

	// TERMINA CUENTADECOBRO GRABAR INSERTAR
	_, err = insForm.Exec(
		tempCuentadecobro.Numero,
		tempCuentadecobro.Fecha.Format(layout),
		tempCuentadecobro.Centro,
		tempCuentadecobro.Tercero,
		tempCuentadecobro.Totalanterior,
		tempCuentadecobro.Totalactual,
		tempCuentadecobro.Total)
	if err != nil {
		panic(err)
	}

	var resultado bool
	resultado = true

	js, err := json.Marshal(SomeStruct{resultado})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	//http.Redirect(w, r, "/CUENTADECOBROLista", 301)
}

func CuentadecobroDatoAgregar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	//var tempCuentadecobro cuentadecobro
	listacuentadecobroDato := []cuentadecobroDato{}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// carga informacion de la CUENTADECOBRO
	err = json.Unmarshal(b, &listacuentadecobroDato)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}



	// INSERTA DETALLE
	for _, x := range listacuentadecobroDato {

		var q string

		q = "update tercero set "
		q += "descuento1 = $2 ,"
		q += "descuento2 = $3, "
		q += "cuotap = $4, "
		q += "cuota1 = $5, "
		q += "cuota2 = $6 "
		q += "where codigo = $1 "

		log.Println("Cadena SQL " + q)
		insForm, err := db.Prepare(q)
		if err != nil {
			panic(err.Error())
		}

		// TERMINA CUENTADECOBRO GRABAR INSERTAR
		_, err = insForm.Exec(
			x.Tercero,
			x.Descuento1,
			x.Descuento2,
			x.Cuotap,
			x.Cuota1,
			x.Cuota2)

		if err != nil {
			panic(err)
		}

		log.Println("Insertar Detalle \n"+ x.Tercero)
	}

	var resultado bool
	resultado = true

	js, err := json.Marshal(SomeStruct{resultado})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	//http.Redirect(w, r, "/CUENTADECOBROLista", 301)
}
// INICIA CUENTADECOBRO EXISTE
func CuentadecobroExiste(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Numero:= mux.Vars(r)["numero"]
	var total int
	row := db.QueryRow("SELECT COUNT(*) FROM cuentadecobro  WHERE  Numero=$1",  Numero)
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

// INICIA CUENTADECOBRO EDITAR
func CuentadecobroDato(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var consulta = ""
	consulta = "select "
	consulta += "tercero.Codigo as Tercero, "
	consulta += "tercero.Nombre as Nombre, "
	consulta += "tercero.Descuento1 as Descuento1, "
	consulta += "tercero.Descuento2 as Descuento2, "
	consulta += "tercero.CuotaP as CuotaP, "
	consulta += "tercero.Cuota1 as Cuota1, "
	consulta += "tercero.Cuota2 as Cuota2 "
	consulta += "from tercero "
	consulta += " order by cast(codigo as integer) "

	log.Println("Cadena SQL " + consulta)

	// traer detalle
	det := []cuentadecobroDato{}
	err2 := db.Select(&det, consulta)
	if err2 != nil {
		fmt.Println(err2)
	}


	//	log.Println("detalle cuentadecobro)
	parametros := map[string]interface{}{
		"cuentadecobroDato":       det,
		"hosting":     ruta,
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cuentadecobro/cuentadecobroDato.html",
		"vista/cuentadecobro/cuentadecobroDatoScript.html",
		)

	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error cuentadecobro nuevo 3")
	miTemplate.Execute(w, parametros)

	//tmp.Execute(w, parametros)
}
// INICIA CUENTADECOBRO EDITAR
func CuentadecobroEditar(w http.ResponseWriter, r *http.Request) {

	Numero:= mux.Vars(r)["numero"]
	//log.Println("inicio cuentadecobro editar" + Documento)
	db := dbConn()

	// traer cuentadecobro
	v := cuentadecobro{}
	err := db.Get(&v, "SELECT * FROM cuentadecobro WHERE  Numero=$1",  Numero)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	det := []cuentadecobrodetalleeditar{}
	err2 := db.Select(&det, CuentadecobroConsultaDetalle(), Numero)
	if err2 != nil {
		fmt.Println(err2)
	}
	// traer tercero
	t := tercero{}
	err1 := db.Get(&t, "SELECT * FROM tercero where codigo=$1", v.Tercero)
	if err1 != nil {
		log.Fatalln(err1)
	}

	//	log.Println("detalle cuentadecobro)
	parametros := map[string]interface{}{
		"cuentadecobro":       v,
		"detalle":     det,
		"hosting":     ruta,
		"tercero":t,
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cuentadecobro/cuentadecobroEditar.html",
		"vista/cuentadecobro/cuentadecobroScript.html",
		"vista/cuentadecobro/autocompletaplandecuentaempresa.html",
		"vista/cuentadecobro/autocompletatercero.html")

	fmt.Printf("%v, %v", miTemplate, err)
	log.Println("Error cuentadecobro nuevo 3")
	miTemplate.Execute(w, parametros)

	//tmp.Execute(w, parametros)
}
// INICIA CUENTADECOBRO BORRAR
func CuentadecobroBorrar(w http.ResponseWriter, r *http.Request) {
	Numero:= mux.Vars(r)["numero"]
	log.Println("inicio cuentadecobro editar" + Numero)
	db := dbConn()

	// traer CUENTADECOBRO
	v := cuentadecobro{}
	err := db.Get(&v, "SELECT * FROM cuentadecobro WHERE Numero=$1", Numero)
	if err != nil {
		log.Fatalln(err)
	}

	// traer detalle
	det := []cuentadecobrodetalleeditar{}
	err2 := db.Select(&det, CuentadecobroConsultaDetalle(), Numero)
	if err2 != nil {
		fmt.Println(err2)
	}


	// traer tercero
	t := tercero{}
	err1 := db.Get(&t, "SELECT * FROM tercero where codigo=$1", v.Tercero)
	if err1 != nil {
		log.Fatalln(err1)
	}
	//	log.Println("detalle producto" + det.Producto+det.ProductoNombre)
	parametros := map[string]interface{}{
		"cuentadecobro":       v,
		"detalle":     det,
		"hosting":     ruta,
		"tercero":t,
		"centro":ListaCentro(),
	}

	miTemplate, err := template.ParseFiles("vista/inicio/Modulo.html",
		"vista/cuentadecobro/cuentadecobroBorrar.html",
		"vista/cuentadecobro/cuentadecobroScript.html",
		"vista/cuentadecobro/autocompletaplandecuentaempresa.html",
		"vista/cuentadecobro/autocompletatercero.html",
		"vista/cuentadecobro/autocompletacentro.html")

	log.Println("Error cuentadecobro nuevo 3")
	miTemplate.Execute(w, parametros)

	//tmp.Execute(w, parametros)
}

// INICIA CUENTADECOBRO ELIMINAR
func CuentadecobroEliminar(w http.ResponseWriter, r *http.Request) {
	log.Println("Inicio Eliminar")
	db := dbConn()
	Numero:= mux.Vars(r)["numero"]

	// borrar CUENTADECOBRO
	delForm, err := db.Prepare("DELETE from cuentadecobro WHERE  Numero=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec( Numero)

	// borar detalle
	delForm1, err := db.Prepare("DELETE from cuentadecobrodetalle WHERE  Numero=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm1.Exec( Numero)

	log.Println("Registro Eliminado" + Numero)
	http.Redirect(w, r, "/CuentadecobroLista", 301)
}

// TERMINA CUENTADECOBRO ELIMINAR

// INICIA CUENTADECOBRO PDF
func CuentadecobroPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Numero := mux.Vars(r)["numero"]
	// TRAER CUENTADECOBRO
	miCuentadecobro := cuentadecobro{}
	err := db.Get(&miCuentadecobro, "SELECT * FROM cuentadecobro where numero=$1", Numero)
	if err != nil {
		log.Fatalln(err)
	}

	// TRAER DETALLE
	miDetalle := []cuentadecobrodetalleeditar{}
	err2 := db.Select(&miDetalle, CuentadecobroConsultaDetalle(), Numero)
	if err2 != nil {
		fmt.Println(err2)
	}

	//var e empresa = ListaEmpresa()

	var buf bytes.Buffer
	var err1 error

	pdf := gofpdf.New("P", "mm", "Letter", cnFontDir)

	CuentadecobroHeader(pdf,miCuentadecobro);
	CuentadecobroFooter(pdf);

	// inicia suma de paginas
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.Ln(1)
	CuentadecobroCabecera(pdf,miCuentadecobro,miDetalle)

	var filas=len(miDetalle)
	// menos de 32
	if(filas<=32){
		for i, miFila := range miDetalle {
			var a = i + 1
			CuentadecobroFilaDetalle(pdf,miFila,a)
		}
		CuentadecobroPieDePagina(pdf,miCuentadecobro)
	}	else {
		// mas de 32 y menos de 73
		if(filas>32 && filas<=73){
			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CuentadecobroFilaDetalle(pdf,miFila,a)
				}
			}
			// segunda pagina
			pdf.AddPage()
			CuentadecobroCabecera(pdf,miCuentadecobro,miDetalle)
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41)	{
					CuentadecobroFilaDetalle(pdf,miFila,a)
				}
			}

			CuentadecobroPieDePagina(pdf,miCuentadecobro)
		} else {
			// mas de 80

			// primera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a<=41)	{
					CuentadecobroFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CuentadecobroCabecera(pdf,miCuentadecobro,miDetalle)
			// segunda pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>41 && a<=82)	{
					CuentadecobroFilaDetalle(pdf,miFila,a)
				}
			}

			pdf.AddPage()
			CuentadecobroCabecera(pdf,miCuentadecobro,miDetalle)
			// tercera pagina
			for i, miFila := range miDetalle {
				var a = i + 1
				if (a>82)	{
					CuentadecobroFilaDetalle(pdf,miFila,a)
				}
			}

			CuentadecobroPieDePagina(pdf,miCuentadecobro)
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
// TERMINA CUENTADECOBRO PDF

// INICIA EMPRESA CUENTA DE COBRO PDF
func CuentadecobroHeader(pdf *gofpdf.Fpdf, miCuentadecobro cuentadecobro){
// ENCABEZADO
	var e empresa = ListaEmpresa()
	pdf.SetHeaderFunc(func() {
		// LOGO
		pdf.SetFont("Arial", "", 10)
		pdf.Image(imageFile("logo.png"), 20, 25, 40, 0, false,
		"", 0, "")

		// EMPRESA

		pdf.SetY(20)
		pdf.SetX(20)
		pdf.SetDrawColor(119,134,153)
		pdf.CellFormat(185, 32, "", "1", 0, "C",
		false, 0, "")

		pdf.SetY(20)
		pdf.CellFormat(190, 10, Mayuscula(e.Nombre), "0", 0,
			"C", false, 0, "")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, "Nit No. "+Coma(e.Codigo)+" - "+e.Dv, "0", 0, "C",
			false, 0, "")
		pdf.Ln(4)
		pdf.CellFormat(190, 10, e.Direccion+" "+e.Telefono1,"0", 0, "C", false, 0,
			"")

	// NOMBRE DEL DOCUMENTO
	pdf.SetFont("Arial", "", 11)
	pdf.Ln(5)
	pdf.SetY(20)
	pdf.SetX(80)
	pdf.CellFormat(190, 10, "CUENTA DE COBRO", "0", 0, "C",
		false, 0, "")
	pdf.Ln(6)
	pdf.SetX(80)
	pdf.CellFormat(190, 10, " No.  " + miCuentadecobro.Numero, "0", 0, "C",
		false, 0, "")

})
}
// TERMINA EMPRESA CUENTA DE COBRO PDF

// INICIA CABECERA
func CuentadecobroCabecera(pdf *gofpdf.Fpdf,miCuentadecobro cuentadecobro, miDetalle []cuentadecobrodetalleeditar ){
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miCuentadecobro.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	var mesletras string
	mesletras=mesLetras(strconv.Itoa(int(miCuentadecobro.Fecha.Month())))
	pdf.Ln(13)
	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Nombre", "", 0,
		"L", false, 0, "")
	pdf.SetX(40)
	pdf.CellFormat(40, 4, miTercero.Nombre, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Nit. No.", "", 0,
		"L", false, 0, "")
	pdf.SetX(125)
	pdf.CellFormat(40, 4, Coma(miTercero.Codigo), "", 0,
		"L", false, 0, "")
	pdf.SetX(158)
	pdf.CellFormat(40, 4, "Fecha", "", 0,
		"L", false, 0, "")
	pdf.SetX(175)
	pdf.CellFormat(40, 4, miCuentadecobro.Fecha.Format("02/01/2006"), "", 0,
		"L", false, 0, "")
	pdf.Ln(5)

	pdf.SetX(21)
	pdf.CellFormat(40, 4, "Direccion", "", 0,
		"L", false, 0, "")
	pdf.SetX(40)
	pdf.CellFormat(40, 4, miTercero.Direccion, "", 0,
		"L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(40, 4, "Telefono", "", 0,
		"L", false, 0, "")
	pdf.SetX(130)
	pdf.CellFormat(40, 4, miTercero.Telefono1, "", 0,
		"L", false, 0, "")
	pdf.Ln(-1)

	// CUADRO TITULO
	pdf.SetFont("Arial", "", 9)
	pdf.SetY(52)
	pdf.SetX(20)
	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)
	pdf.CellFormat(185, 5, "", "1", 0, "C",
		true, 0, "")

	pdf.SetX(20)
	pdf.CellFormat(184, 5, "No.", "0", 0,
		"L", false, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(190, 5, "CUENTA", "0", 0,
		"L", false, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(190, 5, "NOMBRE", "0", 0,
		"L", false, 0, "")
	pdf.SetX(126)
	pdf.CellFormat(190, 5, "ANTERIOR", "0", 0,
		"L", false, 0, "")
	pdf.SetX(165)
	pdf.CellFormat(4, 5, Mayuscula(mesletras), "0", 0,
		"C", false, 0, "")
	pdf.SetX(190)
	pdf.CellFormat(190, 5, "TOTAL", "0", 0,
		"L", false, 0, "")
	pdf.Ln(6)
}
// TERMINA CABECERA

// INICIA DETALLE CUENTADECOBRO PDF
func CuentadecobroFilaDetalle(pdf *gofpdf.Fpdf,miFila cuentadecobrodetalleeditar, a int ){
	if math.Mod(float64(a),2)==0 {
		pdf.SetFillColor(224,231,239)
		pdf.SetTextColor(0,0,0)
	} else{
		pdf.SetFillColor(255,255,255)
		pdf.SetTextColor(0,0,0)
	}

	pdf.SetFont("Arial", "", 10)
	pdf.SetX(21)

	pdf.CellFormat(40, 4, strconv.Itoa(a), "", 0,
		"L", true, 0, "")
	pdf.SetX(30)
	pdf.CellFormat(60, 4, miFila.Cuenta, "", 0,
		"L", true, 0, "")
	pdf.SetX(55)
	pdf.CellFormat(60, 4, Subcadena(miFila.Cuentanombre,0,40), "", 0,
		"L", true, 0, "")
	pdf.SetX(115)
	pdf.CellFormat(30, 4, miFila.Anterior, "", 0,
		"R", true, 0, "")
	pdf.SetX(145)
	pdf.CellFormat(30, 4, miFila.Actual, "", 0,
		"R", true, 0, "")
	pdf.SetX(174)
	pdf.CellFormat(30, 4, miFila.Total,"",0,
		"R", true, 0, "")
	pdf.Ln(4)
}
// TERMINA DETALLE CUENTADECOBRO

// INICIA FINAL DE PAGINA
func CuentadecobroPieDePagina(pdf *gofpdf.Fpdf,miCuentadecobro cuentadecobro ){
	miTercero := tercero{}
	err3 := db.Get(&miTercero, "SELECT * FROM tercero where codigo=$1", miCuentadecobro.Tercero)
	if err3 != nil {
		log.Fatalln(err3)
	}

	var totaldescuento1 string
	var totaldescuento2 string

	totaldescuento1 = FormatoFlotanteEntero(Flotante(miCuentadecobro.Total) - Flotante(miTercero.Descuento1))
	totaldescuento2 = FormatoFlotanteEntero(Flotante(miCuentadecobro.Total) - Flotante(miTercero.Descuento2))

	parametroscontabilidad := configuracioncontabilidad{}
	parametroscontabilidad=TraerParametrosContabilidad()

	pdf.SetFont("Arial", "", 9)
	pdf.Ln(10)
	pdf.SetY(114)
	pdf.SetX(20)
	pdf.CellFormat(40, 10, "", "0", 0,
		"C", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(22)
	pdf.CellFormat(40, 10, "Se omite firma autografo Art. 10 D. R. 836/1991 Art. 1.6.1.12.12 DUR 1625 de 2016", "0", 0, "L",
		false, 0, "")

	pdf.SetFont("Arial", "", 10)

	// RELLENO TOTALES
	pdf.SetY(94)
	pdf.SetX(20)
	pdf.SetFillColor(59,99,146)
	pdf.SetDrawColor(119,134,153)
	pdf.SetTextColor(255,255,255)
	pdf.CellFormat(50, 6, "", "0", 0, "L",
		true, 0, "")

	pdf.CellFormat(135, 6, "TOTALES", "0", 0, "L",
		true, 0, "")
	pdf.SetY(92)
	pdf.SetX(105)
	pdf.CellFormat(40, 10, miCuentadecobro.Totalanterior, "0", 0, "R",
		false, 0, "")
	pdf.SetY(92)
	pdf.SetX(135)
	pdf.CellFormat(40, 10, miCuentadecobro.Totalactual, "0", 0, "R",
		false, 0, "")
	pdf.SetY(92)
	pdf.SetX(164)
	pdf.CellFormat(40, 10, miCuentadecobro.Total, "0", 0, "R",
		false, 0, "")

	// CUADRO DETALLE
	pdf.SetTextColor(0,0,0)

	pdf.SetY(57)
	pdf.SetX(20)
	pdf.SetDrawColor(0,82,165)
	pdf.CellFormat(185, 43, "", "1", 0, "C",
		false, 0, "")

	// CUADRO TOTALES
	pdf.SetY(94)
	pdf.SetX(20)
	pdf.SetDrawColor(0,82,165)
	pdf.CellFormat(185, 6, "", "1", 0, "C",
		false, 0, "")

	// CUADRO PIE
	pdf.SetY(100)
	pdf.SetX(145)
	pdf.SetDrawColor(0,82,165)
	pdf.CellFormat(60, 26, "", "1", 0, "C",false, 0, "")

	// CUADRO AVISO
	pdf.SetY(100)
	pdf.SetX(20)
	pdf.SetDrawColor(0,82,165)
	//pdf.SetFillColor(241,196,15)
	pdf.CellFormat(125, 26, "", "1", 0, "C",false, 0, "")

	pdf.SetFont("Arial", "", 8)
	pdf.SetY(101)
	pdf.SetX(22)
	pdf.CellFormat(121, 6, parametroscontabilidad.Textoaviso1, "0", 0, "C",
		false, 0, "")

	pdf.SetY(105)
	pdf.SetX(22)
	pdf.CellFormat(121, 6, parametroscontabilidad.Textoaviso2, "0", 0, "C",
		false, 0, "")

	pdf.SetY(109)
	pdf.SetX(22)
	pdf.CellFormat(121, 6, parametroscontabilidad.Textoaviso3, "0", 0, "C",
		false, 0, "")

	pdf.SetY(113)
	pdf.SetX(22)
	pdf.CellFormat(121, 6, parametroscontabilidad.Textoaviso4, "0", 0, "C",
		false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetY(100)
	pdf.SetX(147)
	pdf.CellFormat(40, 10, parametroscontabilidad.Textodescuento1, "0", 0, "L",
		false, 0, "")
	pdf.SetY(100)
	pdf.SetX(164)
	pdf.CellFormat(40, 10, totaldescuento1, "0", 0, "R",
		false, 0, "")

	pdf.SetY(104)
	pdf.SetX(147)
	pdf.CellFormat(40, 10, parametroscontabilidad.Textodescuento2, "0", 0, "L",
		false, 0, "")
	pdf.SetY(104)
	pdf.SetX(164)
	pdf.CellFormat(40, 10, totaldescuento2, "0", 0, "R",
		false, 0, "")

	pdf.SetY(108)
	pdf.SetX(147)
	pdf.CellFormat(40, 10, parametroscontabilidad.Textodescuento3, "0", 0, "L",
		false, 0, "")

	pdf.SetY(108)
	pdf.SetX(164)
	pdf.CellFormat(40, 10, miCuentadecobro.Total, "0", 0, "R",
		false, 0, "")
}
// TERMINA FINAL DE PAGINA

// INICIA PIE DE PAGINA
func CuentadecobroFooter(pdf *gofpdf.Fpdf){

	pdf.SetFooterFunc(func() {
		pdf.SetY(118)
		pdf.SetX(147)
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(40, 10, "www.Sadconf.com.co", "",
			0, "L", false, 0, "")
		pdf.SetX(177)
		pdf.CellFormat(30, 10, fmt.Sprintf(" %d de {nb}", pdf.PageNo()), "",
			0, "R", false, 0, "")
	})

}
// TERMINA PIE DE PAGINA


