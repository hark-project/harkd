package driver

import (
	"strings"

	"harkd/util/command"
)

type virtualbox struct {
	command.Runner
}

func (v virtualbox) available() bool {
	// virtualbox works on all supported platforms
	return true
}

func (v virtualbox) installed() bool {
	// Just check for the command in the path
	return v.HaveOnPath("VBoxManage")
}

func (v virtualbox) healthy() bool {
	// Run the command with --version; ignore the output
	res := v.RunSimple("VBoxManage", "--version")
	return res.Error == nil
}

func (v virtualbox) version() string {
	res := v.RunSimple("VBoxManage", "--version")
	if res.Error != nil {
		return ""
	}
	return strings.TrimSpace(string(res.Output))
}
