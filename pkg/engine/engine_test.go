package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjansen/hoggle/pkg/storage"
	"github.com/sjansen/hoggle/pkg/storage/s3"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	for uri, expected := range map[string]storage.Factory{
		"s3://foo": &s3.Factory{
			Bucket: "foo",
		},
		"s3://bar/": &s3.Factory{
			Bucket: "bar",
		},
		"s3://baz/qux": &s3.Factory{
			Bucket: "baz",
			Prefix: "qux",
		},
		"s3://quux/corge/": &s3.Factory{
			Bucket: "quux",
			Prefix: "corge/",
		},
		"s3://grault/garply/waldo?region=us-gov-west-1": &s3.Factory{
			Region: "us-gov-west-1",
			Bucket: "grault",
			Prefix: "garply/waldo",
		},
		"s3+http://example.com:8080/foo": &s3.Factory{
			Bucket:   "foo",
			Endpoint: "http://example.com:8080",
		},
		"s3+http://storage.example.com/bar/": &s3.Factory{
			Bucket:   "bar",
			Endpoint: "http://storage.example.com",
		},
		"s3+https://storage.example.com/baz/qux": &s3.Factory{
			Bucket:   "baz",
			Prefix:   "qux",
			Endpoint: "https://storage.example.com",
		},
		"s3+https://storage.example.com/quux/corge": &s3.Factory{
			Bucket:   "quux",
			Prefix:   "corge",
			Endpoint: "https://storage.example.com",
		},
	} {
		actual, err := parse(uri)
		if assert.NoError(err) {
			assert.Equal(expected, actual)
		}
	}

	for uri, expected := range map[string]string{
		"":                `invalid storage uri: ""`,
		"invalid":         `invalid storage uri: "invalid"`,
		":invalid":        `invalid storage uri: ":invalid"`,
		"s3:bucket":       `invalid storage uri: "s3:bucket"`,
		"s3://":           `invalid storage uri: "s3://"`,
		"s3+ftp://foo":    `unrecognized scheme: "s3+ftp"`,
		"s3+http://host":  `missing bucket: "s3+http://host"`,
		"s3+http://host/": `missing bucket: "s3+http://host/"`,
		"//bucket":        `invalid storage uri: "//bucket"`,
	} {
		_, err := parse(uri)
		if assert.Error(err) {
			actual := err.Error()
			assert.Equal(expected, actual)
		}
	}
}
