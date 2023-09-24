package runner

import (
	"github.com/mski-iksm/home_controller/appliance"
	"github.com/mski-iksm/home_controller/signal"
)

func Send_signal(nature_api_secret string) {
	var appliances []appliance.Appliance = appliance.Build_appliances(nature_api_secret)
	var selected_appliance appliance.Appliance = appliance.Select_applicance(appliances)

	signal.Send_signal(nature_api_secret, selected_appliance)
}
