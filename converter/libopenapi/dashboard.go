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

	"github.com/fgm/jastify/cmd/apischema/libopenapi/index/schemaindexer"
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
	b.Set(converter.Path{}, jData)
	if err := b.Render(w, 0); err != nil {
		return err
	}
	return nil
}

func terraformKeyFromJsonKey(path converter.Path, jk string, props any) (tk string, discriminator string, selector string) {
	type Unit struct{}
	var unit Unit

	plain := map[string]string{
		"formulas":           "formula",
		"formula":            "formula_expression",
		"layout":             "widget_layout",
		"queries":            "query",
		"requests":           "request",
		"template_variables": "template_variable",
		"widgets":            "widget",
	}
	if tk, ok := plain[jk]; ok {
		return tk, discriminator, ""
	}

	// Some keys need a resolution process, e.g. OneOf like widget.definition.
	// Keys which have no plain conversion and no resolvable conversion pass through for robustness.
	if _, ok := map[string]Unit{
		"definition":  unit,
		"data_source": unit,
	}[jk]; !ok {
		return jk, discriminator, ""
	}
	// To resolve, we need to first find the property allowing us to discriminate
	// between the alternatives in the oneOf.
	// TODO handle anyOf, allOf too.

	// This starts by ensuring we have something from which to fetch the discriminator,
	// when we have it. Ensure that first, as it is cheaper than the schemaindexer.Index() call.
	jm, ok := props.(converter.Jmap)
	if !ok {
		log.Fatalf("expected %q props to be a Jmap but got %#v", jk, props)
	}
	path = path.Push(jk)
	s := schemaindexer.Index(path.Slice())
	ref := s.ParentProxy.GetReference()
	discriminator, ok = schemaindexer.Discriminators[ref]
	if !ok {
		log.Fatalf("discriminator for %q not found", ref)
	}
	v, ok := jm[discriminator]
	if !ok {
		log.Fatalf("discriminator %q not found for %q in %v", discriminator, jk, jm)
	}
	selector, ok = v.(string)
	if !ok {
		log.Fatalf("discriminator %q not a string", discriminator)
	}
	// FIXME probably not always "_definition". Check for other oneOfs beyond widgets.
	tk = selector + "_definition"
	return tk, discriminator, selector
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
