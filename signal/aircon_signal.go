package signal

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/mski-iksm/home_controller/appliance"
)

type AirconSettings struct {
	Temperature   float64
	OperationMode string
	AirVolume     string
	AirDirection  string
}

func getCurrentAirconSettings(appliance appliance.Appliance) AirconSettings {
	tempreture, err := strconv.ParseFloat(appliance.Settings.Temp, 64)
	if err != nil {
		fmt.Errorf("温度 %v をfloat64にcastできません", appliance.Settings.Temp)
	}

	return AirconSettings{
		Temperature:   tempreture,
		OperationMode: appliance.Settings.Mode,
		AirVolume:     appliance.Settings.Vol,
		AirDirection:  appliance.Settings.Dir,
	}
}

func parse_temperature(command string, current_settings AirconSettings) float64 {
	for temperature := 15; temperature < 33; temperature++ {
		if strings.Contains(command, strconv.Itoa(temperature)) {
			return float64(temperature)
		}
	}
	return current_settings.Temperature
}

func parse_operation_mode(command string, current_settings AirconSettings) string {
	if strings.Contains(command, "AC") || strings.Contains(command, "ac") || strings.Contains(command, "冷") || strings.Contains(command, "クーラー") || strings.Contains(command, "エアコン") || strings.Contains(command, "cool") {
		return "cool"
	}
	if strings.Contains(command, "dry") || strings.Contains(command, "DRY") || strings.Contains(command, "ドライ") || strings.Contains(command, "湿") {
		return "dry"
	}
	if strings.Contains(command, "warm") || strings.Contains(command, "hot") || strings.Contains(command, "暖") || strings.Contains(command, "温") {
		return "warm"
	}
	return current_settings.OperationMode
}
func parse_air_volume(command string, current_settings AirconSettings) string {
	for volume := 1; volume < 4; volume++ {
		if strings.Contains(command, "風量"+strconv.Itoa(volume)) {
			return strconv.Itoa(volume)
		}
	}
	if strings.Contains(command, "風量AUTO") || strings.Contains(command, "風量auto") || strings.Contains(command, "風量自動") {
		return "auto"
	}
	return current_settings.AirVolume
}
func parse_air_direction(command string, current_settings AirconSettings) string {
	// for direction := 1; direction < 6; direction++ {
	// 	if strings.Contains(command, "風向"+strconv.Itoa(direction)) {
	// 		return strconv.Itoa(direction)
	// 	}
	// }
	// if strings.Contains(command, "風向swing") || strings.Contains(command, "風向SWING") || strings.Contains(command, "風向スイング") {
	// 	return "swing"
	// }
	// if strings.Contains(command, "風向auto") || strings.Contains(command, "風向自動") || strings.Contains(command, "風向オート") || strings.Contains(command, "風向AUTO") {
	// 	return "auto"
	// }
	// return current_settings.AirDirection
	return ""
}

func parse_aircon_command(command string, current_settings AirconSettings) AirconSettings {
	var temperature float64 = parse_temperature(command, current_settings)
	var operation_mode string = parse_operation_mode(command, current_settings)
	var air_volume string = parse_air_volume(command, current_settings)
	var air_direction string = parse_air_direction(command, current_settings)
	return AirconSettings{
		Temperature:   temperature,
		OperationMode: operation_mode,
		AirVolume:     air_volume,
		AirDirection:  air_direction,
	}
}

func select_aircon_settings(appliance appliance.Appliance) AirconSettings {
	var current_settings AirconSettings = getCurrentAirconSettings(appliance)

	fmt.Print("コマンドを入力してください >")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		command := scanner.Text()

		if len(command) > 0 {
			return parse_aircon_command(command, current_settings)
		}

		fmt.Println("入力が不正です。コマンドを入力してください > ")
		continue
	}
}

func PostAirconSignal(api_secret string, aircon_appliance_id string, aircon_settings AirconSettings) int {
	var fullurl string = "https://api.nature.global/1/appliances/" + aircon_appliance_id + "/aircon_settings"

	post_body := url.Values{}
	post_body.Set("temperature", strconv.FormatFloat(aircon_settings.Temperature, 'f', 1, 64))
	post_body.Set("operation_mode", aircon_settings.OperationMode)
	post_body.Set("air_volume", aircon_settings.AirVolume)
	post_body.Set("air_direction", aircon_settings.AirDirection)

	fmt.Printf("温度: %s\n", strconv.FormatFloat(aircon_settings.Temperature, 'f', 1, 64))
	fmt.Printf("モード: %s\n", aircon_settings.OperationMode)
	fmt.Printf("風量: %s\n", aircon_settings.AirVolume)
	fmt.Printf("風向: %s\n", aircon_settings.AirDirection)

	return post_signal(api_secret, fullurl, post_body)
}

func send_aircon_signal(api_secret string, appliance appliance.Appliance) int {
	var aircon_appliance_id string = appliance.ID
	var aircon_settings AirconSettings = select_aircon_settings(appliance)
	return PostAirconSignal(api_secret, aircon_appliance_id, aircon_settings)
}
