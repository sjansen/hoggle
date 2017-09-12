package agent_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/hoggle/pkg/agent"
	"github.com/sjansen/hoggle/pkg/storage"
)

type container struct {
	n       int
	results []error
}
type filesystem struct{ n int }
type readseeker struct{}
type tempfile struct{ n int }

func (c *container) Download(oid string, dst io.WriterAt) (err error) {
	if c.n < len(c.results) {
		err = c.results[c.n]
	}
	c.n += 1
	return err
}

func (c *container) Upload(oid string, src io.ReadSeeker) (err error) {
	if c.n < len(c.results) {
		err = c.results[c.n]
	}
	c.n += 1
	return err
}

func (fs *filesystem) Open(path string) (io.ReadSeeker, error) {
	return &readseeker{}, nil
}

func (fs *filesystem) TempFile() (storage.TempFile, error) {
	fs.n += 1
	return &tempfile{fs.n}, nil
}

func (rs *readseeker) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (rs *readseeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *tempfile) Close() error {
	return nil
}

func (f *tempfile) Name() string {
	return fmt.Sprintf("/tmp/file%d.png", f.n)
}

func (f *tempfile) WriteAt(p []byte, off int64) (n int, err error) {
	return len(p), nil
}

func TestDownload(t *testing.T) {
	require := require.New(t)

	stdin := bytes.NewReader([]byte(
		`{ "event": "init", "operation": "download", "concurrent": true, "concurrenttransfers": 3 }
{ "event": "download", "oid": "bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a", "size": 34623 }
{ "event": "download", "oid": "22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e", "size": 21245 }
{ "event": "terminate" }
`,
	))

	stdout := &bytes.Buffer{}
	agent := &agent.Agent{
		Blobs:  &container{},
		Files:  &filesystem{},
		Stdin:  stdin,
		Stdout: stdout,
	}
	err := agent.Run()
	require.NoError(err)

	// NOTE: if this test proves unstable, compare individual lines using require.JSONEq()
	actual := stdout.String()
	expected :=
		`{ }
{"event":"progress","oid":"bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a","bytesSinceLast":34623,"bytesSoFar":34623}
{"event":"complete","oid":"bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a","path":"/tmp/file1.png"}
{"event":"progress","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e","bytesSinceLast":21245,"bytesSoFar":21245}
{"event":"complete","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e","path":"/tmp/file2.png"}
`
	require.Equal(expected, actual)
}

func TestDownloadFailure(t *testing.T) {
	require := require.New(t)

	stdin := bytes.NewReader([]byte(
		`{ "event": "init", "operation": "download", "concurrent": true, "concurrenttransfers": 3 }
{ "event": "download", "oid": "bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a", "size": 34623 }
{ "event": "download", "oid": "22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e", "size": 21245 }
{ "event": "terminate" }
`,
	))

	stdout := &bytes.Buffer{}
	agent := &agent.Agent{
		Blobs: &container{results: []error{
			errors.New("boom"),
			nil,
		}},
		Files:  &filesystem{},
		Stdin:  stdin,
		Stdout: stdout,
	}
	err := agent.Run()
	require.NoError(err)

	// NOTE: if this test proves unstable, compare individual lines using require.JSONEq()
	actual := stdout.String()
	expected :=
		`{ }
{"event":"complete","oid":"bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a","error":{"code":1,"message":"boom"}}
{"event":"progress","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e","bytesSinceLast":21245,"bytesSoFar":21245}
{"event":"complete","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e","path":"/tmp/file2.png"}
`
	require.Equal(expected, actual)
}

func TestUpload(t *testing.T) {
	require := require.New(t)

	stdin := bytes.NewReader([]byte(
		`{ "event": "init", "operation": "upload", "concurrent": true, "concurrenttransfers": 3 }
{ "event": "upload", "oid": "bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a", "size": 34623, "path": "/path/to/file1.png" }
{ "event": "upload", "oid": "22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e", "size": 21245, "path": "/path/to/file2.png" }
{ "event": "terminate" }
`,
	))

	stdout := &bytes.Buffer{}
	agent := &agent.Agent{
		Blobs:  &container{},
		Files:  &filesystem{},
		Stdin:  stdin,
		Stdout: stdout,
	}
	err := agent.Run()
	require.NoError(err)

	// NOTE: if this test proves unstable, compare individual lines using require.JSONEq()
	actual := stdout.String()
	expected :=
		`{ }
{"event":"progress","oid":"bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a","bytesSinceLast":34623,"bytesSoFar":34623}
{"event":"complete","oid":"bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a"}
{"event":"progress","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e","bytesSinceLast":21245,"bytesSoFar":21245}
{"event":"complete","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e"}
`
	require.Equal(expected, actual)
}

func TestUploadFailure(t *testing.T) {
	require := require.New(t)

	stdin := bytes.NewReader([]byte(
		`{ "event": "init", "operation": "upload", "concurrent": true, "concurrenttransfers": 3 }
{ "event": "upload", "oid": "bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a", "size": 34623, "path": "/path/to/file1.png" }
{ "event": "upload", "oid": "22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e", "size": 21245, "path": "/path/to/file2.png" }
{ "event": "terminate" }
`,
	))

	stdout := &bytes.Buffer{}
	agent := &agent.Agent{
		Blobs: &container{results: []error{
			errors.New("boom"),
			nil,
		}},
		Files:  &filesystem{},
		Stdin:  stdin,
		Stdout: stdout,
	}
	err := agent.Run()
	require.NoError(err)

	// NOTE: if this test proves unstable, compare individual lines using require.JSONEq()
	actual := stdout.String()
	expected :=
		`{ }
{"event":"complete","oid":"bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a","error":{"code":1,"message":"boom"}}
{"event":"progress","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e","bytesSinceLast":21245,"bytesSoFar":21245}
{"event":"complete","oid":"22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e"}
`
	require.Equal(expected, actual)
}
