package converter

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMust(t *testing.T) {
	t.Parallel()
	for _, test := range [...]struct {
		name        string
		err         error
		expecdPanic bool
	}{
		{"success", nil, false},
		{"failure", errors.New(""), true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					if !test.expecdPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()
			actual := Must(func() (any, error) {
				return "foo", test.err
			}())
			if actual != "foo" {
				t.Errorf("expected 'foo', got '%v'", actual)
			}
		})
	}
}

func TestPath_Push(t *testing.T) {
	t.Parallel()
	const one = "one"
	sliceOfOne := []string{one}
	pathOfOne := Path{sliceOfOne}

	for _, test := range [...]struct {
		name   string
		before []string
		after  []string
		result Path
	}{
		{"nil", nil, nil, pathOfOne},
		{"empty", []string{}, []string{}, pathOfOne},
		{"one", sliceOfOne, sliceOfOne, Path{[]string{one, one}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			path := Path{Data: test.before}
			actual := path.Push(one)
			if !cmp.Equal(path, Path{test.after}) {
				t.Errorf("path was modified: %s\n", cmp.Diff(path, Path{test.after}))
			}
			if !cmp.Equal(actual, test.result) {
				t.Errorf("unexpected output: %s\n", cmp.Diff(actual, test.result))
			}
		})
	}
}

func TestPath_Pop(t *testing.T) {
	t.Parallel()
	const one = "one"
	nilPath := Path{}
	emptyPath := Path{make([]string, 0)}
	sliceOfOne := []string{one}
	pathOfOne := Path{sliceOfOne}

	for _, test := range [...]struct {
		name          string
		before        []string
		after         []string
		expectedPath  Path
		expectedValue string
	}{
		{"nil", nil, nil, nilPath, ""},
		{"empty", []string{}, []string{}, emptyPath, ""},
		{"one", sliceOfOne, sliceOfOne, emptyPath, one},
		{"two", []string{one, one}, []string{one, one}, pathOfOne, one},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			before := Path{Data: test.before}
			actualPath, actualValue := before.Pop()
			if !cmp.Equal(before, Path{test.after}) {
				t.Errorf("before was modified: %s\n", cmp.Diff(before, Path{test.after}))
			}
			if !cmp.Equal(actualPath, test.expectedPath) {
				t.Errorf("unexpected output: %s\n", cmp.Diff(actualPath, test.expectedPath))
			}
			if !cmp.Equal(actualValue, test.expectedValue) {
				t.Errorf("unexpected output: %s\n", cmp.Diff(actualValue, test.expectedValue))
			}
		})
	}
}
