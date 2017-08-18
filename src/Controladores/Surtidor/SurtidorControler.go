package SurtidorControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"../../Modulos/Session"

	"../../Modelos/ExpresionesRegulares"
	"../../Modelos/Surtidor"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../Sesiones"
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
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Surtidors.SurtidorsControler.go.IndexGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor

	var Cabecera, Cuerpo string
	numeroRegistros, _ = Surtidor.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Surtidors, err := Surtidor.GetAll()

	if err != nil {
		fmt.Println("Error al obtener surtidores: ", Surtidors)
	}

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidors[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = Surtidor.GeneraTemplatesBusqueda(Surtidors[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)

	Send.SIndex.SResultados = true
	Send.SEstado = true
	ctx.Render("Surtidor/SurtidorIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Surtidor
func IndexPost(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Surtidors.SurtidorsControler.go.IndexPost: POST")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor

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
	ctx.Render("Surtidor/SurtidorIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Surtidor
func AltaGet(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Surtidors.SurtidorsControler.go.AltaGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor
	Send.SEstado = true
	Send.SMsj = "Listo para cargar"
	Send.SurtidorSS.ID = bson.NewObjectId()

	fmt.Println(Send)
	ctx.Render("Surtidor/SurtidorAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Surtidor
func AltaPost(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Surtidors.SurtidorsControler.go.AltaPost: POST")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor

	//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

	ID := MoGeneral.LimpiarCadena(ctx.FormValue("ID"))
	Usuario := MoGeneral.LimpiarCadena(ctx.FormValue("Usuario"))
	CodigoBarra := MoGeneral.LimpiarCadena(ctx.FormValue("CodigoBarra"))

	IDExist, err := Surtidor.ExistOne(ID)
	fmt.Println("Surtidor", ID, " existe: ", IDExist)
	if err != nil {
		Send.SMsj = fmt.Sprintf("No se pudo buscar en la base de datos, error: %v.", err)
		Send.SEstado = false
		fmt.Println("error al comprobar existencia")
	} else {
		if IDExist {
			Send.SMsj = fmt.Sprintf("El ID existe en la base de datos.")
			Send.SEstado = false
			Send.SurtidorSS.ID = bson.NewObjectId()
			fmt.Println("Elemento ya existe en la base de datos")
		} else {
			if MoGeneral.CadenaVacia(ID) || MoGeneral.CadenaVacia(Usuario) || MoGeneral.CadenaVacia(CodigoBarra) {
				Send.SMsj = fmt.Sprintf("Algun dato esta vacio.")
				Send.SEstado = false
				Send.SurtidorSS.ID = bson.ObjectIdHex(ID)
				Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = CodigoBarra
				Send.SurtidorSS.EUsuarioSurtidor.Usuario = Usuario
				fmt.Println("Datos vacios")
			} else {
				var e Surtidor.SurtidorMgo
				e.IDSurtidor = ID
				e.Usuario = Usuario
				e.CodigoBarra = CodigoBarra
				e.ID = bson.ObjectIdHex(ID)

				Categoria := ExpresionesRegulares.ObtenerCategoriaTexto(e.CodigoBarra)

				if Categoria == "Surtidor" {
					fmt.Println("Objeto a insertar: ", e)

					rs, err := e.InsertOne()
					if err != nil {
						Send.SMsj = fmt.Sprintf("Ocurrio un problema al insertar objeto. Intente mas tarde")
						Send.SEstado = false
					} else {
						Send.SMsj = fmt.Sprintf("Objeto insertado Correctamente: %v.", rs)
						Send.SEstado = true
					}
				} else {
					Send.SMsj = fmt.Sprintf("%v no es una categoria válida.", e.CodigoBarra)
					Send.SEstado = false
				}

				Send.SurtidorSS.ID = e.ID
				Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
				Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			}
		}
	}
	fmt.Println(Send)

	ctx.Render("Surtidor/SurtidorDetalle.html", Send)
}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Surtidor
func EditaGet(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Surtidors.SurtidorsControler.go.EditaGet: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor
	//###### TU CÓDIGO AQUÍ PROGRAMADOR
	id := ctx.Param("ID")
	if id != "" {
		e, err := Surtidor.GetOne(id)
		if err != nil {
			Send.SEstado = false
			Send.SurtidorSS.ID = bson.ObjectIdHex(id)
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			Send.SMsj = "No se encontró El usuario..."
			Send.SEstado = false
		} else {
			Send.SEstado = false
			Send.SurtidorSS.ID = bson.ObjectIdHex(id)
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			Send.SMsj = fmt.Sprintf("Usuario %v:%v localizado...", e.IDSurtidor, e.Usuario)
			Send.SEstado = true
		}
	} else {
		ctx.Redirect("/Surtidors", 301)
	}

	ctx.Render("Surtidor/SurtidorEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Surtidor
func EditaPost(ctx *iris.Context) {

	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Surtidors.SurtidorsControler.go.AltaPost: POST")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor

	//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

	ID := MoGeneral.LimpiarCadena(ctx.FormValue("ID"))
	Usuario := MoGeneral.LimpiarCadena(ctx.FormValue("Usuario"))
	CodigoBarra := MoGeneral.LimpiarCadena(ctx.FormValue("CodigoBarra"))

	IDExist, err := Surtidor.ExistOne(ID)
	fmt.Println("Surtidor", ID, " existe: ", IDExist)
	if err != nil {
		Send.SMsj = fmt.Sprintf("No se pudo buscar en la base de datos, error: %v.", err)
		Send.SEstado = false
		fmt.Println("error al comprobar existencia")
	} else {
		if IDExist {

			if MoGeneral.CadenaVacia(ID) || MoGeneral.CadenaVacia(Usuario) || MoGeneral.CadenaVacia(CodigoBarra) {
				Send.SMsj = fmt.Sprintf("Algun dato esta vacio.")
				Send.SEstado = false
				Send.SurtidorSS.ID = bson.ObjectIdHex(ID)
				Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = CodigoBarra
				Send.SurtidorSS.EUsuarioSurtidor.Usuario = Usuario
				fmt.Println("Datos vacios")
			} else {
				var e Surtidor.SurtidorMgo
				e.IDSurtidor = ID
				e.Usuario = Usuario
				e.CodigoBarra = CodigoBarra
				e.ID = bson.ObjectIdHex(ID)
				fmt.Println("Objeto a insertar: ", e)
				rs, err := e.ModifyOne()
				if err != nil {
					Send.SMsj = fmt.Sprintf("Ocurrio un problema al insertar objeto. Intente mas tarde")
					Send.SEstado = false
				} else {
					Send.SMsj = fmt.Sprintf("Objeto Modificado Correctamente: %v.", rs)
					Send.SEstado = true
				}
				Send.SurtidorSS.ID = e.ID
				Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
				Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			}
		} else {
			Send.SMsj = fmt.Sprintf("El ID no existe en la base de datos.")
			Send.SEstado = false
			Send.SurtidorSS.ID = bson.NewObjectId()
			fmt.Println("Elemento ya existe en la base de datos")
		}
	}
	fmt.Println(Send)

	ctx.Render("Surtidor/SurtidorDetalle.html", Send)

}

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
	var Send Surtidor.SSurtidor
	//###### TU CÓDIGO AQUÍ PROGRAMADOR
	id := ctx.Param("ID")
	if id != "" {
		e, err := Surtidor.GetOne(id)
		if err != nil {
			Send.SurtidorSS.ID = bson.NewObjectId()
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = ""
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = ""
			Send.SMsj = "No se encontró la expresion..."
			Send.SEstado = false
		} else {
			Send.SurtidorSS.ID = e.ID
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			rs, err := e.DeleteOne()
			if err != nil {
				Send.SMsj = fmt.Sprintf("Error al eliminar el elemento %v: %v.", id, e.Usuario)
				Send.SEstado = false
			} else {
				Send.SMsj = fmt.Sprintf("Usuario: %v con  ID: %v ha sido eliminado.", e.Usuario, rs)
				Send.SEstado = true
			}
		}
	} else {
		ctx.Redirect("/Expresions", 301)
	}
	ctx.Render("Surtidor/SurtidorDetalle.html", Send)

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
	var Send Surtidor.SSurtidor
	//###### TU CÓDIGO AQUÍ PROGRAMADOR
	id := ctx.Param("ID")
	if id != "" {
		e, err := Surtidor.GetOne(id)
		if err != nil {
			Send.SurtidorSS.ID = bson.NewObjectId()
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = ""
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = ""
			Send.SMsj = "No se encontró la expresion..."
			Send.SEstado = false
		} else {
			Send.SurtidorSS.ID = e.ID
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			rs, err := e.DeleteOne()
			if err != nil {
				Send.SMsj = fmt.Sprintf("Error al eliminar el elemento %v: %v.", id, e.Usuario)
				Send.SEstado = false
			} else {
				Send.SMsj = fmt.Sprintf("Usuario: %v con  ID: %v ha sido eliminado.", e.Usuario, rs)
				Send.SEstado = true
			}
		}
	} else {
		ctx.Redirect("/Expresions", 301)
	}
	ctx.Render("Surtidor/SurtidorDetalle.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor

	// if !sessionUtils.IsStarted(ctx) {
	// 	http.Redirect(ctx.ResponseWriter, ctx.Request, "/Login", 302)
	// }

	//###### TU CÓDIGO AQUÍ PROGRAMADOR
	id := ctx.Param("ID")
	if id != "" {
		e, err := Surtidor.GetOne(id)
		if err != nil {
			Send.SEstado = false
			Send.SurtidorSS.ID = bson.ObjectIdHex(id)
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			Send.SMsj = "No se encontró El usuario..."
			Send.SEstado = false
		} else {
			Send.SEstado = false
			Send.SurtidorSS.ID = bson.ObjectIdHex(id)
			Send.SurtidorSS.ECodigoBarraSurtidor.CodigoBarra = e.CodigoBarra
			Send.SurtidorSS.EUsuarioSurtidor.Usuario = e.Usuario
			Send.SMsj = fmt.Sprintf("Usuario %v:%v localizado...", e.IDSurtidor, e.Usuario)
			Send.SEstado = true
		}
	} else {
		ctx.Redirect("/Surtidors", 301)
	}

	ctx.Render("Surtidor/SurtidorDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	if !sessionUtils.IsStarted(ctx) {
		ctx.Redirect("/Login", 301)
	}
	var Send Surtidor.SSurtidor

	name, nivel, id := Session.GetUserName(ctx.Request)
	Send.SSesion.Name = name
	Send.SSesion.Nivel = nivel
	Send.SSesion.IDS = id

	//###### TU CÓDIGO AQUÍ PROGRAMADOR

	ctx.Render("Surtidor/SurtidorDetalle.html", Send)
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
	var Send Surtidor.SSurtidor
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
		elementos, err := Surtidor.GetRangeInPage(num, limitePorPagina)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Error al conseguir los datos de la página"
		} else {
			Cabecera, Cuerpo := Surtidor.GeneraTemplatesBusqueda(elementos)
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
