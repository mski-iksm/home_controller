package appliance

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

type Signal struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Settings struct {
	Temp      string    `json:"temp"`
	TempUnit  string    `json:"temp_unit"`
	Mode      string    `json:"mode"`
	Vol       string    `json:"vol"`
	Dir       string    `json:"dir"`
	Dirh      string    `json:"dirh"`
	Button    string    `json:"button"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Appliance struct {
	ID     string `json:"id"`
	Device struct {
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
	} `json:"device"`
	Model    interface{} `json:"model"`
	Type     string      `json:"type"`
	Nickname string      `json:"nickname"`
	Image    string      `json:"image"`
	// Settings interface{} `json:"settings"`
	Settings Settings    `json:"settings"`
	Aircon   interface{} `json:"aircon"`
	Signals  []Signal    `json:"signals"`
	Light    struct {
		Buttons []struct {
			Name  string `json:"name"`
			Image string `json:"image"`
			Label string `json:"label"`
		} `json:"buttons"`
		State struct {
			Brightness string `json:"brightness"`
			Power      string `json:"power"`
			LastButton string `json:"last_button"`
		} `json:"state"`
	} `json:"light,omitempty"`
	Tv struct {
		Buttons []struct {
			Name  string `json:"name"`
			Image string `json:"image"`
			Label string `json:"label"`
		} `json:"buttons"`
		State struct {
			Input string `json:"input"`
		} `json:"state"`
	} `json:"tv,omitempty"`
}

func Build_appliances(api_secret string) []Appliance {
	// -------------------------------------------
	// Send http GET request
	// -------------------------------------------

	var url string = "https://api.nature.global/1/appliances"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_secret)
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
	// Decode response to Appliances
	// -------------------------------------------

	var appliances []Appliance
	json.NewDecoder(res.Body).Decode(&appliances)

	return appliances
}

func Select_applicance(appliances []Appliance) Appliance {
	// 機器を選択
	for i, appliance := range appliances {
		fmt.Println(i, appliance.Device.Name, appliance.Nickname)
	}

	fmt.Print("数値もしくは機器名を入力してください > ")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		in := scanner.Text()
		fmt.Println("in: ", in)

		for i, appliance := range appliances {
			if strconv.Itoa(i) == in || appliance.Nickname == in {
				return appliance
			}
		}

		fmt.Print("入力が不正です。数値もしくは機器名を入力してください > ")
		continue
	}
}

func FilterAppliances(appliances []Appliance, device_name string) []Appliance {
	var filtered_appliances []Appliance
	for _, appliance := range appliances {
		if appliance.Device.Name == device_name {
			filtered_appliances = append(filtered_appliances, appliance)
		}
	}
	return filtered_appliances
}
