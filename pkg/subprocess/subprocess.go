package subprocess

import (
	"os/exec"
	"syscall"
)

func Run(name string, args ...string) (rc int, stdout, stderr []byte, err error) {
	cmd := exec.Command(name, args...)

	stdout, err = cmd.Output()
	if err == nil {
		// should be zero
		rc = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	} else {
		if ee, ok := err.(*exec.ExitError); ok {
			err = nil
			rc = ee.Sys().(syscall.WaitStatus).ExitStatus()
			stderr = ee.Stderr
		} else {
			// unable to run subprocess
			rc = 1
			return
		}
	}

	return
}
