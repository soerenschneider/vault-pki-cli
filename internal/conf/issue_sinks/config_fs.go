package issue_sinks

import (
	"errors"
	"fmt"
	"reflect"
)

const FsType = "file"

type FilesystemSink struct {
	CertificateFile string `json:"cert_file"`
	PrivateKeyFile  string `json:"private_key_file"`
	CaFile          string `json:"ca_file,omitempty"`
	FileOwner       string `json:"file_owner,omitempty"`
	FileGroup       string `json:"file_group,omitempty"`
}

func (sink *FilesystemSink) GetType() string {
	return FsType
}

func FsBackendFromMap(args SinkConfig) (*FilesystemSink, error) {
	expectedTypes := map[string]reflect.Kind{
		"cert_file":        reflect.String,
		"private_key_file": reflect.String,
		"ca_file":          reflect.String,
		"file_owner":       reflect.String,
		"file_group":       reflect.String,
	}

	fsBackend := &FilesystemSink{}

	for keyword, valueType := range expectedTypes {
		t, ok := args[keyword]
		if !ok {
			continue
		}

		v := reflect.ValueOf(t)
		if v.Kind() != valueType {
			return nil, fmt.Errorf("expected type %v for keyword '%s'", valueType, keyword)
		}

		switch keyword {
		case "cert_file":
			fsBackend.CertificateFile = args["cert_file"].(string)
		case "private_key_file":
			fsBackend.PrivateKeyFile = args["private_key_file"].(string)
		case "ca_file":
			fsBackend.CaFile = args["ca_file"].(string)
		case "file_owner":
			fsBackend.FileOwner = args["file_owner"].(string)
		case "file_group":
			fsBackend.FileGroup = args["file_group"].(string)
		}
	}

	return fsBackend, nil
}

func (sink *FilesystemSink) Validate() (errs []error) {
	ownerDefined := len(sink.FileOwner) > 0
	groupDefined := len(sink.FileGroup) > 0
	if !ownerDefined && groupDefined {
		errs = append(errs, errors.New("only 'file_group' defined but not 'file_owner'"))
	}
	if ownerDefined && !groupDefined {
		errs = append(errs, errors.New("only 'file_owner' defined but not 'file_group'"))
	}

	emptyPrivateKeyFile := len(sink.PrivateKeyFile) == 0
	if emptyPrivateKeyFile {
		errs = append(errs, errors.New("must provide private key file 'private_key_file'"))
	}

	if len(sink.CertificateFile) == 0 {
		errs = append(errs, errors.New("must provide private key file 'cert_file'"))
	}

	return
}

func (sink *FilesystemSink) PrintConfig() {
	fmt.Println(sink)
}
