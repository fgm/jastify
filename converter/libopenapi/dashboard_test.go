package converter_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	converter2 "github.com/fgm/jastify/converter"
	"github.com/fgm/jastify/converter/legacy"
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

	sb := strings.Builder{}
	if err := converter.GenerateDashboardTerraformCode(&sb, legacy.ResourceName(jm["title"].(string)), jm); err != nil {
		t.Fatal(err)
	}
	actual := sb.String()
	t.Logf("Actual:\n%s\n", actual)
	if actual != expected {
		t.Fatal(cmp.Diff(actual, expected))
	}
}
