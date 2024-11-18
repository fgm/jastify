package legacy

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/sebdah/goldie/v2"

	"github.com/fgm/jastify/converter"
)

//go:embed testdata/monitor.json
var monitorJSON []byte

func Test_generateMonitorTerraformCode(t *testing.T) {
	// TODO Add sad test case.
	for _, test := range [...]struct {
		name    string
		input   []byte
		errorRx string
	}{
		{"monitor", monitorJSON, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			var j = make(converter.Jmap)
			if err := json.Unmarshal(test.input, &j); err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t)
			actual, err := GenerateMonitorTerraformCode(test.name, j)
			switch {
			case test.errorRx == "" && err != nil:
				t.Fatalf("unexpected error %v", err)
			case test.errorRx != "" && err == nil:
				t.Fatalf("unexpected success")
			case test.errorRx != "" && err != nil:
				rx := regexp.MustCompile(test.errorRx)
				if !rx.MatchString(err.Error()) {
					t.Errorf("got error %s but expected to match %s", err, rx.String())
				}
			}
			g.Assert(t, test.name, []byte(actual))
		})
	}
}
