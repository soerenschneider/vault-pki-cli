package pkg

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

func ParseCertPem(content []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, errors.New("could not decode invalid cert data")
	}
	return x509.ParseCertificate(block.Bytes)
}

func GetFormattedSerial(content []byte) (string, error) {
	cert, err := ParseCertPem(content)
	if err != nil {
		return "", fmt.Errorf("could not parse certificate: %v", err)
	}

	return FormatSerial(cert.SerialNumber), nil
}

func FormatSerial(i *big.Int) string {
	hex := fmt.Sprintf("%x", i)
	if len(hex)%2 == 1 {
		hex = "0" + hex
	}
	re := regexp.MustCompile("..")
	return strings.TrimRight(re.ReplaceAllString(hex, "$0:"), ":")
}
