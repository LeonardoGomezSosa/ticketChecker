package MoConexion

import (
	"fmt"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2"
)

//DataM es una estructura que contiene los datos de configuración en el archivo cfg
var DataM = MoVar.CargaSeccionCFG(MoVar.SecMongo)

//check verifica y escribe el error que elige el usuario y grita el error general
func check(err error, mensaje string) {
	if err != nil {
		fmt.Println("##########")
		fmt.Println(mensaje)
		fmt.Println("##########")
		panic(err)
	}
}

//GetColectionMgo regresa una colección específica de Mongo
func GetColectionMgo(coleccion string) (*mgo.Session, *mgo.Collection, error) {
	s, err := mgo.Dial(DataM.Servidor)
	if err != nil {
		return nil, nil, err
	}
	c := s.DB(DataM.NombreBase).C(coleccion)
	return s, c, nil
}

//GetConexionMgo regresa una sesion de mgo y error
func GetConexionMgo() (*mgo.Session, error) {
	session, err := mgo.Dial(DataM.Servidor)
	if err != nil {
		return session, err
	}
	return session, nil
}

//CloseConexionMgo cierra la sesion que se especifica
func CloseConexionMgo(sesion *mgo.Session) {
	sesion.Close()
}

//GetBaseMgo regresa un objeto database de mgo específico de una sesion específica
func GetBaseMgo(base string, sesion *mgo.Session) *mgo.Database {
	sesion.SetMode(mgo.Monotonic, true)
	return sesion.DB(base)
}
