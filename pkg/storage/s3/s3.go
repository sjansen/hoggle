package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	format     string
	prefix     string
	svc        *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

var UnsupportBucketFormatErr = errors.New("unsupported bucket format")

const bucketFormat = "1\n"

func (f *Factory) New() (storage.Container, error) {
	bucket := f.Bucket
	prefix := f.Prefix
	if len(prefix) > 0 && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if f.Profile != "" {
		opts.Profile = f.Profile
	}

	config := &opts.Config
	config.CredentialsChainVerboseErrors = aws.Bool(true)
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

	sess := session.Must(session.NewSessionWithOptions(opts))
	b := &Bucket{
		bucket:     bucket,
		prefix:     prefix,
		svc:        s3.New(sess),
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}
	return b, nil
}

func (b *Bucket) Download(oid string, dst io.WriterAt) error {
	if err := b.getFormat(); err != nil {
		return err
	}
	if b.format != bucketFormat {
		return UnsupportBucketFormatErr
	}

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
	if err := b.getFormat(); err != nil {
		return err
	}
	if b.format == "" {
		if err := b.setFormat(); err != nil {
			return err
		}
	}

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

func (b *Bucket) formatKey() (key string) {
	return b.prefix + "hoggle-format"
}

// TODO refactor bucket format checking
func (b *Bucket) getFormat() (err error) {
	if b.format != "" {
		return nil
	}

	resp, err := b.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(b.formatKey()),
		Range:  aws.String("bytes=0-255"),
	})
	if err != nil {
		if e, ok := err.(awserr.Error); ok && e.Code() == s3.ErrCodeNoSuchKey {
			err = nil
		}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if bytes.Equal(body, []byte(bucketFormat)) {
		b.format = string(body)
	} else {
		err = UnsupportBucketFormatErr
	}
	return
}

func (b *Bucket) oid2key(oid string) (key string) {
	if len(oid) < 4 {
		key = fmt.Sprintf("%sobjects/%s", b.prefix, oid)
	} else {
		key = fmt.Sprintf("%sobjects/%s/%s/%s", b.prefix, oid[0:2], oid[2:4], oid)
	}
	return
}

func (b *Bucket) setFormat() (err error) {
	_, err = b.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(b.formatKey()),
		Body:   aws.ReadSeekCloser(strings.NewReader(bucketFormat)),
	})
	if err == nil {
		b.format = bucketFormat
	}
	return
}
