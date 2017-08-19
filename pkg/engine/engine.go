package engine

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sjansen/hoggle/pkg/storage"
	"github.com/sjansen/hoggle/pkg/storage/s3"
)

func Standalone(uri string) error {
	f, err := parse(uri)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", f)
	return nil
}

func parse(uri string) (f storage.Factory, err error) {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Opaque != "" {
		err = fmt.Errorf("invalid storage uri: %q", uri)
		return
	}

	q := u.Query()
	switch u.Scheme {
	case "s3", "s3+http", "s3+https":
		var endpoint string
		region := q.Get("region")
		bucket := u.Host
		prefix := u.Path
		if len(prefix) >= 1 && prefix[0] == '/' {
			prefix = prefix[1:]
		}
		if len(u.Scheme) > 2 {
			scheme := u.Scheme[3:]
			endpoint = fmt.Sprintf("%s://%s", scheme, u.Host)
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
			err = fmt.Errorf("missing bucket: %q", uri)
			return
		}
		f = &s3.Factory{
			Region:   region,
			Bucket:   bucket,
			Prefix:   prefix,
			Endpoint: endpoint,
		}
		return
	}

	err = fmt.Errorf("unrecognized scheme: %q", u.Scheme)
	return
}
