// This package load the DataDog API V1 schema and extract the models
//
// Find the file at https://github.com/DataDog/datadog-api-client-go/blob/master/.generator/schemas/v1/openapi.yaml
package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-json"
)

//go:embed openapi.yaml
var openapi []byte

func main() {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(openapi)
	if err != nil {
		log.Fatal(err)
	}
	schemas := doc.Components.Schemas
	sr, ok := schemas["Dashboard"]
	if !ok {
		log.Fatal("dashboard schema not found")
	}
	dashboardSchema := sr.Value
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(dashboardSchema)

}
