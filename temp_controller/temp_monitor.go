package temp_controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mski-iksm/home_controller/device"
)

func MonitorTempreture(device device.Device, temptureMaxMinSettings TempretureMaxMinSettings, ntfyUrl string) {
	// deviceから今の気温を取得
	current_tempreture := Get_current_temperature(device)

	// 現在気温が既定値を0.5度以上超えている場合は、緊急アラート
	if current_tempreture.Tempreture >= temptureMaxMinSettings.TooHotThreshold+0.5 {
		req, _ := http.NewRequest("POST", ntfyUrl, strings.NewReader("緊急アラート: 現在の温度が設定値を超えています。\n現在温度: "+fmt.Sprintf("%.1f", current_tempreture.Tempreture)+"\n設定温度: "+fmt.Sprintf("%.1f", temptureMaxMinSettings.TooHotThreshold)))
		req.Header.Set("Priority", "5")
		http.DefaultClient.Do(req)
	}

	if current_tempreture.Tempreture <= temptureMaxMinSettings.TooColdThreshold-0.5 {
		req, _ := http.NewRequest("POST", ntfyUrl, strings.NewReader("緊急アラート: 現在の温度が設定値を下回っています。\n現在温度: "+fmt.Sprintf("%.1f", current_tempreture.Tempreture)+"\n設定温度: "+fmt.Sprintf("%.1f", temptureMaxMinSettings.TooColdThreshold)))
		req.Header.Set("Priority", "5")
		http.DefaultClient.Do(req)
	}
}
