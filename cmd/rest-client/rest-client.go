package main

// Import resty into your code and refer it as `resty`.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

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

// Twitt es la estructura que encapsula los datos a grabar en la base
type Twitt struct {
	Class           string `json:"@class"`
	ID              string `json:"id"`
	Text            string `json:"text"`
	AuthorID        string `json:"author_id"`
	AuthorName      string `json:"author_name"`
	ConversationID  string `json:"conversation_id"`
	InReplyToUserID string `json:"in_reply_to_user_id"`
	CreatedAt       string `json:"created_at"`
}

/*
Habiéndome enterado que este coso no soporta los fundamentos más básicos de programación
orientada a objetos me veo en la obligación de hacer funciones para los accesos a la base
de datos. Accesos que quedaría mucho mejor encapsulados en una clase.
Pero bueno. Todo rompen. Todo.
*/

// OrientDBQuery is Ejecuta una consulta en la base de datos, pasando usuario y contraseña y
// devuelve el resultado en formato JSON
func OrientDBQuery(dbAcc DBAccess, query string, pretty bool, rawjson bool) (result string, statusCode int, status string) {
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
	var res string
	if !rawjson {
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
		res = gjson.Get(result, "result").String()
	} else {
		res = string(resp.Body())
	}

	return res, resp.StatusCode(), resp.Status()
}

func getRandomBook(dbAcc DBAccess) (result string, statusCode int, status string) {

	return OrientDBQuery(dbAcc, "select expand(getRandomRecord('Libro')) as resultado", true, false)
}

func getRandomBookWithAuthor(dbAcc DBAccess) (result string, statusCode int, status string) {

	return OrientDBQuery(dbAcc, "select getRandomLibro() as resultado", true, false)
}

/*
Deprecated: Función utilizada solamente para depuración y pruebas
*/
func insertTwittDirect(dbAcc DBAccess, t Twitt) /*(result string, statusCode int, status string)*/ {
	urlPost := url.URL{
		Scheme: dbAcc.protocol,
		Host:   dbAcc.host + ":" + dbAcc.port,
		Path:   path.Join("document", dbAcc.database),
	}
	//jsonData := fmt.Sprintf(`{'@class' : 'Twitt','id' : '%s','text' : '%s','author_id' : '%s','author_name' : '%s','conversation_id' : '%s','in_reply_to_user_id' : '%s'}`,
	//	t.id, t.text, t.author_id, t.author_name, t.conversation_id, t.in_reply_to_user_id)
	//jsonData := []byte("{'@class' : 'Twitt','id' : '0'}")

	var jsonData []byte
	jsonData, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(jsonData))

	//jsonData, err := json.Marshal(dbAcc)
	log.Println("urlPost:")
	log.Println(urlPost.String())
	log.Println("jsonData:")
	log.Println(jsonData)

	client := &http.Client{}
	// req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(jsonData))
	req, err := http.NewRequest("POST", urlPost.String(), bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip,deflate")
	req.SetBasicAuth(dbAcc.user, dbAcc.password)
	resp, err := client.Do(req)
	//resp, err := http.Post(urlPost.String(), "application/json; charset=utf-8", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)
	/*
		client := resty.New()
		req := client.R().
			EnableTrace().
			SetHeader("Content-Type", "application/json").
			SetBasicAuth(dbAcc.user, dbAcc.password).
			SetBody(jsonData)
		resp, _ := req.Post(urlPost.String())
		log.Println("Request:\n", req.String())
		log.Println("Response:\n", resp.String())
	*/
	//return resp.String(), resp.StatusCode(), resp.Status()
}

// OrientDBBatch is Ejecuta una batch de instrucciones en la base de datos, pasando usuario y contraseña y
// devuelve el resultado en formato JSON
func OrientDBBatch(dbAcc DBAccess, query string) (result string, statusCode int, status string) {
	// Create a Resty Client
	client := resty.New()
	urlBatch := url.URL{
		Scheme: dbAcc.protocol,
		Host:   dbAcc.host + ":" + dbAcc.port,
		Path:   path.Join("batch", dbAcc.database),
	}
	jsonData := fmt.Sprintf(`{ 
		"transaction" : false,
		"operations" : [{
			"type" : "script",
			"language" : "sql",
			"script" : ["%s"]
			}]
		}`, query)
	log.Println("urlQuery:")
	log.Println(urlBatch.String())
	log.Println("jsonData:")
	log.Println(jsonData)

	resp, _ := client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip,deflate").
		SetBasicAuth(dbAcc.user, dbAcc.password).
		SetBody(jsonData).
		Post(urlBatch.String())

	log.Println("Respuesta:\n", string(resp.Body()))
	return string(resp.Body()), resp.StatusCode(), resp.Status()
}

