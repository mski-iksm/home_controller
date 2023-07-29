package signal

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"home_controller/appliance"
)

func select_other_signal(appliance appliance.Appliance) string {
	// 命令を選択
	for i, signal := range appliance.Signals {
		fmt.Println(i, signal.Name)
	}

	fmt.Print("数値もしくは命令を入力してください > ")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		in := scanner.Text()
		fmt.Println("in: ", in)

		for i, signal := range appliance.Signals {
			if strconv.Itoa(i) == in || signal.Name == in {
				return signal.ID
			}
		}

		fmt.Println("入力が不正です。数値もしくは命令を入力してください > ")
		continue
	}
}

func post_other_signal(api_secret string, selected_signal_id string) int {
	var fullurl string = "https://api.nature.global/1/signals/" + selected_signal_id + "/send"

	post_body := url.Values{}
	return post_signal(api_secret, fullurl, post_body)
}

func send_other_signal(api_secret string, appliance appliance.Appliance) int {
	var selected_signal_id string = select_other_signal(appliance)
	return post_other_signal(api_secret, selected_signal_id)
}
