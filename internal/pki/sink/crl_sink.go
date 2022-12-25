package sink

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
)

type CrlSink struct {
	storage pki.StorageImplementation
}

const crlId = "crl"

func CrlSinkFromConfig(storageConfig []map[string]string) (*CrlSink, error) {
	var crlVal string
	for _, conf := range storageConfig {
		val, ok := conf[crlId]
		if ok {
			crlVal = val
		} else {
			log.Info().Msgf("No storage config given for '%s', writing to stdout", crlId)
		}
	}

	builder, err := storage.GetBuilder()
	if err != nil {
		return nil, err
	}
	if len(crlVal) > 0 {
		storageImpl, err := builder.BuildFromUri(crlVal)
		if err != nil {
			return nil, err
		}
		return &CrlSink{storageImpl}, nil
	}

	return &CrlSink{nil}, nil
}

func (out *CrlSink) WriteCrl(crlData []byte) error {
	if out.storage == nil {
		fmt.Println(string(crlData))
		return nil
	}

	return out.storage.Write(crlData)
}
