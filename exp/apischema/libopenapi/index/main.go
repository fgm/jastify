// This examples shows how to access the schema for a given path into
// a JSON dashboard definition, in this case ["widgets", 0, "definition"],
// resolving to the #/components/schemas/WidgetDefinition OneOf.
//
// Do not be confused: it does not analyse a given document, but just digs into
// the OpenAPI schema for the description of a property.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fgm/jastify/exp/apischema/libopenapi/index/schemaindexer"
)

func main() {
	var path = []string{
		"widgets",
		"definition",
		"query_value",
		"requests",
		"queries",
		"metrics",
		"query",
	}
	flag.Usage = func() {
		cmd := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [selector]*\n  With no selector, a default one is used.\n", cmd)
	}
	flag.Parse()
	if args := os.Args[1:]; len(args) > 0 {
		path = args
	} else {
		log.Printf("Demo: using %q example path", strings.Join(path, " "))
	}
	cur := schemaindexer.Index(path)
	fmt.Println(cur.Description)
}
