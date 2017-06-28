package usuario

import (
	"fmt"

	"../../Modulos/Conexiones"
	_ "github.com/lib/pq"
)

//Usuario Estructura para insertar catalogo en elasticsearch
type Usuario struct {
	Usuario   string
	Nombre    string
	Empresa   string
	Correo    string
	Password  string
	Coleccion string
}

// QueryUsrExist consulta si existe un usuario dado en la base de datos
func QueryUsrExist(value string, field string, table string) (bool, Usuario) {
	ptrDB, err := MoConexion.ConexionPsql()
	var usr Usuario

	if err != nil {
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false, usr
	}

	fmt.Println("# Existe al menos uno")
	fmt.Println("# crear stmt")
	stmt := fmt.Sprintf(`SELECT count(*)  FROM public."%v" where "%v"='%v'`, table, field, value)
	fmt.Println("#  stmt: ", stmt)

	row := ptrDB.QueryRow(stmt)
	var existencia int64
	err = row.Scan(&existencia)
	fmt.Println("Existencia: ", existencia)
	if existencia == 0 {
		fmt.Printf("No se encuentra  elemento que contenga  (Campo:%v, Valor: %v) en %v\n", field, value, table)
		return false, usr
	}

	stmt = fmt.Sprintf(`SELECT "Usuario","Correo", "Password"  FROM public."%v" where "%v"='%v'`, table, field, value)
	row = ptrDB.QueryRow(stmt)

	err = row.Scan(&usr.Usuario, &usr.Correo, &usr.Password)

	if err != nil {
		fmt.Println("No se ha podido Recuperar el dato de la consulta : ", err)
		return false, usr
	}

	ptrDB.Close()

	fmt.Printf("Encontrado  (Campo:%v)\n", field)
	return true, usr
}

// QueryFieldValueExist consulta si existe un usuario dado en la base de datos
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

// InsertarUsuarioPostgres funcion que inserta un usuario en la tabla ADMINISTRADORES
func (me *Usuario) InsertarUsuarioPostgres() bool {
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
	INSERT INTO public."ADMINISTRADORES"(
	"Usuario", "Nombre", "Password", "Correo") VALUES(
		'%v','%v','%v','%v'
	)
	`, me.Usuario, me.Nombre, me.Password, me.Correo)
	fmt.Println("Consulta usuario, creacion: \n", queryS)
	stmt, err := tx.Prepare(queryS)
	if err != nil {
		fmt.Println("No se ha podido instanciar  consulta", err)
		return false
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		fmt.Println("Error en InsertarUsuarioPostgres, ejecutar consulta ")
		tx.Rollback() // return an error too, we may want to wrap them
		fmt.Println("No se ha de realizar La creacion del usuario, error:  ", err)
		ptrDB.Close()
		return false
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("No se ha creado el usuario, error:  ", err)
		return false
	}
	fmt.Println("Usuario creado exitosamente")
	return true

}

// EliminarUsuario funcion que elimina un usuario en la tabla ADMINISTRADORES
func (me *Usuario) EliminarUsuario() bool {
	ptrDB, err := MoConexion.ConexionPsql()
	if me.Usuario == "" {
		fmt.Println("No se proporciono usuario ")
		return false
	}
	if err != nil {
		ptrDB.Close()
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false
	}

	fmt.Println("# Existe al menos uno")
	fmt.Println("# crear stmt")
	stmt := fmt.Sprintf(`DELETE FROM public."ADMINISTRADORES" WHERE "Usuario"='%v'`, me.Usuario)
	row := ptrDB.QueryRow(stmt)

	fmt.Println(row)

	ptrDB.Close()
	return true

}

// ActualizaPassUsuario funcion que inserta un usuario en la tabla ADMINISTRADORES
func (me *Usuario) ActualizaPassUsuario(pass string) bool {
	ptrDB, err := MoConexion.ConexionPsql()
	if me.Usuario == "" {
		fmt.Println("No se proporciono usuario ")
		return false
	}
	if err != nil {
		ptrDB.Close()
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false
	}

	fmt.Println("# Existe al menos uno")
	fmt.Println("# crear stmt")
	stmt := fmt.Sprintf(`UPDATE public."ADMINISTRADORES" SET "Password"='%v' WHERE  "Usuario"='%v';`, pass, me.Usuario)
	fmt.Println(stmt)
	ptrDB.QueryRow(stmt)

	ptrDB.Close()
	return true

}
