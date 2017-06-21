package catalogoentrada

import (
	"context"
	"fmt"

	"../../Modulos/Conexiones"
)

//CatalogoEntrada Estructura para insertar catalogo en elasticsearch
type CatalogoEntrada struct {
	SKU         string `json:"clave"`
	Descripcion string `json:"descripcion"`
	ClaveSat    string `json:"claveSAT"`
}

//InsertarElastic Funcion para insertar objeto en elasticsearch
func (me *CatalogoEntrada) InsertarElastic(indice string, tipo string) bool {
	estadooperacion := false
	cliente, err := MoConexion.GetClienteElastic()
	var ctx = context.Background()
	if err != nil {
		fmt.Println("Error al obtener el cliente de elasticsearch", err)
	} else {
		if MoConexion.VerificaIndexName(cliente, indice) {
			Put, err := cliente.Index().Index(indice).Type(tipo).Id(me.SKU).BodyJson(me).Do(ctx)
			if err != nil {
				fmt.Println("No se ha agregado la entrada al indice ElasticSearch", err)
			} else {
				fmt.Printf("\nIndexado en el index %s, con type %s\n", Put.Index, Put.Type)
				estadooperacion = true
			}
		} else {
			// Create an index
			fmt.Println("Si el indice no existe se debe crear")
			_, err = cliente.CreateIndex(indice).Do(ctx)
			if err != nil {
				// Handle error
				fmt.Println("No se pudro crear el indice: ", indice, "\nError:", err)

			} else {
				fmt.Println("indice Creado")
				Put, err := cliente.Index().Index(indice).Type(tipo).Id(me.SKU).BodyJson(me).Do(ctx)
				if err != nil {
					fmt.Println("No se ha agregado la entrada al indice ElasticSearch", err)
				} else {
					fmt.Printf("\nIndexado en el index %s, con type %s\n", Put.Index, Put.Type)
					estadooperacion = true
				}
			}

		}
	}
	return estadooperacion
}
