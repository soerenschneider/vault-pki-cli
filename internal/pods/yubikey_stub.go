//go:build !yubikey

package pods

import (
	"github.com/pkg/errors"
)

type YubikeyPod struct {
}

func NewYubikeyPod(slot uint32, pin string) (*YubikeyPod, error) {
	return nil, errors.New("this build has no yubikey support!")
}

func (pod *YubikeyPod) Read() ([]byte, error) {
	return nil, errors.New("this build has no yubikey support!")
}

func (pod *YubikeyPod) CanRead() error {
	return errors.New("this build has no yubikey support!")
}

func (pod *YubikeyPod) Write(data []byte) error {
	return errors.New("this build has no yubikey support!")
}

func (pod *YubikeyPod) CanWrite() error {
	return errors.New("this build has no yubikey support!")
}

func ValidateSlot(slot uint32) error {
	return nil
}
