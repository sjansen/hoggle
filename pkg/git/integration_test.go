// +build integration

package git_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/hoggle/pkg/git"
	"github.com/sjansen/hoggle/pkg/subprocess"
)

func TestRoundTrip(t *testing.T) {
	require := require.New(t)

	dir, err := ioutil.TempDir("", "git-config-roundtrip")
	require.NoError(err, "tempdir")

	err = os.Chdir(dir)
	require.NoError(err, "chdir")

	rc, _, _, err := subprocess.Run("git", "init")
	require.NoError(err, "git init")
	require.Equal(0, rc, "git init")

	cfg := &git.Config{}
	for name, value := range map[string]string{
		"hoggle.test-value":                          `Spoon!`,
		"hoggle.ultimate":                            `answer=42 # when question="What do you get when you multiply six by seven?"`,
		"hoggle.http://example.com:8080/foo.rfc3092": "bar baz\nqux quux\n",
	} {
		actual, err := cfg.Get(name)
		require.NoError(err, "initial: "+name)
		require.Equal("", actual)

		err = cfg.Set(name, value)
		require.NoError(err, "set: "+name)

		actual, err = cfg.Get(name)
		require.NoError(err, "final: "+name)
		require.Equal(value, actual)
	}
}
