package temp_controller

import (
	"testing"
	"time"
)

func TestDecideTemperatureNotification(t *testing.T) {
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	current := CurrentTempreture{
		Tempreture: 28.0,
		UpdatedAt:  time.Date(2020, 1, 1, 5, 0, 0, 0, tokyo),
	}
	settings := TemperatureNotifySettings{
		TooHotThreshold:  27.5,
		TooColdThreshold: 24.0,
	}

	tests := []struct {
		name     string
		powerOn  bool
		temp     float64
		want     TemperatureNotification
		settings TemperatureNotifySettings
	}{
		{
			name:    "tooHot",
			powerOn: true,
			temp:    28.0,
			want: TemperatureNotification{
				ShouldNotify: true,
				Message:      "気温が暑いです。現在の気温: 28.0",
				Reason:       "temperature is too hot",
			},
			settings: settings,
		},
		{
			name:    "tooCold",
			powerOn: true,
			temp:    23.5,
			want: TemperatureNotification{
				ShouldNotify: true,
				Message:      "気温が寒いです。現在の気温: 23.5",
				Reason:       "temperature is too cold",
			},
			settings: settings,
		},
		{
			name:    "powerOff",
			powerOn: false,
			temp:    28.0,
			want: TemperatureNotification{
				ShouldNotify: false,
				Reason:       "power is off",
			},
			settings: settings,
		},
		{
			name:    "withinThresholds",
			powerOn: true,
			temp:    25.0,
			want: TemperatureNotification{
				ShouldNotify: false,
				Reason:       "temperature is within alert thresholds",
			},
			settings: settings,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecideTemperatureNotification(CurrentTempreture{
				Tempreture: tt.temp,
				UpdatedAt:  current.UpdatedAt,
			}, tt.powerOn, tt.settings)

			if got.ShouldNotify != tt.want.ShouldNotify {
				t.Fatalf("ShouldNotify mismatch. want %v, got %v", tt.want.ShouldNotify, got.ShouldNotify)
			}
			if got.Reason != tt.want.Reason {
				t.Fatalf("Reason mismatch. want %v, got %v", tt.want.Reason, got.Reason)
			}
			if got.Message != tt.want.Message {
				t.Fatalf("Message mismatch. want %q, got %q", tt.want.Message, got.Message)
			}
		})
	}
}
