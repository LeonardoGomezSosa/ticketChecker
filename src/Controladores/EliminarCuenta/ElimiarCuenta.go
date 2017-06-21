package EliminarCuenta

import (
	"context"
	"fmt"

	"../../Modelos/EliminarCuentaModel"
	"../../Modelos/Usuario"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../Sesiones"
	iris "gopkg.in/kataras/iris.v6"
)

//IndexGet renderea al index de Almacen
func IndexGet(ctx *iris.Context) {
	if !sessionUtils.IsStarted(ctx) {
		sessionUtils.DeleteSsn(ctx)
	}

	fmt.Println("EliminarCuenta.EliminarCuenta.go: GET")
	ctx.Render("EliminarCuenta/Eliminar.html", nil)
}

// IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("EliminarCuenta.EliminarCuenta.go: POST")
	var e EliminarModel.SEliminarVista
	if !sessionUtils.IsStarted(ctx) {
		sessionUtils.DeleteSsn(ctx)
	}
	fmt.Println("EliminarCuenta.Eliminar.go: POST")
	operacion := ctx.FormValue("op")
	usr := MoGeneral.LimpiarCadena(ctx.GetCookie("UsuarioKsd"))
	col := MoGeneral.LimpiarCadena(ctx.GetCookie("ColeccionKsd"))

	switch operacion {
	case "0":
		fmt.Println("Eliminar Indice.", col)
		if !MoGeneral.CadenaVacia(col) {
			fmt.Println("Coleccion: ", col)
			e.SEstado, e.Error, e.SMsj = EliminarIndice(col)
			if e.SEstado {
				e.SEstado = true
			}
		} else {
			e.SEstado = false
			e.SMsj = "No es posible Eliminar catalogo: <br> "
		}
		break
	case "1":
		fmt.Println("Eliminar Cuenta:", usr)
		e.SEstado, e.Error, e.SMsj = EliminarUsuario(usr)
		if e.SEstado {
			ctx.Redirect("/Salir", 301)
		}
		break
	}

	fmt.Println(e)
	ctx.Render("EliminarCuenta/Eliminar.html", e)

	fmt.Println("No se puede imprimir")

}

// EliminarIndice funcion que devuelve el cliente  y verifica que exista el indice dado y si no existe lo crea
func EliminarIndice(indice string) (bool, []EliminarModel.Errores, string) {
	var Errores []EliminarModel.Errores
	cliente, err := MoConexion.GetClienteElastic()
	var ctx = context.Background()
	if err != nil {
		fmt.Println("Error al obtener el cliente de elasticsearch", err)
		Errores = append(Errores, EliminarModel.Errores{Error: "Error al obtener el cliente de elasticsearch"})
		return false, Errores, "Error al obtener el cliente"
	}
	if MoConexion.VerificaIndexName(cliente, indice) {
		// Delete an index
		fmt.Println("Si el indice existe se debe Eliminar")
		_, err = cliente.DeleteIndex(indice).Do(ctx)
		if err != nil {
			// Handle error
			fmt.Println("No se pudo eliminar el indice: ", indice, "\nError:", err)
			Errores = append(Errores, EliminarModel.Errores{Error: "No se pudo eliminar el indice: " + indice})
			return false, Errores, "No se pudo eliminar el indice."
		}
		fmt.Println("El indice se ha eliminado de manera exitosa.")
		return true, Errores, "Indice eliminado de manera exitosa."
	}
	fmt.Println("El indice no existe, no hay nada que hacer.")
	Errores = append(Errores, EliminarModel.Errores{Error: "El indice no existe, no hay nada que hacer."})
	return false, Errores, "Indice no existe."
}

// EliminarUsuario elimina la cuenta del usuario
func EliminarUsuario(usr string) (bool, []EliminarModel.Errores, string) {
	var Errores []EliminarModel.Errores
	existe, usu := usuario.QueryUsrExist(usr, "Usuario", "USUARIOS")
	if existe {
		fmt.Println("Usuario encontrado: ", usu.Usuario)
		borrado := usu.EliminarUsuario()
		indiceborrado, indiceError, indicemensaje := EliminarIndice(usu.Coleccion)
		mensaje := ""
		if !indiceborrado {
			for _, err := range indiceError {
				Errores = append(Errores, err)
			}
			mensaje = indicemensaje
		}
		if !borrado {
			mensaje = "Imposible borrar Usuario. " + mensaje
			Errores = append(Errores, EliminarModel.Errores{Error: "Usuario no borrado"})
			return false, Errores, "El usuario no fue borrado"
		}

		return borrado, Errores, "El usuario existe"
	}
	fmt.Println("Usuario no encontrado")
	Errores = append(Errores, EliminarModel.Errores{Error: "Usuario no encontrado"})
	return existe, Errores, "Usuario no existe."
}
