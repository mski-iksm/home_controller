package runner

import (
	"github.com/mski-iksm/home_controller/service"
	"github.com/mski-iksm/home_controller/temp_controller"
)

func NotifyTemp(nature_api_secret string, device_name string, temperatureNotifySettings temp_controller.TemperatureNotifySettings, ntfyUrl string) {
	notifyTempService := service.NewNotifyTempService(nature_api_secret, ntfyUrl)
	_, err := notifyTempService.Run(service.NotifyTempRequest{
		DeviceName: device_name,
		Settings:   temperatureNotifySettings,
	})
	if err != nil {
		errLog.Println(err)
	}
}
