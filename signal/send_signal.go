package signal

import (
	"log"
	"os"

	"home_controller/appliance"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

func is_aircon(appliance appliance.Appliance) bool {
	if appliance.Aircon == nil {
		return false
	}
	return true
}

func is_light(appliance appliance.Appliance) bool {
	if len(appliance.Light.Buttons) == 0 {
		return false
	}
	return true
}

func is_tv(appliance appliance.Appliance) bool {
	if len(appliance.Tv.Buttons) == 0 {
		return false
	}
	return true
}

func Send_signal(api_secret string, selected_appliance appliance.Appliance) int {
	if is_aircon(selected_appliance) {
		return send_aircon_signal(api_secret, selected_appliance)
	}
	if is_light(selected_appliance) {
		return send_light_signal(api_secret, selected_appliance)
	}
	if is_tv(selected_appliance) {
		return send_tv_signal(api_secret, selected_appliance)
	}
	return send_other_signal(api_secret, selected_appliance)
}
