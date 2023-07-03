package pki

import "errors"

var ErrNoCertFound = errors.New("data not found")

// StorageImplementation is a simple wrapper around a key artifact (cert, key, ca, crl, csr). This enables decoupling
// from the actual resource (file-based, kubernetes, network, ..) and make it interchangeable.
type StorageImplementation interface {
	Read() ([]byte, error)
	CanRead() error
	Write([]byte) error
	CanWrite() error
}
