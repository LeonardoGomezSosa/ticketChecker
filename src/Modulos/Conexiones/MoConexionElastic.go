package MoConexion

import (
	"context"
	"fmt"

	"../../Modulos/Variables"
	elastic "gopkg.in/olivere/elastic.v5"
)

//DataE es una estructura que contiene los datos de configuración en el archivo cfg
var DataE = MoVar.CargaSeccionCFG(MoVar.SecElastic)

//ctx Contexto
var ctx = context.Background()

//GetClienteElastic crea un nuevo cliente a elastic
func GetClienteElastic() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(DataE.BaseURL))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return client, nil
}

//VerificaIndex verifica que el indice exista en elastic
func VerificaIndex(client *elastic.Client) bool {
	exists, err := client.IndexExists(MoVar.Index).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if !exists {
		fmt.Printf("\n El Indice: %s no existe ", MoVar.Index)
		return false
	}
	return true
}

//VerificaIndexName verifica que el indice exista en elastic
func VerificaIndexName(client *elastic.Client, name string) bool {
	exists, err := client.IndexExists(name).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if !exists {
		fmt.Printf("\n El Indice: %s no existe ", name)
		return false
	}
	return true
}

//VerificaIndexName2 verifica que el indice exista en elastic
func VerificaIndexName2(client *elastic.Client, name string) error {
	exists, err := client.IndexExists(name).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !exists {
		fmt.Printf("\n El Indice: %s no existe ", name)
		return err
	}
	return nil
}

//InsertaElastic inserta un articulo de minisuper en elastic
func InsertaElastic(Type string, ID string, Data interface{}) bool {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer client.Stop()
	if VerificaIndex(client) {
		Put, err := client.Index().Index(MoVar.Index).Type(Type).Id(ID).BodyJson(Data).Do(ctx)
		if err != nil {
			fmt.Println(err)
			return false
		}
		fmt.Printf("\nIndexado en el index %s, con type %s\n", Put.Index, Put.Type)
		return true
	}
	return false
}

//DeleteElastic elimina un docuemnto de elastic por ID
func DeleteElastic(Type string, ID string) bool {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer client.Stop()

	_, err = client.Get().Index(MoVar.Index).Type(Type).Id(ID).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = client.Delete().Index(MoVar.Index).Type(Type).Id(ID).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//ConsultaElastic elimina un docuemnto de elastic por ID
func ConsultaElastic(Type string, ID string) (*elastic.GetResult, error) {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer client.Stop()

	dato, err := client.Get().Index(MoVar.Index).Type(Type).Id(ID).Do(ctx)
	if err != nil {
		fmt.Println(err)
	}

	return dato, err
}

//BuscaElastic busca documentos por un texto dado
func BuscaElastic(Index, Type string, consulta *elastic.QueryStringQuery) (*elastic.SearchResult, bool) {
	client, err := GetClienteElastic()
	if err != nil {
		return nil, false
	}
	defer client.Stop()
	var docs *elastic.SearchResult

	if VerificaIndexName(client, Index) {
		docs, err = client.Search().Index(Index).Type(Type).Query(consulta).Do(ctx)
		if err != nil {
			fmt.Println(err)
			return nil, false
		}
		return docs, true
	}
	return docs, false
}

//BuscaElasticAvanzada busca documentos por un texto dado
func BuscaElasticAvanzada(Index, Type string, consulta *elastic.QueryStringQuery) (*elastic.SearchResult, error) {
	client, err := GetClienteElastic()
	if err != nil {
		return nil, err
	}
	defer client.Stop()
	var docs *elastic.SearchResult

	err1 := VerificaIndexName2(client, Index)
	if err1 != nil {
		return nil, err1
	}

	docs, err2 := client.Search().Index(Index).Type(Type).Query(consulta).From(0).Size(1000).Do(ctx)
	if err2 != nil {
		return nil, err2
	}

	return docs, nil
}

//FlushElastic hace flush a determinado index de elastic
func FlushElastic() {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.Flush().Index(MoVar.Index).Do(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

//ActualizaElastic  actualiza correctamente un documento en elasticsearch
func ActualizaElastic(Index, Type string, ID string, Data interface{}) error {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println("Error al obtener el cliente elasticSearch: ", err)
		return err
	}
	_, err = client.Update().Index(Index).Type(Type).Id(ID).Doc(Data).DetectNoop(true).Do(context.TODO())
	if err != nil {
		fmt.Println("Error al Actualizar en elasticSearch", err)
		return err
	}
	defer client.Stop()
	return err
}
