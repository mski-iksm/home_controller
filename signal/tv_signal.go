package signal

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/mski-iksm/home_controller/appliance"
)

func select_tv_button(appliance appliance.Appliance) string {
	for i, button := range appliance.Tv.Buttons {
		fmt.Println(i, button.Label)
	}

	fmt.Print("数値もしくはボタン名を入力してください > ")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		in := scanner.Text()
		fmt.Println("in: ", in)

		for i, button := range appliance.Tv.Buttons {
			if strconv.Itoa(i) == in || button.Label == in {
				return button.Name
			}
		}

		fmt.Println("入力が不正です。数値もしくはボタン名を入力してください > ")
		continue
	}
}

func post_tv_signal(api_secret string, tv_appliance_id string, selected_tv_button string) int {
	var fullurl string = "https://api.nature.global/1/appliances/" + tv_appliance_id + "/tv"

	post_body := url.Values{}
	post_body.Set("button", selected_tv_button)

	return post_signal(api_secret, fullurl, post_body)
}

func send_tv_signal(api_secret string, appliance appliance.Appliance) int {
	var tv_appliance_id string = appliance.ID
	var selected_tv_button string = select_tv_button(appliance)
	return post_tv_signal(api_secret, tv_appliance_id, selected_tv_button)
}
