package git

import (
	"bytes"
	"fmt"

	"github.com/sjansen/hoggle/pkg/subprocess"
)

type Config struct{}

func (c *Config) Get(name string) (value string, err error) {
	rc, stdout, _, err := subprocess.Run(
		"git", "config", "--local", "--get-all", "--null", name,
	)
	if err != nil {
		return
	} else if rc == 1 {
		// value is not set
		return
	} else if rc != 0 {
		err = fmt.Errorf(
			"unable to query git config value for %q", name,
		)
		// TODO log stdout & stderr
		return
	} else if bytes.Count(stdout, []byte{0}) > 1 {
		err = fmt.Errorf(
			"git config returned more values than expected for %q", name,
		)
		return
	} else if stdout[len(stdout)-1] != 0 {
		err = fmt.Errorf(
			"git config returned unexpected output for %q", name,
		)
		return
	}

	value = string(stdout[:len(stdout)-1])
	return
}

func (c *Config) Set(name, value string) (err error) {
	// TODO --replace-all ?
	rc, _, _, err := subprocess.Run(
		"git", "config", "--local", name, value,
	)
	if err == nil && rc != 0 {
		err = fmt.Errorf(
			"unable to set git config value for %q", name,
		)
		// TODO log stdout & stderr
		return
	}

	return
}
