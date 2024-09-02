package shape

import (
	"fmt"

	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
)

// CaStorage accepts CA data to write to the configured storage implementation.
type CaStorage struct {
	storage pki.StorageImplementation
}

func NewCaStorage(storage pki.StorageImplementation) (*CaStorage, error) {
	return &CaStorage{
		storage: storage,
	}, nil
}

func (out *CaStorage) WriteCa(certData []byte) error {
	if out.storage == nil {
		fmt.Println(string(certData))
		return nil
	}

	return out.storage.Write(certData)
}
