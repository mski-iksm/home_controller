package main

import (
	"flag"
	"os"

	"github.com/mski-iksm/home_controller/runner"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

var nature_api_secret string
var action_mode string
var device_name string

// temp_control モードで使用する設定 ==========
// 暑すぎる・寒すぎる気温の設定
var tooHotThreshold float64
var tooColdThreshold float64

// 暑すぎ・寒すぎになる前に気温設定を変更する機能。0.0にすると機能しない
var preparationThreshold float64

// 設定可能な気温の下限・上限
var minimumTemperatureSetting float64
var maximumTemperatureSetting float64

// =======================================

// slack設定 ==============
var slackToken string
var slackChannel string

// =======================

func init() {
	flag.StringVar(&nature_api_secret, "nature_api_secret", "", "nature remoのAPI")

	flag.StringVar(&action_mode, "action_mode", "send_signal", "send_signal or temp_control")
	flag.StringVar(&device_name, "device_name", "", "device name")

	flag.Float64Var(&tooHotThreshold, "tooHotThreshold", 27.5, "この気温以上になると暑いと判定し、エアコンの設定温度を下げる")
	flag.Float64Var(&tooColdThreshold, "tooColdThreshold", 24.0, "この気温未満になると暑いと判定し、エアコンの設定温度を上げる")
	flag.Float64Var(&preparationThreshold, "preparationThreshold", 0.0, "暑すぎる・寒すぎる気温になる前に、エアコンの設定温度を変更する機能。0.0にすると機能しない")
	flag.Float64Var(&minimumTemperatureSetting, "minimumTemperatureSetting", 23.0, "エアコンの設定可能温度の下限。安全のためこれ以上下げないようにする。")
	flag.Float64Var(&maximumTemperatureSetting, "maximumTemperatureSetting", 30.0, "エアコンの設定可能温度の上限。安全のためこれ以上上げないようにする。")

	flag.StringVar(&slackToken, "slackToken", "", "slackのtoken")
	flag.StringVar(&slackChannel, "slackChannel", "", "通知を送信するslackのチャンネル名。#から始める。")
}

func main() {
	flag.Parse()

	// nature_api_secretが空の場合には環境変数から読み込む
	if nature_api_secret == "" {
		nature_api_secret = os.Getenv("NATURE_API_SECRET")
	}

	// device_nameが空の場合には環境変数から読み込む
	if device_name == "" {
		device_name = os.Getenv("DEVICE_NAME")
	}

	// slackTokenが空の場合には環境変数から読み込む
	if slackToken == "" {
		slackToken = os.Getenv("SLACK_TOKEN")
	}

	// slackChannelが空の場合には環境変数から読み込む
	if slackChannel == "" {
		slackChannel = os.Getenv("SLACK_CHANNEL")
	}

	slackObject := slack.SlackObject{
		SlackToken:   slackToken,
		SlackChannel: slackChannel,
	}

	if action_mode == "send_signal" {
		runner.Send_signal(nature_api_secret)
		return
	}
	if action_mode == "temp_control" {
		temptureMaxMinSettings := temp_controller.ConstructTempretureMaxMinSettings(tooHotThreshold, tooColdThreshold, preparationThreshold, minimumTemperatureSetting, maximumTemperatureSetting)
		runner.TempControl(nature_api_secret, device_name, *temptureMaxMinSettings, slackObject)
		return
	}
	if action_mode == "notify_temp" {
		temptureMaxMinSettings := temp_controller.ConstructTempretureMaxMinSettings(tooHotThreshold, tooColdThreshold, preparationThreshold, minimumTemperatureSetting, maximumTemperatureSetting)
		runner.NotifyTemp(nature_api_secret, device_name, *temptureMaxMinSettings, slackObject)
		return
	}
}
