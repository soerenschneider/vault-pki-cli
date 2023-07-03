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

const K8sConfigMapScheme = "k8s-cm"

type K8sConfigmapStorage struct {
	*K8sConfig
	client *kubernetes.Clientset
}

func NewK8sConfigmapStorageFromUri(uri string, context *buildContext) (*K8sConfigmapStorage, error) {
	conf, err := K8sConfigFromUri(uri)
	if err != nil {
		return nil, err
	}

	client, err := context.KubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("could not build kubernetes client: %v", err)
	}

	return &K8sConfigmapStorage{
		K8sConfig: conf,
		client:    client,
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
	configMap, err := fs.client.CoreV1().ConfigMaps(fs.Namespace).Get(context.TODO(), fs.Name, meta.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil, pki.ErrNoCertFound
		}
		return nil, err
	}

	cert, ok := configMap.Data[keyCert]
	if !ok {
		return nil, fmt.Errorf("kubernetes configMap '%s' does not contain a certifcate", fs.Name)
	}

	return []byte(cert), nil
}

func (fs *K8sConfigmapStorage) CanRead() error {
	_, err := fs.client.CoreV1().ConfigMaps(fs.Namespace).Get(context.TODO(), fs.Name, meta.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return pki.ErrNoCertFound
		}
		return err
	}

	return nil
}

func (fs *K8sConfigmapStorage) Write(data []byte) error {
	configmap := &v1.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name: fs.Name,
			Labels: map[string]string{
				"app": "k8s-cli-pki",
			},
		},
		Data: map[string]string{
			keyCert: string(data),
		},
	}

	if fs.CanRead() == nil {
		_, err := fs.client.CoreV1().ConfigMaps(fs.Namespace).Update(context.TODO(), configmap, meta.UpdateOptions{})
		return err
	}

	_, err := fs.client.CoreV1().ConfigMaps(fs.Namespace).Create(context.TODO(), configmap, meta.CreateOptions{})
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
