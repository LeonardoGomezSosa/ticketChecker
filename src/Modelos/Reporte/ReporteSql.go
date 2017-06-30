package reporte

import (
	"database/sql"
	"fmt"
	"time"

	"../../Modulos/Conexiones"
	_ "github.com/lib/pq"
)

// Reporte Estructura para almacenar y leer datos de la tabla reporte de postgres
type Reporte struct {
	CodigoBarraTicket   string
	CodigoBarraSurtidor string
	TimeIn              time.Time
	TimeOut             time.Time
	DuracionM           int64
	Respuesta           string
}

// ConsultarTicketExiste consulta si existe una entrada dada de ticket dado un campo
func ConsultarTicketExiste(value string, field string, table string) (bool, error) {
	ptrDB, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false, err
	}

	stmt := fmt.Sprintf(`SELECT count(*)  FROM public."%v" where "%v"='%v'`, table, field, value)
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
	"CodigoBarraTicket", "CodigoBarraSurtidor", "TimeIn", "TimeOut", "DuracionMinutos", "Respuesta" 
	FROM public."%v" where "%v"='%v' ORDER BY "TimeIn" DESC
	LIMIT 1
	`,
		table, field, value)

	row = ptrDB.QueryRow(stmt)
	err = row.Scan(&rep.CodigoBarraTicket, &rep.CodigoBarraSurtidor, &rep.TimeIn, &rep.TimeOut, &rep.DuracionM, &rep.Respuesta)
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
	query := fmt.Sprintf(`INSERT INTO public."%v" VALUES('%v','%v','%v','%v','%v','%v')`, "REPORTE",
		rep.CodigoBarraTicket, rep.CodigoBarraSurtidor,
		rep.TimeIn.Format("2006-01-02 15:04:05 -0700"),
		rep.TimeIn.Format("2006-01-02 15:04:05 -0700"),
		rep.DuracionM, rep.Respuesta)
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
	query := fmt.Sprintf(`UPDATE public."%v" SET "TimeOut"='%v', "DuracionMinutos"=%v , "Respuesta"='%v' WHERE "CodigoBarraTicket"='%v'`, "REPORTE",
		rep.TimeOut.Format("2006-01-02 15:04:05 -0700"), rep.DuracionM, rep.Respuesta, rep.CodigoBarraTicket)
	fmt.Println("Consulta:   ", query)
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql()
	if err != nil {
		fmt.Println("Errores al conectar con postgres: ", err)
		return err
	}
	BasePsql.Exec("set transaction isolation level serializable")

	_, errsql := SesionPsql.Exec(query)
	if errsql != nil {
		SesionPsql.Rollback()
		BasePsql.Close()
		fmt.Println("Error al insertar el Ticket")
		return err
	}
	SesionPsql.Commit()
	BasePsql.Close()
	return nil
}
