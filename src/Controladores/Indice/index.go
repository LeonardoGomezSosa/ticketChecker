package Indice

import (
	"fmt"

	iris "gopkg.in/kataras/iris.v6"

	"../../Modelos/ListaModel"
	"../../Modulos/Conexiones"
	"github.com/tealeg/xlsx"
)

/*
//IndexGet renderea al index de Almacen
func IndexGet(ctx *iris.Context) {
	fmt.Println("Indice.Index.go: GET")
	ctx.Render("Vistas/index.html", nil)
}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {
	fmt.Println("Indice.Index.go: Post")
	file, headerMultiPartFile, err := ctx.FormFile("uploadfile")
	if err != nil {
		fmt.Println("Error al recibir el archivo...", err)
	} else {
		fmt.Println(file, "\n---------------------------------------------------------------------------\n",
			headerMultiPartFile.Filename, headerMultiPartFile.Header)
		_, err := headerMultiPartFile.Open()
		if err != nil {
			fmt.Println("Error al recibir el archivo...", err)
		} else {
			buf := bytes.NewBuffer(nil)
			_, err = io.Copy(buf, file)
			if err != nil {
				fmt.Println("Error al copiar el archivo...", err)
			} else {
				fmt.Println(buf)
				bytesBuf := buf.Bytes()
				bytesString := string(bytesBuf)
				r := csv.NewReader(strings.NewReader(bytesString))
				//preparar insercion
				indice := "catalogosat"
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
						req := elastic.NewBulkIndexRequest().Index(indice).Type(tipo).Id(elemento.SKU).Doc(elemento)
						bulkRequest = bulkRequest.Add(req)
						fmt.Printf("agregado %v a bulk\n", elemento)
					}
					fmt.Println("Total de operaciones en fila:  ", bulkRequest.NumberOfActions())

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
				} else {
					fmt.Println("El cliente o el indice no existen, imposible insertar objetos...")
				}

			}
		}

	}

	ctx.Render("General/index.html", nil)

}

//PrepararInsercion funcion que devuelve el cliente  y verifica que exista el indice dado
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
*/
//DescargaPost metodo para descargar catalogo en formato CSV
func DescargaPost(ctx *iris.Context) {
	//Cargamos los datos de la base de datos
	//Para CSV
	//var data = [][]string{{"10101514", "Primates"}, {"10101515", "Armadillos"}}
	//Para XLSX
	var data2 []ListaModel.CatalogoFinal
	var rowx ListaModel.CatalogoFinal
	idTablaUsuario := ctx.GetCookie("IDUsuario")
	fmt.Println(idTablaUsuario)
	_, cat := GetCatalogo(idTablaUsuario)
	fmt.Println(cat)

	rowx.Sku = "SP-004578"
	rowx.Descripcion = "Armadillo Silvestre"
	rowx.ClaveSat = "10101514"
	data2 = append(data2, rowx)
	////////////////////////
	fmt.Println("Construyendo archivo XLSX")
	var fileX *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	fileX = xlsx.NewFile()
	sheet, err = fileX.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	for _, val := range cat {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = val.Sku
		cell = row.AddCell()
		cell.Value = val.Descripcion
		cell = row.AddCell()
		cell.Value = val.ClaveSat
	}
	err = fileX.Save("../Public/Data/data.csv")
	err = fileX.Save("../Public/Data/data.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Println("Termina de construir archivo XLSX")
	////////////////////////
	ctx.Render("Vistas/descarga.html", nil)
}

//GetCatalogo devuelve el listado total de operaciones pendientes de pago by @melchormendoza
func GetCatalogo(NombreTabla string) (string, []ListaModel.CatalogoFinal) {
	//BasePosGres := conectaPostgresql("192.168.1.110", "ClasificadorCatalogo", "postgres", "12345")
	BasePosGres, err := MoConexion.ConexionPsql()
	//BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()

	Query := fmt.Sprintf(`SELECT "Sku","Descripcion","ClaveSat" FROM public."%v"`, NombreTabla)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	productos := []ListaModel.CatalogoFinal{}
	for resultSet.Next() {
		producto := ListaModel.CatalogoFinal{}
		resultSet.Scan(&producto.Sku, &producto.Descripcion, &producto.ClaveSat)
		productos = append(productos, producto)
	}
	fmt.Println(productos)
	//BasePosGres.Commit()
	resultSet.Close()
	stmt.Close()
	BasePosGres.Close()
	return "_", productos
}
