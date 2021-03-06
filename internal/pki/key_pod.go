package pki

// KeyPod is a simple wrapper around a key (which is just a byte stream itself). This way, we decouple
// the implementation (file-based, memory, network, ..) and make it easily swap- and testable.
type KeyPod interface {
	Read() ([]byte, error)
	CanRead() error
	Write([]byte) error
	CanWrite() error
}
