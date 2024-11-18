package legacy

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/sebdah/goldie/v2"

	"github.com/fgm/jastify/converter"
)

//go:embed testdata/dashboard-good.json
var dashGoodJSON []byte

//go:embed testdata/dashboard-bad.json
var dashBadJSON []byte

//go:embed testdata/dashboard-minimal.json
var dashMiniJSON []byte

//go:embed testdata/screenboard.json
var screenboardJSON []byte

//go:embed testdata/timeboard.json
var timeboardJSON []byte

func Test_generateDashboardTerraformCode(t *testing.T) {
	for _, test := range [...]struct {
		name    string
		input   []byte
		errorRx string
	}{
		{"dashboard-good", dashGoodJSON, ""},
		{"dashboard-minimal", dashMiniJSON, ""},
		{"dashboard-bad", dashBadJSON, "^can't convert key .* with value \".*\"$"},
		{"screenboard", screenboardJSON, ""},
		{"timeboard", timeboardJSON, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			var j = make(converter.Jmap)
			if err := json.Unmarshal(test.input, &j); err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t)
			actual, err := GenerateDashboardTerraformCode(test.name, j)
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
