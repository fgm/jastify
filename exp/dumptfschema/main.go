package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-datadog/datadog"
)

const (
	Limit        = 1e6
	ResourceType = "datadog_dashboard"
)

var (
	count       int
	showCounter bool
)

func indent(level int) string {
	return strings.Repeat("  ", level)
}
func counter() string {
	if !showCounter {
		return ""
	}
	return fmt.Sprintf("%s", counter())
}

func showType(tv schema.ValueType) string {
	switch tv {
	case schema.TypeBool:
		return "bool"
	case schema.TypeInt:
		return "int"
	case schema.TypeFloat:
		return "float"
	case schema.TypeString:
		return "string"
	case schema.TypeList:
		return "list"
	case schema.TypeMap:
		return "map"
	case schema.TypeSet:
		return "set"
	default:
		return tv.String()
	}
}

func describeResource(level int, name string, res *schema.Resource, collection string) {
	if count > Limit {
		return
	}
	count++
	if showCounter {
	}
	if res.Schema == nil {
		res.Schema = res.SchemaFunc()
	}
	if collection != "" {
		fmt.Printf("%s%s- %s (%s[resource]): %s\n", counter(), indent(level), name, collection, res.Description)
	} else {
		fmt.Printf("%s%s- %s (resource): %s\n", counter(), indent(level), name, res.Description)
	}
	keys := make([]string, 0, len(res.Schema))
	for key := range res.Schema {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		describeSchema(level+1, key, res.Schema[key], "")
	}
}

func describeSchema(level int, name string, sch *schema.Schema, collection string) {
	if count > Limit {
		return
	}
	count++
	switch sch.Type {
	case schema.TypeInvalid:
		log.Fatalf("Unexpected invalid type on %s at level %d", name, level)
	case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
		if collection == "" {
			fmt.Printf("%s%s- %s: %s\n", counter(), indent(level), name, showType(sch.Type))
		} else {
			fmt.Printf("%s%s- %s: %s[%s]\n", counter(), indent(level), name, collection, showType(sch.Type))
		}

	case schema.TypeList, schema.TypeMap, schema.TypeSet:
		switch e := sch.Elem.(type) {
		case nil:
			log.Fatalf("Collection type %s without an element type on %s at level %d", showType(sch.Type), name, level)
		case *schema.Resource:
			describeResource(level, name, e, showType(sch.Type))
		case *schema.Schema:
			collection += fmt.Sprintf("%s(min=%d, max=%d)", showType(sch.Type), sch.MinItems, sch.MaxItems)
			describeSchema(level, name, e, collection)
		}
	}
}

func main() {
	p := datadog.Provider()
	ds := p.ResourcesMap[ResourceType]
	describeResource(0, ResourceType, ds, "")
}
