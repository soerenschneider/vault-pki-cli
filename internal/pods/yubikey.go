//go:build yubikey

package pods

import (
	"errors"
	"fmt"
	"github.com/go-piv/piv-go/piv"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"strings"
)

type YubikeyPod struct {
	pin     string
	slot    piv.Slot
	yubikey *piv.YubiKey
}

func NewYubikeyPod(slot uint32, pin string) (*YubikeyPod, error) {
	keySlot, err := translateSlot(slot)
	if err != nil {
		return nil, err
	}

	err = verifyPin(pin)
	if err != nil {
		return nil, err
	}

	// List all smartcards connected to the system.
	cards, err := piv.Cards()
	if err != nil {
		return nil, err
	}

	// Find a YubiKey and open the reader.
	var yk *piv.YubiKey
	for _, card := range cards {
		if strings.Contains(strings.ToLower(card), "yubikey") {
			if yk, err = piv.Open(card); err != nil {
				return nil, err
			}
			return &YubikeyPod{
				pin:     pin,
				slot:    *keySlot,
				yubikey: yk,
			}, nil
		}
	}

	return nil, errors.New("no cards found")
}

func (pod *YubikeyPod) getManagementKey() (*[24]byte, error) {
	m, err := pod.yubikey.Metadata(pod.pin)
	if err != nil {
		return nil, err
	}
	if m.ManagementKey == nil {
		return nil, err
	}
	return m.ManagementKey, nil
}

func ValidateSlot(slot uint32) error {
	_, err := translateSlot(slot)
	return err
}

func translateSlot(slot uint32) (*piv.Slot, error) {
	switch slot {
	case 0x9a:
		return &piv.SlotAuthentication, nil
	case 0x9c:
		return &piv.SlotSignature, nil
	case 0x9e:
		return &piv.SlotCardAuthentication, nil
	case 0x9d:
		return &piv.SlotKeyManagement, nil
	default:
		return nil, errors.New("invalid slot")
	}
}

func verifyPin(pin string) error {
	if len(pin) < 6 {
		return fmt.Errorf("supplied pin shorter than 6 characters: %d", len(pin))
	}

	return nil
}

func (pod *YubikeyPod) Read() ([]byte, error) {
	cert, err := pod.yubikey.Certificate(piv.SlotAuthentication)
	if err != nil {
		return nil, err
	}
	return cert.Raw, nil
}

func (pod *YubikeyPod) CanRead() error {
	_, err := pod.yubikey.Certificate(piv.SlotAuthentication)
	return err
}

func (pod *YubikeyPod) Write(data []byte) error {
	cert, err := pkg.ParseCertPem(data)
	if err != nil {
		return err
	}

	managementKey, err := pod.getManagementKey()
	if err != nil {
		return err
	}
	return pod.yubikey.SetCertificate(*managementKey, pod.slot, cert)
}

func (pod *YubikeyPod) CanWrite() error {
	_, err := pod.getManagementKey()
	return err
}
