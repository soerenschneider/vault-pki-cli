package pods

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"golang.org/x/sys/unix"
)

type FsPod struct {
	FilePath  string
	FileOwner *int
	FileGroup *int
}

func NewFsPod(path, owner, group string) (*FsPod, error) {
	if len(path) == 0 {
		return nil, errors.New("empty path provided")
	}

	var uid, gid *int
	if len(owner) > 0 && len(group) > 0 {
		localUser, err := user.Lookup(owner)
		if err != nil {
			return nil, fmt.Errorf("could not lookup user '%s': %v", owner, err)
		}
		*uid, err = strconv.Atoi(localUser.Uid)
		if err != nil {
			return nil, fmt.Errorf("was expecting a numerical uid, got '%s'", localUser.Uid)
		}

		localGroup, err := user.LookupGroup(group)
		if err != nil {
			return nil, fmt.Errorf("could not lookup group '%s': %v", group, err)
		}
		*gid, err = strconv.Atoi(localGroup.Gid)
		if err != nil {
			return nil, fmt.Errorf("was expecting a numerical gid, got '%s'", localGroup.Gid)
		}
	}

	return &FsPod{
		FilePath:  path,
		FileOwner: uid,
		FileGroup: gid,
	}, nil
}

func (fs *FsPod) Read() ([]byte, error) {
	return ioutil.ReadFile(fs.FilePath)
}

func (fs *FsPod) CanRead() error {
	_, err := os.Stat(fs.FilePath)
	return err
}

func (fs *FsPod) Write(signedData string) error {
	err := ioutil.WriteFile(fs.FilePath, []byte(signedData), 0640)
	if err != nil {
		return fmt.Errorf("could not write file '%s' to disk: %v", fs.FilePath, err)
	}

	if fs.FileOwner != nil && fs.FileGroup != nil {
		err = os.Chown(fs.FilePath, *fs.FileOwner, *fs.FileGroup)
		if err != nil {
			return fmt.Errorf("could not chown file '%s': %v", fs.FilePath, err)
		}
	}

	return nil
}

func (fs *FsPod) CanWrite() error {
	dir := filepath.Dir(fs.FilePath)
	return unix.Access(dir, unix.W_OK)
}
