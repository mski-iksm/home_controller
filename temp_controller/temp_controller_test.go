package temp_controller

import (
	"errors"
	"home_controller/signal"
	"testing"
	"time"
)

func TestBuildNewAirconSettings(t *testing.T) {
	type Args struct {
		current_aircon_setting CurrentAirConSettings
		current_tempreture     CurrentTempreture
	}

	tests := []struct {
		name    string
		args    Args
		want    NewAirConSettings
		wantErr error
	}{
		{
			name: "tooHot",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   27.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Now(),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 30.0,
					UpdatedAt:  time.Now(),
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   26.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: true,
			},
			wantErr: nil,
		},
		{
			name: "tooCold",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   27.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Now(),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 23.0,
					UpdatedAt:  time.Now(),
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   28.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: true,
			},
			wantErr: nil,
		},
		{
			name: "tooColdButNoChange",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   30.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Now(),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 23.0,
					UpdatedAt:  time.Now(),
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   30.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: true,
			},
			wantErr: errors.New("No Change"),
		},
		{
			name: "tooHotButNoChange",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   23.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Now(),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 28.0,
					UpdatedAt:  time.Now(),
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   23.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: true,
			},
			wantErr: errors.New("No Change"),
		},
		{
			name: "dontChangeIfPowerOff",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   27.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Now(),
					PowerOn:   false,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 29.0,
					UpdatedAt:  time.Now(),
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   27.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: false,
			},
			wantErr: errors.New("No Change"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildNewAirconSettings(tt.args.current_aircon_setting, tt.args.current_tempreture)
			if got != tt.want {
				t.Errorf("buildNewAirconSettings mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if (err == nil) != (tt.wantErr == nil) {
				t.Errorf("buildNewAirconSettings error mismatch. Must be %v, got %v\n", tt.wantErr, err)
			}
		})
	}
}
