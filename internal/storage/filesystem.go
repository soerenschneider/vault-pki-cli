package storage

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"

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
	Mode      os.FileMode
}

const (
	FsScheme   = "file"
	ParamChmod = "chmod"
)

var defaultMode os.FileMode = 0700

func NewFilesystemStorageFromUri(uri string) (*FilesystemStorage, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	path := parsed.Path
	if parsed.Host == "~" || parsed.Host == "$HOME" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("tried to expand '%s' but homeDir could not be detected: %v", parsed.Host, err)
		}

		orig := filepath.Join(parsed.Host, path)
		path = filepath.Join(homeDir, orig[len(parsed.Host):])
		log.Info().Msgf("Expanded path '%s' to '%s'", orig, path)
	} else if len(parsed.Host) > 0 {
		return nil, fmt.Errorf("invalid syntax for uri, no host expected: '%s' (did you forget the leading '/'?)", uri)
	}

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

	mode := defaultMode
	params, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("could not parse queries")
	}
	for key, val := range params {
		if key == ParamChmod {
			parsed, err := strconv.ParseInt(val[0], 8, 32)
			if err != nil {
				return nil, fmt.Errorf("could not parse value for 'chmod' param: %v", val)
			}

			mode = os.FileMode(parsed)
			if err != nil {
				return nil, fmt.Errorf("invalid file mode supplied: %v", val[0])
			}
		}
	}

	return newFilesystemStorage(path, username, pass, mode)
}

func newFilesystemStorage(path, owner, group string, mode os.FileMode) (*FilesystemStorage, error) {
	if len(path) == 0 {
		return nil, errors.New("empty path provided")
	}

	var uid, gid *int
	if len(owner) > 0 && len(group) > 0 {
		localUser, err := user.Lookup(owner)
		if err != nil {
			return nil, fmt.Errorf("could not lookup user '%s': %v", owner, err)
		}
		cuid, err := strconv.Atoi(localUser.Uid)
		if err != nil {
			return nil, fmt.Errorf("was expecting a numerical uid, got '%s'", localUser.Uid)
		}
		uid = &cuid

		localGroup, err := user.LookupGroup(group)
		if err != nil {
			return nil, fmt.Errorf("could not lookup group '%s': %v", group, err)
		}
		cgid, err := strconv.Atoi(localGroup.Gid)
		if err != nil {
			return nil, fmt.Errorf("was expecting a numerical gid, got '%s'", localGroup.Gid)
		}
		gid = &cgid
	}

	return &FilesystemStorage{
		FilePath:  path,
		FileOwner: uid,
		FileGroup: gid,
		Mode:      mode,
	}, nil
}

func (fs *FilesystemStorage) Read() ([]byte, error) {
	data, err := os.ReadFile(fs.FilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, pki.ErrNoCertFound
		}
		return nil, err
	}

	return data, nil
}

func (fs *FilesystemStorage) CanRead() error {
	_, err := os.Stat(fs.FilePath)
	return err
}

func (fs *FilesystemStorage) Write(signedData []byte) error {
	err := os.WriteFile(fs.FilePath, signedData, fs.Mode)
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
