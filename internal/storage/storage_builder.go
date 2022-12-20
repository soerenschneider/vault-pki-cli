package storage

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"net/url"
)

func BuildFromUri(uri string) (pki.StorageImplementation, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		log.Error().Msgf("can not build storage from uri '%s': %v", uri, err)
		return nil, err
	}

	switch parsed.Scheme {
	case FsScheme:
		return NewFilesystemStorageFromUri(uri)
	case K8sConfigMapScheme:
		return NewK8sConfigmapStorageFromUri(uri)
	case K8sSecretMapScheme:
		return NewK8sSecretStorageFromUri(uri)
	default:
		return nil, fmt.Errorf("can not build storage impl for unknown scheme '%s'", parsed.Scheme)
	}
}
