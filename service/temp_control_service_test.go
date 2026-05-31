package service

import (
	"testing"
	"time"

	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/signal"
	"github.com/mski-iksm/home_controller/temp_controller"
)

type fakeNatureClient struct {
	devices        []device.Device
	appliances     []appliance.Appliance
	postCount      int
	lastApplianceID string
	lastSettings    signal.AirconSettings
}

func (f *fakeNatureClient) GetDevices() []device.Device {
	return f.devices
}

func (f *fakeNatureClient) GetAppliances() []appliance.Appliance {
	return f.appliances
}

func (f *fakeNatureClient) PostAirconSignal(applianceID string, settings signal.AirconSettings) int {
	f.postCount++
	f.lastApplianceID = applianceID
	f.lastSettings = settings
	return 200
}

type fakeSlackNotifier struct {
	messages []string
}

func (f *fakeSlackNotifier) SendSlack(message string) {
	f.messages = append(f.messages, message)
}

type fakeAlertNotifier struct {
	alerts []temp_controller.TemperatureAlert
}

func (f *fakeAlertNotifier) SendTemperatureAlert(alert temp_controller.TemperatureAlert) {
	f.alerts = append(f.alerts, alert)
}

func TestTempControlService_Run_UsesInjectedInterfaces(t *testing.T) {
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	currentTime := time.Date(2020, 1, 1, 5, 0, 0, 0, tokyo)

	fakeClient := &fakeNatureClient{
		devices: []device.Device{
			{
				Name: "Remo",
				NewestEvents: device.NewestEvents{
					Te: device.Temperature{
						Val:       30.0,
						CreatedAt: currentTime,
					},
				},
			},
		},
		appliances: []appliance.Appliance{
			{
				ID: "aircon-1",
				Device: struct {
					Name              string    `json:"name"`
					ID                string    `json:"id"`
					CreatedAt         time.Time `json:"created_at"`
					UpdatedAt         time.Time `json:"updated_at"`
					MacAddress        string    `json:"mac_address"`
					BtMacAddress      string    `json:"bt_mac_address"`
					SerialNumber      string    `json:"serial_number"`
					FirmwareVersion   string    `json:"firmware_version"`
					TemperatureOffset int       `json:"temperature_offset"`
					HumidityOffset    int       `json:"humidity_offset"`
				}{
					Name: "Remo",
				},
				Aircon: 1,
				Settings: appliance.Settings{
					Temp:      "27.0",
					Mode:      "cool",
					Vol:       "auto",
					Dir:       "auto",
					UpdatedAt: currentTime.Add(-2 * time.Hour),
				},
			},
		},
	}
	fakeSlack := &fakeSlackNotifier{}
	fakeAlert := &fakeAlertNotifier{}

	service := TempControlService{
		NatureClient:             fakeClient,
		SlackNotifier:            fakeSlack,
		TemperatureAlertNotifier: fakeAlert,
	}

	result, err := service.Run(TempControlRequest{
		DeviceName: "Remo",
		Settings: temp_controller.TempretureMaxMinSettings{
			TooHotThreshold:           27.5,
			TooColdThreshold:          24.0,
			PreparationThreshold:      0.0,
			MinimumTemperatureSetting: 23.0,
			MaximumTemperatureSetting: 30.0,
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if fakeClient.postCount != 1 {
		t.Fatalf("expected one aircon post, got %d", fakeClient.postCount)
	}
	if fakeClient.lastApplianceID != "aircon-1" {
		t.Fatalf("expected aircon-1, got %s", fakeClient.lastApplianceID)
	}
	if len(fakeSlack.messages) != 1 {
		t.Fatalf("expected one slack message, got %d", len(fakeSlack.messages))
	}
	if len(fakeAlert.alerts) != 1 {
		t.Fatalf("expected one alert, got %d", len(fakeAlert.alerts))
	}
	if !result.AirconChanged {
		t.Fatalf("expected aircon change result")
	}
}
