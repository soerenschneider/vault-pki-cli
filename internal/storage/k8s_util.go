package storage

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	"net/url"
	"strings"
)

type K8sConfig struct {
	Namespace string
	Name      string
	client    *kubernetes.Clientset
}

type K8sOption func(sink *K8sConfig)

func WithNamespace(namespace string) K8sOption {
	return func(h *K8sConfig) {
		h.Namespace = namespace
	}
}

func WithName(secretName string) K8sOption {
	return func(h *K8sConfig) {
		h.Name = secretName
	}
}

func K8sConfigFromUri(uri string) (*K8sConfig, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	path := parsed.Path
	split := strings.Split(path, "/")
	if len(split) < 3 {
		return nil, errors.New("can't build k8s secret storage, not enough information provided, please supply namespace and name")
	}

	impl := &K8sConfig{}
	impl.Namespace = split[1]
	impl.Name = split[2]

	return impl, err
}
