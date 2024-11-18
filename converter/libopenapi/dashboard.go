package converter

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"

	"github.com/fgm/jastify/cmd/apischema/libopenapi/index/schemaindexer"
	"github.com/fgm/jastify/converter"
	"github.com/fgm/jastify/converter/legacy"
)

func GenerateDashboardTerraformCode(resourceName string, data converter.Jmap) (string, error) {
	sb := strings.Builder{}
	if _, err := sb.WriteString(generate([]string{}, resourceName, data)); err != nil {
		log.Fatal(err)
	}
	return sb.String(), nil
}

func generate(path []string, key string, data converter.Jmap) string {
	var out string
	for k, v := range data {
		sv, ok := v.(string)
		if !ok {
			out += fmt.Sprintf("// Skipping %q, type %T unsupported yet\n", k, v)
			continue // FIXME
		}
		cur := schemaindexer.Index(append(path, k))
		valid := getValidValues(cur)
		if valid != nil && !slices.Contains(valid, sv) {
			log.Fatalf("%s has key %q, not part of %v", key, sv, valid)
		}
		out += legacy.AssignmentString(k, sv)
	}
	return out
}

// Returning nil means any value is valid.
func getValidValues(cur *base.Schema) []string {
	var tl = len(cur.Type)

	switch {
	case cur == nil:
		log.Fatal("expected non-nil schema, got nil")
	case tl > 1:
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
