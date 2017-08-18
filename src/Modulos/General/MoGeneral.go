package MoGeneral

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//########### GENERALES #######################################

//EstaVacio verifica si un objeto está vacío o no
func EstaVacio(object interface{}) bool {
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	if reflect.ValueOf(object).Kind() == reflect.Struct {
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}
	return false
}

//ConstruirCadenas recibe un texto y regresa dos que se utilizarán para buscar en elastic
func ConstruirCadenas(texto string) (string, string) {

	var palabras = []string{}
	var final = []string{}
	var final2 = []string{}
	var cadenafinal string
	var cadenafinal2 string

	nuevacadena := strings.Replace(texto, "/", "\\/", -1)
	nuevacadena = strings.Replace(nuevacadena, "~", "\\~", -1)
	nuevacadena = strings.Replace(nuevacadena, "^", "\\^", -1)
	nuevacadena = strings.Replace(nuevacadena, "+", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "[", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "]", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "{", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "}", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "(", "\\(", -1)
	nuevacadena = strings.Replace(nuevacadena, ")", "\\)", -1)
	nuevacadena = strings.Replace(nuevacadena, "|", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "=", "", -1)
	nuevacadena = strings.Replace(nuevacadena, ">", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "<", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "!", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "&", "", -1)

	palabras = strings.Split(nuevacadena, " ")

	for _, valor := range palabras {
		if valor != "" {
			palabrita := valor + "~2"
			final = append(final, palabrita)
		}
	}

	for _, valor := range palabras {
		if valor != "" {
			palabrita := `+` + `"` + valor + `"`
			final2 = append(final2, palabrita)
		}
	}

	for _, value := range final {
		cadenafinal = cadenafinal + " " + value
	}

	for _, value := range final2 {
		cadenafinal2 = cadenafinal2 + " " + value
	}

	fmt.Println("Primer Cadena: ", cadenafinal)
	fmt.Println("Segunda Cadena: ", cadenafinal2)
	return cadenafinal, cadenafinal2
}

//ValidaCadenas recibe un texto y regresa dos que se utilizarán para buscar en elastic
func ValidaCadenas(texto string) (string, string) {

	var palabras = []string{}
	var final = []string{}
	var final2 = []string{}
	var cadenafinal string
	var cadenafinal2 string

	nuevacadena := strings.Replace(texto, "/", "\\/", -1)
	nuevacadena = strings.Replace(nuevacadena, "~", "\\~", -1)
	nuevacadena = strings.Replace(nuevacadena, "^", "\\^", -1)
	nuevacadena = strings.Replace(nuevacadena, "+", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "[", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "]", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "{", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "}", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "(", "\\(", -1)
	nuevacadena = strings.Replace(nuevacadena, ")", "\\)", -1)
	nuevacadena = strings.Replace(nuevacadena, "|", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "=", "", -1)
	nuevacadena = strings.Replace(nuevacadena, ">", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "<", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "!", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "&", "", -1)

	palabras = strings.Split(nuevacadena, " ")

	for _, valor := range palabras {
		if valor != "" {
			palabrita := valor + "~2"
			final = append(final, palabrita)
		}
	}

	for _, valor := range palabras {
		if valor != "" {
			palabrita := `+` + `"` + valor + `"`
			final2 = append(final2, palabrita)
		}
	}

	for _, value := range final {
		cadenafinal = cadenafinal + " " + value
	}

	for _, value := range final2 {
		cadenafinal2 = cadenafinal2 + " " + value
	}

	fmt.Println("Primer Cadena: ", cadenafinal)
	fmt.Println("Segunda Cadena: ", cadenafinal2)
	return cadenafinal, cadenafinal2
}

//Totalpaginas calcula el número de paginaciones de acuerdo al número
// de resultados encontrados y los que se quieren mostrar en la página.
func Totalpaginas(numeroRegistros int, limitePorPagina int) int {
	NumPagina := float32(numeroRegistros) / float32(limitePorPagina)
	NumPagina2 := int(NumPagina)
	if NumPagina > float32(NumPagina2) {
		NumPagina2++
	}
	return NumPagina2
}

