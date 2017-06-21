package sessionUtils

import (
	"fmt"

	"net/http"

	"../../Modelos/Usuario"
	iris "gopkg.in/kataras/iris.v6"
)

// StartSession Inicia la sesion
func StartSession(ctx *iris.Context, usr usuario.Usuario) {
	galleta := &http.Cookie{Name: "UsuarioKsd", Value: usr.Usuario, MaxAge: 1200}
	galleta2 := &http.Cookie{Name: "ColeccionKsd", Value: usr.Coleccion, MaxAge: 1200}

	ctx.SetCookieKV("IDUsuario", usr.Coleccion)
	http.SetCookie(ctx.ResponseWriter, galleta)
	http.SetCookie(ctx.ResponseWriter, galleta2)
}

// IsStarted verifica si existe la sesion Usuario
func IsStarted(ctx *iris.Context) bool {
	cookie, err := ctx.Request.Cookie("UsuarioKsd")
	if err != nil {
		fmt.Println("la galleta no se puede consumir: ", err)
		return false
	}
	fmt.Println("la galleta: ", cookie)
	if cookie.Value != "" {
		return true
	}
	return false
}

// DeleteSsn Elimina los datos de la sesion
func DeleteSsn(ctx *iris.Context) {
	fmt.Println("Borrar sesion")

	http.SetCookie(ctx.ResponseWriter, nil)
	cookie, err := ctx.Request.Cookie("UsuarioKsd")
	if err != nil {
		fmt.Println("La galleta no se puede consumir (Pobre monstruo come galletas): ", err)
	}

	galleta := &http.Cookie{Name: "UsuarioKsd", Path: "/", MaxAge: -1}
	galleta2 := &http.Cookie{Name: "ColeccionKsd", Path: "/", MaxAge: -1}
	http.SetCookie(ctx.ResponseWriter, galleta)
	http.SetCookie(ctx.ResponseWriter, galleta2)
	fmt.Println("la galleta: ", cookie)
	ctx.Redirect("/", http.StatusFound)

}

// BorrarGalleta  Vino el Monstruo y se la comio
func BorrarGalleta(ctx *iris.Context) bool {
	cookie, err := ctx.Request.Cookie("UsuarioKsd")
	ctx.RemoveCookie("UsuarioKsd")

	if err != nil {
		fmt.Println("la galleta", cookie, " no se puede consumir: ", err)
		return false
	}

	http.SetCookie(ctx.ResponseWriter, nil)
	ctx.Redirect("/", 301)
	return true
}
