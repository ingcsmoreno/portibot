package main

// Import resty into your code and refer it as `resty`.
import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func main() {
	log.Println("Iniciando prueba de REST Client...")

	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get("http://sibila.website:2480/")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())

	log.Println("Fin de la prueba REST Client...")

	log.Println("Iniciando prueba de consulta de OrientDB")
	resp, err = client.R().
		EnableTrace().
		SetBasicAuth("admin", "admin").
		Get("http://sibila.website:2480/query/PPR/sql/select from Concepto")

	fmt.Println("Respuesta:\n", resp)

	log.Println("Fin de prueba de acceso a OrientDB")
}
