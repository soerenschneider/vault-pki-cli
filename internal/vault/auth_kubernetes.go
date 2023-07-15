package vault

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/kubernetes"
	"golang.org/x/net/context"
)

const (
	defaultServiceAccountTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token" // #nosec G101
	defaultMount                   = "kubernetes"
)

type KubernetesAuth struct {
	client                  *api.Client
	role                    string
	serviceAccountTokenFile string
	mount                   string
}

func NewVaultKubernetesAuth(client *api.Client, role string) (*KubernetesAuth, error) {
	if client == nil {
		return nil, errors.New("empty client provided")
	}

	return &KubernetesAuth{
		client:                  client,
		role:                    role,
		mount:                   defaultMount,
		serviceAccountTokenFile: defaultServiceAccountTokenFile,
	}, nil
}

func (t *KubernetesAuth) Cleanup() error {
	path := "auth/token/revoke-self"
	_, err := t.client.Logical().Write(path, map[string]interface{}{})
	return err
}

func (t *KubernetesAuth) Authenticate() (string, error) {

	k8sAuth, err := kubernetes.NewKubernetesAuth(
		t.role,
		kubernetes.WithServiceAccountTokenPath(t.serviceAccountTokenFile),
		kubernetes.WithMountPath(t.mount))

	if err != nil {
		return "", fmt.Errorf("unable to initialize Kubernetes kubernetes method: %w", err)
	}

	authInfo, err := t.client.Auth().Login(context.TODO(), k8sAuth)
	if err != nil {
		return "", fmt.Errorf("unable to log in with Kubernetes kubernetes: %w", err)
	}
	if authInfo == nil {
		return "", fmt.Errorf("no kubernetes info was returned after login")
	}

	return authInfo.Auth.ClientToken, nil
}
