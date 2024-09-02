package backend

import (
	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
)

type BufferPod struct {
	Data []byte
}

func (b *BufferPod) Read() ([]byte, error) {
	if len(b.Data) > 0 {
		return b.Data, nil
	}
	return nil, pki.ErrNoCertFound
}

func (b *BufferPod) CanRead() error {
	if len(b.Data) > 0 {
		return nil
	}

	return pki.ErrNoCertFound
}

func (b *BufferPod) Write(data []byte) error {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	b.Data = data
	return nil
}

func (b *BufferPod) CanWrite() error {
	return nil
}
