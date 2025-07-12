package navigator

import "io"

type Tree interface {
}

type Forest struct {
	trees []Tree
}

func NewTree() *Forest {
	return &Forest{}
}

func (tr *Forest) Parse(name string, r io.Reader) (err error) {
	return
}
