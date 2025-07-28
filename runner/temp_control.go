package runner

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/signal"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

func getTime() time.Time {
	nowTime := time.Now().UTC() // 現在時刻をUTCで取得

	// タイムゾーンからJSTを読み込み
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		errLog.Println(err)
		os.Exit(1)
	}
	timeTokyo := nowTime.In(tokyo)
	return timeTokyo
}

func TempControl(nature_api_secret string, device_name string, temptureMaxMinSettings temp_controller.TempretureMaxMinSettings, slackObject slack.SlackObject) {
	// 時刻を出力
	currentTime := getTime()
	appLog.Printf("時刻: %v\n", currentTime)

	// get devices
	var devices []device.Device = device.Get_devices(nature_api_secret)
	// get appliances
	var appliances []appliance.Appliance = appliance.Build_appliances(nature_api_secret)

	// select device
	selected_device, no_device_err := device.SelectDevice(devices, device_name)
	if no_device_err != nil {
		errLog.Println(no_device_err)
		return
	}

	// filter appliances
	filtered_appliances := appliance.FilterAppliances(appliances, device_name)

	newAirconOrderParameters, no_change_err := temp_controller.BuildNewAirconOrderParameters(filtered_appliances, selected_device, temptureMaxMinSettings)

	if no_change_err != nil {
		errLog.Println(no_change_err)
		return
	}

	signal.PostAirconSignal(nature_api_secret, newAirconOrderParameters.ApplianceId, newAirconOrderParameters.AirconSettings)

	// deviceから今の気温を取得
	current_tempreture := temp_controller.Get_current_temperature(selected_device)

	// slackを送る
	slackMessage := fmt.Sprintf("エアコンの設定を変更しました。\n現在温度: %v\n設定温度: %v\nモード: %v\n風量: %v\n風向: %v\n", current_tempreture.Tempreture, newAirconOrderParameters.AirconSettings.Temperature, newAirconOrderParameters.AirconSettings.OperationMode, newAirconOrderParameters.AirconSettings.AirVolume, newAirconOrderParameters.AirconSettings.AirDirection)
	slackObject.SendSlack(slackMessage)
}
