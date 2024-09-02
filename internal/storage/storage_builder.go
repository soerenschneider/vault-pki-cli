package storage

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg/storage/backend"
	sink2 "github.com/soerenschneider/vault-pki-cli/pkg/storage/shape"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	instance *buildContext
	lock     = &sync.Mutex{}
	once     = sync.Once{}
)

// buildContext is responsible for building storage implementation instances. The struct contains shared resources,
// such as clients that can be shared across multiple instances of storage implementations.
type buildContext struct {
	config           *conf.Config
	kubernetesClient *kubernetes.Clientset
}

func InitBuilder(config *conf.Config) *buildContext {
	once.Do(func() {
		if instance == nil {
			instance = &buildContext{
				config: config,
			}
		}
	})

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
	case backend.FsScheme:
		return backend.NewFilesystemStorageFromUri(uri)
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

func CaStorageFromConfig(storageConfig []map[string]string) (*sink2.CaStorage, error) {
	var caVal string
	for _, conf := range storageConfig {
		val, ok := conf[caId]
		if ok {
			caVal = val
		} else {
			log.Info().Msgf("No storage config given for '%s', writing to stdout", caId)
		}
	}

	builder, err := GetBuilder()
	if err != nil {
		return nil, err
	}
	if len(caVal) > 0 {
		storageImpl, err := builder.BuildFromUri(caVal)
		if err != nil {
			return nil, err
		}
		return sink2.NewCaStorage(storageImpl)
	}

	return sink2.NewCaStorage(nil)
}

const (
	certId = "cert"
	keyId  = "key"
	caId   = "ca"
	csrId  = "csr"
	crlId  = "crl"
)

func CrlStorageFromConfig(storageConfig []map[string]string) (*sink2.CrlStorage, error) {
	var crlVal string
	for _, conf := range storageConfig {
		val, ok := conf[crlId]
		if ok {
			crlVal = val
		} else {
			log.Info().Msgf("No storage config given for '%s', writing to stdout", crlId)
		}
	}

	builder, err := GetBuilder()
	if err != nil {
		return nil, err
	}
	if len(crlVal) > 0 {
		storageImpl, err := builder.BuildFromUri(crlVal)
		if err != nil {
			return nil, err
		}
		return sink2.NewCrlStorage(storageImpl)
	}

	return sink2.NewCrlStorage(nil)
}

func CsrStorageFromConfig(storageConfig []map[string]string) (*sink2.CsrStorage, error) {
	var certVal string
	var csrVal string
	var caVal string
	for _, conf := range storageConfig {
		val, ok := conf[certId]
		if !ok {
			return nil, fmt.Errorf("can not build storage, missing '%s' in storage configuration", certId)
		}
		certVal = val

		val, ok = conf[csrId]
		if !ok {
			return nil, fmt.Errorf("can not build storage, missing '%s' in storage configuration", csrId)
		}
		csrVal = val

		val, ok = conf[caId]
		if ok {
			caVal = val
		}
	}

	builder, err := GetBuilder()
	if err != nil {
		return nil, err
	}
	certStorageImpl, err := builder.BuildFromUri(certVal)
	if err != nil {
		return nil, err
	}

	csrStorageImpl, err := builder.BuildFromUri(csrVal)
	if err != nil {
		return nil, err
	}

	var caStorageImpl pki.StorageImplementation
	if len(caVal) > 0 {
		caStorageImpl, err = builder.BuildFromUri(caVal)
		if err != nil {
			return nil, err
		}
	}

	return sink2.NewCsrStorage(certStorageImpl, csrStorageImpl, caStorageImpl)
}

func KeyPairStorageFromConfig(config *conf.Config) ([]*sink2.KeyPairStorage, error) {
	var sinks []*sink2.KeyPairStorage

	for _, conf := range config.StorageConfig {
		sink, err := buildSink(conf)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, sink)
	}

	return sinks, nil
}

func buildSink(conf map[string]string) (*sink2.KeyPairStorage, error) {
	builder, err := GetBuilder()
	if err != nil {
		return nil, err
	}

	var certVal string
	var keyVal string
	var caVal string

	val, ok := conf[certId]
	if ok {
		certVal = val
	}

	val, ok = conf[keyId]
	if !ok {
		return nil, fmt.Errorf("can not build storage, missing '%s' in storage configuration", keyId)
	}
	keyVal = val

	val, ok = conf[caId]
	if ok {
		caVal = val
	}

	var certSink pki.StorageImplementation
	if len(certVal) > 0 {
		certSink, err = builder.BuildFromUri(certVal)
		if err != nil {
			return nil, err
		}
	}

	keySink, err := builder.BuildFromUri(keyVal)
	if err != nil {
		return nil, err
	}

	var caSink pki.StorageImplementation
	if len(caVal) > 0 {
		caSink, err = builder.BuildFromUri(caVal)
		if err != nil {
			return nil, err
		}
	}

	return sink2.NewKeyPairStorage(certSink, keySink, caSink)
}

func MultiKeyPairStorageFromConfig(config *conf.Config) (*sink2.MultiKeyPairStorage, error) {
	sinks, err := KeyPairStorageFromConfig(config)
	if err != nil {
		return nil, err
	}

	return sink2.NewMultiKeyPairSink(sinks...)
}
