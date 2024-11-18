// This package load the DataDog API V1 schema and extract the models
//
// Find the file at https://github.com/DataDog/datadog-api-client-go/blob/master/.generator/schemas/v1/openapi.yaml
package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-json"

	"github.com/fgm/jastify/exp"
)

type Schema openapi3.Schema
type SchemaRef openapi3.SchemaRef

// MarshalJSON returns the JSON encoding of SchemaRef.
func (x *SchemaRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

func (x *SchemaRef) MarshalYAML() (any, error) {
	sr := (*openapi3.SchemaRef)(x)
	if sr.Value != nil && sr.Ref != "" {
		log.Printf("Describing %s", sr.Ref)
		sr.Ref = "" // Force dumping schema instead of factoring ref.
	}
	for k, v := range sr.Value.Properties {
		dsr := (*SchemaRef)(v)
		if dsr.Value != nil {
			dsr.Ref = "" // Force dumping schema instead of factoring ref.
		}
		sr.Value.Properties[k] = (*openapi3.SchemaRef)(v)
	}
	return sr.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of SchemaRef.
func (x *Schema) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

func (x *Schema) MarshalYAML() (any, error) {
	s := (*openapi3.Schema)(x)
	m1, err := s.MarshalYAML()
	if err != nil {
		return nil, err
	}
	m, ok := m1.(map[string]any)
	if !ok {
		return nil, fmt.Errorf(`"%v" is not map[string]any`, m)
	}
	for k, v := range m {
		sr, ok := v.(*openapi3.SchemaRef)
		if !ok {
			continue
		}
		dsr := (*SchemaRef)(sr)
		if dsr.Value != nil {
			dsr.Ref = "" // Force dumping schema instead of factoring ref.
		}
		m[k] = dsr
	}
	return m, nil
}

func getRoot() *openapi3.Schema {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(exp.OpenAPI)
	if err != nil {
		log.Fatal(err)
	}
	schemas := doc.Components.Schemas
	sr, ok := schemas["Dashboard"]
	if !ok {
		log.Fatal("dashboard schema not found")
	}
	dashboardSchema := sr.Value
	return dashboardSchema
}

func main() {
	root := getRoot()
	ws, ok := root.Properties["widgets"]
	if !ok {
		log.Fatal("widgets schema not found")
	}
	dws := SchemaRef(*ws)
	ds := (*Schema)(dws.Value)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(ds)
	if err != nil {
		log.Fatal(err)
	}

}
