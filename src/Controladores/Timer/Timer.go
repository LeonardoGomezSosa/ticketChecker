package Timer

import (
	"fmt"

	"time"

	"../../Modelos/ExpresionesRegulares"
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
		ctx.Redirect("/Login", 301)
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

	// err = ctx.Request.ParseForm()
	// if err != nil {
	// 	fmt.Println(err)
	// }
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
			existeEnReporte, err := ConsultarTicketExiste(Entrada, "CodigoBarraTicket", "REPORTE")
			rep.CodigoBarraTicket = Entrada
			if err != nil {
				fmt.Println("Error: ", err)
			}
			if existeEnReporte {
				fmt.Println("La matrola existe")

			} else {
				fmt.Println("La matrola no existe, registrar.")
			}

			break
		case "Surtidor":
			rep.CodigoBarraSurtidor = Entrada
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

// modificar posterior para metodo propio de clASE
