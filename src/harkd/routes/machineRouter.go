package routes

import (
	"harkd/context"
	"harkd/core"
	"harkd/services"

	"github.com/ceralena/go-restroute"
)

type machineRouter struct {
	service services.MachineService
	responseWriter
	requestDecoder
	context.Factory
}

func newMachineRouter(ctxFactory context.Factory) machineRouter {
	return machineRouter{
		services.NewMachineService(ctxFactory),
		newResponseWriter(),
		jsonRequestDecoder(),
		ctxFactory,
	}
}

func (mr machineRouter) getRouteMap() restroute.Map {
	return restroute.Map{
		"^/api/machine$": restroute.MethodMap{
			"GET": mr.getMachines,
			"PUT": mr.createMachine,
		},
		`^/api/machine/(?P<machine_id>\w+)$`: restroute.MethodMap{
			"GET": mr.getMachineByID,
		},
	}
}

func (mr machineRouter) getMachines(req restroute.Request) {
	machines, err := mr.service.GetMachines()
	if err != nil {
		mr.WriteResponse(req.W, err)
	} else {
		mr.WriteResponse(req.W, machines)
	}
}

func (mr machineRouter) createMachine(req restroute.Request) {
	// Parse the machine from the request body
	var machine core.Machine
	err := mr.Decode(req.R.Body, &machine)
	if err != nil {
		mr.WriteResponse(req.W, err)
		return
	}

	// Get the service to create the machine
	err = mr.service.CreateMachine(machine)
	if err != nil {
		mr.WriteResponse(req.W, err)
		return
	}

	mr.WriteResponseWithStatus(req.W, 201, nil)
}

func (mr machineRouter) getMachineByID(req restroute.Request) {
	machineID := req.Params["machine_id"]
	m, err := mr.service.GetMachineByID(machineID)
	if err != nil {
		mr.WriteResponse(req.W, err)
	} else {
		mr.WriteResponse(req.W, m)
	}
}
