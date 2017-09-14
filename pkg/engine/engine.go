package engine

import (
	"os"

	"github.com/sjansen/hoggle/pkg/agent"
	"github.com/sjansen/hoggle/pkg/storage/local"
)

func Standalone(uri string) error {
	f, err := parse(uri)
	if err != nil {
		return err
	}
	container, err := f.New()
	if err != nil {
		return err
	}

	agent := &agent.Agent{
		Blobs:  container,
		Files:  &local.Filesystem{},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}
	return agent.Run()
}
