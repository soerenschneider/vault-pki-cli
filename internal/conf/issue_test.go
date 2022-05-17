package conf

import (
	"reflect"
	"testing"
)

func TestIssueArguments_Validate(t *testing.T) {
	tests := []struct {
		name string
		conf IssueArguments
		want []error
	}{
		{
			name: "unitialized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.conf.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
