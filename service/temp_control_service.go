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
	DeviceName string
	Settings   temp_controller.TempretureMaxMinSettings
}

type TempControlResult struct {
	CurrentTemperature temp_controller.CurrentTempreture
	AirconSettings     signal.AirconSettings
	AirconChanged      bool
	TemperatureAlert   temp_controller.TemperatureAlert
}

type NatureClient interface {
	GetDevices() []device.Device
	GetAppliances() []appliance.Appliance
	PostAirconSignal(applianceID string, settings signal.AirconSettings) int
}

type SlackNotifier interface {
	SendSlack(message string)
}

type TemperatureAlertNotifier interface {
	SendTemperatureAlert(alert temp_controller.TemperatureAlert)
}

type NatureRemoClient struct {
	APISecret string
}

func (c NatureRemoClient) GetDevices() []device.Device {
	return device.Get_devices(c.APISecret)
}

func (c NatureRemoClient) GetAppliances() []appliance.Appliance {
	return appliance.Build_appliances(c.APISecret)
}

func (c NatureRemoClient) PostAirconSignal(applianceID string, settings signal.AirconSettings) int {
	return signal.PostAirconSignal(c.APISecret, applianceID, settings)
}

type NtfyNotifier struct {
	URL string
}

func (n NtfyNotifier) SendTemperatureAlert(alert temp_controller.TemperatureAlert) {
	temp_controller.SendTemperatureAlert(n.URL, alert)
}

type TempControlService struct {
	NatureClient             NatureClient
	SlackNotifier            SlackNotifier
	TemperatureAlertNotifier TemperatureAlertNotifier
}

func NewTempControlService(natureAPISecret string, slackObject slack.SlackObject, ntfyURL string) TempControlService {
	return TempControlService{
		NatureClient: NatureRemoClient{
			APISecret: natureAPISecret,
		},
		SlackNotifier:            &slackObject,
		TemperatureAlertNotifier: NtfyNotifier{URL: ntfyURL},
	}
}

func (s TempControlService) Run(request TempControlRequest) (TempControlResult, error) {
	var result TempControlResult

	applianceContext, err := LoadApplianceContext(s.NatureClient, request.DeviceName)
	if err != nil {
		return result, err
	}

	result.CurrentTemperature = applianceContext.CurrentTemperature

	temperatureAlert := temp_controller.DecideTemperatureAlert(applianceContext.CurrentTemperature, request.Settings)
	result.TemperatureAlert = temperatureAlert
	s.TemperatureAlertNotifier.SendTemperatureAlert(temperatureAlert)

	newAirconOrderParameters, err := temp_controller.BuildNewAirconOrderParameters(applianceContext.FilteredAppliances, applianceContext.Device, request.Settings)
	if err != nil {
		return result, err
	}

	s.NatureClient.PostAirconSignal(newAirconOrderParameters.ApplianceId, newAirconOrderParameters.AirconSettings)

	result.AirconSettings = newAirconOrderParameters.AirconSettings
	result.AirconChanged = true

	slackMessage := fmt.Sprintf(
		"エアコンの設定を変更しました。\n現在温度: %v\n設定温度: %v\nモード: %v\n風量: %v\n風向: %v\n",
		applianceContext.CurrentTemperature.Tempreture,
		newAirconOrderParameters.AirconSettings.Temperature,
		newAirconOrderParameters.AirconSettings.OperationMode,
		newAirconOrderParameters.AirconSettings.AirVolume,
		newAirconOrderParameters.AirconSettings.AirDirection,
	)
	s.SlackNotifier.SendSlack(slackMessage)

	return result, nil
}
