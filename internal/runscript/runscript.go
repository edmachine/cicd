package runscript

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

// Script is the structure that holds the script to run, and the output/exit status of the run
type Script struct {
	Path        string
	ExitCode    int
	StandardOut []byte
	StandardErr string
	HasRun      bool
}

// New will return the script and an error
// error is unused for now but can be used in the future as functionality is added
func New(path string) (script Script, err error) {
	script.Path = path
	script.HasRun = false
	return script, nil
}

// Run will execute the script and set the internal data
func (s Script) Run() (err error) {
	if s.HasRun {
		return errors.New("we have already run this command")
	}
	// TODO change this from /bin/echo to just running the script via /bin/bash
	output, err := exec.Command("/bin/echo", s.Path).Output()
	// We don't want to run the same script multiple times
	s.HasRun = true

	s.StandardOut = output
	if err != nil {
		// check if err is an exec.ExitError
		if exiterr, ok := err.(*exec.ExitError); ok {
			s.StandardErr = exiterr.Error()
			// check if the exiterror has a exitStatus
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				s.ExitCode = status.ExitStatus()
			}
		}

		return err
	}

	return nil
}

// StandardOutput converts the byte slice of standard output to a string
func (s Script) StandardOutput() string {
	return fmt.Sprintf("%s", s.StandardOut)
}
