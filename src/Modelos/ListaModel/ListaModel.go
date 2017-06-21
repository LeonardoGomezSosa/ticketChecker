package ListaModel

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	elastic "gopkg.in/olivere/elastic.v5"
)

//ProductosBase objeto que se encarga de tratar los datos obtenidos de elasticsearch
type ProductosBase struct {
	Descripcion        string `json:"descripcion"`
	Unidad             string `json:"unidad"`
	Marca              string `json:"marca"`
	Clave              string `json:"clave"`
	ClaveSAT           string `json:"claveSAT"`
	Descripcionproyser string `json:"descripcionproyser"`
	Idproyser          string `json:"idproyser"`
	Linea              string `json:"linea"`
}

//ProductosSat objeto que se encarga de tratar los datos obtenidos de elasticsearch
type ProductosSat struct {
	Descripcion string `json:"descripcion"`
	Clave       string `json:"claveprodserv"`
}

//ListaVista estructura Auxiliar para enviar listas a la vista
type ListaVista struct {
	ID          string
	SEstado     bool
	SMsj        string
	SIhtml      template.HTML
	SCabecera   template.HTML
	SBody       template.HTML
	SPaginacion template.HTML
	SGrupo      template.HTML
}

//GeneraTemplatesBusquedaBase crea templates de tabla de búsqueda
func GeneraTemplatesBusquedaBase(Productos []ProductosBase, tipo int) (string, string) {
	cuerpo := ``
	cabecera := ``

	switch tipo {
	case 1:

		cabecera = `<tr>
			<th>#</th>
			<th class="text-center">
				(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">
			</th>
			<th><small>SKU</small></th>
			<th><small>Descripción</small></th>			
			</tr>`

		for k, v := range Productos {
			cuerpo += `<tr>`

			cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`

			cuerpo += `<td class="text-center"><input class="ProdExtraido" type="checkbox" id="` + v.Clave + `" value="` + v.Clave + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
			cuerpo += `<td><small>` + v.Clave + `</small></td>`

			cuerpo += `<td><small>` + v.Descripcion + `</small></td>`
			cuerpo += `</tr>`

		}
	case 2:
		cabecera = `<tr>
			<th>#</th>
			<th class="text-center">
				(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" >
			</th>
			<th><small>SKU</small></th>
			<th><small>Código Sat</small></th>
			<th><small>Descripción</small></th>			
			</tr>`

		for k, v := range Productos {
			cuerpo += `<tr>`

			cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`

			cuerpo += `<td class="text-center"><input class="ProdExtraido" type="checkbox" id="` + v.Clave + `" value="` + v.Clave + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
			
cuerpo += `<td><small>` + v.Clave + `</small></td>`
			if v.ClaveSAT != "" {
				cuerpo += `<td><small>` + v.ClaveSAT + `</small></td>`
			} else {
				cuerpo += `<td><small>Sin Código</small></td>`
			}
			cuerpo += `<td><small>` + v.Descripcion + `</small></td>`
			cuerpo += `</tr>`

		}
	case 3:
		cabecera = `<tr>
			<th>#</th>
			<th class="text-center">

				(Todos)<input  class="ProdAsignar" type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">
			</th>
			<th><small>SKU</small></th>
			<th><small>Código Sat</small></th>
			<th><small>Descripción</small></th>			
			</tr>`

		for k, v := range Productos {
			cuerpo += `<tr>`

			cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`

			cuerpo += `<td class="text-center"><input class="ProdExtraido" type="checkbox" id="` + v.Clave + `" value="` + v.Clave + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
			cuerpo += `<td><small>` + v.Clave + `</small></td>`
			if v.ClaveSAT != "" {
				cuerpo += `<td><small>` + v.ClaveSAT + `</small></td>`
			} else {
				cuerpo += `<td><small>Sin Código</small></td>`
			}
			cuerpo += `<td><small>` + v.Descripcion + `</small></td>`
			cuerpo += `</tr>`

		}
	default:
		cabecera = `<tr>
			<th>#</th>
			<th class="text-center">
				(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">
			</th>
			<th><small>SKU</small></th>
			<th><small>Descripción</small></th>			
			</tr>`

		for k, v := range Productos {
			cuerpo += `<tr>`

			cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`

			cuerpo += `<td class="text-center"><input  class="ProdExtraido" type="checkbox" id="` + v.Clave + `" value="` + v.Clave + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
			cuerpo += `<td><small>` + v.Clave + `</small></td>`

			cuerpo += `<td><small>` + v.Descripcion + `</small></td>`
			cuerpo += `</tr>`

		}

	}

	return cabecera, cuerpo
}

