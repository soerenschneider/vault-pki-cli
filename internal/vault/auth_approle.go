package vault

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"golang.org/x/net/context"
)

const (
	KeyRoleId       = "role_id"
	KeySecretId     = "secret_id"
	KeySecretIdFile = "secret_id_file"
)

type AppRoleAuth struct {
	approleMount string
	loginData    map[string]string
}

func NewAppRoleAuth(loginData map[string]string, mountPath string) (*AppRoleAuth, error) {
	return &AppRoleAuth{
		loginData:    loginData,
		approleMount: mountPath,
	}, nil
}

func (t *AppRoleAuth) getLoginTuple() (string, string, error) {
	var roleId, secretId string
	val, ok := t.loginData[KeyRoleId]
	if ok && len(val) > 0 {
		roleId = val
	}

	val, ok = t.loginData[KeySecretId]
	if ok && len(val) > 0 {
		secretId = val
	}

	val, ok = t.loginData[KeySecretIdFile]
	if ok && len(val) > 0 {
		data, err := os.ReadFile(val)
		if err != nil {
			return "", "", fmt.Errorf("could not read secret_id from file '%s': %v", val, err)
		}
		secretId = string(data)
	}

	return roleId, secretId, nil
}

func (t *AppRoleAuth) Cleanup(ctx context.Context, client *api.Client) error {
	path := "auth/token/revoke-self"
	_, err := client.Logical().Write(path, map[string]interface{}{})
	return err
}

func (t *AppRoleAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	roleId, secretId, err := t.getLoginTuple()
	if err != nil {
		return nil, fmt.Errorf("could not get login data: %v", err)
	}

	path := fmt.Sprintf("auth/%s/login", t.approleMount)
	data := map[string]interface{}{
		KeyRoleId:   roleId,
		KeySecretId: secretId,
	}
	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
