package ExpresionesRegulares

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/General"
)

// ExpresionRegular Estructura que almacena un objeto expresion regular
type ExpresionRegular struct {
	ID               string
	Categoria        string
	ExpresionRegular string
}

// ObtenerExpresionesAlmacenadas oBTIENE EL CONJUNTO DE LAS EXPRESIONES REGULARES ALMACENADAS
func ObtenerExpresionesAlmacenadas() ([]ExpresionRegular, error) {

	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	Query := fmt.Sprintf(`SELECT "ID", "Categoria", "ExprReg" FROM public."%v"`, "REGEXTKUS")
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return nil, err
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

	return expresion, err
}

// ObtenerCategoriaTexto Devuelve la categoria de un texto dado comparando con las expresiones regulares almacenadas en sistema
func ObtenerCategoriaTexto(PorValidar string) string {
	exp, err := ObtenerExpresionesAlmacenadas()
	if err != nil {
		return fmt.Sprintf("Error al obtener Categorias: %v. ", err)
	}
	for _, valor := range exp {
		expresion := fmt.Sprintf("%s", string(valor.ExpresionRegular))
		if MoGeneral.ValidaCadenaExpresion(PorValidar, expresion) {
			return valor.Categoria
		}
	}
	return "No se encuentra en catalogo"
}
