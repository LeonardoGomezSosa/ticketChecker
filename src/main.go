package main

import (
	"fmt"

	"./Controladores/EliminarCuenta"
	"./Controladores/Login"
	"./Controladores/Recuperar"
	"./Controladores/Sesiones"
	"./Controladores/Timer"

	"./Modulos/Variables"

	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

func main() {

	//###################### Start ####################################
	app := iris.New()
	app.Adapt(httprouter.New())
	app.Adapt(view.HTML("../Public/Pages", ".html").Reload(true))

	app.Set(iris.OptionCharset("UTF-8"))

	//app.StaticWeb("/icono", "./Recursos/Generales/img")

	app.StaticWeb("/css", "../Public/Resources/css")
	app.StaticWeb("/js", "../Public/Resources/js")
	app.StaticWeb("/Plugins", "../Public/Resources/Plugins")
	app.StaticWeb("/scripts", "../Public/Resources/scripts")
	app.StaticWeb("/img", "../Public/Resources/img")

	//###################### CFG ######################################

	var DataCfg = MoVar.CargaSeccionCFG(MoVar.SecDefault)

	//###################### Ruteo ####################################

	app.Get("/Login", Login.IndexGet)
	app.Post("/Login", Login.IndexPost)

	app.Get("/Recuperar", Recuperar.IndexGet)
	app.Post("/Recuperar", Recuperar.IndexPost)

	app.Get("/Salir", sessionUtils.DeleteSsn)
	app.Post("/Salir", sessionUtils.DeleteSsn)

	app.Get("/Eliminar", EliminarCuenta.IndexGet)
	app.Post("/Eliminar", EliminarCuenta.IndexPost)

	app.Get("/", Timer.IndexGet)
	app.Post("/", Timer.IndexPost)

	app.Get("/RecibirRespuesta", Timer.CapturaRespuestaGet)
	app.Post("/RecibirRespuesta", Timer.CapturaRespuestaPost)

	app.Get("/Surtidors", SurtidorControler.IndexGet)
	app.Post("/Surtidors", SurtidorControler.IndexPost)
	app.Post("/Surtidors/search", SurtidorControler.BuscaPagina)
	app.Post("/Surtidors/agrupa", SurtidorControler.MuestraIndexPorGrupo)

	//Index (Búsqueda)
	app.Get("/Expresions", ExpresionControler.IndexGet)
	app.Post("/Expresions", ExpresionControler.IndexPost)
	app.Post("/Expresions/search", ExpresionControler.BuscaPagina)
	app.Post("/Expresions/agrupa", ExpresionControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Expresions/alta", ExpresionControler.AltaGet)
	app.Post("/Expresions/alta", ExpresionControler.AltaPost)

	//Edicion
	app.Get("/Expresions/edita", ExpresionControler.EditaGet)
	app.Post("/Expresions/edita", ExpresionControler.EditaPost)
	app.Get("/Expresions/edita/:ID", ExpresionControler.EditaGet)
	app.Post("/Expresions/edita/:ID", ExpresionControler.EditaPost)

	//Detalle
	app.Get("/Expresions/detalle", ExpresionControler.DetalleGet)
	app.Post("/Expresions/detalle", ExpresionControler.DetallePost)
	app.Get("/Expresions/detalle/:ID", ExpresionControler.DetalleGet)
	app.Post("/Expresions/detalle/:ID", ExpresionControler.DetallePost)

	//###################### Listen Server #############################

	if DataCfg.Puerto != "" {
		fmt.Println("Ejecutandose en el puerto: ", DataCfg.Puerto)
		fmt.Println("Acceder a la siguiente url: ", DataCfg.BaseURL)
		app.Listen(":" + DataCfg.Puerto)
	} else {
		fmt.Println("Ejecutandose en el puerto: 8080")
		fmt.Println("Acceder a la siguiente url: localhost")
		app.Listen(":8080")
	}

}
