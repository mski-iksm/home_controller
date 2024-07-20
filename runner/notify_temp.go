package runner

import (
	"fmt"

	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

func NotifyTemp(nature_api_secret string, device_name string, temptureMaxMinSettings temp_controller.TempretureMaxMinSettings, slackObject slack.SlackObject) {

	// get devices
	var devices []device.Device = device.Get_devices(nature_api_secret)

	// select device
	selected_device, no_device_err := device.SelectDevice(devices, device_name)
	if no_device_err != nil {
		errLog.Println(no_device_err)
		return
	}

	current_tempreture := temp_controller.Get_current_temperature(selected_device)

	// slackを送る
	if current_tempreture.Tempreture >= temptureMaxMinSettings.TooHotThreshold {
		slackMessage := fmt.Sprintf("気温が暑いです。現在の気温: %v\n", current_tempreture.Tempreture)
		slackObject.SendSlack(slackMessage)
	}
	if current_tempreture.Tempreture < temptureMaxMinSettings.TooColdThreshold {
		slackMessage := fmt.Sprintf("気温が寒いです。現在の気温: %v\n", current_tempreture.Tempreture)
		slackObject.SendSlack(slackMessage)
	}
}
