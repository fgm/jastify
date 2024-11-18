package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fgm/jastify/converter"
	"github.com/fgm/jastify/converter/legacy"
	loa "github.com/fgm/jastify/converter/libopenapi"
)

func main() {
	var (
		jsonData     []byte
		err          error
		resourceName string
	)

	switch len(os.Args) {
	case 1:
		// No args: read from stdin
		jsonData, err = io.ReadAll(os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
			os.Exit(1)
		}
	case 2:
		// One arg: it's a file, read from it.
		jsonData, err = os.ReadFile(os.Args[1])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading file:", err)
			os.Exit(1)
		}
		resourceName = os.Args[1]
	default:
		name := filepath.Base(os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Usage: \n%s < somefile.json\nor\nconvert somefile.json\n", name)
		os.Exit(1)
	}

	var parsedJson converter.Jmap
	if err := json.Unmarshal(jsonData, &parsedJson); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}

	var tf string

	if _, exists := parsedJson["name"]; exists {
		if resourceName == "" {
			resourceName = "monitor_1"
		}
		resourceName = legacy.ResourceName(resourceName)
		tf = converter.Must(loa.GenerateMonitorTerraformCode(resourceName, parsedJson))
	} else {
		if resourceName == "" {
			resourceName = "dashboard_1"
		}
		resourceName = legacy.ResourceName(resourceName)
		tf = converter.Must(loa.GenerateDashboardTerraformCode(resourceName, parsedJson))
	}

	fmt.Print(tf)
}
