package internal

import "testing"

func Test_dumpMetrics(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			want: `# HELP vault_pki_cli_success_bool Boolean that reflects whether the tool ran successful
# TYPE vault_pki_cli_success_bool gauge
vault_pki_cli_success_bool{cn="domain"} 1
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		MetricSuccess.WithLabelValues("domain").Set(1)
		t.Run(tt.name, func(t *testing.T) {
			got, err := dumpMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("dumpMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dumpMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}
