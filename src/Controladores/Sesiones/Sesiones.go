package sessionUtils

import (
	"fmt"
	"net/http"

	"../../Modelos/Reporte"
	"../../Modelos/Usuario"
	"../../Modulos/General"
	"github.com/gorilla/securecookie"
	iris "gopkg.in/kataras/iris.v6"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

// StartSession Inicia la sesion
func StartSession(ctx *iris.Context, usr usuario.Usuario) {
	if encoded, err := cookieHandler.Encode("session", usr); err == nil {
		galleta := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(ctx.ResponseWriter, galleta)
	}
}

// IsStarted verifica si existe la sesion Usuario
func IsStarted(ctx *iris.Context) bool {
	if cookie, err := ctx.Request.Cookie("session"); err == nil {
		cookieValue := usuario.Usuario{}
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			if !MoGeneral.EstaVacio(cookieValue.Usuario) {
				return true
			}
		}
	}
	return false
}

// DeleteSsn Elimina los datos de la sesion
func DeleteSsn(ctx *iris.Context) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(ctx.ResponseWriter, cookie)
	ctx.Redirect("/", 301)
}

// CrearGalletaReporte  Crea una galleta para el monstruo come galletas, especiado con criptografia
func CrearGalletaReporte(ctx *iris.Context, nombre string, valor reporte.Reporte) bool {
	if encoded, err := cookieHandler.Encode(nombre, valor); err == nil {
		galleta := &http.Cookie{
			Name:  nombre,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(ctx.ResponseWriter, galleta)
		return true
	}
	return false
}

// ConsumirGalleta  vino el monstruo come galletas y eliminó la galleta
func ConsumirGalleta(ctx *iris.Context, nombre string) bool {
	cookie := &http.Cookie{
		Name:   nombre,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(ctx.ResponseWriter, cookie)
	return false
}

//LeerGalletaReporte Obtiene los datos de la Galleta de Reporte
func LeerGalletaReporte(ctx *iris.Context, nombre string) *reporte.Reporte {
	cookie, err := ctx.Request.Cookie(nombre)
	if err != nil {
		fmt.Println("no se puede leer cookie: ", err)
		return nil
	}
	fmt.Println("==================================================")
	fmt.Println("Cookie recibida: ", cookie)
	fmt.Println("Cookie Name: ", cookie.Name)
	fmt.Println("Cookie Value: ", cookie.Value)
	fmt.Println("==================================================")
	var reporte reporte.Reporte
	if err = cookieHandler.Decode(cookie.Name, cookie.Value, &reporte); err == nil {
		fmt.Println(reporte)
		return &reporte
	}
	fmt.Println("no se puede decodificar: ", err)
	return nil
}

//LeerGalletaGeneral Obtiene los datos de la Galleta de Reporte
func LeerGalletaGeneral(ctx *iris.Context, nombre string) interface{} {
	cookie, err := ctx.Request.Cookie(nombre)
	if err != nil {
		fmt.Println("no se puede leer cookie: ", err)
		return nil
	}
	fmt.Println("==================================================")
	fmt.Println("Cookie recibida: ", cookie)
	fmt.Println("Cookie Name: ", cookie.Name)
	fmt.Println("Cookie Value: ", cookie.Value)
	fmt.Println("==================================================")
	var reporte interface{}
	if err = cookieHandler.Decode(cookie.Name, cookie.Value, &reporte); err == nil {
		fmt.Println(reporte)
		return &reporte
	}
	fmt.Println("no se puede decodificar: ", err)
	return nil
}
