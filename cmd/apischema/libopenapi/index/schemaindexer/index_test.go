package schemaindexer_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fgm/jastify/cmd/apischema/libopenapi/index/schemaindexer"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		paths    []string
		expected string
	}{
		{[]string{}, "A dashboard is Datadog’s tool for visually tracking, analyzing, and displaying\nkey performance metrics, which enable you to monitor the health of your infrastructure."},
		{[]string{"widgets"}, "List of widgets to display on the dashboard."},
		{[]string{"widgets", "layout"}, "The layout for a widget on a `free` or **new dashboard layout** dashboard."},
		{[]string{"widgets", "layout", "width"}, "The width of the widget. Should be a non-negative integer."},
		{[]string{"widgets", "definition", "slo", "time_windows"}, "Times being monitored."},
		{[]string{"widgets", "definition", "manage_status", "query"}, "Query to filter the monitors with."},
		{[]string{"widgets", "definition", "query_value", "requests", "conditional_formats", "comparator"}, "Comparator to apply."},
		{[]string{"widgets", "definition", "query_value", "requests", "queries", "metrics", "query"}, "Metrics query definition."},
		{[]string{"widgets", "definition", "query_value", "requests", "formulas", "formula"}, "String expression built from queries, formulas, and functions."},
		// number_format exists in exported dashboards but not according to the official client.
		// {[]string{"widgets", "definition", "query_value", "requests", "formulas", "number_format"}, ""},
		{[]string{"widgets", "definition", "timeseries", "requests", "style", "palette"}, "Color palette to apply to the widget."},
	}
	for _, test := range tests {
		name := "Empty path"
		if len(test.paths) != 0 {
			name = strings.Join(test.paths, ",")
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := schemaindexer.Index(test.paths)
			if actual == nil {
				t.Fatal("expected index to return a non-nil schema proxy")
			}
			if actual.Description != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual.Description)
				t.Errorf("%s", cmp.Diff(test.expected, actual.Description))
			}
		})
	}
}
