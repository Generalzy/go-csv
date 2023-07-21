package go_csv

import "errors"

const (
	Tag = "csv"
)

var (
	InvalidTypeError     = errors.New("invalid type")
	UnsupportedTypeError = errors.New("unsupported type")

	CannotSetError = errors.New("cannot set value")

	MissingHeadError = errors.New("must write head firstly")

	NilPointerError = errors.New("nil pointer")
)
