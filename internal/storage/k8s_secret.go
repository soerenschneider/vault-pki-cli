package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sSecretStorage struct {
	*K8sConfig
	client *kubernetes.Clientset
}

const K8sSecretMapScheme = "k8s-sec"

func NewK8sSecretStorageFromUri(uri string, context *buildContext) (*K8sSecretStorage, error) {
	conf, err := K8sConfigFromUri(uri)
	if err != nil {
		return nil, err
	}

	client, err := context.KubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("could not build kubernetes client: %v", err)
	}

	return &K8sSecretStorage{
		K8sConfig: conf,
		client:    client,
	}, nil
}

func NewK8sSecretStorage(client *kubernetes.Clientset, opts ...K8sOption) (*K8sSecretStorage, error) {
	if client == nil {
		return nil, errors.New("empty client provided")
	}

	k8sBackend := &K8sSecretStorage{}
	k8sBackend.client = client

	for _, opt := range opts {
		opt(k8sBackend.K8sConfig)
	}

	return k8sBackend, nil
}

func (fs *K8sSecretStorage) Read() ([]byte, error) {
	if fs.client == nil {
		return nil, errors.New("can't read secret, uninitialized k8s client")
	}

	secret, err := fs.client.CoreV1().Secrets(fs.Namespace).Get(context.TODO(), fs.Name, meta.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil, pki.ErrNoCertFound
		}
		return nil, err
	}

	cert, ok := secret.Data[keyPrivateKey]
	if !ok {
		return nil, fmt.Errorf("kubernetes secret '%s' does not contain a certifcate", fs.Name)
	}

	return cert, nil
}

func (fs *K8sSecretStorage) CanRead() error {
	if fs.client == nil {
		return errors.New("can't read secret, uninitialized k8s client")
	}

	_, err := fs.client.CoreV1().Secrets(fs.Namespace).Get(context.TODO(), fs.Name, meta.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return pki.ErrNoCertFound
		}
		return err
	}

	return nil
}

func (fs *K8sSecretStorage) Write(data []byte) error {
	if fs.client == nil {
		return errors.New("can't read secret, uninitialized k8s client")
	}

	secret := &v1.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: fs.Name,
			Labels: map[string]string{
				"app": "k8s-cli-pki",
			},
		},
		StringData: map[string]string{
			keyPrivateKey: string(data),
		},
	}

	if fs.CanRead() == nil {
		_, err := fs.client.CoreV1().Secrets(fs.Namespace).Update(context.TODO(), secret, meta.UpdateOptions{})
		return err
	}

	_, err := fs.client.CoreV1().Secrets(fs.Namespace).Create(context.TODO(), secret, meta.CreateOptions{})
	return err
}

func (fs *K8sSecretStorage) CanWrite() error {
	if fs.client == nil {
		return errors.New("can't read secret, uninitialized k8s client")
	}

	return nil
}
