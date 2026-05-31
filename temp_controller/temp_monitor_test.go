package temp_controller

import (
	"strings"
	"testing"
	"time"
)

func TestDecideTemperatureAlert(t *testing.T) {
	settings := TempretureMaxMinSettings{
		TooHotThreshold:  27.5,
		TooColdThreshold: 24.0,
	}

	tests := []struct {
		name        string
		tempreture  float64
		wantNotify  bool
		wantReason  string
		wantMessage string
	}{
		{
			name:        "tooHot",
			tempreture:  28.0,
			wantNotify:  true,
			wantReason:  "temperature is too hot",
			wantMessage: "設定値を超えています",
		},
		{
			name:        "tooCold",
			tempreture:  23.5,
			wantNotify:  true,
			wantReason:  "temperature is too cold",
			wantMessage: "設定値を下回っています",
		},
		{
			name:       "withinThresholds",
			tempreture: 25.0,
			wantNotify: false,
			wantReason: "temperature is within alert thresholds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecideTemperatureAlert(CurrentTempreture{
				Tempreture: tt.tempreture,
				UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			}, settings)

			if got.ShouldNotify != tt.wantNotify {
				t.Errorf("ShouldNotify mismatch. Must be %v, got %v\n", tt.wantNotify, got.ShouldNotify)
			}
			if got.Reason != tt.wantReason {
				t.Errorf("Reason mismatch. Must be %v, got %v\n", tt.wantReason, got.Reason)
			}
			if tt.wantMessage != "" && !strings.Contains(got.Message, tt.wantMessage) {
				t.Errorf("Message must contain %v, got %v\n", tt.wantMessage, got.Message)
			}
		})
	}
}
