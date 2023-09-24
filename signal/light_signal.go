package signal

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/mski-iksm/home_controller/appliance"
)

func select_light_button(appliance appliance.Appliance) string {
	for i, button := range appliance.Light.Buttons {
		fmt.Println(i, button.Label)
	}

	fmt.Print("数値もしくはボタン名を入力してください > ")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		in := scanner.Text()
		fmt.Println("in: ", in)

		for i, button := range appliance.Light.Buttons {
			if strconv.Itoa(i) == in || button.Label == in {
				return button.Name
			}
		}

		fmt.Println("入力が不正です。数値もしくはボタン名を入力してください > ")
		continue
	}
}

func post_light_signal(api_secret string, light_appliance_id string, selected_light_button string) int {
	var fullurl string = "https://api.nature.global/1/appliances/" + light_appliance_id + "/light"

	post_body := url.Values{}
	post_body.Set("button", selected_light_button)

	return post_signal(api_secret, fullurl, post_body)
}

func send_light_signal(api_secret string, appliance appliance.Appliance) int {
	var light_appliance_id string = appliance.ID
	var selected_light_button string = select_light_button(appliance)
	return post_light_signal(api_secret, light_appliance_id, selected_light_button)
}
