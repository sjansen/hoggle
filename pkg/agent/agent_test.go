package agent_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/hoggle/pkg/agent"
)

func TestHappyPath(t *testing.T) {
	require := require.New(t)

	stdin := bytes.NewReader([]byte(
		`{ "event": "init", "operation": "upload", "concurrent": true, "concurrenttransfers": 3 }
{ "event": "upload", "oid": "bf3e3e2af9366a3b704ae0c31de5afa64193ebabffde2091936ad2e7510bc03a", "size": 34623, "path": "/path/to/file1.png" }
{ "event": "upload", "oid": "22ab5f63670800cc7be06dbed816012b0dc411e774754c7579467d2536a9cf3e", "size": 21245, "path": "/path/to/file2.png" }
{ "event": "terminate" }
`,
	))

	stdout := &bytes.Buffer{}
	err := agent.Run(stdin, stdout, nil)
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
