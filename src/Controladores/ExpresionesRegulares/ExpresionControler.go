package ExpresionControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"../../Modulos/Session"
	"../Sesiones"

	"../../Modelos/ExpresionesRegulares"
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

//IndexGet renderea al index de Expresion
func IndexGet(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.IndexGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")

	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}

	var Send ExpresionesRegulares.SExpresion

	// var Cabecera, Cuerpo string
	numeroRegistros, err := (ExpresionesRegulares.CountAll())
	if err != nil {
		fmt.Println("Error al contar elementos: ", err)
	}

	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Expresions, err := ExpresionesRegulares.GetAll()
	if err != nil {
		fmt.Println("numeroRegistros: ", numeroRegistros)
	}
	arrIDMgo = []bson.ObjectId{}

	var Cabecera, Cuerpo string
	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(Expresions[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(Expresions[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	// Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	// Paginacion := "MoGeneral.ConstruirPaginacion(paginasTotales, 1)"
	// Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true
	Send.SEstado = true
	Send.SMsj = "Capturar expresion regular"

	ctx.Render("Expresion/ExpresionIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Expresion
func IndexPost(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.IndexPost: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send ExpresionesRegulares.SExpresion

	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Expresion.EVARIABLEExpresion.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := ExpresionesRegulares.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("Expresion/ExpresionIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Expresion
func AltaGet(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.AltaGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	var Send ExpresionesRegulares.SExpresion
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	Send.SEstado = true
	Send.SMsj = "Listo para dar de alta una nueva expresion regular."
	Send.SResultados = false
	Send.Expresion.ID = bson.NewObjectId()
	Send.EIDExpresionExpresion.IDExpresion = ""
	Send.EClaseExpresion.Clase = ""
	Send.EExpresionExpresion.Expresion = ""
	ctx.Render("Expresion/ExpresionAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Expresion
func AltaPost(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.AltaGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	var Send ExpresionesRegulares.SExpresion
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	} else {
		ctx.Render("Expresion/ExpresionAlta.html", Send)
	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Expresion
func EditaGet(ctx *iris.Context) {
	var Send ExpresionesRegulares.SExpresion

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	} else if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true

		//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

		ctx.Render("ExpresionEdita.html", Send)
	} else {
		ctx.Render("IndexDashboard.html", Send)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Expresion
func EditaPost(ctx *iris.Context) {

	var Send ExpresionesRegulares.SExpresion

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	if name == "" {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	} else if nivel == "Administrador" {
		Send.SSesion.IsAdmin = true

		//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

		ctx.Render("ExpresionEdita.html", Send)
	} else {
		ctx.Render("IndexDashboard.html", Send)
	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send ExpresionesRegulares.SExpresion

	if !sessionUtils.IsStarted(ctx) {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	}

	//###### TU CÓDIGO AQUÍ PROGRAMADOR
	id := ctx.Param("ID")
	if id != "" {
		e, err := ExpresionesRegulares.GetOne(id)
		if err != nil {
			Send.Expresion.EIDExpresionExpresion.IDExpresion = ""
			Send.Expresion.EClaseExpresion.Clase = ""
			Send.Expresion.EExpresionExpresion.Expresion = ""
			Send.SMsj = "No se encontró la expresion..."
			Send.SEstado = false
		} else {
			Send.Expresion.ID = e.ID
			Send.Expresion.EIDExpresionExpresion.IDExpresion = e.IDExpresion
			Send.Expresion.EClaseExpresion.Clase = e.Clase
			Send.Expresion.EExpresionExpresion.Expresion = e.Expresion
			Send.SMsj = "Expresion regular localizada..."
			Send.SEstado = true
		}
	} else {
		ctx.Redirect("/Expresions", 301)
	}

	ctx.Render("Expresion/ExpresionDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send ExpresionesRegulares.SExpresion

	if !sessionUtils.IsStarted(ctx) {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	}

	//###### TU CÓDIGO AQUÍ PROGRAMADOR
	id := ctx.Param("ID")
	if id != "" {
		Expresions, err := ExpresionesRegulares.GetOne(id)
		if err != nil {
			Send.Expresion.ID = Expresions.ID
			Send.Expresion.EIDExpresionExpresion.IDExpresion = Expresions.IDExpresion
			Send.Expresion.EClaseExpresion.Clase = Expresions.Clase
			Send.Expresion.EExpresionExpresion.Expresion = Expresions.Expresion
			Send.SMsj = "Expresion regular localizada..."
			Send.SEstado = true
		} else {
			Send.SMsj = "No se encontró la expresion..."
			Send.SEstado = false
		}
	} else {
		ctx.Redirect("/Expresions", 301)
	}

	ctx.Render("Expresion/ExpresionDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send ExpresionesRegulares.SExpresion

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

		Cabecera, Cuerpo := ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrToMongo))
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
	var Send ExpresionesRegulares.SExpresion
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Expresion.ENombreExpresion.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := ExpresionesRegulares.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(ExpresionesRegulares.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
