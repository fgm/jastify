package converter_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	converter2 "github.com/fgm/jastify/converter"
	"github.com/fgm/jastify/converter/legacy"
	converter "github.com/fgm/jastify/converter/libopenapi"
)

func TestGenerateDashboardTerraformCode(t *testing.T) {
	t.Parallel()
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

func TestTFBlock_ResolveConflicts(t *testing.T) {
	t.Parallel()
	s := map[string]*schema.Schema{
		"da":   {ConflictsWith: []string{"nda", "ndab"}, Deprecated: "da"},
		"db":   {ConflictsWith: []string{"ndb", "ndab"}, Deprecated: "db"},
		"ndab": {ConflictsWith: []string{"da", "db"}, Deprecated: ""},
	}
	tests := []struct {
		name              string
		args              []converter.TFArgument // XXX Add same logic for blocks
		expectErr         string
		expectedRemaining []string
	}{
		{
			name: "no conflicts",
			args: []converter.TFArgument{
				{Name: "da"},
				{Name: "db"},
			},
			expectedRemaining: []string{"da", "db"},
		},
		{
			name: "single deprecated vs non-deprecated",
			args: []converter.TFArgument{
				{Name: "db"},
				{Name: "ndab"},
			},
			expectedRemaining: []string{"ndab"},
		},
		{
			name: "multiple deprecated vs non-deprecated",
			args: []converter.TFArgument{
				{Name: "da"},
				{Name: "db"},
				{Name: "ndab"},
			},
			expectedRemaining: []string{"ndab"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			b := converter.TFBlock{
				SchemaMap: s,
				Arguments: test.args,
			}
			b.ResolveConflicts()
			actual := make([]string, 0, len(b.Arguments))
			for _, arg := range b.Arguments {
				actual = append(actual, arg.Name)
			}
			if !cmp.Equal(actual, test.expectedRemaining) {
				t.Fatal(cmp.Diff(actual, test.expectedRemaining))
			}
		})
	}
}
