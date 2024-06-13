package runner

import (
	"fmt"

	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

func Notify(nature_api_secret string, device_name string, temptureMaxMinSettings temp_controller.TempretureMaxMinSettings, slackObject slack.SlackObject) {
	// 時刻を出力
	currentTime := getTime()
	appLog.Printf("時刻: %v\n", currentTime)

	// get devices
	var devices []device.Device = device.Get_devices(nature_api_secret)

	// select device
	selected_device, no_device_err := device.SelectDevice(devices, device_name)
	if no_device_err != nil {
		errLog.Println(no_device_err)
		return
	}

	// get current temperature
	current_temp := temp_controller.GetCurrentTemp(selected_device)

	// check if the current temperature is too hot or too cold
	if current_temp.Tempreture >= temptureMaxMinSettings.TooHotThreshold {
		slackMessage := fmt.Sprintf("暑いです。現在の気温 %v\n", current_temp.Tempreture)
		slackObject.SendSlack(slackMessage)
		return
	}
	if current_temp.Tempreture <= temptureMaxMinSettings.TooColdThreshold {
		slackMessage := fmt.Sprintf("寒いです。現在の気温 %v\n", current_temp.Tempreture)
		slackObject.SendSlack(slackMessage)
		return
	}
}
