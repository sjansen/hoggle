package agent

import (
	"fmt"
	"io"

	"github.com/sjansen/hoggle/pkg/agent/protocol"
	"github.com/sjansen/hoggle/pkg/storage"
)

type Agent struct {
	Blobs  storage.Container
	Files  storage.Filesystem
	Stdin  io.Reader
	Stdout io.Writer
}

func (a *Agent) Run() error {
	session := protocol.NewSession(a.Stdin, a.Stdout)

	init, err := session.Init()
	if err != nil {
		return err
	}

	err = session.Ready()
	if err != nil {
		return err
	}

	switch init.Operation {
	case "download":
		return a.download(session)
	case "upload":
		return a.upload(session)
	}
	return fmt.Errorf("unexpected operation: %q", init.Operation)
}

func (a *Agent) download(s *protocol.Session) error {
	for {
		msg, err := s.StartDownload()
		if err == protocol.Terminate {
			return nil
		} else if err != nil {
			return err
		}

		file, err := a.Files.TempFile()
		if err != nil {
			return err
		}

		err = a.Blobs.Download(msg.Oid, file)
		// TODO report failed download

		s.ReportProgress(msg.Oid, msg.Size, msg.Size)
		s.ReportCompletedDownload(msg.Oid, file.Name())
	}
}

func (a *Agent) upload(s *protocol.Session) error {
	for {
		msg, err := s.StartUpload()
		if err == protocol.Terminate {
			return nil
		} else if err != nil {
			return err
		}

		file, err := a.Files.Open(msg.Path)
		if err != nil {
			return err
		}

		err = a.Blobs.Upload(msg.Oid, file)
		// TODO report failed upload

		s.ReportProgress(msg.Oid, msg.Size, msg.Size)
		s.ReportCompletedUpload(msg.Oid)
	}
}
