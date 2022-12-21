package sink

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
)

// CsrSink accepts CA data to write to the configured storage implementation.
type CaSink struct {
	storage pki.StorageImplementation
}

func CaSinkFromConfig(storageConfig []map[string]string) (*CaSink, error) {
	var caVal string
	for _, conf := range storageConfig {
		val, ok := conf[caId]
		if ok {
			caVal = val
		} else {
			log.Info().Msgf("No storage config given for '%s', writing to stdout", caId)
		}
	}

	builder, err := storage.GetBuilder()
	if err != nil {
		return nil, err
	}
	if len(caVal) > 0 {
		storageImpl, err := builder.BuildFromUri(caVal)
		if err != nil {
			return nil, err
		}
		return &CaSink{storageImpl}, nil
	}

	return &CaSink{nil}, nil
}

func (out *CaSink) WriteCa(certData []byte) error {
	if out.storage == nil {
		fmt.Println(string(certData))
		return nil
	}

	return out.storage.Write(certData)
}
