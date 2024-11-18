package legacy

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fgm/jastify/converter"
)

var (
	j1 = converter.Jmap{"a_cool_key": "a cool value", "another_key": "another value"}
	j2 = converter.Jmap{"a_cool_key": "another cool value", "another_key": "another ANOTHER value"}
)

func Test_assignmentString(t *testing.T) {
	for _, test := range [...]struct {
		name       string
		inputName  string
		inputValue any
		expected   string
	}{
		{"integer value", "i", 7, "i = 7\n"},
		{"nil value", "i", nil, ""},
		{"single-line string", "sls", "foo", "sls = \"foo\"\n"},
		{"multi-line string", "mls", "foo\nbar", "mls = <<EOF\nfoo\nbar\nEOF\n"},
		{"single-line strings list", "slsl", []string{"foo", "bar"}, "slsl = [\"foo\",\"bar\"]\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if actual := AssignmentString(test.inputName, test.inputValue); actual != test.expected {
				t.Errorf(`assignment string %v`, cmp.Diff(actual, test.expected))
			}
		})
	}
}

func Test_block(t *testing.T) {
	input := j1
	const expected = `my_block {
a_cool_key = "a cool value"
another_key = "another value"
}
`

	if actual := block("my_block", input, AssignmentString); actual != expected {
		t.Errorf("block output did not match: %v\n", cmp.Diff(expected, actual))
	}
}

func Test_blockList(t *testing.T) {
	input := converter.Jmaps{j1, j2}
	const expected = `my_block {
a_cool_key = "a cool value"
another_key = "another value"
}
my_block {
a_cool_key = "another cool value"
another_key = "another ANOTHER value"
}
`
	if actual := blockList(input, "my_block", AssignmentString); actual != expected {
		t.Errorf("blockList output did not match: %v\n", cmp.Diff(expected, actual))
	}
}

func Test_convertFromDefinition(t *testing.T) {
	for _, test := range [...]struct {
		name            string
		input           string
		expected        string
		expectedSuccess bool
	}{
		{"existing key", "is_read_only", "is_read_only = false\n", true},
		{"nonexistent key", "not a thing", "", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual, err := convertFromDefinition(DASHBOARD, test.input, false)
			if actual != test.expected {
				t.Errorf("convertFromDefinition output did not match: %v\n", cmp.Diff(actual, test.expected))
			}
			if err == nil != test.expectedSuccess {
				t.Errorf("convertFromDefinition error did not match: %v\n", cmp.Diff(err, test.expectedSuccess))
			}
		})
	}
}

func Test_JmapsFromAny(t *testing.T) {
	var (
		normalInput  = []any{j1, j2}
		normalOutput = converter.Jmaps{j1, j2}
		nj           = []any{"", []any(nil)}
	)

	for _, test := range [...]struct {
		name     string
		input    any
		expected converter.Jmaps
		errorRx  string
	}{
		{"non-slice", "foo", nil, `^items expected as \[\]any but got`},
		{"empty slice of strings", []string{}, nil, `^items expected as \[\]any but got`},
		{"empty slice of any", []any{}, make(converter.Jmaps, 0), ""},
		{"nil", nil, nil, `^items expected as \[\]any but got`},
		{"slice of any non-Jmap", nj, nil, `^item \[[\d]+\] expected as Jmap but got`},
		{"normal Jmaps", normalInput, normalOutput, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual, err := converter.JmapsFromAny(test.input)
			if !cmp.Equal(actual, test.expected) {
				t.Errorf("JmapsFromAny output did not match: %v\n", cmp.Diff(actual, test.expected))
			}
			switch {
			case test.errorRx == "" && err != nil:
				t.Errorf("JmapsFromAny unexpected error: %v\n", err)
			case test.errorRx != "" && err == nil:
				t.Errorf("JmapsFromAny did not expect an error but got %#v\n", err)
			case test.errorRx != "" && err != nil:
				rx := regexp.MustCompile(test.errorRx)
				if !rx.MatchString(err.Error()) {
					t.Errorf("JmapsFromAny expected error to match %q but got %s\n", rx.String(), err)
				}
			}
		})
	}
}

func Test_mapContents(t *testing.T) {
	input := j1
	const expected = `a_cool_key = "a cool value"
another_key = "another value"
`
	if actual := mapContents(input, AssignmentString); actual != expected {
		t.Errorf("mapContents output did not match: %v\n", cmp.Diff(expected, actual))
	}
}

func Test_queryBlock(t *testing.T) {
	input := j1
	const expected = `query {
  q1 {
a_cool_key = "a cool value"
another_key = "another value"
}
}
`
	if actual := queryBlock("q1", input, AssignmentString); actual != expected {
		t.Errorf("queryBlock output did not match: %v\n", cmp.Diff(expected, actual))
	}
}

func Test_queryBlockList(t *testing.T) {
	input := converter.Jmaps{j1, j2}
	const expected = `query {
  metric_query {
a_cool_key = "a cool value"
another_key = "another value"
}
}
query {
  metric_query {
a_cool_key = "another cool value"
another_key = "another ANOTHER value"
}
}
`
	if actual := queryBlockList(input, AssignmentString); actual != expected {
		t.Errorf("queryBlockList output did not match: %v\n", cmp.Diff(expected, actual))
	}
}
