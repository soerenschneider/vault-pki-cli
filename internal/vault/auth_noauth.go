package vault

import (
	"errors"

	"github.com/hashicorp/vault/api"
	"golang.org/x/net/context"
)

type NoAuth struct {
}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (t *NoAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	return nil, errors.New("no auth")
}

func (t *NoAuth) Cleanup(ctx context.Context, client *api.Client) error {
	return nil
}
