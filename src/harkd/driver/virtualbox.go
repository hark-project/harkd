package driver

import (
	"os/exec"
	"strings"
)

func virtualboxAvailable() bool {
	// virtualbox works on all supported platforms
	return true
}

func virtualboxInstalled() bool {
	// Just check for the command in the path
	_, err := exec.LookPath("VBoxManage")
	return err == nil
}

func virtualboxHealthy() bool {
	// Run the command with --version; ignore the output
	c := exec.Command("VBoxManage", "--version")
	if err := c.Start(); err != nil {
		return false
	}
	err := c.Wait()
	return err == nil
}

func virtualboxVersion() string {
	c := exec.Command("VBoxManage", "--version")
	b, err := c.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}
