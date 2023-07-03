package vault

import "errors"

type NoAuth struct {
}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (t *NoAuth) Authenticate() (string, error) {
	return "", errors.New("no auth")
}

func (t *NoAuth) Cleanup() error {
	return nil
}
