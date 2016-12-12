package driver

import (
	"harkd/util/command"
)

// Info is information about a particular driver.
type Info struct {
	DriverName          string `json:"driverName"`
	AvailableOnPlatform bool   `json:"availableOnPlatform"`
	Installed           bool   `json:"installed"`
	Healthy             bool   `json:"healthy"`
	Version             string `json:"version"`
}

// GetDriverInfo returns information on every Driver supported by hark.
func GetDriverInfo(runner command.Runner) []Info {
	vb := virtualbox{runner}
	return []Info{
		{"virtualbox", vb.available(), vb.installed(), vb.healthy(), vb.version()},
	}
}
