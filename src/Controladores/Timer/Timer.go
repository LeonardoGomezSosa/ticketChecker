package Timer

import (
	"fmt"

	"../../Modelos/Usuario"
	"../Sesiones"
	iris "gopkg.in/kataras/iris.v6"
)

// UsuarioPlusErrores Errores interfaz para capturar errores
type UsuarioPlusErrores struct {
	Usr   usuario.Usuario
	Error []Errores
}

//Errores interfaz para capturar errores
type Errores struct {
	Error string
}

//IndexGet renderea al index de Almacen
func IndexGet(ctx *iris.Context) {
	fmt.Println("Timer.Timer.go: GET")
	if sessionUtils.IsStarted(ctx) {
		ctx.Render("Timer/Timer.html", nil)
	} else {
		ctx.Render("Login/login.html", nil)
	}

}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("Timer.Timer.go: POST")
	if !sessionUtils.IsStarted(ctx) {
		fmt.Println("")

		ctx.Redirect("/", 301)
	} else {
		ctx.Render("Timer/Timer.html", nil)
	}

}