//ConstruirPaginacion construtye la paginación en formato html para usarse en la página
func ConstruirPaginacion(paginasTotales int, pag int) string {
	var lt int
	var rt int

	lt = 1
	rt = paginasTotales

	if pag > 2 {
		lt = pag - 1
	}
	if paginasTotales > pag {
		rt = pag + 1
	}

	if pag > paginasTotales {
		pag = paginasTotales
	}
	//inicio
	var templateP string

	templateP = `

    <div class="input-group col-md-8">
      
	  <span class="input-group-btn" onclick="BuscaPagina(1)">
        <button class="btn btn-primary" type="button"><span aria-hidden="true">&laquo;</span></button>
      </span>
      
	  <span class="input-group-btn" onclick="BuscaPagina(` + strconv.Itoa(lt) + `)">
        <button class="btn btn-secondary" type="button"><span aria-hidden="true">&lt;</span></button>
      </span>
 

      <input class="form-control" onkeypress="checkSubmit(this)" id="num" onchange="BuscaPagina(this.value)" value="` + strconv.Itoa(pag) + `" type="number" name="num" min="1" max="` + strconv.Itoa(paginasTotales) + `">
 
      <span class="input-group-btn">
        <button class="btn btn-default" type="button"> de ` + strconv.Itoa(paginasTotales) + `</button>
      </span>

    
      <span class="input-group-btn"  onclick="BuscaPagina(` + strconv.Itoa(rt) + `)">
        <button class="btn btn-secondary" type="button"><span aria-hidden="true">&gt;</span></button>
      </span>
      <span class="input-group-btn" onclick="BuscaPagina(` + strconv.Itoa(paginasTotales) + `)">
        <button class="btn btn-primary" type="button"><span aria-hidden="true">&raquo;</span></button>
      </span>
    </div>
`
	return templateP
}

//ConstruirPaginacion2 construtye la paginación en formato html para usarse en la página
func ConstruirPaginacion2(paginasTotales int, pag int) string {
	var lt int
	var rt int

	lt = 1
	rt = paginasTotales

	if pag > 2 {
		lt = pag - 1
	}
	if paginasTotales > pag {
		rt = pag + 1
	}

	//inicio
	var templateP string

	templateP = `

    <div class="input-group col-md-8">
      
	  <span class="input-group-btn" onclick="BuscaPagina(1)">
        <button class="btn btn-primary" type="button"><span aria-hidden="true">&laquo;</span></button>
      </span>
      
	  <span class="input-group-btn" onclick="BuscaPagina(` + strconv.Itoa(lt) + `)">
        <button class="btn btn-secondary" type="button"><span aria-hidden="true">&lt;</span></button>
      </span>
 

      <input class="form-control" onkeypress="checkSubmit(this)" id="num" onchange="BuscaPagina(this.value)" value="` + strconv.Itoa(pag) + `" type="number" name="paginasat" min="1" max="` + strconv.Itoa(paginasTotales) + `">
 
      <span class="input-group-btn">
        <button class="btn btn-default" type="button"> de ` + strconv.Itoa(paginasTotales) + `</button>
      </span>

    
      <span class="input-group-btn"  onclick="BuscaPagina(` + strconv.Itoa(rt) + `)">
        <button class="btn btn-secondary" type="button"><span aria-hidden="true">&gt;</span></button>
      </span>
      <span class="input-group-btn" onclick="BuscaPagina(` + strconv.Itoa(paginasTotales) + `)">
        <button class="btn btn-primary" type="button"><span aria-hidden="true">&raquo;</span></button>
      </span>
    </div>
`
	return templateP
}

