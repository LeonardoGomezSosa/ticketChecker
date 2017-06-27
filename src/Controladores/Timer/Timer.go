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

<<<<<<< HEAD
=======
type Reporte struct {
	CodigoBarraTicket   string
	CodigoBarraSurtidor string
	TimeIn              time.Time
	TimeOut             time.Time
	DuracionM           int64
	SurtidoCompleto     string
}

type ReporteVista struct {
	CodigoBarraTicket   CodigoBarraTicketVista
	CodigoBarraSurtidor CodigoBarraSurtidorVista
	TimeIn              TimeInVista
	TimeOut             TimeOutVista
	DuracionM           DuracionMVista
	SurtidoCompleto     SurtidoCompletoVista
	Timer               bool
}

type CodigoBarraTicketVista struct {
	CodigoBarraSurtidor string //Valor
	Error               string
	Estado              bool
}

type CodigoBarraSurtidorVista struct {
	CodigoBarraSurtidor string //Valor
	Error               string
	Estado              bool
}
type TimeInVista struct {
	TimeIn time.Time //Valor
	Error  string
	Estado bool
}
type TimeOutVista struct {
	TimeOut time.Time //Valor
	Error   string
	Estado  bool
}
type DuracionMVista struct {
	DuracionM int64 //Valor
	Error     string
	Estado    bool
}
type SurtidoCompletoVista struct {
	SurtidoCompleto string //Valor
	Error           string
	Estado          bool
}

>>>>>>> 4801584225d5814b981d543ce787347dcd7bf76d
//IndexGet renderea al indObtenerExpresionesAlmacenadasex de Almacen
func IndexGet(ctx *iris.Context) {
	fmt.Println("Timer.Timer.go: GET")
	if sessionUtils.IsStarted(ctx) {
		ctx.Render("Timer/Timer.html", nil)
	} else {
		ctx.Redirect("/", 301)
	}

}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("Timer.Timerepr.go: POST")
	var rep Reporte
	exp, _, err := ExpresionesRegulares.ObtenerExpresionesAlmacenadas()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	rep.CodigoBarraTicket = ctx.FormValue("Ticket")
	rep.CodigoBarraSurtidor = ctx.FormValue("Surtidor")

	Entrada := ctx.FormValue("Entrada")

	fmt.Println("Entrada: ", rep)

	if Entrada == "" {
		fmt.Println("Viene vacio")

	} else {
		mensaje := "No se encuentra en alguna categoria."
		for _, valor := range exp {
			expresion := fmt.Sprintf("%s", string(valor.ExpresionRegular))
			if MoGeneral.ValidaCadenaExpresion(Entrada, expresion) {
				mensaje = valor.Categoria
				break
			}
		}

		fmt.Println("Categoria", mensaje)

		switch mensaje {
		case "ticket":
			existeEnReporte, report, err := reporte.ConsultarTicketExisteYRegresarContenidoPorCampo(Entrada, "CodigoBarraTicket", "REPORTE")
			rep.CodigoBarraTicket = Entrada
			if (rep.CodigoBarraTicket != "" && rep.CodigoBarraSurtidor == "") || (rep.CodigoBarraTicket == "" && rep.CodigoBarraSurtidor != "") {
				fmt.Println("Falta 1.")
			}
			if err != nil {
				fmt.Println("Error: ", err)
			}
			if existeEnReporte {
				fmt.Println("La matrola existe")
				if reporte.TimeIn.Before(reporte.TimeOut) {
					fmt.Println("Ya se ha cerrado el Ticket")
					rep.CodigoBarraSurtidor = ""
					rep.CodigoBarraTicket = ""
				} else {
					fmt.Println("InsertarSalida")
					salida := time.Now()
					report.TimeOut = salida
					minutos := int64(reporte.TimeOut.Sub(reporte.TimeIn).Minutes())
					report.DuracionM = minutos
					fmt.Println(report)
					err = ActualizaTicket(report)
					if err != nil {
						fmt.Println("Imposible actualizar de ticket leido con entrada ticket.")
					} else {
						fmt.Printf("Tiempo transcurrido: %v minutos, hora actual: %v", minutos, salida)
					}
					rep.CodigoBarraSurtidor = ""
					rep.CodigoBarraTicket = ""
				}
			} else {
				fmt.Println("Dentro de ticket esta en la opcionde alta.")
				if rep.CodigoBarraSurtidor == "" {
					fmt.Println("No hacer nada hasta que tengas Surtidor.")
				} else {
					fmt.Println("Se tiene previamente un surtidor, valida si existe.")
					existeSurtidor, surt, err := Surtidor.QuerySurtidorExist(rep.CodigoBarraSurtidor, "CodigoBarra", "SURTIDORES")
					if err != nil {
						fmt.Println("Error al buscar Surtidor: ", err)
					} else {
						if existeSurtidor {
							fmt.Println("Si el surtidor existe:")

							rep.TimeIn = time.Now()
							rep.TimeOut = rep.TimeIn
							rep.CodigoBarraSurtidor = surt.CodigoBarra
							rep.DuracionM = 0
							err = InsertarTicket(rep)
							if err != nil {
								fmt.Println(err)
								fmt.Println("No ha sido posible insertar, vuelva a intentar")
							} else {
								fmt.Println("Ha sido posible insertar")
							}
						} else {
							fmt.Println("El surtidor no existe.")
						}
					}
					fmt.Println("Se forza el reinicio de la captura.")
					rep.CodigoBarraTicket = ""
					rep.CodigoBarraSurtidor = ""
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
				fmt.Println("Error al buscar Surtidor: ", err)
			}
			if existeSurtidor {
				fmt.Println("Comprobar si tienes un ticket.")
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
								err = ActualizaTicket(report)
								if err != nil {
									fmt.Println("imposible actualizar")
								} else {

									fmt.Printf("Tiempo transcurrido: %v minutos, hora actual: %v", minutos, salida)
								}
								rep.CodigoBarraSurtidor = ""
								rep.CodigoBarraTicket = ""
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
								err := InsertarTicket(rep)
								if err != nil {
									fmt.Println(err)
								}
								rep.CodigoBarraTicket = ""
								rep.CodigoBarraSurtidor = ""
							}
						}
					} else {
						fmt.Println("Error al consultar ticket: ", err)
						rep.CodigoBarraTicket = ""
						rep.CodigoBarraSurtidor = ""
					}
				} else {
					fmt.Println("No hacer nada hasta recibir alguna entrada valida.")
				}
			} else {
				fmt.Println("No existe el surtidor en la bd.")
			}
			break
		default:
			fmt.Println("No hay nada que hacer")
			rep.CodigoBarraTicket = ""
			rep.CodigoBarraSurtidor = ""
		}

	}
	fmt.Println(rep)
	ctx.Render("Timer/Timer.html", rep)

}
