package main

import (
	"flag"

	"home_controller/runner"
)

var nature_api_secret string
var action_mode string
var device_name string

func init() {
	// Var() 関数の引数は以下の通り
	// 第一引数は束縛する変数のポインタ
	// 第二引数はフラグ名
	// 第三引数はデフォルト値
	// 第四引数はフラグの説明
	flag.StringVar(&nature_api_secret, "nature_api_secret", "", "nature remoのAPI")

	flag.StringVar(&action_mode, "action_mode", "send_signal", "send_signal or temp_control")
	flag.StringVar(&device_name, "device_name", "Remo mini", "device name")
}

func main() {
	flag.Parse()
	if action_mode == "send_signal" {
		runner.Send_signal(nature_api_secret)
		return
	}
	if action_mode == "temp_control" {
		runner.Temp_control(nature_api_secret, device_name)
		return
	}

}
