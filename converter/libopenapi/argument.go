package converter

import (
	"fmt"
	"io"
	"strings"
)

type TFArgument struct {
	Name        string
	Value       any
	RO          bool
	Deprecation string
}

func (arg *TFArgument) Render(w io.Writer, depth int) error {
	switch arg.Value.(type) {
	case bool:
		return arg.renderBool(w, depth)
	case string:
		return arg.renderString(w, depth)
	case []any:
		return arg.renderSlice(w, depth)
	case Unsupported:
		fmt.Fprintf(w, "%s%s%s %T", Indent(depth+1), UnsupportedPrefix, arg.Name, arg.Value)
	}
	return nil
}

func (arg *TFArgument) renderBool(w io.Writer, depth int) error {
	out := Indent(depth)
	if arg.RO {
		out += ReadOnlyPrefix
	}
	out += fmt.Sprintf("%s\t= %t", arg.Name, arg.Value)
	if arg.Deprecation != "" {
		out = out + "\t// " + arg.Deprecation
	}
	_, err := fmt.Fprintln(w, out)
	return err
}

func (arg *TFArgument) renderSlice(w io.Writer, depth int) error {
	out := Indent(depth)
	if arg.RO {
		out += ReadOnlyPrefix
	}
	vs, ok := arg.Value.([]any)
	if !ok {
		return fmt.Errorf("%s is not a slice %T", arg.Name, arg.Value)
	}
	sl := make([]string, len(vs))
	for i, v := range vs {
		sl[i] = fmt.Sprintf("%q", v)
	}
	out += fmt.Sprintf("%s\t= [%s]", arg.Name, strings.Join(sl, ", "))
	if arg.Deprecation != "" {
		out = out + "\t// " + arg.Deprecation
	}
	_, err := fmt.Fprintln(w, out)
	return err
}

func (arg *TFArgument) renderString(w io.Writer, depth int) error {
	out := Indent(depth)
	if arg.RO {
		out += ReadOnlyPrefix
	}
	out += fmt.Sprintf("%s\t= %q", arg.Name, arg.Value)
	if arg.Deprecation != "" {
		out = out + "\t// " + arg.Deprecation
	}
	_, err := fmt.Fprintln(w, out)
	return err
}
