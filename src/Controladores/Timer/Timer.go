package Timer

import (
	"fmt"

	"time"

	"../../Modelos/ExpresionesRegulares"
	"../../Modelos/Reporte"
	"../../Modelos/Surtidor"
	"../../Modulos/General"
	"../Sesiones"
	iris "gopkg.in/kataras/iris.v6"
)

var exp []ExpresionesRegulares.ExpresionRegular

//IndexGet renderea al indObtenerExpresionesAlmacenadasex de Almacen
func IndexGet(ctx *iris.Context) {
	fmt.Println("Timer.Timer.go: GET")
	var vista reporte.ReporteVista
	vista.Estado = true
	vista.Mensaje = "Listo para cargar datos"
	vista.Error = ""
	vista.Timer = false

	if sessionUtils.IsStarted(ctx) {
		ctx.Render("Timer/Timer.html", vista)
	} else {
		ctx.Redirect("/", 301)
	}

}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("Timer.Timerepr.go: POST")
	var rep reporte.Reporte
	var vista reporte.ReporteVista
	vista.Estado = false
	vista.Mensaje = "Listo para cargar datos"

	exp, _, err := ExpresionesRegulares.ObtenerExpresionesAlmacenadas()
	if err != nil {
		fmt.Println("Error: ", err)
		vista.Estado = false
		vista.Mensaje = "Listo para cargar datos"
		vista.Error = "Error al consultar el catalogo de expresiones regulares."
	}
	rep.CodigoBarraTicket = ctx.FormValue("Ticket")
	rep.CodigoBarraSurtidor = ctx.FormValue("Surtidor")

	vista.CodigoBarraTicket.CodigoBarraTicket = rep.CodigoBarraTicket
	vista.CodigoBarraSurtidor.CodigoBarraSurtidor = rep.CodigoBarraSurtidor

	Entrada := ctx.FormValue("Entrada")

	fmt.Println("Entrada: ", rep)

	if Entrada == "" {
		fmt.Println("Viene vacio")
		vista.Estado = true
		vista.Error = "Introducir un dato."

	} else {
		mensaje := "No se encuentra en alguna categoria."
		for _, valor := range exp {
			expresion := fmt.Sprintf("%s", string(valor.ExpresionRegular))
			if MoGeneral.ValidaCadenaExpresion(Entrada, expresion) {
				mensaje = valor.Categoria
				break
			}
		}

		switch mensaje {
		case "ticket":
			existeEnReporte, report, err := reporte.ConsultarTicketExisteYRegresarContenidoPorCampo(Entrada, "CodigoBarraTicket", "REPORTE")
			rep.CodigoBarraTicket = Entrada
			if err != nil {
				vista.Estado = true
				vista.Error = fmt.Sprintf("Error al consultar Ticket: %v.", err)
				vista.CodigoBarraTicket.CodigoBarraTicket = ""
				rep.CodigoBarraTicket = ""
			}
			if (rep.CodigoBarraTicket != "" && rep.CodigoBarraSurtidor == "") || (rep.CodigoBarraTicket == "" && rep.CodigoBarraSurtidor != "") {
				vista.Estado = true
				if rep.CodigoBarraTicket != "" {
					vista.Mensaje = "Tiene 3 segundos para introducir Codigo de Surtidor"
				} else if rep.CodigoBarraSurtidor != "" {
					vista.Mensaje = "Tiene 3 segundos para introducir Codigo de Ticket"
				}

			}
			if existeEnReporte {
				fmt.Println("La matrola existe")
				if report.TimeIn.Before(report.TimeOut) {
					vista.Mensaje = "Ya se ha cerrado el ticket."
					vista.Estado = true
					vista.CodigoBarraTicket.CodigoBarraTicket = ""
					vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
				} else {
					fmt.Println("InsertarSalida")
					salida := time.Now()
					report.TimeOut = salida
					minutos := int64(report.TimeOut.Sub(report.TimeIn).Minutes())
					report.DuracionM = minutos
					fmt.Println(report)
					err = reporte.ActualizaTicket(report)
					if err != nil {
						fmt.Println("Imposible actualizar de ticket leido con entrada ticket.")
						vista.Estado = false
						vista.Error = fmt.Sprintf("Error al Actualizar Reporte del Ticket:\n %v", err)
						vista.Mensaje = ""
					} else {
						vista.Estado = true
						vista.Mensaje = fmt.Sprintf("Tiempo transcurrido: %v minutos, hora actual: %v", minutos, salida)
						vista.Error = ""
					}
					vista.CodigoBarraTicket.CodigoBarraTicket = ""
					vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
				}
			} else {
				fmt.Println("Dentro de ticket esta en la opcionde alta.")
				if rep.CodigoBarraSurtidor == "" {
					fmt.Println("No hacer nada hasta que tengas Surtidor.")
					vista.CodigoBarraTicket.CodigoBarraTicket = rep.CodigoBarraTicket
				} else {
					fmt.Println("Se tiene previamente un surtidor, valida si existe.")
					existeSurtidor, surt, err := Surtidor.QuerySurtidorExist(rep.CodigoBarraSurtidor, "CodigoBarra", "SURTIDORES")
					if err != nil {
						fmt.Println("Error al buscar Surtidor: ", err)
						vista.Estado = false
						vista.Error = fmt.Sprintf("Error al buscar Surtidor:\n %v", err)
						vista.Mensaje = ""
					} else {
						if existeSurtidor {
							fmt.Println("Si el surtidor existe:")
							rep.TimeIn = time.Now()
							rep.TimeOut = rep.TimeIn
							rep.CodigoBarraSurtidor = surt.CodigoBarra
							rep.DuracionM = 0
							err = reporte.InsertarTicket(rep)
							if err != nil {
								vista.Estado = false
								vista.Error = fmt.Sprintf("Error al Insertar Ticket Surtidor:\n %v.", err)
								vista.Mensaje = ""
							}
							vista.Estado = true
							vista.Mensaje = fmt.Sprintf("Insertado el registro, Corre tu tiempo.")
							vista.Error = ""
							vista.CodigoBarraTicket.CodigoBarraTicket = ""
							vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
						} else {
							fmt.Println("El surtidor no existe.")
							vista.Estado = false
							vista.Error = fmt.Sprintf("El surtidor no existe.")
							vista.Mensaje = ""
						}
					}
					fmt.Println("Se forza el reinicio de la captura.")
					vista.CodigoBarraTicket.CodigoBarraTicket = ""
					vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
				}
			}
			break
		case "Surtidor":
			rep.CodigoBarraSurtidor = Entrada
			existeSurtidor := false
			existeSurtidor, surt, err := Surtidor.QuerySurtidorExist(rep.CodigoBarraSurtidor, "CodigoBarra", "SURTIDORES")
			if (rep.CodigoBarraTicket != "" && rep.CodigoBarraSurtidor == "") || (rep.CodigoBarraTicket == "" && rep.CodigoBarraSurtidor != "") {
				fmt.Println("Falta 1.")
			}
			if err != nil {
				vista.Estado = false
				vista.Error = fmt.Sprintf("Error al buscar Surtidor:\n %v.", err)
				fmt.Println("Error al buscar Surtidor: ", err)
			}
			if existeSurtidor {
				fmt.Println("Comprobar si tienes un ticket.")
				vista.CodigoBarraSurtidor.CodigoBarraSurtidor = rep.CodigoBarraSurtidor

				if rep.CodigoBarraTicket != "" {
					existeEnReporte, report, err := reporte.ConsultarTicketExisteYRegresarContenidoPorCampo(rep.CodigoBarraTicket, "CodigoBarraTicket", "REPORTE")
					if err == nil {
						if existeEnReporte {
							fmt.Println("La matrola existe")
							if report.TimeIn.Before(report.TimeOut) {
								fmt.Println("Ya se ha cerrado el Ticket")
								rep.CodigoBarraSurtidor = ""
								rep.CodigoBarraTicket = ""
							} else {
								fmt.Println("InsertarSalida")
								salida := time.Now()
								report.TimeOut = salida
								minutos := int64(report.TimeOut.Sub(report.TimeIn).Minutes())
								report.DuracionM = minutos
								fmt.Println(report)
								err = reporte.ActualizaTicket(report)
								if err != nil {
									fmt.Println("imposible actualizar")
									vista.Estado = false
									vista.Error = fmt.Sprintf("Error al Actualizar Ticket Surtidor:\n %v", err)
								} else {
									vista.Estado = true
									vista.Mensaje = fmt.Sprintf("Tiempo transcurrido: %v minutos, hora actual: %v", minutos, salida)
								}
								vista.CodigoBarraTicket.CodigoBarraTicket = ""
								vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
							}
						} else {
							fmt.Println("La matrola no existe.")
							if rep.CodigoBarraSurtidor == "" {
								fmt.Println("No hacer nada hasta que tengas Surtidor.")
							} else {
								fmt.Println("Registrar timein, ticket y surtidor.")
								rep.TimeIn = time.Now()
								rep.TimeOut = rep.TimeIn
								rep.DuracionM = 0
								rep.CodigoBarraSurtidor = surt.CodigoBarra
								err := reporte.InsertarTicket(rep)
								if err != nil {
									vista.Estado = false
									vista.Error = fmt.Sprintf("Error al Insertar Ticket Surtidor:\n %v", err)
								}
								vista.Estado = true
								vista.Mensaje = fmt.Sprintf("Insertado el registro, Corre tu tiempo.")
								vista.Error = ""
								vista.CodigoBarraTicket.CodigoBarraTicket = ""
								vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
							}
						}
					} else {
						fmt.Println("Error al consultar ticket: ", err)
						vista.Estado = false
						vista.Error = fmt.Sprintf("Error al consultar ticket:\n %v.", err)
						vista.CodigoBarraTicket.CodigoBarraTicket = ""
						vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
					}
				} else {
					fmt.Println("No hacer nada hasta recibir alguna entrada valida.")
					vista.Estado = true
					vista.Mensaje = "Tiene 3 segundos para introducir Codigo de Ticket"
				}
			} else {
				fmt.Println("No existe el surtidor en la BD.")
				vista.Estado = false
				vista.Error = "No existe el surtidor en la BD."
				vista.CodigoBarraTicket.CodigoBarraTicket = ""
				vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
			}
			break
		default:
			fmt.Println("No hay nada que hacer")
			vista.Estado = false
			vista.Error = "No es una categoria de Ticket reconocida"
			vista.Mensaje = ""
			vista.CodigoBarraTicket.CodigoBarraTicket = ""
			vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
		}

	}
	fmt.Println(rep)
	ctx.Render("Timer/Timer.html", vista)

}
