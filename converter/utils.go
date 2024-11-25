package converter

import (
	"errors"
	"strings"
)

var ErrUnimplemented error = errors.New("not implemented")

func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

type Path struct {
	Data []string
}

func (path Path) Push(s string) Path {
	ret := Path{make([]string, len(path.Data), len(path.Data)+1)}
	copy(ret.Data, path.Data)
	ret.Data = append(ret.Data, s)
	return ret
}

func (path Path) Pop() (Path, string) {
	if path.Data == nil {
		return Path{}, ""
	}
	// Since we tested for nil slice, this matches an empty slice.
	if len(path.Data) == 0 {
		return Path{[]string{}}, ""
	}
	max := len(path.Data) - 1
	ret := Path{make([]string, max)}
	copy(ret.Data, path.Data[:max])
	return ret, path.Data[max]
}

func (path Path) String() string {
	return strings.Join(path.Data, "/")
}

func (path Path) Slice() []string {
	return path.Data
}
