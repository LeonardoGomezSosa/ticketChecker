package Indice

import (
	"fmt"

	iris "gopkg.in/kataras/iris.v6"
)

//IndexGet renderea al index de Almacen
func IndexGet(ctx *iris.Context) {
	fmt.Println("Inicio.Index.go: GET")
	fmt.Println(ctx.GetCookie("Usuario"))
	ctx.Render("Vistas/index.html", nil)
}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("Inicio.Index.go: POST")
	ctx.Render("Vistas/index.html", nil)

}
