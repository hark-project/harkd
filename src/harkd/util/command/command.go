package command

import (
	"os/exec"
)

// SimpleResult holds the entire result of a command in memory.
type SimpleResult struct {
	Error      error
	ExitStatus int

	Output []byte
}

// Runner is an interface which can run commands and report whether a command
// is on the path.
type Runner interface {
	HaveOnPath(string) bool
	RunSimple(string, ...string) SimpleResult
}

func NewRunner() Runner {
	return runner{}
}

type runner struct{}

func (r runner) HaveOnPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func (r runner) RunSimple(name string, args ...string) SimpleResult {
	res := SimpleResult{nil, -1, nil}

	cmd := exec.Command(name, args...)

	output, err := cmd.CombinedOutput()
	res.Output = output
	res.Error = err

	// See if this is an ExitError; if so, we can capture the exit status.
	if exitErr, ok := err.(*exec.ExitError); err != nil && ok {
		res.ExitStatus = exitErr.Pid()
	} else {
		res.ExitStatus = 0
	}

	return res
}
