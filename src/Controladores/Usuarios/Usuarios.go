package usuarioxx

import (
	"fmt"

	_ "github.com/lib/pq"

	"../../Modulos/Conexiones"
)

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

func CrearIndice(bson.ObjectID) bool {

}
