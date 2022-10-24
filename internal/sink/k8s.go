package sink

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sSink struct {
	client     *kubernetes.Clientset
	Namespace  string
	SecretName string
}

type K8sOption func(sink *K8sSink)

func WithNamespace(namespace string) K8sOption {
	return func(h *K8sSink) {
		h.Namespace = namespace
	}
}

func WithSecretName(secretName string) K8sOption {
	return func(h *K8sSink) {
		h.SecretName = secretName
	}
}

func NewK8sBackend(opts ...K8sOption) (*K8sSink, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	/*
		config, err := clientcmd.BuildConfigFromFlags("", "/Users/C5338334/.kube/config")
		if err != nil {
			log.Fatal(err)
		}

	*/

	k8sBackend := &K8sSink{client: clientset}
	for _, opt := range opts {
		opt(k8sBackend)
	}

	return k8sBackend, err
}

const (
	keyCa         = "ca"
	keyCert       = "cert"
	keyPrivateKey = "privatekey"
)

func (f *K8sSink) Write(certData *pki.CertData) error {
	data := map[string]string{
		keyCert:       string(certData.Certificate),
		keyPrivateKey: string(certData.PrivateKey),
	}

	if certData.HasCaChain() {
		data[keyCa] = string(certData.CaChain)
	}

	secret := &v1.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: f.SecretName,
			Labels: map[string]string{
				"app": "k8s-cli-pki",
			},
		},
		StringData: data,
	}

	_, err := f.client.CoreV1().Secrets(f.Namespace).Create(context.TODO(), secret, meta.CreateOptions{})
	return err
}

func (f *K8sSink) Read() (*x509.Certificate, error) {
	secret, err := f.client.CoreV1().Secrets(f.Namespace).Get(context.TODO(), f.SecretName, meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	cert, ok := secret.Data[keyCert]
	if !ok {
		return nil, fmt.Errorf("kubernetes secret '%s' does not contain a certifcate", f.SecretName)
	}

	return pkg.ParseCertPem(cert)
}
