package temp_controller

import (
	"testing"
)

func TestConstructTempretureMaxMinSettings(t *testing.T) {
	type Args struct {
		tooHotThreshold           float64
		tooColdThreshold          float64
		preparationThreshold      float64
		minimumTemperatureSetting float64
		maximumTemperatureSetting float64
	}

	tests := []struct {
		name string
		args Args
		want *TempretureMaxMinSettings
	}{
		{
			name: "passBuild",
			args: Args{
				tooHotThreshold:           29.0,
				tooColdThreshold:          22.0,
				preparationThreshold:      0.0,
				minimumTemperatureSetting: 22.0,
				maximumTemperatureSetting: 30.0,
			},
			want: &TempretureMaxMinSettings{
				TooHotThreshold:           29.0,
				TooColdThreshold:          22.0,
				PreparationThreshold:      0.0,
				MinimumTemperatureSetting: 22.0,
				MaximumTemperatureSetting: 30.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstructTempretureMaxMinSettings(tt.args.tooHotThreshold, tt.args.tooColdThreshold, tt.args.preparationThreshold, tt.args.minimumTemperatureSetting, tt.args.maximumTemperatureSetting)
			if (*got).MaximumTemperatureSetting != tt.want.MaximumTemperatureSetting {
				t.Errorf("MaximumTemperatureSetting mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if (*got).MinimumTemperatureSetting != tt.want.MinimumTemperatureSetting {
				t.Errorf("MinimumTemperatureSetting mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if (*got).TooColdThreshold != tt.want.TooColdThreshold {
				t.Errorf("TooColdThreshold mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if (*got).TooHotThreshold != tt.want.TooHotThreshold {
				t.Errorf("TooHotThreshold mismatch. Must be %v, got %v\n", tt.want, got)
			}
		})
	}
}

func TestAssertTresholdSettings(t *testing.T) {
	tests := []struct {
		name  string
		args  TempretureMaxMinSettings
		isErr bool
	}{
		{
			name: "passAssertion",
			args: TempretureMaxMinSettings{
				TooHotThreshold:           29.0,
				TooColdThreshold:          22.0,
				MinimumTemperatureSetting: 22.0,
				MaximumTemperatureSetting: 30.0,
			},
			isErr: false,
		},
		{
			name: "passAssertion",
			args: TempretureMaxMinSettings{
				TooHotThreshold:           45.0,
				TooColdThreshold:          22.0,
				MinimumTemperatureSetting: 22.0,
				MaximumTemperatureSetting: 30.0,
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := assertTresholdSettings(tt.args)
			if tt.isErr != (gotErr != nil) {
				t.Errorf("assertTresholdSettings mismatch. Must be %v, got %v\n", tt.isErr, gotErr)
			}
		})
	}
}
