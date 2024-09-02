package pkg

type SignatureArgs struct {
	CommonName string
	Ttl        string
	IpSans     []string
	AltNames   []string
}

type IssueArgs struct {
	CommonName string
	Ttl        string
	IpSans     []string
	AltNames   []string
}
