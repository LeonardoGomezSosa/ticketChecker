package Login

import (
	"fmt"

	"../../Modelos/Usuario"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
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
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Login.Login.go: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")

	if sessionUtils.IsStarted(ctx) {
		fmt.Println("Se redirige al inicio si ya existe una cookie")
		ctx.Redirect("/admin", 301)
	} else {
		ctx.Render("Login/login.html", nil)
	}

}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("Login.Login.go: POST")
	var VistaRespuesta UsuarioPlusErrores
	VistaRespuesta.Error = []Errores{}

	hasError := false

	Usuario := ctx.FormValue("Usuario")
	Password := ctx.FormValue("Password")
	fmt.Println(Usuario, Password)
	fmt.Printf("u: %v\np:%v\n", Usuario, Password)
	Usuario = MoGeneral.LimpiarCadena(Usuario)
	Password = MoGeneral.LimpiarCadena(Password)

	if MoGeneral.CadenaVacia(Usuario) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Usuario vacio"})
		hasError = true
	}

	if MoGeneral.CadenaVacia(Password) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Password vacio"})
		hasError = true
	}

	if hasError {
		fmt.Println("Errores al validar")
		ctx.Render("Login/login.html", VistaRespuesta)
	} else {
		var usu *usuario.Usuario
		exito, usu := ValidarUsuario(Usuario, Password)
		VistaRespuesta.Usr = *usu
		if !exito {
			VistaRespuesta.Usr = *usu
			fmt.Println(VistaRespuesta.Usr)
			VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Usuario o password incorrectos."})
			ctx.Render("Login/login.html", VistaRespuesta)

		} else {
			fmt.Println("si valida Correo")
			//crear cookie
			if !sessionUtils.IsStarted(ctx) {
				fmt.Println("Se debe instanciar una sesion")
				usu.Password = ""
				usu.Empresa = ""
				usu.Correo = ""
				usu.Coleccion = ""
				sessionUtils.StartSession(ctx, *usu)
				fmt.Println(ctx.GetCookie("Usuario"))
			}
			//redireccionar
			ctx.Redirect("/admin", 301)
		}
	}

}

// ValidarUsuario Valida el ingreso por contraseña de un  usuario
func ValidarUsuario(u string, p string) (bool, *usuario.Usuario) {
	var usr usuario.Usuario
	ptrDB, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println("No se ha podido establecer conexión: ", err)
		return false, nil
	}
	defer ptrDB.Close()

	usr.Usuario = u
	usr.Correo = u
	usr.Password = p

	existeUser, err := usuario.QueryFieldValueExist(u, "Usuario", "ADMINISTRADORES")
	existeEmail, err := usuario.QueryFieldValueExist(u, "Correo", "ADMINISTRADORES")

	if existeUser || existeEmail {
		queryS := fmt.Sprintf(`
	SELECT "Usuario", "Nombre",  "Correo", COUNT(*) as "Coincidencias" FROM public."ADMINISTRADORES" 
	WHERE ("Usuario"='%v' OR "Correo"= '%v') AND  "Password"='%v' GROUP BY "Usuario"
	`, u, u, p)

		fmt.Println(queryS)
		var Coincidencias int64
		row := ptrDB.QueryRow(queryS)
		err = row.Scan(&usr.Usuario, &usr.Nombre, &usr.Correo,
			&Coincidencias)

		if err != nil {
			fmt.Println("No se ha podido realizar la consulta: ", err)
			return false, &usr
		}

		if Coincidencias == 1 {
			fmt.Println("Su usuario ha sido validado ")
			return true, &usr
		}
	}

	return false, &usr

}

func Admin(ctx *iris.Context) {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("Login.Admin.go: GET")
	fmt.Println("=================================")
	fmt.Println("=================================")
	if sessionUtils.IsStarted(ctx) {
		ctx.Render("Admin/index.html", nil)
	} else {
		ctx.Redirect("/Login", 301)
	}
}