func insertTwitt(dbAcc DBAccess, t Twitt) (result string, statusCode int, status string) {
	// SIEMPRE TIENE QUE TENER FECHA y HORA.
	// SI NO SE PASA COMO PARTE DEL TWITT SE PONE FECHA Y HORA ACTUAL
	// DE LA MAQUINA DONDE ESTA CORRIENDO EL BOT
	// EL FORMATO DEL STRING ES: yyyy-mm-dd hh:mi:ss (sin milisegundos)
	var datetime string
	if t.CreatedAt != "" {
		datetime = t.CreatedAt
	} else {
		ahora := time.Now()
		datetime = ahora.Format("2006-01-02 15:04:05")
	}
	query := fmt.Sprintf(`
	BEGIN; 
    LET twitt = SELECT from Twitt where id = '%s';
    if ($twitt.size() = 0) {
        CREATE VERTEX Twitt SET
        id = '%s',
        text = '%s',
        author_id = '%s',
		author_name = '%s',
        conversation_id = '%s',
        in_reply_to_user_id = '%s',
		created_at = '%s';
    }
    COMMIT;`, t.ID, t.ID, t.Text, t.AuthorID, t.AuthorName, t.ConversationID, t.InReplyToUserID, datetime)
	//fmt.Println(query)
	return OrientDBBatch(dbAcc, query)
}

func insertTwittRelation(dbAcc DBAccess, idTwittOrigen string, idTwittDestino string, tipoRelacion string) (result string, statusCode int, status string) {
	var tipoEdge = ""
	if tipoRelacion == "replied_to" {
		tipoEdge = "TwittReply"
	} else if tipoRelacion == "quoted" {
		tipoEdge = "TwittCite"
	} else if tipoRelacion == "retweeted" {
		tipoEdge = "TwittRetweet"
	} else {
		tipoEdge = "E"
	}
	query := fmt.Sprintf(`BEGIN; 
    CREATE EDGE %s from (select from Twitt where id = '%s') to (select from Twitt where id = '%s');
    COMMIT;
	`, tipoEdge, idTwittOrigen, idTwittDestino)
	return OrientDBBatch(dbAcc, query)
}

func main() {
	log.Println("Iniciando prueba de REST Client...")

	// Estructura con los parámetros fijos de acceso al servidor
	acc := DBAccess{
		user:     "admin",
		password: "admin",
		protocol: "http",
		host:     "sibila.website",
		//host:     "localhost",
		port:     "2480",
		database: "portico",
	}

	/*
		result, statusCode, status := insertTwittRelation(acc, "1359892733737984002", "1359889625582551041", "quoted")
		fmt.Println("Response Info (insertTwittRelation):")
		fmt.Println(result)
		fmt.Println(statusCode)
		fmt.Println(status)
	*/

	t := Twitt{
		Class:           "Twitt",
		ID:              "0",
		Text:            "Prueba",
		AuthorID:        "0117",
		AuthorName:      "mcasatti",
		ConversationID:  "001",
		InReplyToUserID: "",
		CreatedAt:       "",
	}
	result, statusCode, status := insertTwitt(acc, t)
	fmt.Println("Response Info (insertTwitt):")
	fmt.Println(result)
	fmt.Println(statusCode)
	fmt.Println(status)

	/*
		result, statusCode, status := getRandomBook(acc)
		fmt.Println("Response Info (getRandomBook):")
		fmt.Println(result)
		fmt.Println(statusCode)
		fmt.Println(status)

		fmt.Println(strings.Repeat("=", 80))

		result, statusCode, status = getRandomBookWithAuthor(acc)
		fmt.Println("Response Info (getRandomBookWithAuthor):")
		fmt.Println(result)
		fmt.Println(statusCode)
		fmt.Println(status)

		fmt.Println(strings.Repeat("=", 80))
	*/
	/* result, statusCode, status = OrientDBQuery(acc, "match {class: Libro, as: l} return $elements", true)
	fmt.Println("Response Info (MATCH):")
	fmt.Println(result)
	fmt.Println(statusCode)
	fmt.Println(status) */
}
