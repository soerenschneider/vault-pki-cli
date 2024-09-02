package shape

import (
	"fmt"

	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
)

type CrlStorage struct {
	storage pki.StorageImplementation
}

func NewCrlStorage(storage pki.StorageImplementation) (*CrlStorage, error) {
	return &CrlStorage{
		storage: storage,
	}, nil
}

func (out *CrlStorage) WriteCrl(crlData []byte) error {
	if out.storage == nil {
		fmt.Println(string(crlData))
		return nil
	}

	return out.storage.Write(crlData)
}
