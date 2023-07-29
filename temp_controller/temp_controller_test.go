package temp_controller

import (
	"errors"
	"home_controller/appliance"
	"home_controller/device"
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

	temptureMaxMinSettings := TempretureMaxMinSettings{
		TooHotThreshold:           27.5,
		TooColdThreshold:          24.0,
		MinimumTemperatureSetting: 23.0,
		MaximumTemperatureSetting: 30.0,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildNewAirconSettings(tt.args.current_aircon_setting, tt.args.current_tempreture, temptureMaxMinSettings)
			if got != tt.want {
				t.Errorf("buildNewAirconSettings mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if (err == nil) != (tt.wantErr == nil) {
				t.Errorf("buildNewAirconSettings error mismatch. Must be %v, got %v\n", tt.wantErr, err)
			}
		})
	}
}

func TestBuildNewAirconOrderParameters(t *testing.T) {
	type Args struct {
		appliances             []appliance.Appliance
		device                 device.Device
		temptureMaxMinSettings TempretureMaxMinSettings
	}

	tests := []struct {
		name    string
		args    Args
		want    AirconOrderParameters
		wantErr error
	}{
		{
			name: "confirmPass",
			args: Args{
				appliances: []appliance.Appliance{
					{
						Aircon: 1,
						Settings: appliance.Settings{
							Mode:      "cool",
							Temp:      "27.0",
							Vol:       "auto",
							Dir:       "auto",
							Button:    "",
							UpdatedAt: time.Now(),
						},
						ID: "1",
					},
				},
				device: device.Device{
					NewestEvents: device.NewestEvents{
						Te: device.Temperature{
							Val:       29.0,
							CreatedAt: time.Now(),
						},
					},
				},
				temptureMaxMinSettings: TempretureMaxMinSettings{
					TooHotThreshold:           27.5,
					TooColdThreshold:          24.0,
					MinimumTemperatureSetting: 24.0,
					MaximumTemperatureSetting: 30.0,
				},
			},
			want: AirconOrderParameters{
				ApplianceId: "1",
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   26.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildNewAirconOrderParameters(tt.args.appliances, tt.args.device, tt.args.temptureMaxMinSettings)
			if got != tt.want {
				t.Errorf("BuildNewAirconOrderParameters mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if (err == nil) != (tt.wantErr == nil) {
				t.Errorf("BuildNewAirconOrderParameters error mismatch. Must be %v, got %v\n", tt.wantErr, err)
			}
		})
	}
}
