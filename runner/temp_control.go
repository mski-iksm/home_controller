package runner

import (
	"log"
	"os"
	"time"

	"github.com/mski-iksm/home_controller/service"
	"github.com/mski-iksm/home_controller/slack"
	"github.com/mski-iksm/home_controller/temp_controller"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

func getTime() time.Time {
	nowTime := time.Now().UTC() // 現在時刻をUTCで取得

	// タイムゾーンからJSTを読み込み
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		errLog.Println(err)
		os.Exit(1)
	}
	timeTokyo := nowTime.In(tokyo)
	return timeTokyo
}

func TempControl(nature_api_secret string, device_name string, temptureMaxMinSettings temp_controller.TempretureMaxMinSettings, slackObject slack.SlackObject, ntfyUrl string) {
	// 時刻を出力
	currentTime := getTime()
	appLog.Printf("時刻: %v\n", currentTime)

	tempControlService := service.NewTempControlService(nature_api_secret, slackObject, ntfyUrl)
	_, err := tempControlService.Run(service.TempControlRequest{
		DeviceName: device_name,
		Settings:   temptureMaxMinSettings,
	})
	if err != nil {
		errLog.Println(err)
		return
	}
}
