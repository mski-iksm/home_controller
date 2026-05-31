package runner

import (
	"fmt"

	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

func NotifyTemp(nature_api_secret string, device_name string, temperatureNotifySettings temp_controller.TemperatureNotifySettings, slackObject slack.SlackObject) {

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

	// エアコンを探して appliance を返す
	aircon_appliance, ac_not_found_err := temp_controller.Find_aircon_appliance(filtered_appliances)
	if ac_not_found_err != nil {
		slackMessage := fmt.Sprintf("エアコンが見つかりませんでした")
		slackObject.SendSlack(slackMessage)
		errLog.Printf("エアコンが見つかりませんでした")
		return
	}

	// 今のエアコンの設定を取得
	current_aircon_setting := temp_controller.GetCurrentAirconSettings(aircon_appliance)

	// エアコンの電源が入っていない場合は何もしない
	current_tempreture := temp_controller.Get_current_temperature(selected_device)
	temperatureNotification := temp_controller.DecideTemperatureNotification(current_tempreture, current_aircon_setting.PowerOn, temperatureNotifySettings)
	if !temperatureNotification.ShouldNotify {
		return
	}

	slackObject.SendSlack(temperatureNotification.Message)
}
