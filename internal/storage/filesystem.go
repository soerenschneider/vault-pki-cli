package storage

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"golang.org/x/sys/unix"
)

type FilesystemStorage struct {
	FilePath  string
	FileOwner *int
	FileGroup *int
}

const FsScheme = "file"

func NewFilesystemStorageFromUri(uri string) (*FilesystemStorage, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	path := parsed.Path

	var username, pass string
	userData := parsed.User
	if userData != nil {
		username = userData.Username()

		var ok bool
		pass, ok = userData.Password()
		if !ok {
			pass = ""
		}
	}

	return NewFilesystemStorage(path, username, pass)
}

func NewFilesystemStorage(path, owner, group string) (*FilesystemStorage, error) {
	if len(path) == 0 {
		return nil, errors.New("empty path provided")
	}

	var uid, gid *int
	if len(owner) > 0 && len(group) > 0 {
		localUser, err := user.Lookup(owner)
		if err != nil {
			return nil, fmt.Errorf("could not lookup user '%s': %v", owner, err)
		}
		conv, err := strconv.Atoi(localUser.Uid)
		if err != nil {
			return nil, fmt.Errorf("was expecting a numerical uid, got '%s'", localUser.Uid)
		}
		uid = &conv

		localGroup, err := user.LookupGroup(group)
		if err != nil {
			return nil, fmt.Errorf("could not lookup group '%s': %v", group, err)
		}
		conv, err = strconv.Atoi(localGroup.Gid)
		if err != nil {
			return nil, fmt.Errorf("was expecting a numerical gid, got '%s'", localGroup.Gid)
		}
		gid = &conv
	}

	return &FilesystemStorage{
		FilePath:  path,
		FileOwner: uid,
		FileGroup: gid,
	}, nil
}

func (fs *FilesystemStorage) Read() ([]byte, error) {
	return os.ReadFile(fs.FilePath)
}

func (fs *FilesystemStorage) CanRead() error {
	_, err := os.Stat(fs.FilePath)
	return err
}

func (fs *FilesystemStorage) Write(signedData []byte) error {
	err := os.WriteFile(fs.FilePath, signedData, 0640)
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

func (fs *FilesystemStorage) CanWrite() error {
	dir := filepath.Dir(fs.FilePath)
	return unix.Access(dir, unix.W_OK)
}
