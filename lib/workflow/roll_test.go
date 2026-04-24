package workflow

import (
	"testing"

	"maprandoseedroller/lib/models"
)

func TestPrepareGameData(t *testing.T) {
	tests := []struct {
		name    string
		req     models.RequestIn
		wantErr bool
	}{
		{
			name: "Valid s4 preset",
			req: models.RequestIn{
				Preset: "s4",
				Flags:  "",
			},
			wantErr: false,
		},
		{
			name: "Invalid preset",
			req: models.RequestIn{
				Preset: "invalid-preset",
				Flags:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := PrepareGameData(tt.req)
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
