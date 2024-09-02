package vault

import (
	"github.com/hashicorp/vault/api"
	"golang.org/x/net/context"
)

type NoAuth struct {
}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (t *NoAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	return nil, nil
}
