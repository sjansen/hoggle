package s3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	require := require.New(t)

	// TODO add more factory test cases
	f := Factory{
		Region: "us-west-2",
		Bucket: "foo",
		Prefix: "hoggle-test",
	}

	expected := Bucket{
		bucket: "foo",
		prefix: "hoggle-test/",
	}

	container, err := f.New()
	require.NoError(err)

	actual := container.(*Bucket)
	require.Equal(expected.prefix, actual.prefix)
}

func TestKeyMunging(t *testing.T) {
	require := require.New(t)

	b := Bucket{prefix: ""}
	for oid, expected := range map[string]string{
		"123":   "objects/123",
		"1234":  "objects/12/34/1234",
		"12345": "objects/12/34/12345",
		"950d608fb3f28c6da83d8a106b7b429734c533b16a35baebea695a4d6b70a233": "objects/95/0d/950d608fb3f28c6da83d8a106b7b429734c533b16a35baebea695a4d6b70a233",
	} {
		actual := b.oid2key(oid)
		require.Equal(expected, actual)
	}
}
