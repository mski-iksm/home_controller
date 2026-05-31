package service

import (
	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
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

	devices := s.NatureClient.GetDevices()
	appliances := s.NatureClient.GetAppliances()

	selectedDevice, err := device.SelectDevice(devices, request.DeviceName)
	if err != nil {
		return result, err
	}

	filteredAppliances := appliance.FilterAppliances(appliances, request.DeviceName)

	airconAppliance, err := temp_controller.Find_aircon_appliance(filteredAppliances)
	if err != nil {
		return result, err
	}

	currentAirconSetting := temp_controller.GetCurrentAirconSettings(airconAppliance)
	currentTemperature := temp_controller.Get_current_temperature(selectedDevice)
	result.CurrentTemperature = currentTemperature

	if !currentAirconSetting.PowerOn {
		return result, nil
	}

	temperatureAlert := temp_controller.DecideTemperatureAlert(
		currentTemperature,
		temp_controller.TempretureMaxMinSettings{
			TooHotThreshold:  request.Settings.TooHotThreshold,
			TooColdThreshold: request.Settings.TooColdThreshold,
		},
	)
	result.TemperatureAlert = temperatureAlert

	s.TemperatureAlertNotifier.SendTemperatureAlert(temperatureAlert)
	return result, nil
}
