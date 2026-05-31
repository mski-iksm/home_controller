package service

import (
	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/temp_controller"
)

type AirconContext struct {
	ApplianceContext
	AirconAppliance       appliance.Appliance
	CurrentAirconSettings temp_controller.CurrentAirConSettings
}

type ApplianceContext struct {
	Device             device.Device
	FilteredAppliances []appliance.Appliance
	CurrentTemperature temp_controller.CurrentTempreture
}

func LoadApplianceContext(client NatureClient, deviceName string) (ApplianceContext, error) {
	var context ApplianceContext

	devices := client.GetDevices()
	appliances := client.GetAppliances()

	selectedDevice, err := device.SelectDevice(devices, deviceName)
	if err != nil {
		return context, err
	}

	filteredAppliances := appliance.FilterAppliances(appliances, deviceName)

	context.Device = selectedDevice
	context.FilteredAppliances = filteredAppliances
	context.CurrentTemperature = temp_controller.Get_current_temperature(selectedDevice)

	return context, nil
}

func LoadAirconContext(client NatureClient, deviceName string) (AirconContext, error) {
	applianceContext, err := LoadApplianceContext(client, deviceName)
	if err != nil {
		return AirconContext{}, err
	}

	airconAppliance, err := temp_controller.Find_aircon_appliance(applianceContext.FilteredAppliances)
	if err != nil {
		return AirconContext{}, err
	}

	context := AirconContext{
		ApplianceContext:      applianceContext,
		AirconAppliance:       airconAppliance,
		CurrentAirconSettings: temp_controller.GetCurrentAirconSettings(airconAppliance),
	}

	return context, nil
}
