package renew_strategy

import (
	"crypto/x509"
	"errors"
	"fmt"
	"math"
	"time"
)

type Percentage struct {
	PercentageMinThreshold float32
}

func NewPercentage(percentage float32) (*Percentage, error) {
	if percentage < 20 || percentage > 80 {
		return nil, fmt.Errorf("percentage must be between [20, 80]")
	}

	return &Percentage{PercentageMinThreshold: percentage}, nil
}

func (p *Percentage) Renew(cert *x509.Certificate) (bool, error) {
	if cert == nil {
		return true, errors.New("empty certificate provided")
	}

	percentage := GetPercentage(*cert)
	return percentage <= p.PercentageMinThreshold, nil
}

func GetPercentage(cert x509.Certificate) float32 {
	from := cert.NotBefore
	expiry := cert.NotAfter

	secondsTotal := expiry.Sub(from).Seconds()
	durationUntilExpiration := time.Until(expiry)

	return float32(math.Max(0, durationUntilExpiration.Seconds()*100./secondsTotal))
}
