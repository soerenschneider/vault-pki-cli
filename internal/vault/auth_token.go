package vault

type TokenAuth struct {
	token string
}

func NewTokenAuth(token string) (*TokenAuth, error) {
	return &TokenAuth{token}, nil
}

func (t *TokenAuth) Authenticate() (string, error) {
	return t.token, nil
}
