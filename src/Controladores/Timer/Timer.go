package Timer

import (
	"database/sql"
	"fmt"

	"time"

	"../../Modelos/ExpresionesRegulares"
	"../../Modelos/Surtidor"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../Sesiones"
	iris "gopkg.in/kataras/iris.v6"
)

var exp []ExpresionesRegulares.ExpresionRegular

type Reporte struct {
	CodigoBarraTicket   string
	CodigoBarraSurtidor string
	TimeIn              time.Time
	TimeOut             time.Time
	DuracionM           int64
}

type ReporteVista struct {
	CodigoBarraTicket   CodigoBarraTicketVista
	CodigoBarraSurtidor CodigoBarraSurtidorVista
	TimeIn              TimeInVista
	TimeOut             TimeOutVista
	DuracionM           DuracionMVista
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
			existeEnReporte, reporte, err := ConsultarTicketExisteYRegresarContenidoPorCampo(Entrada, "CodigoBarraTicket", "REPORTE")
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
					reporte.TimeOut = salida
					minutos := int64(reporte.TimeOut.Sub(reporte.TimeIn).Minutes())
					reporte.DuracionM = minutos
					fmt.Println(reporte)
					err = ActualizaTicket(reporte)
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
					existeEnReporte, reporte, err := ConsultarTicketExisteYRegresarContenidoPorCampo(rep.CodigoBarraTicket, "CodigoBarraTicket", "REPORTE")
					if err == nil {
						if existeEnReporte {
							fmt.Println("La matrola existe")
							if reporte.TimeIn.Before(reporte.TimeOut) {
								fmt.Println("Ya se ha cerrado el Ticket")
								rep.CodigoBarraSurtidor = ""
								rep.CodigoBarraTicket = ""
							} else {
								fmt.Println("InsertarSalida")
								salida := time.Now()
								reporte.TimeOut = salida
								minutos := int64(reporte.TimeOut.Sub(reporte.TimeIn).Minutes())
								reporte.DuracionM = minutos
								fmt.Println(reporte)
								err = ActualizaTicket(reporte)
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

// ConsultarTicketExiste consulta si existe una entrada dada de ticket dado un campo
func ConsultarTicketExiste(value string, field string, table string) (bool, error) {
	ptrDB, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false, err
	}

	fmt.Println("# Querying")
	stmt := fmt.Sprintf(`SELECT count(*)  FROM public."%v" where "%v"='%v'`, table, field, value)
	fmt.Println(stmt)
	row := ptrDB.QueryRow(stmt)

	var existencia int64
	err = row.Scan(&existencia)

	ptrDB.Close()
	if err != nil {
		fmt.Println("No se ha podido Recuperar el dato de la consulta : ", err)
		return false, err
	}
	if existencia == 0 {
		fmt.Printf("No se encuentra  elemento que contenga  (Campo:%v, Valor: %v) en %v\n", field, value, table)
		return false, err
	}

	fmt.Printf("Encontrado  (Campo:%v, Valor: %v) en %v\n", field, value, table)
	return true, nil

}

// ConsultarTicketExisteYRegresarContenidoPorCampo Consulta si el indice existe y regresa el elemento en caso de que exista tal
func ConsultarTicketExisteYRegresarContenidoPorCampo(value string, field string, table string) (bool, Reporte, error) {
	ptrDB, err := MoConexion.ConexionPsql()
	var rep Reporte

	if err != nil {
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false, rep, err
	}
	stmt := fmt.Sprintf(`SELECT count(*)  FROM public."%v" where "%v"='%v'`, table, field, value)
	fmt.Println(stmt)
	row := ptrDB.QueryRow(stmt)
	var existencia int64
	err = row.Scan(&existencia)

	if err != nil {
		fmt.Println("No se ha podido Recuperar el dato de la consulta : ", err)
		ptrDB.Close()
		return false, rep, err
	}
	if existencia == 0 {
		fmt.Printf("No se encuentra  elemento que contenga  (Campo:%v, Valor: %v) en %v\n", field, value, table)
		ptrDB.Close()
		return false, rep, nil
	}

	stmt = fmt.Sprintf(`
	SELECT 
	"CodigoBarraTicket", "CodigoBarraSurtidor", "TimeIn", "TimeOut", "DuracionMinutos" 
	FROM public."%v" where "%v"='%v' ORDER BY "TimeIn" DESC
	LIMIT 1
	`,
		table, field, value)

	row = ptrDB.QueryRow(stmt)
	err = row.Scan(&rep.CodigoBarraTicket, &rep.CodigoBarraSurtidor, &rep.TimeIn, &rep.TimeOut, &rep.DuracionM)
	if err != nil {
		fmt.Println("No se ha podido Recuperar el dato de la consulta : ", err)
		ptrDB.Close()
		return false, rep, err
	}
	ptrDB.Close()
	return true, rep, nil

}

// InsertarTicket inserta una entrada en la tabla reporte
func InsertarTicket(rep Reporte) error {
	var SesionPsql *sql.Tx
	var err error
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql()
	if err != nil {
		fmt.Println("Errores al conectar con postgres: ", err)
		return err
	}
	BasePsql.Exec("set transaction isolation level serializable")
	query := fmt.Sprintf(`INSERT INTO public."%v" VALUES('%v','%v','%v','%v','%v')`, "REPORTE",
		rep.CodigoBarraTicket, rep.CodigoBarraSurtidor,
		rep.TimeIn.Format("2006-01-02 15:04:05 -0700"),
		rep.TimeIn.Format("2006-01-02 15:04:05 -0700"),
		rep.DuracionM)
	_, errsql := SesionPsql.Exec(query)
	if errsql != nil {
		SesionPsql.Rollback()
		BasePsql.Close()
		fmt.Println("Error al insertar el Ticket")
		fmt.Println(query)
		return err
	}
	SesionPsql.Commit()
	BasePsql.Close()
	return err
}

// ActualizaTicket inserta una entrada en la tabla reporte
func ActualizaTicket(rep Reporte) error {
	var SesionPsql *sql.Tx
	var err error
	query := fmt.Sprintf(`UPDATE public."%v" SET "TimeOut"='%v', "DuracionMinutos"=%v WHERE "CodigoBarraTicket"='%v'`, "REPORTE",
		rep.TimeOut.Format("2006-01-02 15:04:05 -0700"), rep.DuracionM, rep.CodigoBarraTicket)
	fmt.Println(query)
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql()
	if err != nil {
		fmt.Println("Errores al conectar con postgres: ", err)
		return err
	}
	fmt.Println("1")
	BasePsql.Exec("set transaction isolation level serializable")

	_, errsql := SesionPsql.Exec(query)
	fmt.Println("2")
	if errsql != nil {
		fmt.Println("3")
		SesionPsql.Rollback()
		BasePsql.Close()
		fmt.Println("Error al insertar el Ticket")
		fmt.Println(query)
		return err
	}
	fmt.Println("4")
	SesionPsql.Commit()
	BasePsql.Close()
	return nil
}
