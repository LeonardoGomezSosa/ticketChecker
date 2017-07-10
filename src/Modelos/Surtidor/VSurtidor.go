package Surtidor

import (
	"html/template"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ECodigoBarraSurtidor Estructura de campo de Surtidor
type ECodigoBarraSurtidor struct {
	CodigoBarra string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EUsuarioSurtidor Estructura de campo de Surtidor
type EUsuarioSurtidor struct {
	Usuario  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//SurtidorSS estructura de Surtidors mongo
type SurtidorSS struct {
	ID bson.ObjectId
	ECodigoBarraSurtidor
	EUsuarioSurtidor
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

//SSurtidor estructura de Surtidores para la vista
type SSurtidor struct {
	SEstado bool
	SMsj    string
	SurtidorSS
	SIndex
	SSesion
}
