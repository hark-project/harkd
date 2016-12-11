package services

import (
	"harkd/driver"
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
	return systemService{}
}

type systemService struct{}

func (sc systemService) GetStatus() Status {
	return Status{true}
}

func (sc systemService) GetDriverInfo() []driver.Info {
	return driver.GetDriverInfo()
}
