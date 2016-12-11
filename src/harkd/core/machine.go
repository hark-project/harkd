package core

import (
	"harkd/errors"
)

// Machine is the core hark data structure for a single machine.
type Machine struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Validate validates the machine.
func (m Machine) Validate() error {
	if m.ID == "" {
		return errors.ErrEntityInvalid("machine id cannot be empty")
	}
	if m.Name == "" {
		return errors.ErrEntityInvalid("machine name cannot be empty")
	}
	return nil
}
