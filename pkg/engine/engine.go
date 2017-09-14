package engine

import (
	"os"

	"github.com/sjansen/hoggle/pkg/agent"
	"github.com/sjansen/hoggle/pkg/config"
	"github.com/sjansen/hoggle/pkg/git"
	"github.com/sjansen/hoggle/pkg/storage/local"
)

func Init(uri string) error {
	// validate uri
	_, err := parse(uri)
	if err != nil {
		return err
	}

	return config.Init(&git.Config{}, uri)
}

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
