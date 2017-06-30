package Timer

import (
	"fmt"
	"strconv"

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
	vista.TimerOn = false
	vista.Concluido = false

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
	vista.TimerOn = true
	vista.Concluido = false

	exp, _, err := ExpresionesRegulares.ObtenerExpresionesAlmacenadas()
	if err != nil {
		fmt.Println("Error: ", err)
		vista.Estado = false
		vista.Mensaje = "Listo para cargar datos"
		vista.Error = "Error al consultar el catalogo de expresiones regulares."
	}
	rep.CodigoBarraTicket = ctx.FormValue("Ticket")
	rep.CodigoBarraSurtidor = ctx.FormValue("Surtidor")
	concluido := ctx.FormValue("Concluido")

	vista.CodigoBarraTicket.CodigoBarraTicket = rep.CodigoBarraTicket
	vista.CodigoBarraSurtidor.CodigoBarraSurtidor = rep.CodigoBarraSurtidor
	vista.Concluido, err = strconv.ParseBool(concluido)

	if err != nil {
		fmt.Println("Error al parsear Concluido", err)
		vista.Concluido = false
	}

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
			if !vista.Concluido {
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
						vista.CodigoBarraTicket.CodigoBarraTicket = report.CodigoBarraTicket
						vista.CodigoBarraSurtidor.CodigoBarraSurtidor = report.CodigoBarraSurtidor
						vista.TimeIn.TimeIn = report.TimeIn
						vista.TimeOut.TimeOut = report.TimeOut
						vista.DuracionM.DuracionM = report.DuracionM
						err = reporte.ActualizaTicket(report)
						if err != nil {
							fmt.Println("Imposible actualizar de ticket leido con entrada ticket.")
							vista.Estado = false
							vista.Error = fmt.Sprintf("Error al Actualizar Reporte del Ticket:\n %v", err)
							vista.Mensaje = ""
							vista.CodigoBarraTicket.CodigoBarraTicket = ""
							vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
						} else {
							vista.Estado = true
							vista.Mensaje = fmt.Sprintf("Tiempo transcurrido: %v minutos, hora actual: %v. \n ¿Ha Surtido el pedido Completo?", minutos, salida.Format("2006-01-02 15:04:05 -0700"))
							nombre := "vista"
							galletaCreada := sessionUtils.CrearGalletaReporte(ctx, nombre, report)
							if galletaCreada {
								ctx.Redirect("/queestapasando", 301)
								// sessionUtils.ConsumirGalleta(ctx, nombre)
							}
							vista.Error = ""
							vista.TimerOn = false
							vista.Concluido = true
						}
					}
				} else {
					fmt.Println("Dentro de ticket esta en la opcionde alta.")
					if rep.CodigoBarraSurtidor == "" {
						fmt.Println("No hacer nada hasta que tengas Surtidor.")
						vista.CodigoBarraTicket.CodigoBarraTicket = rep.CodigoBarraTicket
						vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
					} else {
						fmt.Println("Se tiene previamente un surtidor, valida si existe.")
						existeSurtidor, surt, err := Surtidor.QuerySurtidorExist(rep.CodigoBarraSurtidor, "CodigoBarra", "SURTIDORES")
						if err != nil {
							fmt.Println("Error al buscar Surtidor: ", err)
							vista.Estado = false
							vista.Error = fmt.Sprintf("Error al buscar Surtidor:\n %v", err)
							vista.Mensaje = ""
							vista.CodigoBarraTicket.CodigoBarraTicket = ""
							vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
						} else {
							if existeSurtidor {
								fmt.Println("Si el surtidor existe:")
								rep.TimeIn = time.Now()
								rep.TimeOut = rep.TimeIn
								rep.CodigoBarraSurtidor = surt.CodigoBarra
								rep.DuracionM = 0
								rep.Respuesta = ""
								err = reporte.InsertarTicket(rep)
								if err != nil {
									vista.Estado = false
									vista.Error = fmt.Sprintf("Error al Insertar Ticket Surtidor:\n %v.", err)
									vista.Mensaje = ""
									vista.CodigoBarraTicket.CodigoBarraTicket = ""
									vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
								} else {
									vista.Estado = true
									vista.Mensaje = fmt.Sprintf("Insertado el registro, Corre tu tiempo.")
									vista.Error = ""
									vista.CodigoBarraTicket.CodigoBarraTicket = ""
									vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
								}
							} else {
								fmt.Println("El surtidor no existe.")
								vista.Estado = false
								vista.Error = fmt.Sprintf("El surtidor no existe.")
								vista.Mensaje = ""
								vista.CodigoBarraTicket.CodigoBarraTicket = ""
								vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
							}
						}
					}
				}
			}
			break
		case "Surtidor":
			if !vista.Concluido {
				rep.CodigoBarraSurtidor = Entrada
				existeSurtidor := false
				existeSurtidor, surt, err := Surtidor.QuerySurtidorExist(rep.CodigoBarraSurtidor, "CodigoBarra", "SURTIDORES")
				if (rep.CodigoBarraTicket != "" && rep.CodigoBarraSurtidor == "") || (rep.CodigoBarraTicket == "" && rep.CodigoBarraSurtidor != "") {
					fmt.Println("Falta 1.")
				}
				if err != nil {
					vista.Estado = false
					vista.Error = fmt.Sprintf("Error al buscar Surtidor:\n %v.", err)
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
									vista.CodigoBarraTicket.CodigoBarraTicket = report.CodigoBarraTicket
									vista.CodigoBarraSurtidor.CodigoBarraSurtidor = report.CodigoBarraSurtidor
									vista.TimeIn.TimeIn = report.TimeIn
									vista.TimeOut.TimeOut = report.TimeOut
									vista.DuracionM.DuracionM = report.DuracionM
									vista.Concluido = true
									err = reporte.ActualizaTicket(report)
									if err != nil {
										fmt.Println("Imposible actualizar de ticket leido con entrada ticket.")
										vista.Estado = false
										vista.Error = fmt.Sprintf("Error al Actualizar Reporte del Ticket:\n %v", err)
										vista.Mensaje = ""
										vista.CodigoBarraTicket.CodigoBarraTicket = ""
										vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
									} else {
										vista.Estado = true
										vista.Mensaje = fmt.Sprintf("Tiempo transcurrido: %v minutos, hora actual: %v. \n ¿Ha Surtido el pedido Completo?", minutos, salida.Format("2006-01-02 15:04:05 -0700"))
										nombre := "vista"
										galletaCreada := sessionUtils.CrearGalletaReporte(ctx, nombre, report)
										fmt.Println("La galleta fue creada? ", galletaCreada)
										if galletaCreada {
											ctx.Redirect("/queestapasando", 301)
											// sessionUtils.ConsumirGalleta(ctx, nombre)
										}
										vista.Error = ""
										vista.TimerOn = false
										vista.Concluido = true
									}
								}
							} else {
								fmt.Println("La matrola no existe.")
								if rep.CodigoBarraTicket == "" {
									fmt.Println("No hacer nada hasta que tengas Ticket.")
									vista.CodigoBarraTicket.CodigoBarraTicket = ""
									vista.CodigoBarraSurtidor.CodigoBarraSurtidor = rep.CodigoBarraSurtidor
								} else {
									fmt.Println("Registrar timein, ticket y surtidor.")
									rep.TimeIn = time.Now()
									rep.TimeOut = rep.TimeIn
									rep.DuracionM = 0
									rep.CodigoBarraSurtidor = surt.CodigoBarra
									rep.Respuesta = ""
									vista.CodigoBarraTicket.CodigoBarraTicket = report.CodigoBarraTicket
									vista.CodigoBarraSurtidor.CodigoBarraSurtidor = report.CodigoBarraSurtidor
									vista.TimeIn.TimeIn = report.TimeIn
									vista.TimeOut.TimeOut = report.TimeOut
									vista.DuracionM.DuracionM = report.DuracionM
									err := reporte.InsertarTicket(rep)
									if err != nil {
										vista.Estado = false
										vista.Error = fmt.Sprintf("Error al Insertar Ticket Surtidor:\n %v", err)

									} else {
										vista.Estado = true
										vista.Mensaje = fmt.Sprintf("Insertado el registro, Corre tu tiempo.")
										vista.Error = ""
									}
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
			}
			break
		default:
			if !vista.Concluido {
				fmt.Println("No hay nada que hacer")
				vista.Estado = false
				vista.Error = "No es una categoria de Ticket reconocida"
				vista.Mensaje = ""
				vista.CodigoBarraTicket.CodigoBarraTicket = ""
				vista.CodigoBarraSurtidor.CodigoBarraSurtidor = ""
			} else {
				vista.Concluido = true
			}
		}

	}
	fmt.Println(vista)
	ctx.Render("Timer/Timer.html", vista)

}

