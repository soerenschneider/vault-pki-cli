package storage

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const K8sConfigMapScheme = "k8s-cm"

type K8sConfigmapStorage struct {
	*K8sConfig
}

func NewK8sConfigmapStorageFromUri(uri string) (*K8sConfigmapStorage, error) {
	conf, err := K8sConfigFromUri(uri)
	if err != nil {
		return nil, err
	}

	return &K8sConfigmapStorage{
		conf,
	}, nil
}

func NewK8sConfigmapStorage(client *kubernetes.Clientset, name string, opts ...K8sOption) (*K8sConfigmapStorage, error) {
	if client == nil {
		return nil, errors.New("empty client provided")
	}

	k8sBackend := &K8sConfigmapStorage{}
	k8sBackend.client = client
	for _, opt := range opts {
		opt(k8sBackend.K8sConfig)
	}

	return k8sBackend, nil
}

func (fs *K8sConfigmapStorage) Read() ([]byte, error) {
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

func (fs *K8sConfigmapStorage) CanRead() error {
	_, err := fs.client.CoreV1().Secrets(fs.Namespace).Get(context.TODO(), fs.Name, meta.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (fs *K8sConfigmapStorage) Write(data []byte) error {
	secret := &v1.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name: fs.Name,
			Labels: map[string]string{
				"app": "k8s-cli-pki",
			},
		},
		Data: map[string]string{
			"data": string(data),
		},
	}

	_, err := fs.client.CoreV1().ConfigMaps(fs.Namespace).Create(context.TODO(), secret, meta.CreateOptions{})
	return err
}

func (fs *K8sConfigmapStorage) CanWrite() error {
	return nil
}

const (
	keyCa         = "ca"
	keyCert       = "cert"
	keyPrivateKey = "privatekey"
)
