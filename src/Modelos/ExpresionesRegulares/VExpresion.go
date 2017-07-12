package ExpresionesRegulares

import (
	"html/template"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EIDExpresionExpresion Estructura de campo de Expresion
type EIDExpresionExpresion struct {
	IDExpresion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EClaseExpresion Estructura de campo de Expresion
type EClaseExpresion struct {
	Clase    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EExpresionExpresion Estructura de campo de Expresion
type EExpresionExpresion struct {
	Expresion string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Expresion estructura de Expresions mongo
type Expresion struct {
	ID bson.ObjectId
	EIDExpresionExpresion
	EClaseExpresion
	EExpresionExpresion
}

//SSesion estructura de variables de sesion de Usuarios del sistema
type SSesion struct {
	Name    string
	Nivel   string
	IDS     string
	IsAdmin bool
}

//SIndex estructura de variables de index
type SIndex struct {
	SResultados bool
	SRMsj       string
	SCabecera   template.HTML
	SBody       template.HTML
	SPaginacion template.HTML
	SGrupo      template.HTML
}

//SExpresion estructura de Expresiones para la vista
type SExpresion struct {
	SEstado bool
	SMsj    string
	Expresion
	SIndex
	SSesion
}
