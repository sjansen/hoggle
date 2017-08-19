// +build integration

package s3_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	"github.com/sjansen/hoggle/pkg/storage/s3"
)

const battlecry = "Spoon!"

func TestRoundTrip(t *testing.T) {
	require := require.New(t)

	bucket := os.Getenv("HOGGLE_TEST_S3_BUCKET")
	prefix := os.Getenv("HOGGLE_TEST_S3_PREFIX")
	region := os.Getenv("HOGGLE_TEST_S3_REGION")
	endpoint := os.Getenv("HOGGLE_TEST_S3_ENDPOINT")
	if bucket == "" {
		t.Skip("$HOGGLE_TEST_S3_BUCKET not set")
	}
	if endpoint != "" {
		require.NotEmpty(
			region,
			"$HOGGLE_TEST_S3_REGION must be set when $HOGGLE_TEST_S3_ENDPOINT is set",
		)
	}

	f := &s3.Factory{
		Region:   region,
		Bucket:   bucket,
		Prefix:   prefix,
		Endpoint: endpoint,
	}
	b, err := f.New()
	require.NoError(err)

	src := bytes.NewReader([]byte(battlecry))
	err = b.Upload("battlecry", src)
	require.NoError(err)

	buf := make([]byte, 0)
	dst := aws.NewWriteAtBuffer(buf)
	err = b.Download("battlecry", dst)
	require.NoError(err)

	require.Equal([]byte(battlecry), dst.Bytes())
}
