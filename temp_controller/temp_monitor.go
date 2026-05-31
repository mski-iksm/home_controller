package temp_controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mski-iksm/home_controller/device"
)

func DecideTemperatureAlert(current_tempreture CurrentTempreture, temptureMaxMinSettings TempretureMaxMinSettings) TemperatureAlert {
	if current_tempreture.Tempreture >= temptureMaxMinSettings.TooHotThreshold+0.5 {
		return TemperatureAlert{
			ShouldNotify: true,
			Message:      "緊急アラート: 現在の温度が設定値を超えています。\n現在温度: " + fmt.Sprintf("%.1f", current_tempreture.Tempreture) + "\n設定温度: " + fmt.Sprintf("%.1f", temptureMaxMinSettings.TooHotThreshold),
			Priority:     "5",
			Reason:       "temperature is too hot",
		}
	}

	if current_tempreture.Tempreture <= temptureMaxMinSettings.TooColdThreshold-0.5 {
		return TemperatureAlert{
			ShouldNotify: true,
			Message:      "緊急アラート: 現在の温度が設定値を下回っています。\n現在温度: " + fmt.Sprintf("%.1f", current_tempreture.Tempreture) + "\n設定温度: " + fmt.Sprintf("%.1f", temptureMaxMinSettings.TooColdThreshold),
			Priority:     "5",
			Reason:       "temperature is too cold",
		}
	}

	return TemperatureAlert{
		ShouldNotify: false,
		Reason:       "temperature is within alert thresholds",
	}
}

func SendTemperatureAlert(ntfyUrl string, alert TemperatureAlert) {
	if !alert.ShouldNotify {
		return
	}

	req, _ := http.NewRequest("POST", ntfyUrl, strings.NewReader(alert.Message))
	req.Header.Set("Priority", alert.Priority)
	http.DefaultClient.Do(req)
}

func DecideTemperatureNotification(current_tempreture CurrentTempreture, powerOn bool, temperatureNotifySettings TemperatureNotifySettings) TemperatureNotification {
	if !powerOn {
		return TemperatureNotification{
			ShouldNotify: false,
			Reason:       "power is off",
		}
	}

	if current_tempreture.Tempreture >= temperatureNotifySettings.TooHotThreshold {
		return TemperatureNotification{
			ShouldNotify: true,
			Message:      "気温が暑いです。現在の気温: " + fmt.Sprintf("%.1f", current_tempreture.Tempreture),
			Reason:       "temperature is too hot",
		}
	}

	if current_tempreture.Tempreture < temperatureNotifySettings.TooColdThreshold {
		return TemperatureNotification{
			ShouldNotify: true,
			Message:      "気温が寒いです。現在の気温: " + fmt.Sprintf("%.1f", current_tempreture.Tempreture),
			Reason:       "temperature is too cold",
		}
	}

	return TemperatureNotification{
		ShouldNotify: false,
		Reason:       "temperature is within alert thresholds",
	}
}

func MonitorTempreture(device device.Device, temptureMaxMinSettings TempretureMaxMinSettings, ntfyUrl string) {
	// deviceから今の気温を取得
	current_tempreture := Get_current_temperature(device)

	alert := DecideTemperatureAlert(current_tempreture, temptureMaxMinSettings)
	if alert.Reason == "temperature is too hot" {
		println("Sending hot alert to ntfy")
	}
	if alert.Reason == "temperature is too cold" {
		println("Sending cold alert to ntfy")
	}
	SendTemperatureAlert(ntfyUrl, alert)
}
