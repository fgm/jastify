package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"

	"github.com/fgm/jastify/exp"
)

const (
	Indent       = "  "
	MaxIndent    = 30
	MaxStringLen = 60 // In runes
)

func indent(level int) string {
	if level > MaxIndent {
		log.Fatal("Too deep indentation")
	}
	return fmt.Sprintf("%3d", level) + strings.Repeat(Indent, level)
}

func dumpProperties(level int, m *orderedmap.Map[string, *base.SchemaProxy]) {
	level++
	if m == nil || m.Len() == 0 {
		return
	}
	fmt.Printf("%sproperties:\n", indent(level))
	for k := range m.KeysFromOldest() {
		v := m.GetOrZero(k)
		if v == nil {
			return
		}
		dump(level, v, "")
	}
}

func dumpBool(level int, name string, b *bool) {
	if b == nil {
		return
	}
	level++
	fmt.Printf("%s%s: %t\n", indent(level), name, *b)
}

func dumpFloat(level int, name string, f *float64) {
	if f == nil {
		return
	}
	level++
	fmt.Printf("%s%s: %f\n", indent(level), name, *f)
}

func dumpInt(level int, name string, n *int64) {
	if n == nil {
		return
	}
	level++
	fmt.Printf("%s%s: %d\n", indent(level), name, *n)
}

func dumpString(level int, name, s string) {
	if s == "" {
		return
	}
	level++
	if utf8.RuneCountInString(s) > MaxStringLen {
		rs := []rune(s)
		rs = rs[:MaxStringLen]
		s = string(rs) + "..."
	}
	s = regexp.MustCompile("\n").ReplaceAllString(s, "\\n")
	fmt.Printf("%s%s: %s\n", indent(level), name, s)
}

func dumpItems(level int, v *base.DynamicValue[*base.SchemaProxy, bool]) {
	if v == nil {
		return
	}
	var x any
	if v.IsA() {
		x = v.A
	} else if v.IsB() {
		x = v.B
	} else {
		log.Fatalf("DynamicValue has N = %d", v.N)
	}
	sp, ok := x.(*base.SchemaProxy)
	if !ok {
		log.Fatalf("DynamicValue is not a SchemaProxy: %Td", x)
	}
	dump(level, sp, "")
}

func dumpVariants(level int, name string, ssp []*base.SchemaProxy) {
	if len(ssp) == 0 {
		return
	}
	level++
	fmt.Printf("%s%s:\n", indent(level), name)
	for _, sp := range ssp {
		name := sp.Schema().Title
		if name == "" {
			if sp.IsReference() {
				name = sp.GetReference()
			} else {
				log.Fatalf("Schema reference has no name of reference: %#v", sp)
			}
		}
		dump(level, sp, name)
	}
}

func dumpType(level int, typ []string) {
	switch len(typ) {
	case 0:
		return
	case 1:
		fmt.Printf("%stype: %s\n", indent(level+1), typ[0])
	default:
		fmt.Printf("%stype: %v\n", indent(level+1), typ)
	}
}

func dump(level int, sp *base.SchemaProxy, nameOverride string) {
	level++
	kn := sp.GetSchemaKeyNode()
	name := kn.Value
	s, err := sp.BuildSchema()
	if err != nil {
		log.Fatalf("building schema at level %d for %v: %v", level, name, err)
	}
	if nameOverride != "" {
		if name == "" {
			name = nameOverride
		} else {
			name = nameOverride + "(" + name + ")"
		}
	}
	fmt.Printf("%s%s:\n", indent(level), name)
	dumpString(level, "title", s.Title)
	dumpString(level, "description", s.Description)
	dumpString(level, "anchor", s.Anchor)
	dumpString(level, "format", s.Format)
	dumpString(level, "pattern", s.Pattern)
	dumpString(level, "schema_type_ref", s.SchemaTypeRef)

	dumpType(level, s.Type)
	dumpProperties(level, s.Properties)
	dumpItems(level, s.Items)

	dumpVariants(level, "one_of", s.OneOf)
	dumpVariants(level, "all_of", s.AllOf)
	dumpVariants(level, "any_of", s.AnyOf)

	dumpInt(level, "max_contains", s.MaxContains)
	dumpInt(level, "min_contains", s.MinContains)
	dumpInt(level, "min_items", s.MinItems)
	dumpInt(level, "maxItems", s.MaxItems)
	dumpInt(level, "min_length", s.MinLength)
	dumpInt(level, "max_length", s.MaxLength)
	dumpInt(level, "min_properties", s.MinProperties)
	dumpInt(level, "max_properties", s.MaxProperties)
	dumpFloat(level, "minimum", s.Minimum)
	dumpFloat(level, "maximum", s.Maximum)

	dumpBool(level, "deprecated", s.Deprecated)
	dumpBool(level, "read_only", s.ReadOnly)
	dumpBool(level, "nullable", s.Nullable)
	dumpBool(level, "unique_items", s.UniqueItems)
	dumpBool(level, "write_only", s.WriteOnly)
}

func main() {
	rawDoc, err := libopenapi.NewDocument(exp.OpenAPI)
	if err != nil {
		log.Fatal(err)
	}
	doc, errs := rawDoc.BuildV3Model()
	if len(errs) != 0 {
		for _, err := range errs {
			log.Println(err)
		}
		os.Exit(1)
	}
	const name = "Dashboard"
	dsp := doc.Model.Components.Schemas.Value(name)
	dump(0, dsp, "")
}
