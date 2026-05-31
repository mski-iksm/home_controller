package service

import (
	"testing"
	"time"

	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/device"
	"github.com/mski-iksm/home_controller/temp_controller"
)

func TestNotifyTempService_Run_SendsSlackForHotTemperature(t *testing.T) {
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	currentTime := time.Date(2020, 1, 1, 5, 0, 0, 0, tokyo)

	fakeClient := &fakeNatureClient{
		devices: []device.Device{
			{
				Name: "Remo",
				NewestEvents: device.NewestEvents{
					Te: device.Temperature{
						Val:       28.0,
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
					Button:    "power-on",
					UpdatedAt: currentTime.Add(-2 * time.Hour),
				},
			},
		},
	}
	fakeSlack := &fakeSlackNotifier{}

	service := NotifyTempService{
		NatureClient:  fakeClient,
		SlackNotifier: fakeSlack,
	}

	result, err := service.Run(NotifyTempRequest{
		DeviceName: "Remo",
		Settings: temp_controller.TemperatureNotifySettings{
			TooHotThreshold:  27.5,
			TooColdThreshold: 24.0,
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(fakeSlack.messages) != 1 {
		t.Fatalf("expected one slack message, got %d", len(fakeSlack.messages))
	}
	if !result.TemperatureNotification.ShouldNotify {
		t.Fatalf("expected notification result")
	}
	if result.TemperatureNotification.Reason != "temperature is too hot" {
		t.Fatalf("unexpected reason: %s", result.TemperatureNotification.Reason)
	}
}

func TestNotifyTempService_Run_SkipsSlackWhenWithinThreshold(t *testing.T) {
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	currentTime := time.Date(2020, 1, 1, 5, 0, 0, 0, tokyo)

	fakeClient := &fakeNatureClient{
		devices: []device.Device{
			{
				Name: "Remo",
				NewestEvents: device.NewestEvents{
					Te: device.Temperature{
						Val:       25.0,
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
					Button:    "power-on",
					UpdatedAt: currentTime.Add(-2 * time.Hour),
				},
			},
		},
	}
	fakeSlack := &fakeSlackNotifier{}

	service := NotifyTempService{
		NatureClient:  fakeClient,
		SlackNotifier: fakeSlack,
	}

	result, err := service.Run(NotifyTempRequest{
		DeviceName: "Remo",
		Settings: temp_controller.TemperatureNotifySettings{
			TooHotThreshold:  27.5,
			TooColdThreshold: 24.0,
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(fakeSlack.messages) != 0 {
		t.Fatalf("expected no slack messages, got %d", len(fakeSlack.messages))
	}
	if result.TemperatureNotification.ShouldNotify {
		t.Fatalf("expected no notification result")
	}
}
