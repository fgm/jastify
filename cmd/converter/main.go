package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fgm/jastify/converter"
	legacy "github.com/fgm/jastify/converter/legacy"
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
		_, _ = fmt.Fprintln(os.Stderr, "Usage: \nconvert < somefile.json\nor\nconvert somefile.json")
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
		tf = converter.Must(legacy.GenerateMonitorTerraformCode(resourceName, parsedJson))
	} else {
		if resourceName == "" {
			resourceName = "dashboard_1"
		}
		resourceName = legacy.ResourceName(resourceName)
		tf = converter.Must(legacy.GenerateDashboardTerraformCode(resourceName, parsedJson))
	}

	fmt.Print(tf)
}