//GeneraTemplatesBusquedaSat crea templates de tabla de búsqueda
func GeneraTemplatesBusquedaSat(Productos []ProductosSat) (string, string) {

	cuerpo := ``

	cabecera := `<tr>
			<th>#</th>
			<th><small>Selecciona</small></th>
			<th><small>Código</small></th>
			<th><small>Descripción</small></th>			
			</tr>`

	for k, v := range Productos {
		cuerpo += `<tr>`

		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`

		cuerpo += `<td class="text-center"><input class="ProdAsignar"  type="radio" id="` + v.Clave + `" value="` + v.Clave + `" name="mismo" onClick="AsignarClaveSat(this.id)"></td>`
		cuerpo += `<td><small>` + v.Clave + `</small></td>`
		cuerpo += `<td><small>` + strings.ToUpper(v.Descripcion) + `</small></td>`
		cuerpo += `</tr>`

	}

	return cabecera, cuerpo
}

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(index, tipo, texto string, ccode int) *[]ProductosBase {
	var Productos []ProductosBase

	queryTilde := elastic.NewQueryStringQuery(texto)
	queryTilde = queryTilde.Field("descripcion")
	queryTilde = queryTilde.Field("marca")

	var docs *elastic.SearchResult
	var err error

	docs, err = MoConexion.BuscaElasticAvanzada(index, tipo, queryTilde)
	if err != nil {
		fmt.Println(err)
		return &Productos
	}

	var Producto ProductosBase

	if docs.Hits.TotalHits > 0 {
		for _, v := range docs.Hits.Hits {
			Producto = ProductosBase{}
			err := json.Unmarshal(*v.Source, &Producto)
			if err == nil {

				if ccode == 1 {
					if Producto.ClaveSAT == "" {
						Productos = append(Productos, Producto)
					}
				} else if ccode == 2 {
					if Producto.ClaveSAT != "" {
						Productos = append(Productos, Producto)
					}
				} else if ccode == 3 {
					Productos = append(Productos, Producto)
				} else {
					if Producto.ClaveSAT == "" {
						Productos = append(Productos, Producto)
					}
				}

			} else {
				fmt.Println("Unmarchall Error: ", err)
			}
		}
	}
	return &Productos

}

//CheckConCodigo regresa sólo los que no tienen código
func CheckConCodigo(Productos *[]ProductosBase) *[]ProductosBase {
	var Products []ProductosBase

	for _, v := range *Productos {
		if v.ClaveSAT != "" {
			Products = append(Products, v)
		}

	}
	return &Products
}

//CheckSinCodigo regresa sólo los que no tienen código
func CheckSinCodigo(Productos *[]ProductosBase) *[]ProductosBase {
	var Products []ProductosBase

	for _, v := range *Productos {
		if v.ClaveSAT == "" {
			Products = append(Products, v)
		}

	}
	return &Products
}

//BuscarEnElasticSat busca el texto solicitado en los campos solicitados
func BuscarEnElasticSat(index, tipo, texto string) *[]ProductosSat {
	var Productos []ProductosSat

	queryTilde := elastic.NewQueryStringQuery(texto)
	queryTilde = queryTilde.Field("descripcion")
	queryTilde = queryTilde.Field("claveprodserv")

	docs, err := MoConexion.BuscaElasticAvanzada(index, tipo, queryTilde)
	if err != nil {
		fmt.Println(err)
		return &Productos
	}

	var Producto ProductosSat

	if docs.Hits.TotalHits > 0 {
		for _, v := range docs.Hits.Hits {
			Producto = ProductosSat{}
			err := json.Unmarshal(*v.Source, &Producto)
			if err == nil {
				Productos = append(Productos, Producto)
			} else {
				fmt.Println("Unmarchall Error: ", err)
			}
		}
	}

	return &Productos
}

//AsignarRelacionClaveSku actualiza en elastic la relacion de la clave sat y los skus
func AsignarRelacionClaveSku(IndiceElastic, tipoElastic, NombreTabla, claveSat string, arraySkus []string) error {
	var err error
	for _, value := range arraySkus {
		var Producto ProductosBase
		resultado, err := MoConexion.ConsultaElastic(MoVar.Tipo, value)
		if err != nil {
			fmt.Println("Error al consultar el id en elastic")
			return err
		} else {
			valores := resultado.Source
			err := json.Unmarshal(*valores, &Producto)
			if err != nil {
				log.Fatalln("error:", err)
				return err
			}
			Producto.ClaveSAT = claveSat
			err = MoConexion.ActualizaElastic(IndiceElastic, tipoElastic, value, Producto)
			if err != nil {
				fmt.Println("Ha ocurrido un error al actualizar el objeto en elastic: ", err)
				return err
			} else {
				err := MoConexion.InsertaOActualizaRelacion(NombreTabla, value, Producto.Descripcion, claveSat)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

//CatalogoFinal Estructura Provisional POSTGRESQL/ELASTICSEARCH @melchormendoza
type CatalogoFinal struct {
	Sku         string `bson:"Sku"`
	Descripcion string `bson:"Descripcion"`
	ClaveSat    string `bson:"ClaveSat"`
}

//GetCatalogo devuelve el listado total de operaciones pendientes de pago by @melchormendoza
func GetCatalogo(NombreTabla string) (string, []CatalogoFinal) {

	resultSet, err := MoConexion.ConsultaCatalogo(NombreTabla)
	fmt.Println(err)
	productos := []CatalogoFinal{}
	for resultSet.Next() {
		producto := CatalogoFinal{}
		resultSet.Scan(&producto.Sku, &producto.Descripcion, &producto.ClaveSat)
		productos = append(productos, producto)
	}
	fmt.Println(productos)

	return "_", productos
}
