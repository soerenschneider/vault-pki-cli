package storage

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func getLiteral(x int) *int {
	return &x
}

func getOsDependendGroup() string {
	if runtime.GOOS == "linux" {
		return "root"
	}
	return "wheel"
}

func TestNewFilesystemStorageFromUri(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    *FilesystemStorage
		wantErr bool
	}{
		{
			name: "Simple",
			uri:  "file:///home/soeren/.certs/cert.pem",
			want: &FilesystemStorage{
				FilePath: "/home/soeren/.certs/cert.pem",
			},
			wantErr: false,
		},
		{
			name: "With user and group",
			uri:  fmt.Sprintf("file://root:%s@/home/soeren/.certs/cert.pem", getOsDependendGroup()),
			want: &FilesystemStorage{
				FilePath:  "/home/soeren/.certs/cert.pem",
				FileOwner: getLiteral(0),
				FileGroup: getLiteral(0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFilesystemStorageFromUri(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilesystemStorageFromUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilesystemStorageFromUri() got = %v, want %v", got, tt.want)
			}
		})
	}
}