//EliminarEspaciosInicioFinal Elimina los espacios en blanco Al inicio y final de una cadena:
//recibe cadena, regresa cadena limpia de espacios al inicio o final o "" si solo contiene espacios
func EliminarEspaciosInicioFinal(cadena string) string {
	var cadenalimpia string
	cadenalimpia = cadena
	re := regexp.MustCompile("(^\\s+|\\s+$)")
	cadenalimpia = re.ReplaceAllString(cadenalimpia, "")
	return cadenalimpia
}

//EliminarMultiplesEspaciosIntermedios Elimina los espacios en blanco de una cadena:
//recibe cadena, regresa cadena limpia  si solo contiene espacios
func EliminarMultiplesEspaciosIntermedios(cadena string) string {
	var cadenalimpia string
	cadenalimpia = cadena
	re := regexp.MustCompile("[\\s]+")
	cadenalimpia = re.ReplaceAllString(cadenalimpia, " ")
	return cadenalimpia
}

//LimpiarCadena Elimina los espacios en blanco de una cadena:
//recibe cadena, regresa cadena limpia o "" si solo contiene espacios
func LimpiarCadena(cadena string) string {
	var cadenalimpia string
	cadenalimpia = EliminarMultiplesEspaciosIntermedios(cadena)
	cadenalimpia = EliminarEspaciosInicioFinal(cadenalimpia)
	return cadenalimpia
}

//ValidaCadenaExpresion Sirve para reconocer si una cadena valida  :
//recibe cadena, patron,  regresa true si la cadena coincide con el patron, false en otro caso
func ValidaCadenaExpresion(cadena string, patron string) bool {
	re := regexp.MustCompile("^" + patron + "$")
	return re.MatchString(cadena)
}

//RFCValido Sirve para reconocer si un RFC es Válido:
//recibe cadena, regresa true si es un RFC válido, false en otro caso
func RFCValido(rfc string) bool {
	re := regexp.MustCompile("^([a-zA-Z]{3}|[a-zA-Z]{4})\\d{6}[a-zA-Z0-9]{3}$")
	return re.MatchString(rfc)
}

//CPValido Sirve para reconocer si un CP es Válido:
//recibe cadena, regresa true si es un CP válido, false en otro caso
func CPValido(CP string) bool {
	re := regexp.MustCompile("^[0-9]{5}$")
	return re.MatchString(CP)
}

//TelOCelValido Sirve para reconocer si un Telefono o Celular es Válido:
//recibe cadena, regresa true si es un el valor recibido es un conjunto de 10 digitos, false en otro caso
func TelOCelValido(CP string) bool {
	re := regexp.MustCompile("^([0-9]{10}$")
	return re.MatchString(CP)
}

//CSVValido Sirve para reconocer si un Telefono o Celular es Válido:
//recibe cadena, regresa true si es un el valor recibido es un conjunto de 10 digitos, false en otro caso
func CSVValido(archivo string) bool {
	re := regexp.MustCompile("^.*(\\.csv)$")
	return re.MatchString(archivo)
}

//CorreoValido Sirve para reconocer si un email Válido:
//recibe cadena, regresa true si es un el valor recibido es un email(fair enough), false en otro caso
func CorreoValido(correo string) bool {
	re := regexp.MustCompile("^(\\w[-._\\w]*\\w@\\w[-._\\w]*\\w\\.\\w{1,3})$")
	return re.MatchString(correo)
}

//CadenaVacia Sirve para reconocer si un RFC es Válido:
//recibe cadena, regresa true si es un RFC válido, false en otro caso
func CadenaVacia(cadena string) bool {
	if cadena == "" {
		return true
	}
	return false
}

//CargaComboMostrarEnIndex carga las opciones de mostrar en el index
func CargaComboMostrarEnIndex(Muestra int) string {
	var Cantidades = []int{5, 10, 15, 20}
	templ := ``

	for _, v := range Cantidades {
		if Muestra == v {
			templ += `<option value="` + strconv.Itoa(v) + `" selected>` + strconv.Itoa(v) + `</option>`
		} else {
			templ += `<option value="` + strconv.Itoa(v) + `">` + strconv.Itoa(v) + `</option>`
		}
	}
	return templ
}
