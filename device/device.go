package device

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

type Device struct {
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
	Users             []struct {
		ID        string `json:"id"`
		Nickname  string `json:"nickname"`
		Superuser bool   `json:"superuser"`
	} `json:"users"`
	NewestEvents struct {
		Te struct {
			Val       float64   `json:"val"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"te"`
	} `json:"newest_events"`
}

func Get_devices(nature_api_secret string) []Device {
	// -------------------------------------------
	// Send http GET request
	// -------------------------------------------

	var url string = "https://api.nature.global/1/devices"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+nature_api_secret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		errLog.Println(err)
		return nil
	}
	defer res.Body.Close()

	// -------------------------------------------
	// Check http status code
	// -------------------------------------------

	if res.StatusCode != http.StatusOK {
		errLog.Printf("http status code: %d", res.StatusCode)
		return nil
	}

	// -------------------------------------------
	// Decode response to Device
	// -------------------------------------------

	var devices []Device
	json.NewDecoder(res.Body).Decode(&devices)

	return devices
}

func SelectDevice(devices []Device, device_name string) (Device, error) {
	for _, device := range devices {
		if device.Name == device_name {
			return device, nil
		}
	}
	return Device{}, errors.New("No device named " + device_name)
}
