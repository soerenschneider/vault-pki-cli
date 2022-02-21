package issue_strategies

import (
	"crypto/x509"
	"errors"
	"fmt"
	"math"
	"time"

	log "github.com/rs/zerolog/log"
)

type Percentage struct {
	PercentageMinThreshold float64
}

func NewPercentage(percentage float64) (*Percentage, error) {
	if percentage < 10 || percentage > 80 {
		return nil, fmt.Errorf("percentage must be between [10, 80]")
	}

	return &Percentage{PercentageMinThreshold: percentage}, nil
}

func (p *Percentage) Renew(cert *x509.Certificate) (bool, error) {
	if cert == nil {
		return true, errors.New("empty certificate provided")
	}

	from := cert.NotBefore
	expiry := cert.NotAfter

	secondsTotal := expiry.Sub(from).Seconds()
	durationUntilExpiration := time.Until(expiry)

	percentage := math.Max(0, durationUntilExpiration.Seconds()*100./secondsTotal)
	log.Info().Msgf("Lifetime at %.2f%%, %s left (valid from '%v', until '%v')", percentage, durationUntilExpiration.Round(time.Second), from, expiry)

	return percentage <= p.PercentageMinThreshold, nil
}
