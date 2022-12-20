package storage

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sSecretStorage struct {
	*K8sConfig
}

const K8sSecretMapScheme = "k8s-sec"

func (fs *K8sSecretStorage) Read() ([]byte, error) {
	secret, err := fs.client.CoreV1().Secrets(fs.Namespace).Get(context.TODO(), fs.Name, meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	cert, ok := secret.Data[keyCert]
	if !ok {
		return nil, fmt.Errorf("kubernetes secret '%s' does not contain a certifcate", fs.Name)
	}

	return cert, nil
}

func (fs *K8sSecretStorage) CanRead() error {
	_, err := fs.client.CoreV1().Secrets(fs.Namespace).Get(context.TODO(), fs.Name, meta.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (fs *K8sSecretStorage) Write(data []byte) error {
	secret := &v1.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: fs.Name,
			Labels: map[string]string{
				"app": "k8s-cli-pki",
			},
		},
		StringData: map[string]string{
			"data": string(data),
		},
	}

	_, err := fs.client.CoreV1().Secrets(fs.Namespace).Create(context.TODO(), secret, meta.CreateOptions{})
	return err
}

func (fs *K8sSecretStorage) CanWrite() error {
	return nil
}

func NewK8sSecretStorageFromUri(uri string) (*K8sSecretStorage, error) {
	conf, err := K8sConfigFromUri(uri)
	if err != nil {
		return nil, err
	}

	return &K8sSecretStorage{
		conf,
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
