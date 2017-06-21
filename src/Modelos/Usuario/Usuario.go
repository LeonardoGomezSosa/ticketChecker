package usuario

import (
	"database/sql"
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

	stmt = fmt.Sprintf(`SELECT "Usuario", "Coleccion", "Correo", "Password"  FROM public."%v" where "%v"='%v'`, table, field, value)
	row = ptrDB.QueryRow(stmt)

	err = row.Scan(&usr.Usuario, &usr.Coleccion, &usr.Correo, &usr.Password)

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

	fmt.Println("# Querying")
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

// InsertarUsuarioPostgres funcion que inserta un usuario en la tabla USUARIOS
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
	fmt.Println("# Querying")
	queryS := fmt.Sprintf(`
	INSERT INTO public."USUARIOS"(
	"Usuario", "Nombre", "Empresa", "Password", "Correo", "Coleccion") VALUES(
		'%v','%v','%v','%v','%v','%v'
	)
	`, me.Usuario, me.Nombre, me.Empresa, me.Password, me.Correo, me.Coleccion)
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

	queryS = fmt.Sprintf(`
			CREATE TABLE public."%v"
			(
				"Sku" character varying(50) COLLATE pg_catalog."default" NOT NULL,
				"Descripcion" character varying(500) COLLATE pg_catalog."default" NOT NULL,
				"ClaveSat" character varying(8) COLLATE pg_catalog."default",
				CONSTRAINT key_%v PRIMARY KEY ("Sku")
			)

		`, me.Coleccion, me.Coleccion)
	stmt, err = tx.Prepare(queryS)

	if _, err := stmt.Exec(); err != nil {
		fmt.Println("Error en Crear coleccion, ejecutar consulta ")
		tx.Rollback() // return an error too, we may want to wrap them
		fmt.Println("No se ha de realizar La creacion del indice correspondiente al catalogo usuario, error:  ", err)
		stmt.Close()
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

// CrearIndiceElastic funcion que inserta un usuario en la tabla USUARIOS
func (me *Usuario) CrearIndiceElastic() bool {
	return true
}

//sku, descripcion, clave sat

// CrearIndicePostgres funcion que inserta un usuario en la tabla USUARIOS
func (me *Usuario) CrearIndicePostgres(ptrDB *sql.DB) bool {
	// ptrDB, err := MoConexion.ConexionPsql()
	// if err != nil {
	// 	fmt.Println("No se ha podido establecer conexión: ", err)
	// 	return false
	// }
	// tx2, err := ptrDB.Begin()
	// if err != nil {
	// 	fmt.Println("No se ha podido iniciar transaccion: ", err)
	// 	return false
	// }
	// queryS := fmt.Sprintf(`
	// 		CREATE TABLE public."%v"
	// 		(
	// 			"SKU" character varying(50) COLLATE pg_catalog."default" NOT NULL,
	// 			"DESCRIPCION" character varying(500) COLLATE pg_catalog."default" NOT NULL,
	// 			"CLAVESAT" character varying(8) COLLATE pg_catalog."default",
	// 			CONSTRAINT key_%v PRIMARY KEY ("SKU")
	// 		)

	// 	`, me.Coleccion, me.Coleccion)

	// fmt.Println("Consulta indice usuario, crear tabla: \n", queryS)

	// stmt2, err := tx2.Prepare(queryS)
	// if err != nil {
	// 	fmt.Println("Error al ejecutar stmt")
	// 	return false
	// }
	// defer stmt2.Close()

	// if _, err := stmt2.Exec(me.Coleccion); err != nil {
	// 	tx2.Rollback() // return an error too, we may want to wrap them
	// 	fmt.Println("No se ha de realizar operacion, error:  ", err)

	// 	return false
	// }
	// err = tx2.Commit()
	// if err != nil {
	// 	fmt.Println("No se ha realizado la  operacion, error:  ", err)
	// 	return false
	// }
	// fmt.Println("Tabla creada exitosamente")
	// return true
	return true
}

// InsertarUsuarioPostgres funcion que inserta un usuario en la tabla USUARIOS
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
	stmt := fmt.Sprintf(`DELETE FROM public."USUARIOS" WHERE "Usuario"='%v'`, me.Usuario)
	row := ptrDB.QueryRow(stmt)

	fmt.Println(row)

	ptrDB.Close()
	return true

}

// ActualizaPassUsuario funcion que inserta un usuario en la tabla USUARIOS
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
	stmt := fmt.Sprintf(`UPDATE public."USUARIOS" SET "Password"='%v' WHERE  "Usuario"='%v';`, pass, me.Usuario)
	fmt.Println(stmt)
	ptrDB.QueryRow(stmt)

	ptrDB.Close()
	return true

}
