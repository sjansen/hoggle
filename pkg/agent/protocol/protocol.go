package protocol

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	Terminate = errors.New("terminate")

	mustInitErr    = errors.New("must Init() first")
	mustReadyErr   = errors.New("must Ready() first")
	notDownloadErr = errors.New("not a download session")
	notUploadErr   = errors.New("not an upload session")
	terminatedErr  = errors.New("session terminated")
)

type Session struct {
	stdin  *bufio.Scanner
	stdout *bufio.Writer
	// state
	down  bool
	init  bool
	ready bool
	term  bool

	files map[string]bool
}

type InitMsg struct {
	Event               string `json:"event"`
	Operation           string `json:"operation"`
	Concurrent          bool   `json:"concurrent"`
	ConcurrentTransfers int    `json:"concurrenttransfers"`
}

type DownloadMsg struct {
	Event string `json:"event"`
	Oid   string `json:"oid"`
	Size  int64  `json:"size"`
}

type UploadMsg struct {
	Event string `json:"event"`
	Oid   string `json:"oid"`
	Size  int64  `json:"size"`
	Path  string `json:"path"`
}

type errorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type progressMsg struct {
	Event          string `json:"event"`
	Oid            string `json:"oid"`
	BytesSinceLast int64  `json:"bytesSinceLast"`
	BytesSoFar     int64  `json:"bytesSoFar"`
}

type downloadCompleteMsg struct {
	Event string    `json:"event"`
	Oid   string    `json:"oid"`
	Path  string    `json:"path,omitempty"`
	Error *errorMsg `json:"error,omitempty"`
}

type uploadCompleteMsg struct {
	Event string    `json:"event"`
	Oid   string    `json:"oid"`
	Error *errorMsg `json:"error,omitempty"`
}

func NewSession(stdin io.Reader, stdout io.Writer) *Session {
	return &Session{
		stdin:  bufio.NewScanner(stdin),
		stdout: bufio.NewWriter(stdout),
		files:  make(map[string]bool),
	}
}

func (s *Session) Init() (msg *InitMsg, err error) {
	if ok := s.stdin.Scan(); !ok {
		if err = s.stdin.Err(); err == nil {
			err = fmt.Errorf("expected init message")
		}
		return
	}

	line := s.stdin.Text()
	tmp := &InitMsg{}
	if err = json.Unmarshal([]byte(line), tmp); err != nil {
		return
	}

	if tmp.Event != "init" {
		err = fmt.Errorf("unexpected message: %q", line)
		return
	}

	switch tmp.Operation {
	case "download":
		s.down = true
	case "upload":
		s.down = false
	default:
		err = fmt.Errorf("unexpected operation: %q", tmp.Operation)
		return
	}

	msg = tmp
	s.init = true
	return
}

func (s *Session) Ready() (err error) {
	if !s.init {
		return mustInitErr
	}

	if _, err = s.stdout.WriteString("{ }\n"); err != nil {
		return
	}
	if err = s.stdout.Flush(); err != nil {
		return
	}

	s.ready = true
	return
}

func (s *Session) complete(oid string, b []byte) (err error) {
	if complete, ok := s.files[oid]; !ok {
		return fmt.Errorf("invalid oid: %q", oid)
	} else if complete {
		return fmt.Errorf("already complete: %q", oid)
	}

	if _, err = s.stdout.Write(b); err != nil {
		return
	}
	if err = s.stdout.WriteByte('\n'); err != nil {
		return
	}
	err = s.stdout.Flush()

	s.files[oid] = true
	return
}

func (s *Session) ReportCompletedDownload(oid, path string) (err error) {
	b, err := json.Marshal(&downloadCompleteMsg{
		Event: "complete",
		Oid:   oid,
		Path:  path,
	})
	if err != nil {
		return
	}
	return s.complete(oid, b)
}

func (s *Session) ReportCompletedUpload(oid string) (err error) {
	b, err := json.Marshal(&uploadCompleteMsg{
		Event: "complete",
		Oid:   oid,
	})
	if err != nil {
		return
	}
	return s.complete(oid, b)
}

func (s *Session) ReportFailedDownload(oid, msg string) (err error) {
	b, err := json.Marshal(&downloadCompleteMsg{
		Event: "complete",
		Oid:   oid,
		Error: &errorMsg{
			Code:    1,
			Message: msg,
		},
	})
	if err != nil {
		return
	}
	return s.complete(oid, b)
}

func (s *Session) ReportFailedUpload(oid, msg string) (err error) {
	b, err := json.Marshal(&uploadCompleteMsg{
		Event: "complete",
		Oid:   oid,
		Error: &errorMsg{
			Code:    1,
			Message: msg,
		},
	})
	if err != nil {
		return
	}
	return s.complete(oid, b)
}

func (s *Session) ReportProgress(oid string, change, total int64) (err error) {
	if complete, ok := s.files[oid]; !ok {
		return fmt.Errorf("invalid oid: %q", oid)
	} else if complete {
		return fmt.Errorf("already complete: %q", oid)
	}

	b, err := json.Marshal(&progressMsg{
		Event:          "progress",
		Oid:            oid,
		BytesSinceLast: change,
		BytesSoFar:     total,
	})
	if err != nil {
		return
	}

	if _, err = s.stdout.Write(b); err != nil {
		return
	}
	if err = s.stdout.WriteByte('\n'); err != nil {
		return
	}
	err = s.stdout.Flush()

	return
}

func (s *Session) StartDownload() (msg *DownloadMsg, err error) {
	if !s.ready {
		err = mustReadyErr
		return
	} else if !s.down {
		err = notDownloadErr
		return
	} else if s.term {
		err = terminatedErr
		return
	}

	if ok := s.stdin.Scan(); !ok {
		err = errors.New("unexpected end of input")
		return
	}

	line := s.stdin.Text()
	tmp := &DownloadMsg{}
	if err = json.Unmarshal([]byte(line), tmp); err != nil {
		return
	}

	switch tmp.Event {
	case "terminate":
		err = Terminate
	case "download":
		msg = tmp
		s.files[msg.Oid] = false
	default:
		err = fmt.Errorf("unexpected message: %q", line)
	}

	return
}

func (s *Session) StartUpload() (msg *UploadMsg, err error) {
	if !s.ready {
		err = mustReadyErr
		return
	} else if s.down {
		err = notUploadErr
		return
	} else if s.term {
		err = terminatedErr
		return
	}

	if ok := s.stdin.Scan(); !ok {
		err = errors.New("unexpected end of input")
		return
	}

	line := s.stdin.Text()
	tmp := &UploadMsg{}
	if err = json.Unmarshal([]byte(line), tmp); err != nil {
		return
	}

	switch tmp.Event {
	case "terminate":
		err = Terminate
	case "upload":
		msg = tmp
		s.files[msg.Oid] = false
	default:
		err = fmt.Errorf("unexpected message: %q", line)
	}

	return
}
