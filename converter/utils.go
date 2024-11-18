package converter

import "errors"

var ErrUnimplemented error = errors.New("not implemented")

func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

