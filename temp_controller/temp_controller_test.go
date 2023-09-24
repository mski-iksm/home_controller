package temp_controller

import (
	"errors"
	"testing"
	"time"

	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/signal"
)

func TestBuildNewAirconSettings(t *testing.T) {
	type Args struct {
		current_aircon_setting CurrentAirConSettings
		current_tempreture     CurrentTempreture
	}

	tokyoTZ, _ := time.LoadLocation("Asia/Tokyo")
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
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, tokyoTZ),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 30.0,
					UpdatedAt:  time.Date(2020, 1, 1, 5, 0, 0, 0, tokyoTZ),
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
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, tokyoTZ),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 23.0,
					UpdatedAt:  time.Date(2020, 1, 1, 5, 0, 0, 0, tokyoTZ),
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
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, tokyoTZ),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 23.0,
					UpdatedAt:  time.Date(2020, 1, 1, 5, 0, 0, 0, tokyoTZ),
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
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, tokyoTZ),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 28.0,
					UpdatedAt:  time.Date(2020, 1, 1, 5, 0, 0, 0, tokyoTZ),
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
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, tokyoTZ),
					PowerOn:   false,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 29.0,
					UpdatedAt:  time.Date(2020, 1, 1, 5, 0, 0, 0, tokyoTZ),
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
		{
			name: "closeToTooCold",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   27.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 24.2,
					UpdatedAt:  time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
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
			name: "closeToTooColdButNoChange",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   27.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Now().Add(-3 * time.Minute),
					// 現在時刻から53分前だと変更されない
					PowerOn: true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 24.2,
					UpdatedAt:  time.Now().Add(40 * time.Hour),
					// 測定時刻は前回の設定変更よりも十分先にセット
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   27.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: true,
			},
			wantErr: errors.New("No Change"),
		},
		{
			name: "noChangeBecauseTempretureMeasurementIsLessThen10MinutesFromLastChange",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   27.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, tokyoTZ),
					PowerOn:   true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 23.0,
					UpdatedAt:  time.Date(2020, 1, 1, 0, 5, 0, 0, tokyoTZ),
					// 5分後
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   27.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: true,
			},
			wantErr: errors.New("No Change"),
		},
		{
			name: "noChangeBecauseIsLessThen20MinutesFromLastChange",
			args: Args{
				current_aircon_setting: CurrentAirConSettings{
					AirconSettings: signal.AirconSettings{
						OperationMode: "cool",
						Temperature:   27.0,
						AirVolume:     "auto",
						AirDirection:  "auto",
					},
					UpdatedAt: time.Now().Add(-15 * time.Minute),
					// 現在より15分前だと変更されない
					PowerOn: true,
				},
				current_tempreture: CurrentTempreture{
					Tempreture: 23.0,
					UpdatedAt:  time.Now().Add(40 * time.Hour),
					// 測定時刻は前回の設定変更よりも十分先にセット
				},
			},
			want: NewAirConSettings{
				AirconSettings: signal.AirconSettings{
					OperationMode: "cool",
					Temperature:   27.0,
					AirVolume:     "auto",
					AirDirection:  "auto",
				},
				PowerOn: true,
			},
			wantErr: errors.New("No Change"),
		},
	}

	temptureMaxMinSettings := TempretureMaxMinSettings{
		TooHotThreshold:           27.5,
		TooColdThreshold:          24.0,
		PreparationThreshold:      0.5,
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
	tokyoTZ, _ := time.LoadLocation("Asia/Tokyo")

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
							UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, tokyoTZ),
						},
						ID: "1",
					},
				},
				device: device.Device{
					NewestEvents: device.NewestEvents{
						Te: device.Temperature{
							Val:       29.0,
							CreatedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, tokyoTZ),
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

func TestConvertUTCToJST(t *testing.T) {
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	tests := []struct {
		name    string
		utcTime time.Time
		want    time.Time
	}{
		{
			name:    "confirmPass",
			utcTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			want:    time.Date(2020, 1, 1, 9, 0, 0, 0, tokyo),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertUTCToJST(tt.utcTime)
			if got.Year() != tt.want.Year() {
				t.Errorf("ConvertUTCToJST Year mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if got.Month() != tt.want.Month() {
				t.Errorf("ConvertUTCToJST Month mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if got.Day() != tt.want.Day() {
				t.Errorf("ConvertUTCToJST Day mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if got.Hour() != tt.want.Hour() {
				t.Errorf("ConvertUTCToJST Hour mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if got.Minute() != tt.want.Minute() {
				t.Errorf("ConvertUTCToJST Minute mismatch. Must be %v, got %v\n", tt.want, got)
			}
			if got.Second() != tt.want.Second() {
				t.Errorf("ConvertUTCToJST Second mismatch. Must be %v, got %v\n", tt.want, got)
			}
		})
	}
}
