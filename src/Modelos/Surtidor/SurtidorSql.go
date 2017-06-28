package Surtidor

import (
	"fmt"

	"../../Modulos/Conexiones"
	_ "github.com/lib/pq" // SE UTILIZA PARA VINCULAR CON POSTGRES
)

//Surtidor Estructura para insertar catalogo en elasticsearch
type Surtidor struct {
	Surtidor    string
	CodigoBarra string
}

// QuerySurtidorExist consulta si existe un Surtidor dado en la base de datos
func QuerySurtidorExist(value string, field string, table string) (bool, Surtidor, error) {
	ptrDB, err := MoConexion.ConexionPsql()
	var usr Surtidor

	if err != nil {
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false, usr, err
	}

	stmt := fmt.Sprintf(`SELECT count(*)  FROM public."%v" where "%v"='%v'`, table, field, value)

	row := ptrDB.QueryRow(stmt)
	var existencia int64
	err = row.Scan(&existencia)
	if existencia == 0 {
		fmt.Printf("No se encuentra  elemento que contenga  (Campo:%v, Valor: %v) en %v\n", field, value, table)
		return false, usr, nil
	}

	stmt = fmt.Sprintf(`SELECT "Usuario", "CodigoBarra" FROM public."%v" where "%v"='%v'`, table, field, value)
	row = ptrDB.QueryRow(stmt)

	err = row.Scan(&usr.Surtidor, &usr.CodigoBarra)

	if err != nil {
		fmt.Println("No se ha podido Recuperar el dato de la consulta : ", err)
		return false, usr, err
	}

	ptrDB.Close()

	fmt.Printf("Encontrado  (Campo:%v)\n", field)
	return true, usr, nil
}

// QueryFieldValueExist consulta si existe un Surtidor dado en la base de datos
func QueryFieldValueExist(value string, field string, table string) (bool, error) {
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
		return false, nil
	}

	fmt.Printf("Encontrado  (Campo:%v, Valor: %v) en %v\n", field, value, table)
	return true, nil
}

// InsertarSurtidorPostgres funcion que inserta un Surtidor en la tabla SURTIDORES
func (me *Surtidor) InsertarSurtidorPostgres() bool {
	ptrDB, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false
	}
	tx, err := ptrDB.Begin()
	if err != nil {
		fmt.Println("No se ha podido iniciar transaccion: ", err)
		return false
	}
	queryS := fmt.Sprintf(`
	INSERT INTO public."SURTIDORES"(
	"Surtidor", "CodigoBarra") VALUES(
		'%v','%v'
	)
	`, me.Surtidor, me.CodigoBarra)
	stmt, err := tx.Prepare(queryS)
	if err != nil {
		fmt.Println("No se ha podido instanciar  consulta", err)
		return false
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		fmt.Println("Error en InsertarSurtidorPostgres, ejecutar consulta ")
		tx.Rollback() // return an error too, we may want to wrap them
		fmt.Println("No se ha de realizar La creacion del Surtidor, error:  ", err)
		ptrDB.Close()
		return false
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("No se ha creado el Surtidor, error:  ", err)
		tx.Rollback()
		ptrDB.Close()
		return false
	}
	fmt.Println("Surtidor creado exitosamente")
	return true

}

// ActualizaNombreSurtidor funcion que inserta un Surtidor en la tabla SURTIDORES
func (me *Surtidor) ActualizaNombreSurtidor(Nombre string) bool {
	ptrDB, err := MoConexion.ConexionPsql()
	if me.Surtidor == "" {
		fmt.Println("No se proporciono Surtidor ")
		return false
	}
	if err != nil {
		ptrDB.Close()
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false
	}

	stmt := fmt.Sprintf(`UPDATE public."SURTIDORES" SET "Usuario"='%v' WHERE  "CodigoBarra"='%v';`, Nombre, me.CodigoBarra)
	ptrDB.QueryRow(stmt)

	ptrDB.Close()
	return true

}
