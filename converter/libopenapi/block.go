package converter

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/fgm/jastify/converter"
)

type TFBlock struct {
	// The schema map describes the schema implemented by the TF Provider.
	SchemaMap TFSchemaMap
	Arguments []TFArgument
	Blocks    []TFBlock

	// Render-able fields.
	Type   string
	Labels []string
	Body   string
}

func (b *TFBlock) Render(w io.Writer, depth int) error {
	_, err := fmt.Fprintf(w, "%s%s", Indent(depth), b.Type)
	if err != nil {
		return err
	}
	for _, label := range b.Labels {
		_, err := fmt.Fprintf(w, " %q", label)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, " {\n%s", b.Body)
	if err != nil {
		return err
	}
	sort.Slice(b.Arguments, func(i, j int) bool {
		return b.Arguments[i].Name < b.Arguments[j].Name
	})
	for _, arg := range b.Arguments {
		if err := arg.Render(w, depth+1); err != nil {
			return err
		}
	}

	_, err = fmt.Fprintln(w, "}")
	return err
}

func (b *TFBlock) Set(jm converter.Jmap) {
	for jk, v := range jm {
		if jk == "id" {
			continue
		}
		tk := terraformKeyFromJsonKey(jk)
		tvs, known := b.SchemaMap[tk]
		if !known {
			tvs = &schema.Schema{Type: schema.TypeInvalid}
		}
	retry:
		switch tvs.Type {
		case schema.TypeInvalid:
			arg := TFArgument{Name: tk, Value: Unsupported(v)}
			b.Arguments = append(b.Arguments, arg)
		case schema.TypeBool:
			bv, ok := v.(bool)
			if !ok {
				log.Fatalf("key %q (JSON: %q) is not a bool", tk, jk)
			}
			arg := TFArgument{Name: tk, Value: bv, RO: tvs.Computed, Deprecation: tvs.Deprecated}
			b.Arguments = append(b.Arguments, arg)
		// case schema.TypeInt:
		// case schema.TypeFloat:
		case schema.TypeString:
			sv, ok := v.(string)
			if !ok {
				log.Fatalf("key %q (JSON: %q) is not a string", tk, jk)
			}
			arg := TFArgument{Name: tk, Value: sv, RO: tvs.Computed, Deprecation: tvs.Deprecated}
			b.Arguments = append(b.Arguments, arg)

		case schema.TypeList, schema.TypeMap, schema.TypeSet:
			switch t := tvs.Elem.(type) {
			case *schema.Resource:
				// Provide a TFBlock
			case *schema.Schema:
				// Provide a TFArgument
				if reflect.ValueOf(v).Kind() != reflect.Slice {
					log.Fatalf("key %q (JSON: %q) is not a slice", tk, jk)
				}
				arg := TFArgument{Name: tk, Value: v, RO: tvs.Computed, Deprecation: tvs.Deprecated}
				b.Arguments = append(b.Arguments, arg)
			default:
				log.Fatalf("unexpected type %T in schema map for key %q", t, tk)
			}
		default:
			tvs.Type = schema.TypeInvalid
			goto retry
		}
	}
}
