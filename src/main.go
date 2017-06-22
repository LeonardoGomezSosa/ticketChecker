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

	app.Get("/", Login.IndexGet)
	app.Post("/", Login.IndexPost)

	app.Get("/Recuperar", Recuperar.IndexGet)
	app.Post("/Recuperar", Recuperar.IndexPost)

	app.Get("/Salir", sessionUtils.DeleteSsn)
	app.Post("/Salir", sessionUtils.DeleteSsn)

	app.Get("/Eliminar", EliminarCuenta.IndexGet)
	app.Post("/Eliminar", EliminarCuenta.IndexPost)

	app.Get("/Timer", Timer.IndexGet)
	app.Post("/Timer", Timer.IndexPost)

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
