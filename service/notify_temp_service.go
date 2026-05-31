package service

import (
	"github.com/mski-iksm/home_controller/temp_controller"
)

type NotifyTempRequest struct {
	DeviceName string
	Settings   temp_controller.TemperatureNotifySettings
}

type NotifyTempResult struct {
	CurrentTemperature temp_controller.CurrentTempreture
	TemperatureAlert   temp_controller.TemperatureAlert
}

type NotifyTempService struct {
	NatureClient             NatureClient
	TemperatureAlertNotifier TemperatureAlertNotifier
}

func NewNotifyTempService(natureAPISecret string, ntfyURL string) NotifyTempService {
	return NotifyTempService{
		NatureClient: NatureRemoClient{
			APISecret: natureAPISecret,
		},
		TemperatureAlertNotifier: NtfyNotifier{URL: ntfyURL},
	}
}

func (s NotifyTempService) Run(request NotifyTempRequest) (NotifyTempResult, error) {
	var result NotifyTempResult

	airconContext, err := LoadAirconContext(s.NatureClient, request.DeviceName)
	if err != nil {
		return result, err
	}

	result.CurrentTemperature = airconContext.CurrentTemperature

	if !airconContext.CurrentAirconSettings.PowerOn {
		return result, nil
	}

	temperatureAlert := temp_controller.DecideTemperatureAlert(
		airconContext.CurrentTemperature,
		temp_controller.TempretureMaxMinSettings{
			TooHotThreshold:  request.Settings.TooHotThreshold,
			TooColdThreshold: request.Settings.TooColdThreshold,
		},
	)
	result.TemperatureAlert = temperatureAlert

	s.TemperatureAlertNotifier.SendTemperatureAlert(temperatureAlert)
	return result, nil
}
