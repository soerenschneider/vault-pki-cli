package issue_strategies

import "crypto/x509"

type IssueStrategy interface {
	Renew(cert *x509.Certificate) (bool, error)
}
