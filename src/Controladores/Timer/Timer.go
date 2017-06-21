package Timer

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/dbatbold/beep"
	"github.com/rakyll/portmidi"

	"../../Modelos/Usuario"
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
	fmt.Println("Timer.Timer.go: GET")
	if sessionUtils.IsStarted(ctx) {
		ctx.Render("Timer/Timer.html", nil)
	} else {
		ctx.Redirect("/Login", 301)
	}

}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("Timer.Timer.go: POST")
	portmidi.Initialize()

	err := ctx.Request.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	Entrada := ctx.FormValue("Entrada")
	// Entrada := "TS-001"
	// Entrada = MoGeneral.EliminarEspaciosInicioFinal(Entrada)

	fmt.Println("Entrada: ", Entrada)
	music := beep.NewMusic("") // output can be file as "music.wav"
	volume := 100

	if err := beep.OpenSoundDevice("default"); err != nil {
		log.Fatal(err)
	}
	if err := beep.InitSoundDevice(); err != nil {
		log.Fatal(err)
	}
	beep.PrintSheet = true
	defer beep.CloseSoundDevice()

	musicScore := `
        VP SA8 SR9
        A9HRDE cc DScszs|DEc 
    `
	if Entrada == "" {
		fmt.Println("Viene vacio")

	} else {
		if MoGeneral.ValidaCadenaExpresion(Entrada, "^(SURTIDOR)(\\d{1}|\\d{2})$") {
			fmt.Println("Surtidor")
			musicScore = `
        VP SA8 SR9
        A9HRDE cc DScszs|DEc 
    `
		} else if MoGeneral.ValidaCadenaExpresion(Entrada, "^((TR)|(TS))-\\d{4}$") {
			fmt.Println("Ticket")
			musicScore = `
        VP SA8 SR9
       A3HLDE [n z,    |cHRq HLz, |[n
    `
		} else {
			fmt.Println("No Pela")
		}

	}

	reader := bufio.NewReader(strings.NewReader(musicScore))
	go music.Play(reader, volume)
	music.Wait()
	beep.FlushSoundBuffer()

	ctx.Render("Timer/Timer.html", nil)

}
