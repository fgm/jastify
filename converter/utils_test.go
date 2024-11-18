package converter

import (
	"errors"
	"testing"
)

func TestMust(t *testing.T) {
	for _, test := range [...]struct {
		name        string
		err         error
		expecdPanic bool
	}{
		{"success", nil, false},
		{"failure", errors.New(""), true},
	} {
		t.Run(test.name, func(t *testing.T) {
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
