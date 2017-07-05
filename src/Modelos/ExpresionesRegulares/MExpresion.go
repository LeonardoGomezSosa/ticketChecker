package ExpresionesRegulares

import (
	"fmt"
	"strconv"

	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"github.com/leekchan/accounting"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//ExpresionMgo estructura de Expresions mongo
type ExpresionMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	IDExpresion int           `bson:"IDExpresion"`
	Clase       string        `bson:"Clase"`
	Expresion   string        `bson:"Expresion"`
}

//ExpresionElastic estructura de Expresions para insertar en Elastic
type ExpresionElastic struct {
	IDExpresion int    `json:"IDExpresion"`
	Clase       string `json:"Clase"`
	Expresion   string `json:"Expresion"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []ExpresionMgo {
	var result []ExpresionMgo
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}
	err = Expresions.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Expresions.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) ExpresionMgo {
	var result ExpresionMgo
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}
	err = Expresions.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
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
	Expresions := GetAll()

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
	floats := accounting.Accounting{Symbol: "", Precision: 2}
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
		cuerpo += `<td>` + string(v.IDExpresion) + `</td>`

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
	var err bool

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
