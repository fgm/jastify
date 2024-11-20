package legacy

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/fgm/jastify/converter"
)

var (
	nameExclusionTx = regexp.MustCompile(`[^-[:word:]]+`)
)

// Generates an assignment string for a key-value pair
func AssignmentString(key string, value any) string {
	if value == nil {
		return ""
	}
	displayValue := LiteralString(value)
	return fmt.Sprintf("%s\t= %s\n", key, displayValue)
}

// Creates a block with a name and converted contents
func block(name string, contents converter.Jmap, converter func(string, any) string) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s {\n", name))
	keys := make([]string, 0, len(contents))
	for k := range contents {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		result.WriteString(converter(k, contents[k]))
	}
	result.WriteString("}\n")
	return result.String()
}

// Generates a list of blocks
func blockList(array converter.Jmaps, blockName string, contentConverter func(string, any) string) string {
	var result strings.Builder
	for _, elem := range array {
		result.WriteString(block(blockName, elem, contentConverter))
	}
	return result.String()
}

// Converts a key-value pair using a definition set
func convertFromDefinition(definitionSet map[string]stringFunc, name string, v any) (string, error) {
	if converter, exists := definitionSet[name]; exists {
		return converter(v), nil
	}
	return "", fmt.Errorf("can't convert key '%s' with value %#v", name, v)
}

func ConvertSlice[T any](v []T) string {
	var result strings.Builder
	result.WriteString("[")
	max := len(v) - 1
	for i, elem := range v {
		result.WriteString(LiteralString(elem))
		if i != max {
			result.WriteString(",")
		}
	}
	result.WriteString("]")
	return result.String()
}

// Converts a string or list of strings value to a literal string representation
func LiteralString(value any) string {
	switch v := value.(type) {
	case string:
		if strings.Contains(value.(string), "\n") {
			return fmt.Sprintf("<<EOF\n%s\nEOF", v)
		}
		return fmt.Sprintf("\"%s\"", v)
	case []any:
		return ConvertSlice(v)
	case []string:
		return ConvertSlice(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Maps over contents and applies a converter function
func mapContents(contents converter.Jmap, converter func(string, any) string) string {
	var result strings.Builder
	keys := make([]string, 0, len(contents))
	for k := range contents {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		result.WriteString(converter(k, contents[k]))
	}
	return result.String()
}

// Creates a query block with a name and converted contents
func queryBlock(name string, contents converter.Jmap, converter func(string, any) string) string {
	return fmt.Sprintf("query {\n  %s {\n%s}\n}\n", name, mapContents(contents, converter))
}

// Generates a list of query blocks
func queryBlockList(array converter.Jmaps, contentConverter func(string, any) string) string {
	var result []string
	for _, elem := range array {
		result = append(result, queryBlock("metric_query", elem, contentConverter))
	}
	return strings.Join(result, "")
}

// ResourceName converts any string to a version suitable for a Terraform resource name.
func ResourceName(s string) string {
	return nameExclusionTx.ReplaceAllString(s, "_")
}
