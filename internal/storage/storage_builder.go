package storage

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/url"
	"sync"
)

var (
	instance *buildContext
	lock     = &sync.Mutex{}
)

// buildContext is responsible for building storage implementation instances. The struct contains shared resources,
// such as clients that can be shared across multiple instances of storage implementations.
type buildContext struct {
	config           *conf.Config
	kubernetesClient *kubernetes.Clientset
}

func InitBuilder(config *conf.Config) *buildContext {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &buildContext{
				config: config,
			}
		}
	}

	return instance
}

func GetBuilder() (*buildContext, error) {
	if instance == nil {
		return nil, errors.New("buildContext not initialized yet")
	}

	return instance, nil
}

func (context *buildContext) BuildFromUri(uri string) (pki.StorageImplementation, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		log.Error().Msgf("can not build storage from uri '%s': %v", uri, err)
		return nil, err
	}

	switch parsed.Scheme {
	case FsScheme:
		return NewFilesystemStorageFromUri(uri)
	case K8sConfigMapScheme:
		return NewK8sConfigmapStorageFromUri(uri, context)
	case K8sSecretMapScheme:
		return NewK8sSecretStorageFromUri(uri, context)
	default:
		return nil, fmt.Errorf("can not build storage impl for unknown scheme '%s'", parsed.Scheme)
	}
}

func (context *buildContext) KubernetesClient() (*kubernetes.Clientset, error) {
	if context.kubernetesClient == nil {
		lock.Lock()
		defer lock.Unlock()
		if context.kubernetesClient == nil {
			log.Info().Msg("Building kubernetes client")
			config, err := rest.InClusterConfig()
			if err != nil {
				return nil, fmt.Errorf("could not build kubernetes client config: %v", err)
			}
			client, err := kubernetes.NewForConfig(config)
			if err != nil {
				return nil, fmt.Errorf("could not build kubernetes client: %v", err)
			}
			context.kubernetesClient = client
		}
	}

	return context.kubernetesClient, nil
}
