package pods

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

type FsPod struct {
	FilePath string
}

func (fs *FsPod) Read() ([]byte, error) {
	return ioutil.ReadFile(fs.FilePath)
}

func (fs *FsPod) CanRead() error {
	_, err := os.Stat(fs.FilePath)
	return err
}

func (fs *FsPod) Write(signedData string) error {
	return ioutil.WriteFile(fs.FilePath, []byte(signedData), 0640)
}

func (fs *FsPod) CanWrite() error {
	dir := filepath.Dir(fs.FilePath)
	return unix.Access(dir, unix.W_OK)
}
