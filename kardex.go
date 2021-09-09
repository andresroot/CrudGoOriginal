package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"math"
	"net/http"

	"time"
)

const operacionInventarioInicial = "1"
const operacionCompra = "2"
const operacionSoporte = "3"
const operacionDevolucionVenta = "4"
const operacionTrasladoEntrada = "5"
const operacionDevolucionCompra = "6"
const operacionDevolucionSoporte = "7"
const operacionVenta = "8"
const operacionTrasladoSalida = "9"


// CENTRO TABLA
type kardex struct {
	Fecha string `json:"Fecha"`
	//time.Time `json:"Fecha"`
	Mes  			string  `json:"Mes"`
	Filas           string  `json:"Filas"`
	Producto        string  `json:"Producto"`
	Tipo            string  `json:"Tipo"`
	Operacion       string  `json:"Operacion"`
	Codigo          string  `json:"Documento"`
	Bodega          string  `json:"Bodega"`
	Cantidadentrada float64 `json:"CantidadE"`
	Precioentrada   float64 `json:"PrecioE"`
	Totalentrada    float64 `json:"TotalE"`
	Cantidadsalida  float64 `json:"Cantidad"`
	Preciosalida    float64 `json:"Precio"`
	Totalsalida     float64 `json:"Total"`
	Cantidadsaldo   float64 `json:"CantidadT"`
	Preciosaldo     float64 `json:"PrecioT"`
	Totalsaldo      float64 `json:"TotalT"`
}


