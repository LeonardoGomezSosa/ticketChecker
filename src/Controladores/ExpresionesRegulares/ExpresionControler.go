package ExpresionControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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
var limitePorPagina = 10

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrID []string
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

	var Cabecera, Cuerpo string
	numeroRegistros, _ := (ExpresionesRegulares.CountAll())
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Expresions, _ := ExpresionesRegulares.GetAll()

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(Expresions[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = ExpresionesRegulares.GeneraTemplatesBusqueda(Expresions[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)

	Send.SIndex.SResultados = true
	Send.SEstado = true

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
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send ExpresionesRegulares.SExpresion
	// if !sessionUtils.IsStarted(ctx) {
	// 	ctx.Redirect("/Login", 301)
	// }
	Send.SEstado = true
	Send.SMsj = "Listo para dar de alta una nueva expresion regular, se genera un ID Aleatorio."
	Send.SResultados = false
	Send.Expresion.ID = bson.NewObjectId()
	Send.EIDExpresionExpresion.IDExpresion = Send.Expresion.ID.Hex()
	Send.EClaseExpresion.Clase = ""
	Send.EExpresionExpresion.Expresion = ""
	ctx.Render("Expresion/ExpresionAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Expresion
func AltaPost(ctx *iris.Context) {
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.AltaGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	var Send ExpresionesRegulares.SExpresion
	id := MoGeneral.LimpiarCadena(ctx.FormValue("IDExpresion"))
	Clase := MoGeneral.LimpiarCadena(ctx.FormValue("Clase"))
	Expresion := MoGeneral.LimpiarCadena(ctx.FormValue("Expresion"))

	IDExist, err := ExpresionesRegulares.ExistOne(id)

	if err != nil {
		Send.SMsj = fmt.Sprintf("No se pudo buscar en la base de datos, error: %v.", err)
		Send.SEstado = false
	} else {
		if IDExist {
			Send.SMsj = fmt.Sprintf("El ID existe en la base de datos.")
			Send.SEstado = false
		} else {
			if MoGeneral.CadenaVacia(id) || MoGeneral.CadenaVacia(Clase) || MoGeneral.CadenaVacia(Expresion) {
				Send.SMsj = fmt.Sprintf("Algun dato esta vacio.")
				Send.SEstado = false
			} else {
				var e ExpresionesRegulares.ExpresionMgo
				e.IDExpresion = id
				e.Clase = Clase
				e.Expresion = Expresion
				e.ID = bson.ObjectIdHex(id)

				rs, err := e.InsertOne()
				if err != nil {
					Send.SMsj = fmt.Sprintf("Ocurrio un problema al insertar objeto. Intente mas tarde")
					Send.SEstado = false
				} else {
					Send.SMsj = fmt.Sprintf("Objeto insertado Correctamente: %v.", rs)
					Send.SEstado = true
				}
				Send.Expresion.ID = e.ID
				Send.Expresion.EIDExpresionExpresion.IDExpresion = e.IDExpresion
				Send.Expresion.EClaseExpresion.Clase = e.Clase
				Send.Expresion.EExpresionExpresion.Expresion = e.Expresion
			}
		}
	}

	ctx.Render("Expresion/ExpresionDetalle.html", Send)

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Expresion
func EditaGet(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.EditaGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send ExpresionesRegulares.SExpresion
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
	ctx.Render("Expresion/ExpresionEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Expresion
func EditaPost(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.EditaPost: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send ExpresionesRegulares.SExpresion
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
			Clase := MoGeneral.LimpiarCadena(ctx.FormValue("Clase"))
			Expresion := MoGeneral.LimpiarCadena(ctx.FormValue("Expresion"))

			if MoGeneral.CadenaVacia(Clase) || MoGeneral.CadenaVacia(Expresion) {
				Send.SMsj = "No se puede modificar el elemento indicado, no se aceptan campos Vacios"
				Send.Expresion.EIDExpresionExpresion.IDExpresion = e.IDExpresion
				Send.Expresion.EClaseExpresion.Clase = e.Clase
				Send.Expresion.EExpresionExpresion.Expresion = e.Expresion
				Send.SEstado = false
			} else {
				Send.Expresion.ID = e.ID
				Send.Expresion.EIDExpresionExpresion.IDExpresion = e.IDExpresion
				Send.Expresion.EClaseExpresion.Clase = Clase
				Send.Expresion.EExpresionExpresion.Expresion = Expresion
				Send.SMsj = "Modificar Expresion regular localizada..."
				e.Clase = Clase
				e.Expresion = Expresion
				elem, err := e.ModifyOne()
				if err != nil {
					Send.SMsj = fmt.Sprintf("Error al actualizar el elemento: %v.", err)
					Send.SEstado = false
				} else {
					Send.SMsj = fmt.Sprintf("Actualizado el elemento %v.", elem)
					Send.SEstado = true
					// ctx.Redirect(fmt.Sprintf("/Expresions/detalle/%v", e.IDExpresion), 301)
				}
			}
		}
	} else {
		ctx.Redirect("/Expresions", 301)
	}
	ctx.Render("Expresion/ExpresionEdita.html", Send)

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

// eliminar

//EliminaGet renderea a la edición de Expresion
func EliminaGet(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.EliminaGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send ExpresionesRegulares.SExpresion
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
			rs, err := e.DeleteOne()
			if err != nil {
				Send.SMsj = fmt.Sprintf("Error al eliminar el elemento %v: %v.", id, e.Expresion)
				Send.SEstado = false
			} else {
				Send.SMsj = fmt.Sprintf("Expresion regular: %v con  ID: %v ha sido eliminada.", e.Expresion, rs)
				Send.SEstado = true
			}
		}
	} else {
		ctx.Redirect("/Expresions", 301)
	}
	ctx.Render("Expresion/ExpresionDetalle.html", Send)

}

//EliminaPost renderea a la edición de Expresion
func EliminaPost(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.EliminaPost: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send ExpresionesRegulares.SExpresion
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
			rs, err := e.DeleteOne()
			if err != nil {
				Send.SMsj = fmt.Sprintf("Error al eliminar el elemento %v: %v.", id, e.Expresion)
				Send.SEstado = false
			} else {
				Send.SMsj = fmt.Sprintf("Expresion regular: %v con  ID: %v ha sido eliminada.", e.Expresion, rs)
				Send.SEstado = true
			}
		}
	} else {
		ctx.Redirect("/Expresions", 301)
	}
	ctx.Render("Expresion/ExpresionDetalle.html", Send)

}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Expresionesregulares.ExpresionControler.go.BuscaPagina: Pos")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send ExpresionesRegulares.SExpresion
	Pagina := MoGeneral.LimpiarCadena(ctx.FormValue("Pag"))
	if Pagina != "" {
		num, _ := strconv.Atoi(Pagina)
		fmt.Println("Pagina: ", num)
		if num > paginasTotales {
			num = paginasTotales
		}
		if num == 0 {
			num = 1
		}

		fmt.Println("Pagina: ", num)

		NumPagina = num
		elementos, err := ExpresionesRegulares.GetRangeInPage(num, limitePorPagina)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Error al conseguir los datos de la página"
		} else {
			Cabecera, Cuerpo := ExpresionesRegulares.GeneraTemplatesBusqueda(elementos)
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
		}

		Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, num)
		Send.SIndex.SPaginacion = template.HTML(Paginacion)
		// Send.SIndex.SRMsj = string(num)
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
