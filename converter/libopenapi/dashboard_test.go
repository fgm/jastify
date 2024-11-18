package converter_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	converter2 "github.com/fgm/jastify/converter"
	converter "github.com/fgm/jastify/converter/libopenapi"
)

func TestGenerateDashboardTerraformCode(t *testing.T) {
	bs, err := os.ReadFile("../legacy/testdata/dashboard-good.json")
	if err != nil {
		t.Fatal(err)
	}
	jm := make(converter2.Jmap)
	if err := json.Unmarshal(bs, &jm); err != nil {
		t.Fatal(err)
	}
	bs, err = os.ReadFile("../legacy/testdata/dashboard-good.golden")
	if err != nil {
		t.Fatal(err)
	}
	expected := string(bs)

	actual, err := converter.GenerateDashboardTerraformCode("dashboard_1", jm)
	t.Logf("Actual:\n%s\n", actual)
	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatal(cmp.Diff(actual, expected))
	}
}
