package ExpresionesRegulares

import (
	"database/sql"
	"fmt"

	"../../Modulos/Conexiones"
)

// ExpresionRegular Estructura que almacena un objeto expresion regular
type ExpresionRegular struct {
	ID               int64
	Categoria        string
	ExpresionRegular string
}

// ObtenerExpresionesAlmacenadas oBTIENE EL CONJUNTO DE LAS EXPRESIONES REGULARES ALMACENADAS
func ObtenerExpresionesAlmacenadas() ([]ExpresionRegular, *sql.Rows, error) {
	db, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer db.Close()

	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	Query := fmt.Sprintf(`SELECT * FROM public."%v"`, "REGEXTKUS")
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	var aux ExpresionRegular
	var expresion []ExpresionRegular

	for resultSet.Next() {
		err := resultSet.Scan(&aux.ID, &aux.Categoria, &aux.ExpresionRegular)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			expresion = append(expresion, aux)
		}

	}
	resultSet.Close()
	stmt.Close()
	BasePosGres.Close()

	return expresion, resultSet, err
}
