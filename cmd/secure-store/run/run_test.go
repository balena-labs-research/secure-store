package run

import (
	"testing"
)

func TestExecuteArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Test execution of a command",
			args:    []string{"echo", "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecuteArgs(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
