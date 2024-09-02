package pkg

type Signature struct {
	Certificate []byte
	CaData      []byte
	Serial      string
}

func (cert *Signature) HasCaData() bool {
	return len(cert.CaData) > 0
}
