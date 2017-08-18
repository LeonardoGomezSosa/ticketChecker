package ExpresionesRegulares

import (
	"fmt"
	"strconv"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//ExpresionMgo estructura de Expresions mongo
type ExpresionMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	IDExpresion string        `bson:"IDExpresion"`
	Clase       string        `bson:"Clase"`
	Expresion   string        `bson:"Expresion"`
}

//ExpresionElastic estructura de Expresions para insertar en Elastic
type ExpresionElastic struct {
	IDExpresion string `json:"IDExpresion"`
	Clase       string `json:"Clase"`
	Expresion   string `json:"Expresion"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() ([]ExpresionMgo, error) {
	var aux ExpresionMgo
	var expresion []ExpresionMgo
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return expresion, err
	}

	Query := fmt.Sprintf(`SELECT * FROM public."%v"`, "REGEXTKUS")
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return expresion, err
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return expresion, err
	}

	for resultSet.Next() {
		err := resultSet.Scan(&aux.IDExpresion, &aux.Clase, &aux.Expresion)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			aux.ID = bson.ObjectIdHex(aux.IDExpresion)
			expresion = append(expresion, aux)
		}

	}
	resultSet.Close()
	stmt.Close()
	BasePosGres.Close()

	return expresion, nil

}

// GetRangeInPage Regresa el rango de eleentos entre paginas
func GetRangeInPage(pagina int, elementosPorPagina int) ([]ExpresionMgo, error) {
	var aux ExpresionMgo
	var expresion []ExpresionMgo
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return expresion, err
	}
	saltar := (pagina - 1) * elementosPorPagina
	Query := fmt.Sprintf(`SELECT * FROM public."%v" LIMIT %v OFFSET %v`, "REGEXTKUS", elementosPorPagina, saltar)
	fmt.Println(Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return expresion, err
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return expresion, err
	}

	for resultSet.Next() {
		err := resultSet.Scan(&aux.IDExpresion, &aux.Clase, &aux.Expresion)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			aux.ID = bson.ObjectIdHex(aux.IDExpresion)
			expresion = append(expresion, aux)
		}

	}
	resultSet.Close()
	stmt.Close()
	BasePosGres.Close()

	return expresion, nil

}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() (int, error) {

	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	Query := fmt.Sprintf(`SELECT count(*) FROM public."%v"`, "REGEXTKUS")
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println("Error preparando sentencia", err)
		return 0, err
	}
	resultSet := stmt.QueryRow()

	var total int
	err = resultSet.Scan(&total)
	if err != nil {
		fmt.Println("Error recuperando numero de registros: ", err)
		return 0, err
	}

	stmt.Close()
	BasePosGres.Close()

	return total, nil
}

//GetOne Regresa un documento específico de SQL (Por Coleccion)
func GetOne(ID string) (ExpresionMgo, error) {
	var result ExpresionMgo
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`SELECT *, count(*) FROM public."%v" WHERE "ID"='%v' GROUP BY "ID"`, "REGEXTKUS", ID)
	fmt.Println("Query One : ", Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println("Error preparando sentencia", err)
		return result, err
	}
	resultSet := stmt.QueryRow()
	i := 0
	err = resultSet.Scan(&result.IDExpresion, &result.Clase, &result.Expresion, &i)
	if err != nil {
		fmt.Println("Error recuperando numero de registros: ", err)
		return result, err
	}

	result.ID = bson.ObjectIdHex(result.IDExpresion)
	stmt.Close()
	BasePosGres.Close()
	fmt.Println("Finaliza GetOne")
	return result, nil
}

//ExistOne Regresa un documento específico de SQL (Por Coleccion)
func ExistOne(ID string) (bool, error) {
	var result bool
	result = false
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`SELECT COUNT(*) FROM public."%v" WHERE "ID"='%v'`, "REGEXTKUS", ID)
	fmt.Println("Query One : ", Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println("Error preparando sentencia", err)
		return result, err
	}

	resultSet := stmt.QueryRow()
	i := 0
	err = resultSet.Scan(&i)
	if i > 0 {
		result = true
	}
	if err != nil {
		fmt.Println("Error recuperando numero de registros: ", err)
		return result, err
	}

	stmt.Close()
	BasePosGres.Close()
	fmt.Println("Finaliza GetOne")
	return result, nil
}

//InsertOne Inserta un documento específico de SQL (Por Coleccion)
func (me *ExpresionMgo) InsertOne() (string, error) {
	var result string
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`INSERT INTO public."%v"(
            "ID", "Categoria", "ExprReg")
    VALUES ('%v', '%v', '%v') RETURNING "ID"`,
		"REGEXTKUS", me.IDExpresion, me.Clase, me.Expresion)

	fmt.Println("Query Update One : ", Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println("Error preparando sentencia", err)
		return result, err
	}
	resultSet := stmt.QueryRow()

	err = resultSet.Scan(&result)
	if err != nil {
		fmt.Println("Error recuperando ultimo actualizado: ", err)
		return result, err
	}

	stmt.Close()
	BasePosGres.Close()
	fmt.Println("Finaliza InsertOne", me)
	return result, nil
}

//ModifyOne Modifica un documento específico de SQL (Por Coleccion)
func (me *ExpresionMgo) ModifyOne() (string, error) {
	var result string
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`UPDATE public."%v"
							SET  "Categoria"='%v', "ExprReg"='%v'
							WHERE "ID"='%v' RETURNING "ID"`,
		"REGEXTKUS", me.Clase, me.Expresion, me.IDExpresion)

	fmt.Println("Query Update One : ", Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println("Error preparando sentencia", err)
		return result, err
	}
	resultSet := stmt.QueryRow()

	err = resultSet.Scan(&result)
	if err != nil {
		fmt.Println("Error recuperando ultimo actualizado: ", err)
		return result, err
	}

	stmt.Close()
	BasePosGres.Close()
	fmt.Println("Finaliza ModifyOne", me)
	return result, nil
}

//DeleteOne Elimina un documento específico de SQL (Por Coleccion)
func (me *ExpresionMgo) DeleteOne() (string, error) {
	var result string
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`DELETE FROM  public."%v"
							WHERE "ID"='%v' RETURNING "ID"`,
		"REGEXTKUS", me.IDExpresion)

	fmt.Println("Query DeleteOne : ", Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println("Error preparando sentencia", err)
		return result, err
	}
	resultSet := stmt.QueryRow()

	err = resultSet.Scan(&result)
	if err != nil {
		fmt.Println("Error recuperando ultimo eliminado: ", err)
		return result, err
	}

	stmt.Close()
	BasePosGres.Close()
	fmt.Println("Finaliza DeleteOne", me)
	return result, nil
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []ExpresionMgo {
	var result []ExpresionMgo
	var aux ExpresionMgo
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = ExpresionMgo{}
		Expresions.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) ExpresionMgo {
	var result ExpresionMgo
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)

	if err != nil {
		fmt.Println(err)
	}
	err = Expresions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result ExpresionMgo
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}
	err = Expresions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboExpresions regresa un combo de Expresion de mongo
func CargaComboExpresions(ID string) string {
	Expresions, err := GetAll()
	if err != nil {
		fmt.Println("Error al cargar combo expresiones")
		return ``
	}
	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Expresions {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Expresion + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Expresion + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Expresions []ExpresionMgo) (string, string) {
	// floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
			<th>#</th>
			
				<th>IDExpresion</th>					
				
				<th>Clase</th>					
				
				<th>Expresion</th>					
				</tr>`

	for k, v := range Expresions {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Expresions/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.IDExpresion + `</td>`

		cuerpo += `<td>` + v.Clase + `</td>`

		cuerpo += `<td>` + v.Expresion + `</td>`

		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("IDExpresion")
	queryQuotes = queryQuotes.Field("IDExpresion")

	queryTilde = queryTilde.Field("Clase")
	queryQuotes = queryQuotes.Field("Clase")

	queryTilde = queryTilde.Field("Expresion")
	queryQuotes = queryQuotes.Field("Expresion")

	var docs *elastic.SearchResult
	// var err bool

	// docs, err = MoConexion.BuscaElastic(MoVar.TipoExpresion, queryTilde)
	// if err {
	// 	fmt.Println("No Match 1st Try")
	// }

	// if docs.Hits.TotalHits == 0 {
	// 	docs, err = MoConexion.BuscaElastic(MoVar.TipoExpresion, queryQuotes)
	// 	if err {
	// 		fmt.Println("No Match 2nd Try")
	// 	}
	// }

	return docs
}
