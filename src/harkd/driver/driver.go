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
	virtualbox := Virtualbox{runner}
	return []Info{
		{"virtualbox", virtualbox.available(), virtualbox.installed(), virtualbox.healthy(), virtualbox.version()},
	}
}
