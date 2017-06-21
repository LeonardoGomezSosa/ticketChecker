package Recuperar

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"../../Modelos/RecuperarModel"
	"../../Modulos/General"

	"../../Modelos/Usuario"
	"../Sesiones"
	iris "gopkg.in/kataras/iris.v6"
)

type Envio struct {
	Correo   string
	Password string
	Servidor string
	Puerto   string
}

// Envios slice de clase envio
type Envios []Envio

// ListaDistribucion lista de correos origen de distribucion
var ListaDistribucion = Envios{
	{Correo: "test.miderp@hotmail.com", Password: "D3m0st3n3s", Servidor: "smtp-mail.outlook.com", Puerto: "587"},
	{Correo: "seguimientoaclientes2@outlook.com", Password: "A2345678@", Servidor: "smtp-mail.outlook.com", Puerto: "587"},
	{Correo: "seguimientoaclientes3@outlook.com", Password: "A2345678@", Servidor: "smtp-mail.outlook.com", Puerto: "587"},
}

//IndexGet renderea al index de Almacen
func IndexGet(ctx *iris.Context) {
	fmt.Println("Recuperar.Recuperar.go: GET")
	if sessionUtils.IsStarted(ctx) {
		fmt.Println("Se redirige al inicio si ya existe una cookie")
		ctx.Redirect("/Login", 301)
	} else {
		ctx.Render("Recuperar/Recuperar.html", nil)
	}

}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	var r RecuperarModel.SRecuperarVista
	var usr usuario.Usuario
	fmt.Println("Recuperar.Recuperar.go: POST")
	val := MoGeneral.LimpiarCadena(ctx.FormValue("Usuario"))
	if val != "" {
		if MoGeneral.CorreoValido(val) {
			//verifica si usuario existe
			fmt.Println("Verificar Correo")
			r.SEstado, usr = usuario.QueryUsrExist(val, "Correo", "USUARIOS")
		} else {
			//verifica si correo existe
			r.SEstado, usr = usuario.QueryUsrExist(val, "Usuario", "USUARIOS")
		}
		if r.SEstado {
			pass := RandStringRunes(8)
			oldPass := usr.Password
			fmt.Println("pass:= ", pass)
			if usr.ActualizaPassUsuario(pass) {
				mensaje := fmt.Sprintf(`
					<h1>Matching Catalogo - Catalogo SAT<br><small>Notificacion de cambio de contrase&ntilde;a</small></h1>
					<p>
						Estimado usuario <em>%v</em>,
						Se el informa que su contraseña ha sido reestablecida a solicitud de la p&aacute;gina de recuperaci&oacute;n de contrase&ntilde;a de nuesto sistema.
					</p>
					<p>
						Su nueva contraseña es: <h3><em>%v</em></h3>
					</p>
					<p>
						Agradecemos su preferencia
					</p>
					<h3>Equipo Kore <small>Ingenieria de negocios.</small></h3>
				`, usr.Usuario, pass)
				fmt.Println(mensaje)
				if SendEmail(usr.Correo, mensaje) {
					r.SMsj = "Clave actualizada."
				} else {
					usr.ActualizaPassUsuario(oldPass)
					r.SEstado = false
					r.Error = append(r.Error, RecuperarModel.ErroresR{Error: "No fue posible enviar el correo, intente mas tarde o contacte a su proveedor de sistema."})
				}

			} else {
				r.SEstado = false
				r.Error = append(r.Error, RecuperarModel.ErroresR{Error: "No fue posible cambiar la contraseña, intente mas tarde."})
			}
			fmt.Println(usr, pass)
		} else {
			r.SEstado = false
			r.Error = append(r.Error, RecuperarModel.ErroresR{Error: "No se Encontraron los datos solicitados. Verifique sus datos o intente mas tarde."})
		}
	} else {
		r.SEstado = false
		r.Error = append(r.Error, RecuperarModel.ErroresR{Error: "Proporcione sus datos."})
	}

	fmt.Println(r.SEstado)

	//    si no existe indicar que no se encuentra el usuario
	ctx.Render("Recuperar/Recuperar.html", r)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SendEmail(mail string, mensaje string) bool {

	if mail != "" {
		body := mensaje
		turno := 0
		elemento := ListaDistribucion[turno]
		Destino := fmt.Sprintf("%v", mail)
		uname := elemento.Correo
		pass := elemento.Password
		serv := elemento.Servidor

		servPort := fmt.Sprintf("%v:%v", serv, elemento.Puerto)
		fmt.Println("Enviado desde : ", servPort)
		msg2 := "To: " + Destino + "\r\n" +
			"Content-type: text/html" + "\r\n" +
			"Subject: Notificacion modificacion de contrase&ntilde;a." + "\r\n\r\n" +
			body + "\r\n"

		auth2 := smtp.PlainAuth(
			"",
			uname,
			pass,
			serv,
		)
		fmt.Println("Auth:  ", auth2, ".  Server: ", " Server: ", servPort, " Destino: ", Destino)

		err := smtp.SendMail(
			servPort,
			auth2,
			uname,
			[]string{"test.miderp@hotmail.com", Destino},
			[]byte(msg2),
		)
		if err != nil {
			fmt.Println("Ha ocurrido algun error: ", err)
			fmt.Println("Correo no ha sido enviado ")
			return false

		}
		fmt.Println("Solicitud ha sido realizada ")
		fmt.Println("Correo ha sido enviado ")
		return true
	}
	fmt.Println("No existe ningun correo.")
	return false

}
