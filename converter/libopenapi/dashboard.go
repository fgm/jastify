// Package converter implements dashboard and monitor resource conversion
// from JSON to Terraform in HCL.
//
// "Apps Hungarian" is used here in this way:
//   - any identifier prefixed with "j" or "J" belongs to the JSON realm.
//   - any identifier prefixed with "t(f)" or "T(F)" belongs to the Terraform realm.
//
// Exception: the JMap and JMaps types may appear in both realms.
package converter

import (
	"io"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/terraform-providers/terraform-provider-datadog/datadog"

	"github.com/fgm/jastify/converter"
)

const (
	IndentSize            = 2
	ReadOnlyPrefix        = "// (readonly) "
	UnsupportedPrefix     = "// (unsupported) "
	ResourceTypeDashboard = "datadog_dashboard"
	ResourceTypeMonitor   = "datadog_monitor"
)

var (
	// tfs is the Terraform schema for a dashboard.
	trm = datadog.Provider().ResourcesMap
)

func init() {
	if trm == nil {
		log.Fatalf("datadog_provider has no resource named %q", ResourceTypeDashboard)
	}
}

type TFSchemaMap map[string]*schema.Schema

type Unsupported struct {
	wrapped any
}

func Indent(level int) string {
	return strings.Repeat(" ", level*IndentSize)
}

func GenerateDashboardTerraformCode(w io.Writer, tfResourceName string, jData converter.Jmap) error {
	b := TFBlock{
		SchemaMap: trm[ResourceTypeDashboard].SchemaMap(),
		Type:      "resource",
		Labels:    []string{ResourceTypeDashboard, tfResourceName},
	}
	b.Set(jData)
	if err := b.Render(w, 0); err != nil {
		return err
	}
	return nil
}

func terraformKeyFromJsonKey(jk string) string {
	mapping := map[string]string{
		"template_variables": "template_variable",
		"widgets":            "widget",
	}
	if tk, ok := mapping[jk]; ok {
		return tk
	}
	return jk
}

// Returning nil means any value is valid.
func GetValidValues(cur *base.Schema) []string {
	switch {
	case cur == nil:
		log.Fatal("expected non-nil schema, got nil")
	case len(cur.Type) > 1:
		log.Fatalf("expected 1 type, got %v", cur.Type)
	case cur.Type[0] == "string":
		if cur.Enum == nil {
			return nil
		}
		sl := make([]string, 0, len(cur.Enum))
		for _, yn := range cur.Enum {
			sl = append(sl, yn.Value)
		}
		return sl
	default:
		log.Fatal(converter.ErrUnimplemented)
	}
	return nil
}
