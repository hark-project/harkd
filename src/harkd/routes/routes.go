package routes

import (
	"net/http"

	"harkd/context"

	"github.com/ceralena/go-restroute"
)

// Router implements http.Handler.
type Router http.Handler

// New provides a new Router.
func New(ctxFactory context.Factory) (Router, error) {
	m := restroute.Merge(
		newSystemRouter(ctxFactory).getRouteMap(),
		newMachineRouter(ctxFactory).getRouteMap(),
	)
	return m.Compile()
}
