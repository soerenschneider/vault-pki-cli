package vault

import (
	"github.com/hashicorp/vault/api"
	"golang.org/x/net/context"
)

type TokenAuth struct {
	token string
}

func NewTokenAuth(token string) (*TokenAuth, error) {
	return &TokenAuth{token}, nil
}

func (t *TokenAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	ret := &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: t.token,
		},
	}

	return ret, nil
}

func (t *TokenAuth) Cleanup(ctx context.Context, client *api.Client) error {
	return nil
}
