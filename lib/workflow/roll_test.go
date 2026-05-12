package workflow

import (
	"testing"
)

func TestPrepareGameData(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Valid s4 preset",
			data:    "s4",
			wantErr: false,
		},
		{
			name:    "Default values",
			data:    "",
			wantErr: false,
		},
		{
			name:    "Invalid preset",
			data:    "invalid-preset",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := PrepareGameData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareGameData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("PrepareGameData() returned empty data for valid request")
			}
		})
	}
}
