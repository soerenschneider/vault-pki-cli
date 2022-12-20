package storage

import (
	"fmt"
)

type BufferPod struct {
	Data []byte
}

func (b *BufferPod) Read() ([]byte, error) {
	if len(b.Data) > 0 {
		return b.Data, nil
	}
	return nil, fmt.Errorf("empty buffer")
}

func (b *BufferPod) CanRead() error {
	if len(b.Data) > 0 {
		return nil
	}
	return fmt.Errorf("empty buffer")
}

func (b *BufferPod) Write(data []byte) error {
	b.Data = data
	return nil
}

func (b *BufferPod) CanWrite() error {
	return nil
}
