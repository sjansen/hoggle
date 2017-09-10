package agent

import (
	"fmt"
	"io"

	"github.com/sjansen/hoggle/pkg/agent/protocol"
	"github.com/sjansen/hoggle/pkg/storage"
)

func Run(r io.Reader, w io.Writer, c storage.Container) error {
	session := protocol.NewSession(r, w)

	init, err := session.Init()
	if err != nil {
		return err
	}

	err = session.Ready()
	if err != nil {
		return err
	}

	switch init.Operation {
	case "upload":
		return upload(session)
	}
	return fmt.Errorf("unexpected operation: %q", init.Operation)
}

func upload(s *protocol.Session) error {
	for {
		msg, err := s.StartUpload()
		if err == protocol.Terminate {
			return nil
		} else if err != nil {
			return err
		}

		s.ReportProgress(msg.Oid, msg.Size, msg.Size)
		s.ReportCompletedUpload(msg.Oid)
	}
}
