package MoVar

import (
	config "github.com/robfig/config"
)

const (

	//############# ARCHIVOS LOCALES ######################################

	//FileConfigName contiene el nombre del archivo CFG
	FileConfigName = "Acfg.cfg"

	//################ SECCIONES CFG  ######################################

	//SecDefault nombre de la seccion default del servidor en CFG
	SecDefault = "DEFAULT"
	//SecMongo nombre de la seccion de mongo en CFG
	SecMongo = "CONFIG_DB_MONGO"
	//SecPsql nombre de la seccion de postgresql en cfg
	SecPsql = "CONFIG_DB_POSTGRES"
	//SecElastic nombre de la seccion de postgresql en cfg
	SecElastic = "CONFIG_DB_ELASTIC"

	//############# COLUMNAS DE ALMACEN PSQL ######################################

	//############# COLECCIONES MONGO ######################################

	//################# DATOS ELASTIC ######################################
	//Elastic---------------> Documento
	Index        = "clasificadorvisorus"
	Tipo         = "clasificadorvisorus"
	IndexElastic = "CatalogoProductos"
	//ColeccionOperacion Nombre de IndexElasticla coleccion que almacena las operaciones Generales
	ProductosServicios = "ProductosServicios"

	//ColeccionSurtidor nombre de la coleccion de Surtidor en mongo
	ColeccionSurtidor = "Surtidor"

	//Elastic---------------> Surtidor

	//TipoSurtidor tipo a manejar en elastic
	TipoSurtidor = "Surtidor"

	ColeccionExpresion = "Expresion"

	//Elastic---------------> Expresion

	//TipoExpresion tipo a manejar en elastic
	TipoExpresion = "Expresion"

	//TipoOperacion establecido temporalmente para minimizar los errores, posteriormente se obtendrá de algun catalogo
	TipoOperacion = "Venta"
)

//DataCfg estructura de datos del entorno
type DataCfg struct {
	BaseURL    string
	Servidor   string
	Puerto     string
	Usuario    string
	Pass       string
	Protocolo  string
	NombreBase string
}

//#################<Funciones Generales>#######################################

//CargaSeccionCFG rellena los datos de la seccion a utilizar
func CargaSeccionCFG(seccion string) DataCfg {
	var d DataCfg
	var FileConfig, err = config.ReadDefault(FileConfigName)
	if err == nil {
		if FileConfig.HasOption(seccion, "baseurl") {
			d.BaseURL, _ = FileConfig.String(seccion, "baseurl")
		}
		if FileConfig.HasOption(seccion, "servidor") {
			d.Servidor, _ = FileConfig.String(seccion, "servidor")
		}
		if FileConfig.HasOption(seccion, "puerto") {
			d.Puerto, _ = FileConfig.String(seccion, "puerto")
		}
		if FileConfig.HasOption(seccion, "usuario") {
			d.Usuario, _ = FileConfig.String(seccion, "usuario")
		}
		if FileConfig.HasOption(seccion, "pass") {
			d.Pass, _ = FileConfig.String(seccion, "pass")
		}
		if FileConfig.HasOption(seccion, "protocolo") {
			d.Protocolo, _ = FileConfig.String(seccion, "protocolo")
		}
		if FileConfig.HasOption(seccion, "base") {
			d.NombreBase, _ = FileConfig.String(seccion, "base")
		}
	}
	return d
}
