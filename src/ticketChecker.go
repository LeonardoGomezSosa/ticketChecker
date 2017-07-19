package main

import (
	"fmt"

	"./Controladores/EliminarCuenta"
	"./Controladores/ExpresionesRegulares"
	"./Controladores/Login"
	"./Controladores/Recuperar"
	"./Controladores/Sesiones"
	"./Controladores/Surtidor"
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

	app.Get("/admin", Login.Admin)

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

	//Alta
	app.Get("/Surtidors/alta", SurtidorControler.AltaGet)
	app.Post("/Surtidors/alta", SurtidorControler.AltaPost)

	//Detalle
	app.Get("/Surtidors/detalle", SurtidorControler.DetalleGet)
	app.Post("/Surtidors/detalle", SurtidorControler.DetallePost)
	app.Get("/Surtidors/detalle/:ID", SurtidorControler.DetalleGet)
	app.Post("/Surtidors/detalle/:ID", SurtidorControler.DetallePost)

	//Edicion
	app.Get("/Surtidors/edita", SurtidorControler.EditaGet)
	app.Post("/Surtidors/edita", SurtidorControler.EditaPost)
	app.Get("/Surtidors/edita/:ID", SurtidorControler.EditaGet)
	app.Post("/Surtidors/edita/:ID", SurtidorControler.EditaPost)
	//Elimina
	app.Get("/Surtidors/Elimina/", SurtidorControler.EliminaGet)
	app.Post("/Surtidors/Elimina/", SurtidorControler.EliminaPost)
	app.Get("/Surtidors/Elimina/:ID", SurtidorControler.EliminaGet)
	app.Post("/Surtidors/Elimina/:ID", SurtidorControler.EliminaPost)

	app.Post("/Surtidors/search", SurtidorControler.BuscaPagina)
	// app.Post("/Surtidors/agrupa", SurtidorControler.MuestraIndexPorGrupo)

	//Index (Búsqueda)
	app.Get("/Expresions", ExpresionControler.IndexGet)
	app.Post("/Expresions", ExpresionControler.IndexPost)

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

	//Elimina
	app.Get("/Expresions/Elimina", ExpresionControler.EliminaGet)
	app.Post("/Expresions/Elimina", ExpresionControler.EliminaPost)
	app.Get("/Expresions/Elimina/:ID", ExpresionControler.EliminaGet)
	app.Post("/Expresions/Elimina/:ID", ExpresionControler.EliminaPost)

	app.Post("/Expresions/search", ExpresionControler.BuscaPagina)
	// app.Post("/Expresions/agrupa", ExpresionControler.MuestraIndexPorGrupo)
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
