package MoConexion

import (
	"database/sql"
	"fmt"

	"../../Modulos/Variables"
	_ "github.com/lib/pq"
)

//DataP es una estructura que contiene los datos de configuración en el archivo cfg
var DataP = MoVar.CargaSeccionCFG(MoVar.SecPsql)

//ConexionPsql abre una conexión a PostgreSql
func ConexionPsql() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DataP.Servidor, DataP.Usuario, DataP.Pass, DataP.NombreBase)
	db, err := sql.Open("postgres", dbinfo)
	return db, err
}

//IniciaSesionEspecificaPsql regresa una base de datos con su sesion especificada en el parametro para commit y rollback
func IniciaSesionEspecificaPsql() (*sql.DB, *sql.Tx, error) {
	Psql, err := ConexionPsql()
	if err != nil {
		return nil, nil, err
	}

	tx, err := Psql.Begin()
	if err != nil {
		return nil, nil, err
	}
	return Psql, tx, nil
}

//InsertaOActualizaRelacion verifica la tabla relacion e inserta el sku en caso de que no se encuentre; caso contrario actualiza la clave del sat
func InsertaOActualizaRelacion(tabla, sku, descripcion, claveSat string) error {
	var SesionPsql *sql.Tx
	var err error
	BasePsql, SesionPsql, err := IniciaSesionEspecificaPsql()
	if err != nil {
		fmt.Println("Errores al conectar con postgres: ", err)
		return err
	}
	BasePsql.Exec("set transaction isolation level serializable")

	Query := fmt.Sprintf(`SELECT "Sku" FROM "%v"  WHERE "Sku" ='%v'`, tabla, sku)
	Elemento, err := BasePsql.Query(Query)
	if err != nil {
		fmt.Println("Error al consultar el sku: ", err, Query)
		return err
	}
	var encontrado bool
	for Elemento.Next() {
		var skuEnc string
		err = Elemento.Scan(&skuEnc)
		if err != nil {
			fmt.Println("Error al consultar el sku: (2)", err)
			return err
		}
		if sku == skuEnc {
			encontrado = true
		}
	}
	if encontrado {
		Query = fmt.Sprintf(`UPDATE  public."%v"  SET  "ClaveSat" = '%v' WHERE "Sku" ='%v'`, tabla, claveSat, sku)
		_, err = SesionPsql.Exec(Query)
		if err != nil {
			fmt.Println("Ha ocurrido un error en la actualizacion", err)
			SesionPsql.Rollback()
			BasePsql.Close()
			return err
		}
	} else {
		query := fmt.Sprintf(`INSERT INTO public."%v" VALUES('%v','%v','%v')`, tabla, sku, descripcion, claveSat)
		_, errsql := SesionPsql.Exec(query)
		if errsql != nil {
			SesionPsql.Rollback()
			BasePsql.Close()
			fmt.Println("Error al insertar el producto")
			fmt.Println(query)
			return err
		}
	}
	SesionPsql.Commit()
	BasePsql.Close()
	return err
}

//CatalogoFinal Estructura Provisional POSTGRESQL/ELASTICSEARCH @melchormendoza
type CatalogoFinal struct {
	Sku         string
	Descripcion string
	ClaveSat    string
}

/*
//GetCatalogo devuelve el listado total de operaciones pendientes de pago by @melchormendoza
func GetCatalogo() (string, []CatalogoFinal) {
	BasePosGres, err := ConexionPsql()

	//BasePosGres := conectaPostgresql("192.168.1.110", "ClasificadorCatalogo", "postgres", "12345")
	//BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()

	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	Query := fmt.Sprintf(`SELECT "Sku","Descripcion","ClaveSat" FROM public."CatalogoFinal"`)
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

	productos := []CatalogoFinal{}
	for resultSet.Next() {
		producto := CatalogoFinal{}
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
*/
//GetCatalogo devuelve el listado total de operaciones pendientes de pago by @melchormendoza
func ConsultaCatalogo(NombreTabla string) (*sql.Rows, error) {
	BasePosGres, err := ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	Query := fmt.Sprintf(`SELECT "Sku","Descripcion","ClaveSat" FROM public."%v"`, NombreTabla)
	stmt, err := BasePosGres.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resultSet.Close()
	stmt.Close()
	BasePosGres.Close()
	return resultSet, err
}
