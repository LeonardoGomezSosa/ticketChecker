package ListaControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	//"log"
	"strconv"

	"../../Modelos/ListaModel"
	"../../Modulos/General"
	"gopkg.in/kataras/iris.v6"
)

//##########< Variables Generales > ############

var cadenaBusqueda string
var cadenaBusquedaSat string
var numeroRegistros int
var numeroRegistrosSat int
var paginasTotales int
var paginasTotalesSat int

//NumPagina especifica el numero de página en la que se cargarán los registros
var NumPagina int

var tipoConsultaBase = 1

//limitePorPagina limite de registros a mostrar en la pagina
var limitePorPaginaBase = 10
var limitePorPaginaSat = 10

var productoElastic ListaModel.ProductosBase
var arrProductoElastic []ListaModel.ProductosBase

var productoElasticSat ListaModel.ProductosSat
var arrProductoElasticSat []ListaModel.ProductosSat

var productoVista ListaModel.ProductosBase
var arrProductoVista []ListaModel.ProductosBase

var productoVistaSat ListaModel.ProductosSat
var arrProductoVistaSat []ListaModel.ProductosSat

//####################< LISTA >###########################

//IndexGet renderea al index de Catalogo
func IndexGet(ctx *iris.Context) {
	var Send ListaModel.ListaVista

	//Valida Usuario y obten ID del USUARIO
	Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaBase))

	Send.ID = "USUARIO"
	ctx.Render("Vistas/index.html", Send)

}

//IndexPost renderea al index de Catalogo //Aquí se aplica la Petición
func IndexPost(ctx *iris.Context) {
	ctx.Render("Vistas/index.html", nil)
}

