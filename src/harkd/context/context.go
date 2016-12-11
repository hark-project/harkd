package context

import (
	"harkd/dal"
)

// Context is an interface that can provide a context for interacting with hark
// state and the hark VM runtime.
type Context interface {
	GetDal() dal.Dal
}

type dirContext struct {
	dir string
	dal dal.Dal
}

func (d dirContext) GetDal() dal.Dal {
	return d.dal
}
