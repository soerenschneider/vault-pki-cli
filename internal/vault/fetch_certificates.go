package vault

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
)

func FetchCert(vaultAddress, pkiMount string, binary bool) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/v1/%s/ca", vaultAddress, pkiMount)
	if !binary {
		endpoint = fmt.Sprintf("%s/v1/%s/ca/pem", vaultAddress, pkiMount)
	}

	return fetchResource(endpoint)
}

func FetchCertChain(vaultAddress, pkiMount string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/v1/%s/ca_chain", vaultAddress, pkiMount)
	return fetchResource(endpoint)
}

func FetchCrl(vaultAddress, pkiMount string, binary bool) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/v1/%s/crl", vaultAddress, pkiMount)
	if !binary {
		endpoint += "/pem"
	}
	return fetchResource(endpoint)
}

func fetchResource(url string) ([]byte, error) {
	client := retryablehttp.NewClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %v", err)
	}

	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	return read, nil
}
