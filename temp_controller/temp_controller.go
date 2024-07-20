package temp_controller

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/signal"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

func ConvertUTCToJST(utcTime time.Time) time.Time {
	// タイムゾーンからJSTを読み込み
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		errLog.Println(err)
		os.Exit(1)
	}
	timeTokyo := utcTime.In(tokyo)
	return timeTokyo
}

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
		UpdatedAt: ConvertUTCToJST(appliance.Settings.UpdatedAt),
		PowerOn:   len(appliance.Settings.Button) == 0 || appliance.Settings.Button == "power-on",
		// 空 もしくは power-on なら電源オン
	}

	appLog.Printf("設定温度: %v\n", current_aircon_settings.AirconSettings.Temperature)
	appLog.Printf("設定時刻: %v\n", current_aircon_settings.UpdatedAt)
	appLog.Printf("電源: %v\n", current_aircon_settings.PowerOn)

	return current_aircon_settings
}

func Get_current_temperature(device device.Device) CurrentTempreture {
	current_tempreture := CurrentTempreture{
		Tempreture: device.NewestEvents.Te.Val,
		UpdatedAt:  ConvertUTCToJST(device.NewestEvents.Te.CreatedAt),
	}

	appLog.Printf("現在温度: %v\n", current_tempreture.Tempreture)
	appLog.Printf("温度更新時刻: %v\n", current_tempreture.UpdatedAt)

	return current_tempreture
}

// 新settingを作る
func buildNewAirconSettings(current_aircon_setting CurrentAirConSettings, current_tempreture CurrentTempreture, temptureMaxMinSettings TempretureMaxMinSettings) (NewAirConSettings, error) {
	var no_error error = nil
	currentTimeJST := ConvertUTCToJST(time.Now())

	// 電源がオフなら No Change
	if !current_aircon_setting.PowerOn {
		return NewAirConSettings{
			AirconSettings: current_aircon_setting.AirconSettings,
			PowerOn:        false,
		}, errors.New("No Change")
	}

	// 気温を測定した時刻が前回の変更から10分以内なら No Change
	if current_tempreture.UpdatedAt.Sub(current_aircon_setting.UpdatedAt).Minutes() < 10.0 {
		return NewAirConSettings{
			AirconSettings: current_aircon_setting.AirconSettings,
			PowerOn:        true,
		}, errors.New("No Change")
	}

	// 前回変更から20分以内なら No Change
	if currentTimeJST.Sub(current_aircon_setting.UpdatedAt).Minutes() < 20.0 {
		return NewAirConSettings{
			AirconSettings: current_aircon_setting.AirconSettings,
			PowerOn:        true,
		}, errors.New("No Change")
	}

	// too hot なら温度下げる
	if current_tempreture.Tempreture >= temptureMaxMinSettings.TooHotThreshold && current_aircon_setting.AirconSettings.Temperature > temptureMaxMinSettings.MinimumTemperatureSetting {
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
	if current_tempreture.Tempreture <= temptureMaxMinSettings.TooColdThreshold && current_aircon_setting.AirconSettings.Temperature < temptureMaxMinSettings.MaximumTemperatureSetting {
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

	// too hot に近づき、かつ1時間以上設定変更されていなければ温度下げる
	if current_tempreture.Tempreture >= temptureMaxMinSettings.TooHotThreshold-temptureMaxMinSettings.PreparationThreshold && current_aircon_setting.AirconSettings.Temperature > temptureMaxMinSettings.MinimumTemperatureSetting {
		time_diff := currentTimeJST.Sub(current_aircon_setting.UpdatedAt)
		if time_diff.Minutes() > 58.0 {
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
	}

	// too cold に近づき、かつ1時間以上設定変更されていなければ温度上げる
	if current_tempreture.Tempreture <= temptureMaxMinSettings.TooColdThreshold+temptureMaxMinSettings.PreparationThreshold && current_aircon_setting.AirconSettings.Temperature < temptureMaxMinSettings.MaximumTemperatureSetting {
		time_diff := currentTimeJST.Sub(current_aircon_setting.UpdatedAt)
		if time_diff.Minutes() > 58.0 {
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
	}

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

func BuildNewAirconOrderParameters(appliances []appliance.Appliance, device device.Device, temptureMaxMinSettings TempretureMaxMinSettings) (AirconOrderParameters, error) {
	// エアコンを探して appliance を返す
	aircon_appliance, ac_not_found_err := find_aircon_appliance(appliances)
	if ac_not_found_err != nil {
		errLog.Printf("エアコンが見つかりませんでした")
		no_appliance_error := errors.New("No Change")
		return AirconOrderParameters{}, no_appliance_error
	}

	// 今の設定を取得
	current_aircon_setting := getCurrentAirconSettings(aircon_appliance)

	// deviceから今の気温を取得
	current_tempreture := Get_current_temperature(device)

	// new settingを作る
	new_aircon_settings, no_settings_change_error := buildNewAirconSettings(current_aircon_setting, current_tempreture, temptureMaxMinSettings)
	if no_settings_change_error != nil {
		errLog.Printf("設定変更なし")
		return AirconOrderParameters{}, no_settings_change_error
	}

	// setting から order_parameter を作る
	newAirconOrderParameters := AirconOrderParameters{
		ApplianceId:    aircon_appliance.ID,
		AirconSettings: new_aircon_settings.AirconSettings,
	}
	appLog.Printf("newAirconOrderParameters: %v", newAirconOrderParameters)

	return newAirconOrderParameters, nil
}
