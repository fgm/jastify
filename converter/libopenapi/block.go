package converter

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"slices"
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
			arg := TFArgument{Name: tk, Value: Unsupported{v}}
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
				tvs.Type = schema.TypeInvalid
				goto retry
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
	b.ResolveConflicts()
}

// ResolveConflicts removes conflicting arguments, under multiple assumptions:
//
//   - only Arguments can be deprecated: needs to be reconsidered in the general case,
//     but in the Datadog provider, this only applies to the agentRule in securityMonitoringRule,
//     so we do not care in that version.
//   - provider schema is correct, not conflicting two deprecated declarations,
//     nor two non-deprecated declarations
//
// TODO verify the second assumption.
// TODO use the result to expose the removed keys in the generated output.
func (b *TFBlock) ResolveConflicts() []string {
	removed := make([]string, 0)
	for _, arg := range b.Arguments {
		s, ok := b.SchemaMap[arg.Name]
		if !ok {
			log.Fatalf("argument key %q is not in schema. Should not happen", arg.Name)
		}

		// Only delete deprecated keys when they conflict with a non-deprecated
		// key that is actually in use on the block.
		if s.Deprecated != "" {
			for _, c := range s.ConflictsWith {
				if slices.ContainsFunc(b.Arguments, func(arg TFArgument) bool {
					return arg.Name == c
				}) {
					removed = append(removed, arg.Name)
				}
			}
		}
	}
	// These two steps should be redundant, as arguments should be unique,
	// but this makes the code more resilient.
	slices.Sort(removed)
	slices.Compact(removed)

	for _, doit := range removed {
		b.Arguments = slices.DeleteFunc(b.Arguments, func(arg TFArgument) bool {
			return arg.Name == doit
		})
	}
	return removed
}
