package ExpresionesRegulares

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IExpresion interface con los métodos de la clase
type IExpresion interface {
	InsertaMgo() bool
	InsertaElastic() bool

	ActualizaMgo(campos []string, valores []interface{}) bool
	ActualizaElastic(campos []string, valores []interface{}) bool //Reemplaza No Actualiza

	ReemplazaMgo() bool
	ReemplazaElastic() bool

	ConsultaExistenciaByFieldMgo(field string, valor string)

	ConsultaExistenciaByIDMgo() bool
	ConsultaExistenciaByIDElastic() bool

	EliminaByIDMgo() bool
	EliminaByIDElastic() bool
}

//################################################<<METODOS DE GESTION >>################################################################

//##################################<< INSERTAR >>###################################

//InsertaMgo es un método que crea un registro en Mongo
func (p ExpresionMgo) InsertaMgo() bool {
	result := false
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}

	err = Expresions.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p ExpresionMgo) InsertaElastic() bool {
	var ExpresionE ExpresionElastic

	ExpresionE.IDExpresion = p.IDExpresion
	ExpresionE.Clase = p.Clase
	ExpresionE.Expresion = p.Expresion
	insert := MoConexion.InsertaElastic(MoVar.TipoExpresion, p.ID.Hex(), ExpresionE)
	if !insert {
		fmt.Println("Error al insertar Expresion en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p ExpresionMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Expresions.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p ExpresionMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoExpresion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Expresion en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Expresion en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p ExpresionMgo) ReemplazaMgo() bool {
	result := false
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	err = Expresions.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Expresion en elastic
func (p ExpresionMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoExpresion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Expresion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoExpresion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Expresion en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p ExpresionMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Expresions.Find(bson.M{field: valor}).Count()
	if e != nil {
		fmt.Println(e)
	}
	if n > 0 {
		result = true
	}
	s.Close()
	return result
}

//ConsultaExistenciaByIDMgo es un método que encuentra un registro en Mongo buscándolo por ID
func (p ExpresionMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Expresions.Find(bson.M{"_id": p.ID}).Count()
	if e != nil {
		fmt.Println(e)
	}
	if n > 0 {
		result = true
	}
	s.Close()
	return result
}

//ConsultaExistenciaByIDElastic es un método que encuentra un registro en Mongo buscándolo por ID
// func (p ExpresionMgo) ConsultaExistenciaByIDElastic() bool {
// 	result := MoConexion.ConsultaElastic(MoVar.TipoExpresion, p.ID.Hex())
// 	return result
// }

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p ExpresionMgo) EliminaByIDMgo() bool {
	result := false
	s, Expresions, err := MoConexion.GetColectionMgo(MoVar.ColeccionExpresion)
	if err != nil {
		fmt.Println(err)
	}
	e := Expresions.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p ExpresionMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoExpresion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Expresion en Elastic")
		return false
	}
	return true
}
