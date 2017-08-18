package Surtidor

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

//SurtidorMgo estructura de Surtidors mongo
type SurtidorMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	IDSurtidor  string        `bson:"ID,omitempty"`
	CodigoBarra string        `bson:"CodigoBarra"`
	Usuario     string        `bson:"Usuario"`
}

//SurtidorElastic estructura de Surtidors para insertar en Elastic
type SurtidorElastic struct {
	CodigoBarra string `json:"CodigoBarra"`
	Usuario     string `json:"Usuario"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() ([]SurtidorMgo, error) {
	var result []SurtidorMgo
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`SELECT "ID", "CodigoBarra", "Usuario" FROM public."%v"`, "SURTIDORES")
	fmt.Println("Query Surtidores GetAll: ", Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	var aux SurtidorMgo
	for resultSet.Next() {
		err := resultSet.Scan(&aux.IDSurtidor, &aux.CodigoBarra, &aux.Usuario)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			aux.ID = bson.ObjectIdHex(aux.IDSurtidor)
			result = append(result, aux)
		}

	}

	resultSet.Close()
	stmt.Close()
	BasePosGres.Close()

	return result, nil
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() (int, error) {

	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	Query := fmt.Sprintf(`SELECT count(*) FROM public."%v"`, "SURTIDORES")
	stmt, err := BasePosGres.Prepare(Query)
	fmt.Println("Query Count(*) Surtidores: ", Query)

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
	fmt.Println("Total Surtidores: ", total)
	return total, nil
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID string) (SurtidorMgo, error) {
	var result SurtidorMgo
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`SELECT * FROM public."%v" WHERE "ID"='%v' GROUP BY "ID"`, "SURTIDORES", ID)
	fmt.Println("Query GetOne : ", Query)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println("Error preparando sentencia", err)
		return result, err
	}
	resultSet := stmt.QueryRow()
	err = resultSet.Scan(&result.IDSurtidor, &result.Usuario, &result.CodigoBarra)
	if err != nil {
		fmt.Println("Error recuperando numero de registros: ", err)
		return result, err
	}

	result.ID = bson.ObjectIdHex(result.IDSurtidor)
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

	Query := fmt.Sprintf(`SELECT COUNT(*) FROM public."%v" WHERE "ID"='%v'`, "SURTIDORES", ID)
	fmt.Println("Query ExistOne: ", Query)
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
	fmt.Println("Finaliza ExistOne")
	return result, nil
}

//InsertOne Inserta un documento específico de SQL (Por Coleccion)
func (me *SurtidorMgo) InsertOne() (string, error) {
	var result string
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`INSERT INTO public."%v"(
            "ID", "Usuario", "CodigoBarra")
    VALUES ('%v', '%v', '%v') RETURNING "ID"`,
		"SURTIDORES", me.IDSurtidor, me.Usuario, me.CodigoBarra)

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
func (me *SurtidorMgo) ModifyOne() (string, error) {
	var result string
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`UPDATE public."%v"
							SET  "Usuario"='%v', "CodigoBarra"='%v'
							WHERE "ID"='%v' RETURNING "ID"`,
		"SURTIDORES", me.Usuario, me.CodigoBarra, me.IDSurtidor)

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
func (me *SurtidorMgo) DeleteOne() (string, error) {
	var result string
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	Query := fmt.Sprintf(`DELETE FROM  public."%v"
							WHERE "ID"='%v' RETURNING "ID"`,
		"SURTIDORES", me.IDSurtidor)

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

// GetRangeInPage Regresa el rango de eleentos entre paginas
func GetRangeInPage(pagina int, elementosPorPagina int) ([]SurtidorMgo, error) {
	var aux SurtidorMgo
	var expresion []SurtidorMgo
	BasePosGres, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return expresion, err
	}
	saltar := (pagina - 1) * elementosPorPagina
	Query := fmt.Sprintf(`SELECT * FROM public."%v" LIMIT %v OFFSET %v`, "SURTIDORES", elementosPorPagina, saltar)
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
		err := resultSet.Scan(&aux.IDSurtidor, &aux.Usuario, &aux.CodigoBarra)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			aux.ID = bson.ObjectIdHex(aux.IDSurtidor)
			expresion = append(expresion, aux)
		}

	}
	resultSet.Close()
	stmt.Close()
	BasePosGres.Close()

	return expresion, nil

}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []SurtidorMgo {
	var result []SurtidorMgo
	var aux SurtidorMgo
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = SurtidorMgo{}
		Surtidors.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) SurtidorMgo {
	var result SurtidorMgo
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)

	if err != nil {
		fmt.Println(err)
	}
	err = Surtidors.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result SurtidorMgo
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}
	err = Surtidors.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboSurtidors regresa un combo de Surtidor de mongo
func CargaComboSurtidors(ID string) string {
	Surtidors, _ := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Surtidors {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Usuario + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Usuario + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Surtidors []SurtidorMgo) (string, string) {
	cuerpo := ``

	cabecera := `<tr>
			<th>#</th>	
				<th>CodigoBarra</th>					
				<th>Usuario</th>					
				</tr>`

	for k, v := range Surtidors {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Surtidors/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.CodigoBarra + `</td>`

		cuerpo += `<td>` + v.Usuario + `</td>`

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

	queryTilde = queryTilde.Field("CodigoBarra")
	queryQuotes = queryQuotes.Field("CodigoBarra")

	queryTilde = queryTilde.Field("Usuario")
	queryQuotes = queryQuotes.Field("Usuario")

	var docs *elastic.SearchResult
	// var err bool

	// docs, err = MoConexion.BuscaElastic(MoVar.TipoSurtidor, queryTilde)
	// if err {
	// 	fmt.Println("No Match 1st Try")
	// }

	// if docs.Hits.TotalHits == 0 {
	// 	docs, err = MoConexion.BuscaElastic(MoVar.TipoSurtidor, queryQuotes)
	// 	if err {
	// 		fmt.Println("No Match 2nd Try")
	// 	}
	// }

	return docs
}
