package runner

import (
	"home_controller/appliance"
	"home_controller/device"
	"home_controller/signal"
	"home_controller/temp_controller"
	"log"
	"os"
	"time"
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

func TempControl(nature_api_secret string, device_name string, temptureMaxMinSettings temp_controller.TempretureMaxMinSettings) {
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

	// slackを送る
}
