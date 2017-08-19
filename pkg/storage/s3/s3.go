package s3

import (
	"io"

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
}

type Bucket struct {
	bucket     string
	prefix     string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func (f *Factory) New() (storage.Storage, error) {
	config := &aws.Config{}
	if f.Endpoint != "" {
		config.Region = aws.String(f.Region)
		config.Endpoint = aws.String(f.Endpoint)
		config.DisableSSL = aws.Bool(true)
		config.S3ForcePathStyle = aws.Bool(true)
	}
	sess := session.Must(session.NewSession(config))
	b := &Bucket{
		bucket:     f.Bucket,
		prefix:     f.Prefix,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}
	return b, nil
}

func (b *Bucket) Download(oid string, dst io.WriterAt) error {
	params := &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(b.prefix + oid),
	}
	if _, err := b.downloader.Download(dst, params); err != nil {
		return err
	} else {
		return nil
	}
}

func (b *Bucket) Upload(oid string, src io.ReadSeeker) error {
	params := &s3manager.UploadInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(b.prefix + oid),
		Body:   src,
	}
	if _, err := b.uploader.Upload(params); err != nil {
		return err
	} else {
		return nil
	}
}
