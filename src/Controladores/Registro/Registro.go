package Registro

import (
	"fmt"

	"../../Modelos/Usuario"
	"../../Modulos/General"

	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//UsuarioPlusErrores interfaz para capturar errores y rellenar los datos obtenidos del formulario
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
	var VistaRespuesta UsuarioPlusErrores
	coleccion := bson.NewObjectId()
	VistaRespuesta.Usr.Coleccion = coleccion.Hex()
	fmt.Println("Coleccion: ", VistaRespuesta.Usr.Coleccion)
	fmt.Println("Registro.Registro.go: GET")
	ctx.Render("Login/registro.html", VistaRespuesta)
}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	var VistaRespuesta UsuarioPlusErrores
	VistaRespuesta.Error = []Errores{}
	fmt.Println("Registro.Registro.go: POST")
	Usuario := ctx.FormValue("Usuario")
	Nombre := ctx.FormValue("Nombre")
	Empresa := ctx.FormValue("Empresa")
	Correo := ctx.FormValue("Correo")
	Password := ctx.FormValue("Password")
	Confirma := ctx.FormValue("Confirma")
	Coleccion := ctx.FormValue("Coleccion")

	fmt.Println("Usuario: ", Usuario)
	fmt.Println("Nombre:", Nombre)
	fmt.Println("Empresa:", Empresa)
	fmt.Println("Correo:", Correo)
	fmt.Println("Password:", Password)
	fmt.Println("Confirma:", Confirma)

	Usuario = MoGeneral.LimpiarCadena(Usuario)
	Nombre = MoGeneral.LimpiarCadena(Nombre)
	Empresa = MoGeneral.LimpiarCadena(Empresa)
	Correo = MoGeneral.LimpiarCadena(Correo)
	Password = MoGeneral.LimpiarCadena(Password)
	Confirma = MoGeneral.LimpiarCadena(Confirma)
	Coleccion = MoGeneral.LimpiarCadena(Coleccion)
	fmt.Println("La coleccion recibida", Coleccion)
	if MoGeneral.CadenaVacia(Usuario) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Usuario vacio"})
	}

	if MoGeneral.CadenaVacia(Nombre) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Nombre vacio"})
	}
	if MoGeneral.CadenaVacia(Empresa) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Empresa vacio"})
	}
	if MoGeneral.CadenaVacia(Correo) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Correo vacio"})
	} else {
		if !MoGeneral.CorreoValido(Correo) {
			VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Correo no Valido"})
		}
	}
	if MoGeneral.CadenaVacia(Password) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Password vacio"})
	}
	if MoGeneral.CadenaVacia(Confirma) {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "Confirma vacio"})
	}

	if Password != Confirma {
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "No coincide Password con Confirmaci&oactute;n."})
	}

	existeUsuario, err := usuario.QueryFieldValueExist(Usuario, "Usuario", "USUARIOS")
	existeEmpresa, err := usuario.QueryFieldValueExist(Empresa, "Empresa", "USUARIOS")
	existeCorreo, err := usuario.QueryFieldValueExist(Correo, "Correo", "USUARIOS")

	UnicoError := ""
	var usr usuario.Usuario
	usr.Usuario = Usuario
	usr.Nombre = Nombre
	usr.Empresa = Empresa
	usr.Correo = Correo
	usr.Password = Password
	usr.Coleccion = Coleccion
	if existeCorreo || existeEmpresa || existeUsuario || err != nil {
		if err != nil {
			UnicoError = "No se ha podido validar los datos que solicita."
		} else {
			UnicoError = "Los datos proporcionados ya existen en la base de datos."
		}
		VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: UnicoError})
		fmt.Println("No se puede crear el usuario solicitado: ", UnicoError)
	} else {

		fmt.Println("Estamos listos para intentar crear un usuario nuevo!!!", usr)
	}
	VistaRespuesta.Usr = usr
	// fmt.Println(VistaRespuesta.Error)
	if len(VistaRespuesta.Error) > 0 {
		ctx.Render("Login/registro.html", VistaRespuesta)
	} else {
		fmt.Println("Verificar si no hay datos duplicados!!!", usr)
		if !existeUsuario && !existeEmpresa && !existeCorreo {
			fmt.Println("Intentando insertar...")
			if VistaRespuesta.Usr.InsertarUsuarioPostgres() {
				fmt.Println("Insertado!!!", usr)
				ctx.Render("Login/registroconcluido.html", nil)
			} else {
				fmt.Println("Ocurrio un error al insertar!!!", usr)
				VistaRespuesta.Error = append(VistaRespuesta.Error, Errores{Error: "No se ha podido crear el usuario."})
				ctx.Render("Login/registro.html", VistaRespuesta)
			}
		} else {
			ctx.Render("Login/registro.html", VistaRespuesta)
		}

	}

}

//