//ConsultaElasticBase recupera la cedena de texto a buscar y regresa los datos con la paginación
func ConsultaElasticBase(ctx *iris.Context) {
	var Send ListaModel.ListaVista

	filtro := ctx.FormValue("Filtro")
	if filtro != "" {
		fil, _ := strconv.Atoi(filtro)
		tipoConsultaBase = fil
	}

	grupo := ctx.FormValue("GrupoBase")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPaginaBase = gru
	}

	cadena := ctx.FormValue("Cadena")

	index := ctx.GetCookie("IDUsuario")
	if index == "" {
		Send.SEstado = false
		Send.SMsj = "No se encuentran datos de sesion."
	} else {
		tipo := "productoservicio"
		fmt.Println("Consulta en base: ", cadena, " con grupo: ", grupo, " y filtro: ", filtro)
		if cadena != "" {

			arrProductoElastic2 := ListaModel.BuscarEnElastic(index, tipo, cadena, tipoConsultaBase)

			numeroRegistros = len(*arrProductoElastic2)

			fmt.Println("Total de registros: ", numeroRegistros)
			if numeroRegistros == 0 {
				Send.SEstado = false
				Send.SMsj = "No se encontró ningún registro para mostrar."
				jData, _ := json.Marshal(Send)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}
			arrProductoElastic = *arrProductoElastic2
			arrProductoVista = []ListaModel.ProductosBase{}
			if numeroRegistros <= limitePorPaginaBase {
				for _, v := range arrProductoElastic[0:numeroRegistros] {
					arrProductoVista = append(arrProductoVista, v)
				}
			} else if numeroRegistros >= limitePorPaginaBase {
				for _, v := range arrProductoElastic[0:limitePorPaginaBase] {
					arrProductoVista = append(arrProductoVista, v)
				}
			}

			Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaBase(arrProductoVista, tipoConsultaBase)
			Send.SCabecera = template.HTML(Cabecera)
			Send.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPaginaBase)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SPaginacion = template.HTML(Paginacion)

			Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaBase))

			Send.SEstado = true
		} else {
			Send.SEstado = false
			Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
		}
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//ConsultaElasticSat recupera la cedena de texto a buscar y regresa los datos con la paginación
func ConsultaElasticSat(ctx *iris.Context) {
	var Send ListaModel.ListaVista

	grupo := ctx.FormValue("GrupoSat")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPaginaSat = gru
	}

	cadena := ctx.FormValue("Cadena")
	index := "clasificadorsat"
	tipo := "clasificadorsat"

	fmt.Println("Consulta en Sat: ", cadena)
	if cadena != "" {

		arrProductoElasticSat2 := ListaModel.BuscarEnElasticSat(index, tipo, cadena)
		numeroRegistrosSat = len(*arrProductoElasticSat2)

		fmt.Println("Total de registros: ", numeroRegistrosSat)

		if numeroRegistrosSat == 0 {
			Send.SEstado = false
			Send.SMsj = "No se encontró ningún registro para mostrar."
			jData, _ := json.Marshal(Send)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		arrProductoElasticSat = *arrProductoElasticSat2

		arrProductoVistaSat = []ListaModel.ProductosSat{}
		if numeroRegistrosSat <= limitePorPaginaSat {
			for _, v := range arrProductoElasticSat[0:numeroRegistrosSat] {
				arrProductoVistaSat = append(arrProductoVistaSat, v)
			}
		} else if numeroRegistrosSat >= limitePorPaginaSat {
			for _, v := range arrProductoElasticSat[0:limitePorPaginaSat] {
				arrProductoVistaSat = append(arrProductoVistaSat, v)
			}
		}

		Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaSat(arrProductoVistaSat)
		Send.SCabecera = template.HTML(Cabecera)
		Send.SBody = template.HTML(Cuerpo)

		paginasTotalesSat = MoGeneral.Totalpaginas(numeroRegistrosSat, limitePorPaginaSat)
		Paginacion := MoGeneral.ConstruirPaginacion2(paginasTotalesSat, 1)
		Send.SPaginacion = template.HTML(Paginacion)

		Send.SEstado = true
	} else {
		Send.SEstado = false
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send ListaModel.ListaVista

	filtro := ctx.FormValue("Filtro")
	if filtro != "" {
		fil, _ := strconv.Atoi(filtro)
		tipoConsultaBase = fil
	}

	Pagina := ctx.FormValue("Pag")

	if Pagina != "" {
		num, _ := strconv.Atoi(Pagina)
		if num > 0 {
			if numeroRegistros > 0 {
				NumPagina = num
				skip := limitePorPaginaBase * (NumPagina - 1)
				limite := skip + limitePorPaginaBase

				arrProductoVista = []ListaModel.ProductosBase{}
				if NumPagina == paginasTotales {
					final := int(numeroRegistros) % limitePorPaginaBase
					if final == 0 {
						for _, v := range arrProductoElastic[skip:limite] {
							arrProductoVista = append(arrProductoVista, v)
						}
					} else {
						for _, v := range arrProductoElastic[skip : skip+final] {
							arrProductoVista = append(arrProductoVista, v)
						}
					}

				} else {
					for _, v := range arrProductoElastic[skip:limite] {
						arrProductoVista = append(arrProductoVista, v)
					}
				}

				if tipoConsultaBase == 2 {
					arrProductoVista = *ListaModel.CheckConCodigo(&arrProductoVista)
				} else if tipoConsultaBase == 1 {
					arrProductoVista = *ListaModel.CheckSinCodigo(&arrProductoVista)
				}

				Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaBase(arrProductoVista, tipoConsultaBase)
				Send.SCabecera = template.HTML(Cabecera)
				Send.SBody = template.HTML(Cuerpo)

				paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPaginaBase)
				Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, num)
				Send.SPaginacion = template.HTML(Paginacion)

				Send.SEstado = true
			} else {
				Send.SMsj = "No se pueden mostrar páginas, no existen registros que mostrar."
				Send.SEstado = false
			}
		} else {
			Send.SMsj = "No hay páginas Cero."
			Send.SEstado = false
		}

	} else {
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
		Send.SEstado = false
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//BuscaPaginaSat regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPaginaSat(ctx *iris.Context) {
	var Send ListaModel.ListaVista
	Pagina := ctx.FormValue("Pag")

	if Pagina != "" {
		num, _ := strconv.Atoi(Pagina)
		if num > 0 {
			if numeroRegistrosSat > 0 {
				NumPagina = num
				skip := limitePorPaginaSat * (NumPagina - 1)
				limite := skip + limitePorPaginaSat

				arrProductoVistaSat = []ListaModel.ProductosSat{}
				if NumPagina == paginasTotalesSat {
					final := int(numeroRegistrosSat) % limitePorPaginaSat
					if final == 0 {
						for _, v := range arrProductoElasticSat[skip:limite] {
							arrProductoVistaSat = append(arrProductoVistaSat, v)
						}
					} else {
						for _, v := range arrProductoElasticSat[skip : skip+final] {
							arrProductoVistaSat = append(arrProductoVistaSat, v)
						}
					}

				} else {
					for _, v := range arrProductoElasticSat[skip:limite] {
						arrProductoVistaSat = append(arrProductoVistaSat, v)
					}
				}

				Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaSat(arrProductoVistaSat)
				Send.SCabecera = template.HTML(Cabecera)
				Send.SBody = template.HTML(Cuerpo)

				paginasTotalesSat = MoGeneral.Totalpaginas(numeroRegistrosSat, limitePorPaginaSat)
				Paginacion := MoGeneral.ConstruirPaginacion2(paginasTotalesSat, num)
				Send.SPaginacion = template.HTML(Paginacion)

				Send.SEstado = true
			} else {
				Send.SMsj = "No se pueden mostrar páginas, no existen registros que mostrar."
				Send.SEstado = false
			}
		} else {
			Send.SMsj = "No hay páginas Cero."
			Send.SEstado = false
		}

	} else {
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
		Send.SEstado = false
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//MuestraIndexPorGrupoB regresa template de busqueda y paginacion de acuerdo a la agrupacion solicitada
func MuestraIndexPorGrupoB(ctx *iris.Context) {
	var Send ListaModel.ListaVista

	filtro := ctx.FormValue("Filtro")
	if filtro != "" {
		fil, _ := strconv.Atoi(filtro)
		tipoConsultaBase = fil
	}

	grupo := ctx.FormValue("GrupoBase")
	fmt.Println(grupo)
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPaginaBase = gru
	}

	cadenaBusqueda = ctx.FormValue("Cadena")

	index := ctx.GetCookie("IDUsuario")
	if index == "" {
		Send.SEstado = false
		Send.SMsj = "No se encuentran datos de sesion."
	} else {
		tipo := "productoservicio"
		fmt.Println("Consulta en base: ", cadenaBusqueda, " con grupo: ", grupo, " y filtro: ", filtro)

		if cadenaBusqueda != "" {

			arrProductoElastic2 := ListaModel.BuscarEnElastic(index, tipo, cadenaBusqueda, tipoConsultaBase)
			numeroRegistros = len(*arrProductoElastic2)

			if numeroRegistros == 0 {
				Send.SEstado = false
				Send.SMsj = "No se encontró ningún registro para mostrar."
				jData, _ := json.Marshal(Send)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}
			arrProductoElastic = *arrProductoElastic2

			arrProductoVista = []ListaModel.ProductosBase{}
			if numeroRegistros <= limitePorPaginaBase {
				for _, v := range arrProductoElastic[0:numeroRegistros] {
					arrProductoVista = append(arrProductoVista, v)
				}
			} else if numeroRegistros >= limitePorPaginaBase {
				for _, v := range arrProductoElastic[0:limitePorPaginaBase] {
					arrProductoVista = append(arrProductoVista, v)
				}
			}

			Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaBase(arrProductoVista, tipoConsultaBase)
			Send.SCabecera = template.HTML(Cabecera)
			Send.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPaginaBase)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SPaginacion = template.HTML(Paginacion)

			Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaBase))

			Send.SEstado = true

		} else {
			if len(arrProductoElastic) > 0 {

				arrProductoVista = []ListaModel.ProductosBase{}
				if numeroRegistros <= limitePorPaginaBase {
					for _, v := range arrProductoElastic[0:numeroRegistros] {
						arrProductoVista = append(arrProductoVista, v)
					}
				} else if numeroRegistros >= limitePorPaginaBase {
					for _, v := range arrProductoElastic[0:limitePorPaginaBase] {
						arrProductoVista = append(arrProductoVista, v)
					}
				}

				if tipoConsultaBase == 2 {
					arrProductoVista = *ListaModel.CheckConCodigo(&arrProductoVista)
				} else if tipoConsultaBase == 1 {
					arrProductoVista = *ListaModel.CheckSinCodigo(&arrProductoVista)
				}

				Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaBase(arrProductoVista, tipoConsultaBase)
				Send.SCabecera = template.HTML(Cabecera)
				Send.SBody = template.HTML(Cuerpo)

				paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPaginaBase)
				Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
				Send.SPaginacion = template.HTML(Paginacion)

				Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaBase))

			}
		}
	}
	Send.SEstado = true
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//MuestraIndexPorGrupoS regresa template de busqueda y paginacion de acuerdo a la agrupacion solicitada
func MuestraIndexPorGrupoS(ctx *iris.Context) {
	var Send ListaModel.ListaVista

	grupo := ctx.FormValue("GrupoSat")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPaginaSat = gru
	}

	cadenaBusqueda = ctx.FormValue("Cadena")
	index := "clasificadorsat"
	tipo := "clasificadorsat"

	if cadenaBusqueda != "" {

		arrProductoElasticSat2 := ListaModel.BuscarEnElasticSat(index, tipo, cadenaBusqueda)

		numeroRegistrosSat = len(*arrProductoElasticSat2)

		if numeroRegistrosSat == 0 {
			Send.SEstado = false
			Send.SMsj = "No se encontró ningún registro para mostrar."
			jData, _ := json.Marshal(Send)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}
		arrProductoElasticSat = *arrProductoElasticSat2

		arrProductoVistaSat = []ListaModel.ProductosSat{}
		if numeroRegistrosSat <= limitePorPaginaSat {
			for _, v := range arrProductoElasticSat[0:numeroRegistrosSat] {
				arrProductoVistaSat = append(arrProductoVistaSat, v)
			}
		} else if numeroRegistrosSat >= limitePorPaginaSat {
			for _, v := range arrProductoElasticSat[0:limitePorPaginaSat] {
				arrProductoVistaSat = append(arrProductoVistaSat, v)
			}
		}

		Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaSat(arrProductoVistaSat)
		Send.SCabecera = template.HTML(Cabecera)
		Send.SBody = template.HTML(Cuerpo)

		paginasTotalesSat = MoGeneral.Totalpaginas(numeroRegistrosSat, limitePorPaginaSat)
		Paginacion := MoGeneral.ConstruirPaginacion2(paginasTotalesSat, 1)
		Send.SPaginacion = template.HTML(Paginacion)

		Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaSat))

	} else {
		if len(arrProductoElasticSat) > 0 {

			arrProductoVistaSat = []ListaModel.ProductosSat{}
			if numeroRegistrosSat <= limitePorPaginaSat {
				for _, v := range arrProductoElasticSat[0:numeroRegistrosSat] {
					arrProductoVistaSat = append(arrProductoVistaSat, v)
				}
			} else if numeroRegistrosSat >= limitePorPaginaSat {
				for _, v := range arrProductoElasticSat[0:limitePorPaginaSat] {
					arrProductoVistaSat = append(arrProductoVistaSat, v)
				}
			}

			Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaSat(arrProductoVistaSat)
			Send.SCabecera = template.HTML(Cabecera)
			Send.SBody = template.HTML(Cuerpo)

			paginasTotalesSat = MoGeneral.Totalpaginas(numeroRegistrosSat, limitePorPaginaSat)
			Paginacion := MoGeneral.ConstruirPaginacion2(paginasTotalesSat, 1)
			Send.SPaginacion = template.HTML(Paginacion)

			Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaSat))
		}

	}

	Send.SEstado = true
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//AsignarRelacionCatalogo toma datos del formulario y los asigna al catalogo ...xxxxx
func AsignarRelacionCatalogo(ctx *iris.Context) {

	var send ListaModel.SListaVista
	send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaBase))
	ctx.Request.ParseForm()
	SkuSeleccionados := ctx.Request.Form["SkuSeleccionados"][0]
	ClaveSatSeleccionado := ctx.FormValue("ClaveSatSeleccionado")
	if len(SkuSeleccionados) > 0 && ClaveSatSeleccionado != "" {
		var arraySkus []string
		err := json.Unmarshal([]byte(SkuSeleccionados), &arraySkus)
		if err != nil {
			log.Fatal(err)
		}
		idTablaUsuario := ctx.GetCookie("IDUsuario")
		tipo := "productoservicio"
		err = ListaModel.AsignarRelacionClaveSku(idTablaUsuario, tipo, idTablaUsuario, ClaveSatSeleccionado, arraySkus)
		if err != nil {
			send.SEstado = false
			send.SMsj = "Ocurrieron errores al guardar su relacion. Intentelo mas tarde. : " + err.Error()
			//ctx.Render("Vistas/index.html", send)
			ctx.MustRender("Vistas/index.html", send)
		} else {
			send.SEstado = true
			send.SMsj = "Sus datos se han registrado correctamente"
			ctx.Render("Vistas/index.html", send)
		}
	} else {
		send.SEstado = false
		send.SMsj = "Seleccione al menos un elemento de cada categoria"
		ctx.Render("Vistas/index.html", send)
		//ctx.Writef("Debera seleccionar una categoria")
	}
}
