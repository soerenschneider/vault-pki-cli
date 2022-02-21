package pods

import (
	"reflect"
	"testing"
)

func TestFsPod_Read(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     []byte
		wantErr  bool
	}{
		{
			name:     "read hello file",
			filePath: "../../assets/tests/hello.txt",
			want:     []byte("hello"),
			wantErr:  false,
		},
		{
			name:     "read unexistent file",
			filePath: "../../assets/tests/hello-im-not-here.txt",
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FsPod{
				FilePath: tt.filePath,
			}
			got, err := fs.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("FsPod.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FsPod.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFsPod_CanRead(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "read hello file",
			filePath: "../../assets/tests/hello.txt",
			wantErr:  false,
		},
		{
			name:     "read unexistent file",
			filePath: "../../assets/tests/hello-im-not-here.txt",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FsPod{
				FilePath: tt.filePath,
			}
			if err := fs.CanRead(); (err != nil) != tt.wantErr {
				t.Errorf("FsPod.CanRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
