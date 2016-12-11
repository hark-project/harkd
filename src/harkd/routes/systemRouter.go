package routes

import (
	"harkd/context"
	"harkd/services"

	"github.com/ceralena/go-restroute"
)

type systemRouter struct {
	service services.SystemService
	responseEncoder
	context.Factory
}

func newSystemRouter(ctxFactory context.Factory) systemRouter {
	return systemRouter{
		services.NewSystemService(),
		jsonResponseEncoder(),
		ctxFactory,
	}
}

func (sr systemRouter) getRouteMap() restroute.Map {
	return restroute.Map{
		"^/api/system/status$": restroute.MethodMap{
			"GET": sr.getStatus,
		},
		"^/api/system/driver$": restroute.MethodMap{
			"GET": sr.getDriverInfo,
		},
	}
}

func (sr systemRouter) getStatus(req restroute.Request) {
	status := sr.service.GetStatus()
	sr.Encode(req.W, status)
}

func (sr systemRouter) getDriverInfo(req restroute.Request) {
	driverInfo := sr.service.GetDriverInfo()
	sr.Encode(req.W, driverInfo)
}
