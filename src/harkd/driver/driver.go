package driver

// Info is information about a particular driver.
type Info struct {
	DriverName          string `json:"driverName"`
	AvailableOnPlatform bool   `json:"availableOnPlatform"`
	Installed           bool   `json:"installed"`
	Healthy             bool   `json:"healthy"`
	Version             string `json:"version"`
}

// GetDriverInfo returns information on every Driver supported by hark.
func GetDriverInfo() []Info {
	return []Info{
		{"virtualbox", virtualboxAvailable(), virtualboxInstalled(), virtualboxHealthy(), virtualboxVersion()},
	}
}
