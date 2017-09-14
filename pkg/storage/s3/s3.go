package s3

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/sjansen/hoggle/pkg/storage"
)

type Factory struct {
	Region   string
	Bucket   string
	Prefix   string
	Endpoint string
	Profile  string
}

type Bucket struct {
	bucket     string
	prefix     string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func (f *Factory) New() (storage.Container, error) {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	config := &opts.Config
	if f.Endpoint != "" {
		config.Region = aws.String(f.Region)
		config.Endpoint = aws.String(f.Endpoint)
		if strings.HasPrefix(f.Endpoint, "http://") {
			config.DisableSSL = aws.Bool(true)
		}
		config.S3ForcePathStyle = aws.Bool(true)
	} else if f.Region != "" {
		config.Region = aws.String(f.Region)
	}
	if f.Profile != "" {
		opts.Profile = f.Profile
	}
	sess := session.Must(session.NewSessionWithOptions(opts))
	prefix := f.Prefix
	if len(prefix) > 0 && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}
	b := &Bucket{
		bucket:     f.Bucket,
		prefix:     prefix,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}
	return b, nil
}

func (b *Bucket) Download(oid string, dst io.WriterAt) error {
	key := b.oid2key(oid)
	params := &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	}
	if _, err := b.downloader.Download(dst, params); err != nil {
		return err
	} else {
		return nil
	}
}

func (b *Bucket) Upload(oid string, src io.ReadSeeker) error {
	key := b.oid2key(oid)
	params := &s3manager.UploadInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
		Body:   src,
	}
	if _, err := b.uploader.Upload(params); err != nil {
		return err
	} else {
		return nil
	}
}

func (b *Bucket) oid2key(oid string) (key string) {
	if len(oid) < 4 {
		key = fmt.Sprintf("%sobjects/%s", b.prefix, oid)
	} else {
		key = fmt.Sprintf("%sobjects/%s/%s/%s", b.prefix, oid[0:2], oid[2:4], oid)
	}
	return
}
