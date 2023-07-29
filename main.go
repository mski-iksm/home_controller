package main

import (
	"flag"

	"home_controller/runner"
	"home_controller/temp_controller"
)

var nature_api_secret string
var action_mode string
var device_name string

// temp_control モードで使用する設定 ==========
// 暑すぎる・寒すぎる気温の設定
var tooHotThreshold float64  // = 27.5
var tooColdThreshold float64 //= 24.0

// 設定可能な気温の下限・上限
var minimumTemperatureSetting float64 //= 23.0
var maximumTemperatureSetting float64 //= 30.0
// =======================================

func init() {
	flag.StringVar(&nature_api_secret, "nature_api_secret", "", "nature remoのAPI")

	flag.StringVar(&action_mode, "action_mode", "send_signal", "send_signal or temp_control")
	flag.StringVar(&device_name, "device_name", "Remo mini", "device name")

	flag.Float64Var(&tooHotThreshold, "tooHotThreshold", 27.5, "この気温以上になると暑いと判定し、エアコンの設定温度を下げる")
	flag.Float64Var(&tooColdThreshold, "tooColdThreshold", 24.0, "この気温未満になると暑いと判定し、エアコンの設定温度を上げる")
	flag.Float64Var(&minimumTemperatureSetting, "minimumTemperatureSetting", 23.0, "エアコンの設定可能温度の下限。安全のためこれ以上下げないようにする。")
	flag.Float64Var(&maximumTemperatureSetting, "maximumTemperatureSetting", 30.0, "エアコンの設定可能温度の上限。安全のためこれ以上上げないようにする。")
}

func main() {
	flag.Parse()
	if action_mode == "send_signal" {
		runner.Send_signal(nature_api_secret)
		return
	}
	if action_mode == "temp_control" {
		temptureMaxMinSettings := temp_controller.TempretureMaxMinSettings{
			TooHotThreshold:           tooHotThreshold,
			TooColdThreshold:          tooColdThreshold,
			MinimumTemperatureSetting: minimumTemperatureSetting,
			MaximumTemperatureSetting: maximumTemperatureSetting,
		}
		runner.TempControl(nature_api_secret, device_name, temptureMaxMinSettings)
		return
	}
}
