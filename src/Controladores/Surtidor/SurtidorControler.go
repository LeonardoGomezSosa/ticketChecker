package SurtidorControler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"../../Modulos/Session"

	"../../Modelos/Surtidor"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//##########< Variables Generales > ############

var cadenaBusqueda string
var numeroRegistros int
var paginasTotales int

//NumPagina especifica el numero de página en la que se cargarán los registros
var NumPagina int

//limitePorPagina limite de registros a mostrar en la pagina
var limitePorPagina = 5

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Surtidor
func IndexGet(ctx *iris.Context) {

	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	}

	if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true
	}

	var Cabecera, Cuerpo string
	numeroRegistros = Surtidor.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Surtidors := Surtidor.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Surtidors {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidors[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidors[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	// Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("SurtidorIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Surtidor
func IndexPost(ctx *iris.Context) {

	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	}

	if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true
	}

	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Surtidor.EVARIABLESurtidor.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := Surtidor.BuscarEnElastic(cadenaBusqueda)

		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)

			arrToMongo = []bson.ObjectId{}
			if numeroRegistros <= limitePorPagina {
				for _, v := range arrIDElastic[0:numeroRegistros] {
					arrToMongo = append(arrToMongo, v)
				}
			} else if numeroRegistros >= limitePorPagina {
				for _, v := range arrIDElastic[0:limitePorPagina] {
					arrToMongo = append(arrToMongo, v)
				}
			}

			MoConexion.FlushElastic()

			Cabecera, Cuerpo := Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrIDMgo[0:limitePorPagina]))
			}

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

			Send.SIndex.SRMsj = "No se encontraron resultados para: " + cadenaBusqueda + " ."
		}

		Send.SEstado = true

	} else {
		Send.SEstado = false
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
		Send.SResultados = false
	}
	// Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	ctx.Render("SurtidorIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Surtidor
func AltaGet(ctx *iris.Context) {

	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	} else if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true

		//####   TÚ CÓDIGO PARA CARGAR DATOS A LA VISTA DE ALTA----> PROGRAMADOR

		ctx.Render("SurtidorAlta.html", Send)
	} else {
		ctx.Render("IndexDashboard.html", Send)
	}

}

//AltaPost regresa la petición post que se hizo desde el alta de Surtidor
func AltaPost(ctx *iris.Context) {

	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	} else if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true

		//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

		ctx.Render("SurtidorAlta.html", Send)
	} else {
		ctx.Render("IndexDashboard.html", Send)
	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Surtidor
func EditaGet(ctx *iris.Context) {

	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	} else if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true

		//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

		ctx.Render("SurtidorEdita.html", Send)
	} else {
		ctx.Render("IndexDashboard.html", Send)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Surtidor
func EditaPost(ctx *iris.Context) {

	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	} else if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true

		//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

		ctx.Render("SurtidorEdita.html", Send)
	} else {
		ctx.Render("IndexDashboard.html", Send)
	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	}

	if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true
	}

	//###### TU CÓDIGO AQUÍ PROGRAMADOR

	ctx.Render("SurtidorDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	}

	if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true
	}

	//###### TU CÓDIGO AQUÍ PROGRAMADOR

	ctx.Render("SurtidorDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send Surtidor.SSurtidor

	Pagina := ctx.FormValue("Pag")
	if Pagina != "" {
		num, _ := strconv.Atoi(Pagina)
		if num == 0 {
			num = 1
		}
		NumPagina = num
		skip := limitePorPagina * (NumPagina - 1)
		limite := skip + limitePorPagina

		arrToMongo = []bson.ObjectId{}
		if NumPagina == paginasTotales {
			final := int(numeroRegistros) % limitePorPagina
			if final == 0 {
				for _, v := range arrIDElastic[skip:limite] {
					arrToMongo = append(arrToMongo, v)
				}
			} else {
				for _, v := range arrIDElastic[skip : skip+final] {
					arrToMongo = append(arrToMongo, v)
				}
			}

		} else {
			for _, v := range arrIDElastic[skip:limite] {
				arrToMongo = append(arrToMongo, v)
			}
		}

		Cabecera, Cuerpo := Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrToMongo))
		Send.SIndex.SCabecera = template.HTML(Cabecera)
		Send.SIndex.SBody = template.HTML(Cuerpo)

		Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, NumPagina)
		Send.SIndex.SPaginacion = template.HTML(Paginacion)

	} else {
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."

	}

	// Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Send.SEstado = true

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//MuestraIndexPorGrupo regresa template de busqueda y paginacion de acuerdo a la agrupacion solicitada
func MuestraIndexPorGrupo(ctx *iris.Context) {
	var Send Surtidor.SSurtidor
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Surtidor.ENombreSurtidor.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := Surtidor.BuscarEnElastic(cadenaBusqueda)

		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)

			arrToMongo = []bson.ObjectId{}
			if numeroRegistros <= limitePorPagina {
				for _, v := range arrIDElastic[0:numeroRegistros] {
					arrToMongo = append(arrToMongo, v)
				}
			} else if numeroRegistros >= limitePorPagina {
				for _, v := range arrIDElastic[0:limitePorPagina] {
					arrToMongo = append(arrToMongo, v)
				}
			}

			Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrIDMgo[0:limitePorPagina]))
			}

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

			Send.SIndex.SRMsj = "No se encontraron resultados para: " + cadenaBusqueda + " ."
		}

	} else {

		if numeroRegistros <= limitePorPagina {
			Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidor.GetEspecifics(arrIDMgo[0:limitePorPagina]))
		}

		Send.SIndex.SCabecera = template.HTML(Cabecera)
		Send.SIndex.SBody = template.HTML(Cuerpo)

		paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
		Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
		Send.SIndex.SPaginacion = template.HTML(Paginacion)

		Send.SIndex.SRMsj = "No se encontraron resultados para: " + cadenaBusqueda + " ."
	}
	// Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Send.SEstado = true

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}
