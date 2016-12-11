package services

import (
	"harkd/context"
	"harkd/core"
	"harkd/dal"
)

// MachineService is a http service for working with machines.
type MachineService interface {
	GetMachineByID(id string) (core.Machine, error)
	GetMachines() ([]core.Machine, error)

	CreateMachine(core.Machine) error
}

// NewMachineService provides a MachineService.
func NewMachineService(ctxFactory context.Factory) MachineService {
	return machineService{ctxFactory, ctxFactory.GetContext().GetDal()}
}

type machineService struct {
	context.Factory
	dal dal.Dal
}

// GetMachineID looks up a machine by ID. It returns an error if the machine does not exist.
func (mc machineService) GetMachineByID(id string) (core.Machine, error) {
	return mc.dal.GetMachineByID(id)
}

// GetMachines looks up all of the machines managed by hark.
func (mc machineService) GetMachines() ([]core.Machine, error) {
	return mc.dal.GetMachines()
}

// CreateMachine creates a new Machine and saves it to the state.
func (mc machineService) CreateMachine(m core.Machine) error {
	return mc.dal.SaveMachine(m)
}