func KardexDatosTodos(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	//CodigoParamerto := mux.Vars(r)["codigo"]


	Discriminar := mux.Vars(r)["discriminar"]
	FechaInicial := mux.Vars(r)["fechainicial"]
	FechaFinal := mux.Vars(r)["fechafinal"]

	dateinicial, err := time.Parse("2006-01-02", FechaInicial)
	//datefinal, err := time.Parse("2006-01-02", fechaFinal)

	if err == nil {
		fmt.Println("Fecha Inicial suma"+dateinicial.String())
	}

	log.Println("fecha Inicial : " + FechaInicial)
	args := []interface{}{}
	BodegaParametro := mux.Vars(r)["bodega"]
	TipoParametro := mux.Vars(r)["tipo"]

	res := []inventario{}
	listadokardex := []kardex{}
	listadokardexfinal := []kardex{}
	listadokardexresumen := []kardex{}

	var cadena string
	var cadenaProducto string

		cadenaProducto = "select codigo,nombre from producto"
		if BodegaParametro == "Todas" {
			cadena = "SELECT fecha,EXTRACT(MONTH FROM  fecha) as mes, tipo, codigo, bodega, producto," +
				" cantidad, precio, operacion FROM inventario where" +
				"  (fecha>=$1  AND fecha <=$2)" +
				" ORDER BY producto,Fecha,operacion "

			cadenaProducto="SELECT  distinct inventario.producto, producto.nombre,producto.iva FROM inventario " +
				" inner join producto on producto.codigo=inventario.producto where" +
				"  (fecha>=$1  AND fecha <=$2)" +
				" group by inventario.producto, producto.nombre,producto.iva order by inventario.producto"

			args = append(args, &FechaInicial)
			args = append(args, &FechaFinal)

		} else {

			cadena = "SELECT fecha,EXTRACT(MONTH FROM  fecha) as mes, tipo, codigo, bodega, producto," +
				" cantidad, precio, operacion FROM inventario where" +
				" and (fecha>=$1  AND fecha <=$2)  and bodega=$3 " +
				" ORDER BY producto,Fecha,operacion "

			cadenaProducto="SELECT  distinct inventario.producto, producto.nombre,producto.iva FROM inventario " +
				" inner join producto on producto.codigo=inventario.producto where " +
				" and (fecha>=$1  AND fecha <=$2)  and bodega=$3 " +
				" group by inventario.producto, producto.nombre,producto.iva order by inventario.producto"
			args = append(args, &FechaInicial)
			args = append(args, &FechaFinal)
			args = append(args, &BodegaParametro)
		}
	var siexisteproducto bool
	selDB1, err := db.Query(cadenaProducto, args...)
	switch err {
	//resltadvaa
	case nil:
		log.Printf("Datos Kardex existe")
		siexisteproducto = true
	case sql.ErrNoRows:
		log.Println("Datos Kardex no encontrados")
	default:
		log.Printf("tercero error: %s\n", err)
	}

	var simueve bool
	simueve=false

	if siexisteproducto==true{
		var codigoProducto string
		var nombreProducto string
		var ivaProduco string
		args=nil

		for selDB1.Next() {
			// recorrer productos
			args=nil
			simueve=true

			err = selDB1.Scan(&codigoProducto, &nombreProducto,&ivaProduco)

			log.Println("Datos Kardex producto"+codigoProducto)
			if err != nil {
				panic(err.Error())
			}

			if BodegaParametro == "Todas" {
				cadena = "SELECT fecha,EXTRACT(MONTH FROM  fecha) as mes, tipo, codigo, bodega, producto," +
					" cantidad, precio, operacion FROM inventario where" +
					" producto=$1 " +
					" and (fecha <=$2)" +
					" ORDER BY producto,Fecha,operacion "
				args = append(args, &codigoProducto)
				//args = append(args, &FechaInicial)
				args = append(args, &FechaFinal)

			} else {

				cadena = "SELECT fecha,EXTRACT(MONTH FROM  fecha) as mes, tipo, codigo, bodega, producto," +
					" cantidad, precio, operacion FROM inventario where" +
					" producto=$1" +
					" and ( fecha <=$2)  and bodega=$3 " +
					" ORDER BY producto,Fecha,operacion "

				args = append(args, &codigoProducto)
				//args = append(args, &FechaInicial)
				args = append(args, &FechaFinal)
				args = append(args, &BodegaParametro)
			}

			selDB, err := db.Query(cadena, args...)
			log.Println("sql : " + cadena)
			var siexiste bool
			siexiste = false
			switch err {
			//resltadvaa
			case nil:
				log.Printf("Datos Kardex existe")
				siexiste = true
			case sql.ErrNoRows:
				log.Println("Datos Kardex no encontrados")
			default:
				log.Printf("tercero error: %s\n", err)
			}

			if siexiste == true {
				var saldo float64
				var costo float64
				//var precosto float64
				var total float64
				saldo = 0
				costo = 0
				total = 0
				var primero bool
				primero = false
				var totalanterior float64
				totalanterior=0
				var Fecha time.Time
				var Tipo string
				var mes string
				var Codigo string
				var Bodega string
				var Producto string
				var Cantidad string
				var Precio string
				var Operacion string
				var Cantidadentrada float64
				var Precioentrada float64
				var Totalentrada float64
				var Cantidadsalida float64
				var Preciosalida float64
				var Totalsalida float64
				var Cantidadsaldo float64
				var Preciosaldo float64
				var Totalsaldo float64
				for selDB.Next() {



					Cantidadentrada = 0
					Precioentrada = 0
					Totalentrada = 0
					Cantidadsalida = 0
					Preciosalida = 0
					Totalsalida = 0
					Cantidadsaldo = 0
					Preciosaldo = 0
					Totalsaldo = 0

					err = selDB.Scan(&Fecha,&mes, &Tipo, &Codigo, &Bodega, &Producto, &Cantidad,
						&Precio, &Operacion)
					if err != nil {
						panic(err.Error())
					}

					// OPERACION DE ASIGNACION AL COSTO O PRECIO DE ENTRADAS

					if Operacion == operacionInventarioInicial || Operacion == operacionCompra || Operacion == operacionSoporte || Operacion == operacionDevolucionVenta || Operacion == operacionTrasladoEntrada{

						Cantidadentrada = Flotante(Cantidad)
						Precioentrada = Flotante(Precio)
						Totalentrada = Flotante(Cantidad) * Flotante(Precio)

						if Operacion == operacionDevolucionVenta ||  Operacion == operacionTrasladoEntrada{
							Cantidadentrada = Flotante(Cantidad)
							Precioentrada = costo
							Totalentrada = Flotante(Cantidad) * costo
						} else {
							Cantidadentrada = Flotante(Cantidad)
							Precioentrada = Flotante(Precio)
							Totalentrada = Flotante(Cantidad) * Flotante(Precio)
						}

					} else {

						if Operacion == operacionDevolucionCompra ||  Operacion == operacionDevolucionSoporte{
							Cantidadsalida = Flotante(Cantidad)
							Preciosalida = Flotante(Precio)
							Totalsalida = Flotante(Cantidad) * Flotante(Precio)
						} else {

							Cantidadsalida = Flotante(Cantidad)
							Preciosalida = costo
							Totalsalida = Flotante(Cantidad) * costo
						}

					}

					if primero == false ||saldo==0{
						saldo = Flotante(Cantidad)
						costo = Flotante(Precio)
						total = Flotante(Cantidad) * Flotante(Precio)
						primero = true
					} else {
						log.Println("cantidad")
						log.Println(Cantidadentrada + saldo)
						log.Println("total")
						log.Println(((Flotante(Cantidad) * Flotante(Precio)) + total))

						// OPERACION DE ASIGNACION CALCULOS DE ENTRADA Y SALIDA
						if Operacion == operacionInventarioInicial || Operacion == operacionCompra || Operacion == operacionSoporte || Operacion == operacionDevolucionVenta || Operacion == operacionTrasladoEntrada{

							if Operacion == operacionDevolucionVenta ||  Operacion == operacionTrasladoEntrada{
								costo = ((Flotante(Cantidad) * costo) + total) / (Cantidadentrada - Cantidadsalida + saldo)
							} else {
								costo = ((Flotante(Cantidad) * Flotante(Precio)) + total) / (Cantidadentrada - Cantidadsalida + saldo)
							}

						} else {
							costo = (total - (Flotante(Cantidad) * (costo))) / (Cantidadentrada - Cantidadsalida + saldo)

							if Operacion == operacionDevolucionCompra ||  Operacion == operacionDevolucionSoporte{

								costo = (total - (Flotante(Cantidad) * Flotante(Precio))) / (Cantidadentrada - Cantidadsalida + saldo)
							} else {

								costo = (total - (Flotante(Cantidad) * (costo))) / (Cantidadentrada - Cantidadsalida + saldo)

							}
						}
						log.Println("costo")
						log.Println(costo)
						saldo = saldo + Cantidadentrada - Cantidadsalida
						total = saldo * costo
						costo = Redondear(costo, 2)

					}
					if saldo == 0{
						//totalanterior=0
						Cantidadsaldo=0
						Preciosaldo=0
						Totalsaldo=0

						if Cantidadentrada>0{
							total=Totalentrada
						}else{
							total=Totalsalida
							Totalsalida=totalanterior
						}
						totalanterior=0
						//total
						Cantidadsaldo = saldo
						Preciosaldo = 0
						//costo
						Totalsaldo = 0

					}else{

						if Cantidadentrada>0{
							total=totalanterior+Totalentrada
						}else{
							total=totalanterior-Totalsalida
						}
						totalanterior=total
						Cantidadsaldo = saldo
						Preciosaldo = costo
						Totalsaldo = total
					}
					//if saldo == 0{
					//	totalanterior=0
					//	Cantidadsaldo=0
					//	Preciosaldo=0
					//	Totalsaldo=0
					//
					//}else{
					//
					//	if Cantidadentrada>0{
					//		total=totalanterior+Totalentrada
					//	}else{
					//		total=totalanterior-Totalsalida
					//	}
					//	totalanterior=total
					//	Cantidadsaldo = saldo
					//	Preciosaldo = costo
					//	Totalsaldo = total
					//}


					// operaciones
					if Cantidadsaldo<0{

						log.Println("Saldo Negativo "+Producto+" "+Operacion+" "+Tipo+" "+Codigo)
					}

					res = append(res, inventario{Fecha, Tipo, Codigo,
						Bodega, Producto, Cantidad, Precio, Operacion})

					listadokardex = append(listadokardex, kardex{Fecha.Format("02/01/06"),
						mes,
						"",
						Producto,
						Tipo,
						Operacion,
						Codigo,
						Bodega,
						Cantidadentrada,
						Precioentrada,
						Totalentrada,
						Cantidadsalida,
						Preciosalida,
						Totalsalida,
						Cantidadsaldo,
						Preciosaldo,
						Totalsaldo})

				}

				listadokardexresumen = append(listadokardexresumen, kardex{Fecha.Format("02/01/06"),
					mes,
					"",
					Producto,
					Tipo,
					Operacion,
					Codigo,
					Bodega,
					Cantidadentrada,
					Precioentrada,
					Totalentrada,
					Cantidadsalida,
					Preciosalida,
					Totalsalida,
					Cantidadsaldo,
					Preciosaldo,
					Totalsaldo})


			}
		}
	}

	if simueve == false {
		var slice []string
		slice = make([]string, 0)
		data, _ := json.Marshal(slice)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {



		if Discriminar=="SI"{
				if TipoParametro=="Todos"{
					listadokardexfinal=listadokardex
				}else{
					var totalEntrada float64
					var totalSalida float64

					totalEntrada=0
					totalSalida=0


					for _, miFila := range listadokardex {
						dateFila, err := time.Parse("02/01/06", miFila.Fecha)
						//datefinal, err := time.Parse("2006-01-02", fechaFinal)

						if err == nil {
							fmt.Println("Fecha Inicial suma"+dateinicial.String())
						}

						// fecha igual o mayor a la inicial
						if dateFila.After(dateinicial) || dateFila==dateinicial	{

							if miFila.Operacion==TipoParametro{
								totalEntrada+=miFila.Totalentrada
								totalSalida+=miFila.Totalsalida
								listadokardexfinal=append(listadokardexfinal,miFila)

							}
						}


					}

					// suma totales
					listadokardexfinal=append(listadokardexfinal, kardex{
						"",
						"",
						"",
						"",
						"Total",
						"",
						"",
						"",
						0,
						0,
						totalEntrada,
						0,
						0,
						totalSalida,
						0,
						0,
						0})
				}
		} else {

			listadokardexfinal=listadokardexresumen
		}

		data, _ := json.Marshal(listadokardexfinal)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func KardexDatos(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	CodigoParamerto := mux.Vars(r)["codigo"]
	FechaInicial := mux.Vars(r)["fechainicial"]
	FechaFinal := mux.Vars(r)["fechafinal"]

	log.Println("fecha Inicial : " + FechaInicial)
	args := []interface{}{}
	BodegaParametro := mux.Vars(r)["bodega"]
	TipoParametro := mux.Vars(r)["tipo"]

	//res := []inventario{}
	listadokardex := []kardex{}
	listadokardexfinal := []kardex{}
	var cadena string


		//cadenaProducto = "select codigo,nombre from producto where producto=$1"
		if BodegaParametro == "Todas" {
			cadena = "SELECT fecha,EXTRACT(MONTH FROM  fecha) as mes, tipo, codigo, bodega, producto," +
				" cantidad, precio, operacion FROM inventario where" +
				" producto=$1 " +
				" and (fecha>=$2  AND fecha <=$3)" +
				" ORDER BY producto,Fecha,operacion "
			args = append(args, &CodigoParamerto)
			args = append(args, &FechaInicial)
			args = append(args, &FechaFinal)

		} else {

			cadena = "SELECT fecha,EXTRACT(MONTH FROM  fecha) as mes, tipo, codigo, bodega, producto," +
				" cantidad, precio, operacion FROM inventario where" +
				" producto=$1" +
				" and (fecha>=$2  AND fecha <=$3)  and bodega=$4 " +
				" ORDER BY producto,Fecha,operacion "

			args = append(args, &CodigoParamerto)
			args = append(args, &FechaInicial)
			args = append(args, &FechaFinal)
			args = append(args, &BodegaParametro)
		}



	selDB, err := db.Query(cadena, args...)
	log.Println("sql : " + cadena)

	var siexiste bool
	siexiste = false

	switch err {
	//resltadvaa
	case nil:
		log.Printf("Datos Kardex existe")
		siexiste = true
	case sql.ErrNoRows:
		log.Println("Datos Kardex no encontrados")
	default:
		log.Printf("tercero error: %s\n", err)
	}

	if siexiste == true {
		var saldo float64
		var costo float64
		//var precosto float64
		var total float64
		saldo = 0
		costo = 0
		total = 0
		var primero bool
		primero = false
		var totalanterior float64
		totalanterior = 0

		for selDB.Next() {
			var Fecha time.Time
			var Tipo string
			var mes string
			var Codigo string
			var Bodega string
			var Producto string
			var Cantidad string
			var Precio string
			var Operacion string
			var Cantidadentrada float64
			var Precioentrada float64
			var Totalentrada float64
			var Cantidadsalida float64
			var Preciosalida float64
			var Totalsalida float64
			var Cantidadsaldo float64
			var Preciosaldo float64
			var Totalsaldo float64

			Cantidadentrada = 0
			Precioentrada = 0
			Totalentrada = 0
			Cantidadsalida = 0
			Preciosalida = 0
			Totalsalida = 0
			Cantidadsaldo = 0
			Preciosaldo = 0
			Totalsaldo = 0

			err = selDB.Scan(&Fecha,&mes, &Tipo, &Codigo, &Bodega, &Producto, &Cantidad,
				&Precio, &Operacion)
			if err != nil {
				panic(err.Error())
			}



			if Operacion == operacionInventarioInicial || Operacion == operacionCompra || Operacion == operacionSoporte || Operacion == operacionDevolucionVenta || Operacion == operacionTrasladoEntrada{

				Cantidadentrada = Flotante(Cantidad)
				Precioentrada = Flotante(Precio)
				Totalentrada = Flotante(Cantidad) * Flotante(Precio)

				if Operacion == operacionDevolucionVenta ||  Operacion == operacionTrasladoEntrada{
					Cantidadentrada = Flotante(Cantidad)
					Precioentrada = costo
					Totalentrada = Flotante(Cantidad) * costo
				} else {
					Cantidadentrada = Flotante(Cantidad)
					Precioentrada = Flotante(Precio)
					Totalentrada = Flotante(Cantidad) * Flotante(Precio)
				}

			} else {

				if Operacion == operacionDevolucionCompra ||  Operacion == operacionDevolucionSoporte{
					Cantidadsalida = Flotante(Cantidad)
					Preciosalida = Flotante(Precio)
					Totalsalida = Flotante(Cantidad) * Flotante(Precio)
				} else {

					Cantidadsalida = Flotante(Cantidad)
					Preciosalida = costo
					Totalsalida = Flotante(Cantidad) * costo
				}

			}

			if primero == false ||saldo == 0{
				saldo = Flotante(Cantidad)
				costo = Flotante(Precio)
				total = Flotante(Cantidad) * Flotante(Precio)
				primero = true
			} else {
				log.Println("cantidad")
				log.Println(Cantidadentrada + saldo)
				log.Println("total")
				log.Println(((Flotante(Cantidad) * Flotante(Precio)) + total))

				// OPERACION DE ASIGNACION CALCULOS DE ENTRADA Y SALIDA
				if Operacion == operacionInventarioInicial || Operacion == operacionCompra || Operacion == operacionSoporte || Operacion == operacionDevolucionVenta || Operacion == operacionTrasladoEntrada{
					if Operacion == operacionDevolucionVenta ||  Operacion == operacionTrasladoEntrada{
						costo = ((Flotante(Cantidad) * costo) + total) / (Cantidadentrada - Cantidadsalida + saldo)
					} else {
						costo = ((Flotante(Cantidad) * Flotante(Precio)) + total) / (Cantidadentrada - Cantidadsalida + saldo)
					}
				} else {
					costo = (total - (Flotante(Cantidad) * (costo))) / (Cantidadentrada - Cantidadsalida + saldo)
					if Operacion == operacionDevolucionCompra ||  Operacion == operacionDevolucionSoporte{
						costo = (total - (Flotante(Cantidad) * Flotante(Precio))) / (Cantidadentrada - Cantidadsalida + saldo)
					} else {
						costo = (total - (Flotante(Cantidad) * (costo))) / (Cantidadentrada - Cantidadsalida + saldo)
					}
				}
				log.Println("costo")
				log.Println(costo)
					saldo = saldo + Cantidadentrada - Cantidadsalida
					total = saldo * costo
					costo = Redondear(costo, 2)
			}


			if saldo == 0{
				//totalanterior=0
				Cantidadsaldo=0
				Preciosaldo=0
				Totalsaldo=0

				if Cantidadentrada>0{
					total=Totalentrada
				}else{
					total=Totalsalida
					Totalsalida=totalanterior
				}
				totalanterior=0
					//total
				Cantidadsaldo = saldo
				Preciosaldo = 0
					//costo
				Totalsaldo = 0

			}else{

				if Cantidadentrada>0{
					total=totalanterior+Totalentrada
				}else{
					total=totalanterior-Totalsalida
				}
				totalanterior=total
				Cantidadsaldo = saldo
				Preciosaldo = costo
				Totalsaldo = total
			}


			//if Cantidadentrada>0{
			//	total=totalanterior+Totalentrada
			//}else{
			//	total=totalanterior-Totalsalida
			//}
			//totalanterior=total
			//
			//
			//
			//Cantidadsaldo = saldo
			//Preciosaldo = costo
			//Totalsaldo = total

			// operaciones

			//res = append(res, inventario{Fecha, Tipo, Codigo,
			///	Bodega, Producto, Cantidad, Precio, Operacion})

			listadokardex = append(listadokardex, kardex{
				Fecha.Format("02/01/06"),
				mes,
				"",
				Producto,
				Tipo,
				Operacion,
				Codigo,
				Bodega,
				Cantidadentrada,
				Precioentrada,
				Totalentrada,
				Cantidadsalida,
				Preciosalida,
				Totalsalida,
				Cantidadsaldo,
				Preciosaldo,
				Totalsaldo})

		}

	}

	if siexiste == false {
		var slice []string
		slice = make([]string, 0)
		data, _ := json.Marshal(slice)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		if TipoParametro=="Todos"{
			listadokardexfinal=listadokardex
		}else{
			var totalEntrada float64
			var totalSalida float64

			totalEntrada=0
			totalSalida=0


			for _, miFila := range listadokardex {

				if miFila.Operacion==TipoParametro{
					totalEntrada+=miFila.Totalentrada
					totalSalida+=miFila.Totalsalida
					listadokardexfinal=append(listadokardexfinal,miFila)

				}

			}

			// suma totales
			listadokardexfinal=append(listadokardexfinal, kardex{
				"",
				"",
				"",
				"",
				"Total",
				"",
				"",
				"",
				0,
				0,
				totalEntrada,
				0,
				0,
				totalSalida,
				0,
				0,
				0})
		}
		data, _ := json.Marshal(listadokardexfinal)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}


// CENTRO KARDEX
func KardexLista(w http.ResponseWriter, r *http.Request) {
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

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func Redondear(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// INICIA CENTRO PDF
func kardexPdf(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	Fecha := mux.Vars(r)["Fecha"]
	t := inventario{}
	var e empresa = ListaEmpresa()
	var c ciudad = TraerCiudad(e.Ciudad)
	err := db.Get(&t, "SELECT * FROM inventario where fecha=$1", Fecha)
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
		pdf.CellFormat(190, 10, "Datos Centro de Costos", "0", 0,
			"C", false, 0, "")
		pdf.Ln(10)
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Fecha", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Fecha.Format("02/01/2006"), "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Tipo:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Tipo, "", 0,
		"", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.CellFormat(40, 4, "Documento No.:", "", 0,
		"", false, 0, "")
	pdf.CellFormat(40, 4, t.Codigo, "", 0,
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

// TERMINA CENTRO PDF
