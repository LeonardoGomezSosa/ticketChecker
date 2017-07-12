package Surtidor

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//ISurtidor interface con los métodos de la clase
type ISurtidor interface {
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
func (p SurtidorMgo) InsertaMgo() bool {
	result := false
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}

	err = Surtidors.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p SurtidorMgo) InsertaElastic() bool {
	var SurtidorE SurtidorElastic

	SurtidorE.CodigoBarra = p.CodigoBarra
	SurtidorE.Usuario = p.Usuario
	insert := MoConexion.InsertaElastic(MoVar.TipoSurtidor, p.ID.Hex(), SurtidorE)
	if !insert {
		fmt.Println("Error al insertar Surtidor en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p SurtidorMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Surtidors.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p SurtidorMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoSurtidor, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Surtidor en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Surtidor en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p SurtidorMgo) ReemplazaMgo() bool {
	result := false
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	err = Surtidors.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Surtidor en elastic
func (p SurtidorMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoSurtidor, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Surtidor en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoSurtidor, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Surtidor en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p SurtidorMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Surtidors.Find(bson.M{field: valor}).Count()
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
func (p SurtidorMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Surtidors.Find(bson.M{"_id": p.ID}).Count()
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
// func (p SurtidorMgo) ConsultaExistenciaByIDElastic() bool {
// 	result := MoConexion.ConsultaElastic(MoVar.TipoSurtidor, p.ID.Hex())
// 	return result
// }

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p SurtidorMgo) EliminaByIDMgo() bool {
	result := false
	s, Surtidors, err := MoConexion.GetColectionMgo(MoVar.ColeccionSurtidor)
	if err != nil {
		fmt.Println(err)
	}
	e := Surtidors.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p SurtidorMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoSurtidor, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Surtidor en Elastic")
		return false
	}
	return true
}
