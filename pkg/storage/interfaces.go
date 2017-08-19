package storage

import "io"

type Factory interface {
	New() (Storage, error)
}

type Storage interface {
	Download(oid string, dst io.WriterAt) error
	Upload(oid string, src io.ReadSeeker) error
}
