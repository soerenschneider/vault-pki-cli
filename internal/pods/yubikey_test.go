package pods

import (
	"github.com/go-piv/piv-go/piv"
	"reflect"
	"testing"
)

func Test_getSlot(t *testing.T) {
	tests := []struct {
		name    string
		slot    uint32
		want    *piv.Slot
		wantErr bool
	}{
		{
			name:    "invalid slot",
			slot:    0,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid slot",
			slot:    0x9a,
			want:    &piv.SlotAuthentication,
			wantErr: false,
		},
		{
			name:    "different valid slot",
			slot:    0x9c,
			want:    &piv.SlotSignature,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TranslateSlot(tt.slot)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranslateSlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TranslateSlot() got = %v, want %v", got, tt.want)
			}
		})
	}
}
