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

func terraformKeyFromJsonKey(path converter.Path, jk string) string {
	type Unit struct{}
	var unit Unit

	plain := map[string]string{
		"layout":             "widget_layout",
		"template_variables": "template_variable",
		"widgets":            "widget",
	}
	if tk, ok := plain[jk]; ok {
		return tk
	}

	// Some keys need a resolution process, e.g. OneOf like widget.definition.
	// Keys which have no plain conversion and no resolvable conversion pass through for robustness.
	if _, ok := map[string]Unit{
		"definition": unit,
	}[jk]; !ok {
		return jk
	}
	path = path.Push(jk)
	s := schemaindexer.Index(path.Slice())
	ref := s.ParentProxy.GetReference()
	discriminator, ok := schemaindexer.Discriminators[ref]
	if !ok {
		log.Fatalf("discriminator for %q not found", ref)
	}
	_ = schemaindexer.Index(path.Push(discriminator).Slice())
	log.Printf("[INFO] discriminator for %q found as %q", ref, discriminator)
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
