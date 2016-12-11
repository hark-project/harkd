package services

import (
	"harkd/driver"
	"harkd/util/command"
)

// SystemService is a controller for getting system-level diagnostics
// and status information.
type SystemService interface {
	GetStatus() Status
	GetDriverInfo() []driver.Info
}

// Status represents the current overall status of the hark service.
type Status struct {
	Healthy bool `json:"healthy"`
}

// NewSystemService constructs a SystemService.
func NewSystemService() SystemService {
	return systemService{command.NewRunner()}
}

type systemService struct {
	command.Runner
}

func (sc systemService) GetStatus() Status {
	return Status{true}
}

func (sc systemService) GetDriverInfo() []driver.Info {
	return driver.GetDriverInfo(sc.Runner)
}
