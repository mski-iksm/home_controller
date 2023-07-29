package temp_controller

import (
	"errors"
	"fmt"
	"home_controller/appliance"
	"home_controller/device"
	"home_controller/signal"
	"log"
	"os"
	"strconv"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

var TOO_HOT_THRESHOLD float64 = 27.5
var TOO_COLD_THRESHOLD float64 = 24.0

// 設定可能な再下限・最上限
var MINIMUM_TEMPERATURE_SETTING float64 = 23.0
var MAXIMUM_TEMPERATURE_SETTING float64 = 30.0

func getCurrentAirconSettings(appliance appliance.Appliance) CurrentAirConSettings {
	tempreture, err := strconv.ParseFloat(appliance.Settings.Temp, 64)
	if err != nil {
		fmt.Errorf("温度 %v をfloat64にcastできません", appliance.Settings.Temp)
	}

	current_aircon_settings := CurrentAirConSettings{
		AirconSettings: signal.AirconSettings{
			Temperature:   tempreture,
			OperationMode: appliance.Settings.Mode,
			AirVolume:     appliance.Settings.Vol,
			AirDirection:  appliance.Settings.Dir,
		},
		UpdatedAt: appliance.Settings.UpdatedAt,
		PowerOn:   len(appliance.Settings.Button) == 0,
	}

	appLog.Printf("設定温度: %v\n", current_aircon_settings.AirconSettings.Temperature)
	appLog.Printf("設定時刻: %v\n", current_aircon_settings.UpdatedAt)
	appLog.Printf("電源: %v\n", current_aircon_settings.PowerOn)

	return current_aircon_settings
}

func get_current_temperature(device device.Device) CurrentTempreture {
	current_tempreture := CurrentTempreture{
		Tempreture: device.NewestEvents.Te.Val,
		UpdatedAt:  device.NewestEvents.Te.CreatedAt,
	}

	appLog.Printf("現在温度: %v\n", current_tempreture.Tempreture)
	appLog.Printf("温度更新時刻: %v\n", current_tempreture.UpdatedAt)

	return current_tempreture
}

// 新settingを作る
func buildNewAirconSettings(current_aircon_setting CurrentAirConSettings, current_tempreture CurrentTempreture) (NewAirConSettings, error) {
	var no_error error = nil

	// 電源がオフなら No Change
	if !current_aircon_setting.PowerOn {
		return NewAirConSettings{
			AirconSettings: current_aircon_setting.AirconSettings,
			PowerOn:        false,
		}, errors.New("No Change")
	}

	// too hot なら温度下げる
	if current_tempreture.Tempreture >= TOO_HOT_THRESHOLD && current_aircon_setting.AirconSettings.Temperature > MINIMUM_TEMPERATURE_SETTING {
		new_aircon_setting := NewAirConSettings{
			AirconSettings: signal.AirconSettings{
				OperationMode: current_aircon_setting.AirconSettings.OperationMode,
				Temperature:   current_aircon_setting.AirconSettings.Temperature - 1.0,
				AirVolume:     current_aircon_setting.AirconSettings.AirVolume,
				AirDirection:  current_aircon_setting.AirconSettings.AirDirection,
			},
			PowerOn: true,
		}
		return new_aircon_setting, no_error
	}

	// too cold なら温度上げる
	if current_tempreture.Tempreture <= TOO_COLD_THRESHOLD && current_aircon_setting.AirconSettings.Temperature < MAXIMUM_TEMPERATURE_SETTING {
		new_aircon_setting := NewAirConSettings{
			AirconSettings: signal.AirconSettings{
				OperationMode: current_aircon_setting.AirconSettings.OperationMode,
				Temperature:   current_aircon_setting.AirconSettings.Temperature + 1.0,
				AirVolume:     current_aircon_setting.AirconSettings.AirVolume,
				AirDirection:  current_aircon_setting.AirconSettings.AirDirection,
			},
			PowerOn: true,
		}
		return new_aircon_setting, no_error
	}

	// 気温がOK範囲の場合、setting is too high をチェック -> 一旦不要

	// 特に問題なしなら 現状のまま を返す
	return NewAirConSettings{
		AirconSettings: current_aircon_setting.AirconSettings,
		PowerOn:        true,
	}, errors.New("No Change")
}

func find_aircon_appliance(appliances []appliance.Appliance) (appliance.Appliance, error) {
	var no_error error = nil

	for _, appliance := range appliances {
		if appliance.Aircon != nil {
			return appliance, no_error
		}
	}
	return appliance.Appliance{}, errors.New("No AC found")
}

func BuildNewAirconAppliance(appliances []appliance.Appliance, device device.Device) (AirconAppliance, error) {
	// エアコンを探して appliance を返す
	aircon_appliance, ac_not_found_err := find_aircon_appliance(appliances)
	if ac_not_found_err != nil {
		errLog.Printf("エアコンが見つかりませんでした")
		no_appliance_error := errors.New("No Change")
		return AirconAppliance{}, no_appliance_error
	}

	// 今の設定を取得
	current_aircon_setting := getCurrentAirconSettings(aircon_appliance)

	// deviceから今の気温を取得
	current_tempreture := get_current_temperature(device)

	// new settingを作る
	new_aircon_settings, no_settings_change_error := buildNewAirconSettings(current_aircon_setting, current_tempreture)
	if no_settings_change_error != nil {
		errLog.Printf("設定変更なし")
		return AirconAppliance{}, no_settings_change_error
	}

	// setting から appliance を作る
	new_aircon_appliance := AirconAppliance{
		ApplianceId:    aircon_appliance.ID,
		AirconSettings: new_aircon_settings.AirconSettings,
	}
	appLog.Printf("new_aircon_appliance: %v", new_aircon_appliance)

	return new_aircon_appliance, nil
}
