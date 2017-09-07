package agent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sjansen/hoggle/pkg/storage"
)

type errorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type initMsg struct {
	Event               string `json:"event"`
	Operation           string `json:"operation"`
	Concurrent          bool   `json:"concurrent"`
	ConcurrentTransfers int    `json:"concurrenttransfers"`
}

type progressMsg struct {
	Event          string `json:"event"`
	Oid            string `json:"oid"`
	BytesSinceLast int64  `json:"bytesSinceLast"`
	BytesSoFar     int64  `json:"bytesSoFar"`
}

type uploadMsg struct {
	Event string `json:"event"`
	Oid   string `json:"oid"`
	Size  int64  `json:"size"`
	Path  string `json:"path"`
}

type uploadCompleteMsg struct {
	Event string    `json:"event"`
	Oid   string    `json:"oid"`
	Error *errorMsg `json:"error,omitempty"`
}

func Run(r io.Reader, w io.Writer, c storage.Container) error {
	stdin := bufio.NewScanner(r)
	stdout := bufio.NewWriter(w)

	msg, err := readInit(stdin)
	if err != nil {
		return err
	}
	stdout.WriteString("{ }\n")
	stdout.Flush()

	fmt.Println(msg)
	if msg.Operation == "upload" {
		for stdin.Scan() {
			req, err := readUpload(stdin)
			if err != nil {
				return err
			}
			fmt.Println(req)
			if req.Event == "terminate" {
				return nil
			}
			writeProgress(stdout, req.Oid, req.Size, req.Size)
			writeUploadComplete(stdout, req.Oid)
		}
	}
	return fmt.Errorf("unexpected end of input")
}

func readInit(stdin *bufio.Scanner) (msg *initMsg, err error) {
	if ok := stdin.Scan(); !ok {
		if err = stdin.Err(); err == nil {
			err = fmt.Errorf("expected init message")
		}
		return
	}

	line := stdin.Text()
	msg = &initMsg{}
	if err = json.Unmarshal([]byte(line), msg); err != nil {
		msg = nil
	} else if msg.Event != "init" {
		err = fmt.Errorf("unexpected message: %q", line)
		msg = nil
	} else if msg.Operation != "download" && msg.Operation != "upload" {
		err = fmt.Errorf("unexpected operation: %q", msg.Operation)
		msg = nil
	}

	return
}

func readUpload(stdin *bufio.Scanner) (msg *uploadMsg, err error) {
	line := stdin.Text()
	msg = &uploadMsg{}
	if err = json.Unmarshal([]byte(line), msg); err != nil {
		msg = nil
	} else if msg.Event != "upload" && msg.Event != "terminate" {
		err = fmt.Errorf("unexpected message: %q", line)
		msg = nil
	}

	return
}

func writeProgress(stdout *bufio.Writer, oid string, change, total int64) error {
	b, err := json.Marshal(&progressMsg{
		Event:          "progress",
		Oid:            oid,
		BytesSinceLast: change,
		BytesSoFar:     total,
	})
	if err != nil {
		return err
	}

	stdout.Write(b)
	stdout.WriteByte('\n')
	stdout.Flush()

	return nil
}

func writeUploadComplete(stdout *bufio.Writer, oid string) error {
	b, err := json.Marshal(&uploadCompleteMsg{
		Event: "complete",
		Oid:   oid,
	})
	if err != nil {
		return err
	}

	stdout.Write(b)
	stdout.WriteByte('\n')
	stdout.Flush()

	return nil
}
