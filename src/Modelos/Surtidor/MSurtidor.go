package Surtidor

import (
	"fmt"
	"strconv"

	"github.com/leekchan/accounting"

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
func GetAll() []SurtidorMgo {
	var result []SurtidorMgo
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}
	err = Surtidors.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Surtidors.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) SurtidorMgo {
	var result SurtidorMgo
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}
	err = Surtidors.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
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
	Surtidors := GetAll()

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
	floats := accounting.Accounting{Symbol: "", Precision: 2}
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
	var err bool

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