//CapturaRespuestaGet regresa la peticon get que se hi
func CapturaRespuestaGet(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Timer.CapturaRespuestaGet")
	fmt.Println("=================================")
	fmt.Println("=================================")
	var vista reporte.ReporteVista
	V := sessionUtils.LeerGalletaReporte(ctx, "vista")

	if V != nil {
		fmt.Println()
		vista.CodigoBarraTicket.CodigoBarraTicket = V.CodigoBarraTicket
		vista.CodigoBarraSurtidor.CodigoBarraSurtidor = V.CodigoBarraSurtidor
		vista.TimeIn.TimeIn = V.TimeIn
		vista.TimeOut.TimeOut = V.TimeOut
		vista.DuracionM.DuracionM = V.DuracionM
		vista.Respuesta.Respuesta = V.Respuesta
		// vista.Estado = true
		// vista.Mensaje = "Confirma Operacion"
		// vista.TimerOn = true
	} else {
		ctx.Redirect("/queestapasando", 301)
	}

	ctx.Render("Timer/TimerResponse.html", vista)
}

//CapturaRespuestaPost regresa la peticon post que se hizo
func CapturaRespuestaPost(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Timer.CapturaRespuestaPost")
	fmt.Println("=================================")
	fmt.Println("=================================")

	fmt.Println("Timer.CapturaRespuestaPost")
	var vista reporte.ReporteVista
	V := sessionUtils.LeerGalletaReporte(ctx, "vista")
	fmt.Println("vista: ", V)

	if V != nil {
		vista.CodigoBarraTicket.CodigoBarraTicket = V.CodigoBarraTicket
		vista.CodigoBarraSurtidor.CodigoBarraSurtidor = V.CodigoBarraSurtidor
		vista.TimeIn.TimeIn = V.TimeIn
		vista.TimeOut.TimeOut = V.TimeOut
		vista.DuracionM.DuracionM = V.DuracionM
		vista.Respuesta.Respuesta = V.Respuesta
		// vista.Estado = true
		// vista.TimerOn = false
	} else {
		sessionUtils.ConsumirGalleta(ctx, "vista")
		ctx.Redirect("/Timer", 301)
	}

	ctx.Render("Timer/TimerResponse.html", vista)

}
