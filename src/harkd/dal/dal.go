package dal

import (
	"harkd/core"
)

// Dal is the interface for reading and persisting Hark state.
type Dal interface {
	GetMachines() ([]core.Machine, error)
	GetMachineByID(string) (core.Machine, error)

	SaveMachine(core.Machine) error
}
