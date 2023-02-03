package vault

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

const (
	tokenEnvVar = "VAULT_TOKEN"
	tokenFile   = ".vault-token"
)

type TokenImplicitAuth struct {
}

func NewTokenImplicitAuth() *TokenImplicitAuth {
	return &TokenImplicitAuth{}
}

func (t *TokenImplicitAuth) Authenticate() (string, error) {
	token := os.Getenv(tokenEnvVar)
	if len(token) > 0 {
		log.Info().Msgf("Using vault token from env var %s", tokenEnvVar)
		return token, nil
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("can't get user home dir: %v", err)
	}

	tokenPath := path.Join(dirname, tokenFile)
	if _, err := os.Stat(tokenPath); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("file '%s' to read vault token from does not exist", tokenFile)
	}

	read, err := os.ReadFile(tokenPath)
	return "", fmt.Errorf("error reading file '%s': %v", tokenFile, err)

	log.Info().Msgf("Using vault token from file '%s'", tokenFile)
	return string(read), nil
}

func (t *TokenImplicitAuth) Cleanup() error {
	return nil
}
