package main

// Import resty into your code and refer it as `resty`.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// DBAccess is Estructura que engloba los parámetros para utilizar en el acceso a una base de
// datos por medio de API REST. Una vez creada, se utiliza en todas las llamadas
// a funciones de acceso a la base de datos
type DBAccess struct {
	user     string
	password string
	protocol string
	host     string
	port     string
	database string
}

/*
Habiéndome enterado que este coso no soporta los fundamentos más básicos de programación
orientada a objetos me veo en la obligación de hacer funciones para los accesos a la base
de datos. Accesos que quedaría mucho mejor encapsulados en una clase.
Pero bueno. Todo rompen. Todo.
*/

// OrientDBQuery is Ejecuta una consulta en la base de datos, pasando usuario y contraseña y
// devuelve el resultado en formato JSON
func OrientDBQuery(dbAcc DBAccess, query string, pretty bool) (result string, statusCode int, status string) {
	// Create a Resty Client
	client := resty.New()
	urlQuery := url.URL{
		Scheme: dbAcc.protocol,
		Host:   dbAcc.host + ":" + dbAcc.port,
		Path:   path.Join("query", dbAcc.database, "sql", query),
	}
	log.Println("urlQuery:")
	log.Println(urlQuery.String())
	resp, _ := client.R().
		EnableTrace().
		SetBasicAuth(dbAcc.user, dbAcc.password).
		Get(urlQuery.String())

	//log.Println("Respuesta:\n", string(resp.Body()))

	if pretty {
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, resp.Body(), "", "\t")
		if error != nil {
			log.Println("JSON parse error: ", error)
		}
		result = string(prettyJSON.Bytes())
	} else {
		result = string(resp.Body())
	}
	// Antes de devolver result, extraer especificamente el campo 'result' que contiene los datos
	// El otro campo es EXPLAIN que contiene el plan de ejecucion de la consulta.
	res := gjson.Get(result, "result")

	return res.String(), resp.StatusCode(), resp.Status()
}

func main() {
	log.Println("Iniciando prueba de REST Client...")

	// Estructura con los parámetros fijos de acceso al servidor
	acc := DBAccess{
		user:     "admin",
		password: "admin",
		protocol: "http",
		host:     "sibila.website",
		port:     "2480",
		database: "portico",
	}

	result, statusCode, status := OrientDBQuery(acc, "select from Pelicula", true)
	fmt.Println("Response Info (SELECT):")
	fmt.Println(result)
	fmt.Println(statusCode)
	fmt.Println(status)

	fmt.Println(strings.Repeat("=", 80))

	result, statusCode, status = OrientDBQuery(acc, "match {class: Libro, as: l} return $elements", true)
	fmt.Println("Response Info (MATCH):")
	fmt.Println(result)
	fmt.Println(statusCode)
	fmt.Println(status)
}
