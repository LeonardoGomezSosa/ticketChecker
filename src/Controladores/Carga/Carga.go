package Carga

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"../../Modelos/CargaModel"
	"../../Modelos/CatalogoEntrada"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../Sesiones"
	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/olivere/elastic.v5"
)

//IndexGet renderea al index de Almacen
func IndexGet(ctx *iris.Context) {
	if !sessionUtils.IsStarted(ctx) {
		sessionUtils.DeleteSsn(ctx)
	}

	fmt.Println("Carga.Index.go: GET")
	ctx.Render("Carga/index.html", nil)
}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {

	var send CargaModel.SCargaVista
	//ctx.Redirect("Login/login.html", 300)

	if !sessionUtils.IsStarted(ctx) {
		sessionUtils.DeleteSsn(ctx)
		var Errores []CargaModel.Errores
		var Error CargaModel.Errores
		Error.Error = "La sesion a expirado"
		Errores = append(Errores, Error)
		send.Error = Errores
		_, _, _ = ctx.FormFile("uploadfile")
		ctx.Render("Login/login.html", send)
		fmt.Println("--++--++--++--++--")
		return
	}
	fmt.Println("Carga.Index.go: POST")
	file, headerMultiPartFile, err := ctx.FormFile("uploadfile")
	if err != nil {
		fmt.Println("Error al recibir el archivo...", err)
		send.SEstado = false
		send.SMsj = "Selecionar un archivo csv"
		ctx.Render("Carga/index.html", send)

	} else {
		fmt.Println(file, "\n---------------------------------------------------------------------------\n", headerMultiPartFile.Filename)
		_, err := headerMultiPartFile.Open()
		if err != nil {
			send.SEstado = false
			send.SMsj = "Error al recibir el archivo..." + err.Error()
			ctx.Render("Carga/index.html", send)
		} else {
			buf := bytes.NewBuffer(nil)
			_, err = io.Copy(buf, file)
			if err != nil {
				fmt.Println("Error al copiar el archivo...", err)
				send.SEstado = false
				send.SMsj = "Error al recibir el archivo..." + err.Error()
				ctx.Render("Carga/index.html", send)
			} else {
				if MoGeneral.CSVValido(headerMultiPartFile.Filename) {
					fmt.Println("Si se puede Parsear: ")
					bytesBuf := buf.Bytes()
					bytesString := string(bytesBuf)
					r := csv.NewReader(strings.NewReader(bytesString))
					indice := ctx.GetCookie("ColeccionKsd")
					//preparar insercion
					if indice != "" {

						tipo := "productoservicio"
						cliente, exito := PrepararInsercion(indice)
						bulkRequest := cliente.Bulk()
						if exito {
							var contexto = context.Background()
							fmt.Println("El cliente y el indice existen, se procede a insertar objetos...")
							startTime := time.Now()
							fmt.Println("Inicia insercion: ", startTime.Format(time.Stamp))

							for {
								record, err := r.Read()
								if err == io.EOF {
									fmt.Println("Se ha llegado al fin de archivo...")
									break
								}
								if err != nil {
									fmt.Println("Error al realizar la codificacion del csv...")
								}
								var elemento catalogoentrada.CatalogoEntrada
								elemento.SKU = string(record[0])
								elemento.Descripcion = record[1]
								elemento.ClaveSat = ""
								req := elastic.NewBulkIndexRequest().Index(indice).Type(tipo).Id(elemento.SKU).Doc(elemento)
								bulkRequest = bulkRequest.Add(req)
								fmt.Printf("agregado %v a bulk\n", elemento)
							}

							i := bulkRequest.NumberOfActions()
							send.SMsj = "Numero de campos insertado: " + strconv.Itoa(i)

							bulkResponse, err := bulkRequest.Do(contexto)
							if err != nil {
								fmt.Println(err)
							}
							if bulkResponse != nil {
								fmt.Println("took ", bulkResponse.Took, " milis")
								fmt.Println("items ", bulkResponse.Items, " objetos")
							}

							endTime := time.Now()

							fmt.Println("Inicia insercion: ", endTime.Format(time.Stamp))
							cliente.CloseIndex(indice)

							send.SEstado = true
							send.SMsj = send.SMsj + "  Insercion de catalogo de manera  exitosa"
							ctx.Render("Carga/index.html", send)

						} else {
							send.SEstado = false
							send.SMsj = "No se puede Parsear"
							ctx.Render("Carga/index.html", send)
							fmt.Println("No se puede Parsear")
						}

					} else {
						var Errores []CargaModel.Errores
						var Error CargaModel.Errores
						Error.Error = "La sesion a expirado"
						Errores = append(Errores, Error)
						send.Error = Errores
						ctx.Render("Login/login.html", send)
					}

				} else {
					send.SEstado = false
					send.SMsj = "El formato del archivo es incorreco selecione el archivo en formato csv..."
					ctx.Render("Carga/index.html", send)
					fmt.Println("El formato del archivo es incorreco selecione el archivo en formato csv")
				}

			}
		}

	}
}

//PrepararInsercion funcion que devuelve el cliente  y verifica que exista el indice dado y si no existe lo crea
func PrepararInsercion(indice string) (*elastic.Client, bool) {
	exito := false
	cliente, err := MoConexion.GetClienteElastic()
	var ctx = context.Background()
	if err != nil {
		fmt.Println("Error al obtener el cliente de elasticsearch", err)
		cliente = nil
	} else {
		if !MoConexion.VerificaIndexName(cliente, indice) {
			// Create an index
			fmt.Println("Si el indice no existe se debe crear")
			_, err = cliente.CreateIndex(indice).Do(ctx)
			if err != nil {
				// Handle error
				fmt.Println("No se pudro crear el indice: ", indice, "\nError:", err)
				cliente = nil
			} else {
				fmt.Println("El indice SE ha creado de manera exitosa.")
				exito = true
			}
		} else {
			fmt.Println("El indice existe, no hay nada que hacer, Exito.")
			exito = true
		}

	}
	return cliente, exito
}
