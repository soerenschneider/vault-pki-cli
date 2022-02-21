package issue_strategies

import "crypto/x509"

type StaticRenewal struct {
	Decision bool
}

func NewStaticRenewal(renew bool) (*StaticRenewal, error) {
	return &StaticRenewal{renew}, nil
}

func (p *StaticRenewal) Renew(cert *x509.Certificate) (bool, error) {
	return p.Decision, nil
}
