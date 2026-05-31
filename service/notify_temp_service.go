package service

import (
	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

type NotifyTempRequest struct {
	DeviceName string
	Settings   temp_controller.TemperatureNotifySettings
}

type NotifyTempResult struct {
	CurrentTemperature      temp_controller.CurrentTempreture
	TemperatureNotification temp_controller.TemperatureNotification
}

type NotifyTempService struct {
	NatureClient  NatureClient
	SlackNotifier SlackNotifier
}

func NewNotifyTempService(natureAPISecret string, slackObject slack.SlackObject) NotifyTempService {
	return NotifyTempService{
		NatureClient: NatureRemoClient{
			APISecret: natureAPISecret,
		},
		SlackNotifier: &slackObject,
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
		s.SlackNotifier.SendSlack("エアコンが見つかりませんでした")
		return result, err
	}

	currentAirconSetting := temp_controller.GetCurrentAirconSettings(airconAppliance)
	currentTemperature := temp_controller.Get_current_temperature(selectedDevice)
	result.CurrentTemperature = currentTemperature

	temperatureNotification := temp_controller.DecideTemperatureNotification(
		currentTemperature,
		currentAirconSetting.PowerOn,
		request.Settings,
	)
	result.TemperatureNotification = temperatureNotification

	if !temperatureNotification.ShouldNotify {
		return result, nil
	}

	s.SlackNotifier.SendSlack(temperatureNotification.Message)
	return result, nil
}
