package pkg

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

func ParseCertPem(data []byte) (*x509.Certificate, error) {
	if data == nil {
		return nil, errors.New("emtpy data provided")
	}

	var der *pem.Block
	rest := data
	for {
		der, rest = pem.Decode(rest)
		if der == nil {
			return nil, errors.New("invalid pem provided")
		}

		if strings.Contains(der.Type, "PRIVATE KEY") {
			continue
		}

		cert, err := x509.ParseCertificate(der.Bytes)
		if err != nil {
			return nil, err
		}

		if !cert.IsCA {
			return cert, nil
		}
	}
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

func IsCertExpired(cert x509.Certificate) bool {
	return time.Now().After(cert.NotAfter)
}
