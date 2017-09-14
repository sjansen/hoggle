package engine

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sjansen/hoggle/pkg/storage"
	"github.com/sjansen/hoggle/pkg/storage/s3"
)

type urlParser func(url *url.URL) (f storage.Factory, err error)

var schemes = map[string]urlParser{
	"s3":       parseS3,
	"s3+http":  parseS3,
	"s3+https": parseS3,
}

func parse(uri string) (f storage.Factory, err error) {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Opaque != "" {
		err = fmt.Errorf("invalid storage uri: %q", uri)
		return
	}

	if parseURL, ok := schemes[u.Scheme]; ok {
		return parseURL(u)
	}

	err = fmt.Errorf("unrecognized scheme: %q", u.Scheme)
	return
}

func parseS3(url *url.URL) (f storage.Factory, err error) {
	var endpoint string
	q := url.Query()
	profile := q.Get("profile")
	region := q.Get("region")
	bucket := url.Host
	prefix := url.Path
	if len(prefix) >= 1 && prefix[0] == '/' {
		prefix = prefix[1:]
	}
	if len(url.Scheme) > 2 {
		scheme := url.Scheme[3:]
		endpoint = fmt.Sprintf("%s://%s", scheme, url.Host)
		if len(prefix) < 1 {
			bucket = ""
			prefix = ""
		} else if idx := strings.Index(prefix, "/"); idx < 0 {
			bucket = prefix
			prefix = ""
		} else {
			bucket = prefix[0:idx]
			prefix = prefix[idx+1:]
		}
	}
	if len(bucket) < 1 {
		err = fmt.Errorf("missing bucket: %q", url.String())
		return
	}
	f = &s3.Factory{
		Region:   region,
		Bucket:   bucket,
		Prefix:   prefix,
		Endpoint: endpoint,
		Profile:  profile,
	}
	return
}
