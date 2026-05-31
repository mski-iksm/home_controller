package service

import (
	"fmt"

	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/signal"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

type TempControlRequest struct {
	NatureAPISecret string
	DeviceName      string
	Settings        temp_controller.TempretureMaxMinSettings
	SlackObject     slack.SlackObject
	NtfyURL         string
}

type TempControlResult struct {
	CurrentTemperature temp_controller.CurrentTempreture
	AirconSettings     signal.AirconSettings
	AirconChanged      bool
	TemperatureAlert   temp_controller.TemperatureAlert
}

type TempControlService struct{}

func NewTempControlService() TempControlService {
	return TempControlService{}
}

func (s TempControlService) Run(request TempControlRequest) (TempControlResult, error) {
	var result TempControlResult

	devices := device.Get_devices(request.NatureAPISecret)
	appliances := appliance.Build_appliances(request.NatureAPISecret)

	selectedDevice, err := device.SelectDevice(devices, request.DeviceName)
	if err != nil {
		return result, err
	}

	filteredAppliances := appliance.FilterAppliances(appliances, request.DeviceName)

	currentTemperature := temp_controller.Get_current_temperature(selectedDevice)
	result.CurrentTemperature = currentTemperature

	temperatureAlert := temp_controller.DecideTemperatureAlert(currentTemperature, request.Settings)
	result.TemperatureAlert = temperatureAlert
	temp_controller.SendTemperatureAlert(request.NtfyURL, temperatureAlert)

	newAirconOrderParameters, err := temp_controller.BuildNewAirconOrderParameters(filteredAppliances, selectedDevice, request.Settings)
	if err != nil {
		return result, err
	}

	signal.PostAirconSignal(request.NatureAPISecret, newAirconOrderParameters.ApplianceId, newAirconOrderParameters.AirconSettings)

	result.AirconSettings = newAirconOrderParameters.AirconSettings
	result.AirconChanged = true

	slackMessage := fmt.Sprintf(
		"エアコンの設定を変更しました。\n現在温度: %v\n設定温度: %v\nモード: %v\n風量: %v\n風向: %v\n",
		currentTemperature.Tempreture,
		newAirconOrderParameters.AirconSettings.Temperature,
		newAirconOrderParameters.AirconSettings.OperationMode,
		newAirconOrderParameters.AirconSettings.AirVolume,
		newAirconOrderParameters.AirconSettings.AirDirection,
	)
	slackObject := request.SlackObject
	slackObject.SendSlack(slackMessage)

	return result, nil
}
