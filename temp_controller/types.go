package temp_controller

import (
	"home_controller/signal"
	"time"
)

type CurrentAirConSettings struct {
	AirconSettings signal.AirconSettings
	PowerOn        bool
	UpdatedAt      time.Time
}

type NewAirConSettings struct {
	AirconSettings signal.AirconSettings
	PowerOn        bool
}

type CurrentTempreture struct {
	Tempreture float64
	UpdatedAt  time.Time
}

type AirconOrderParameters struct {
	ApplianceId    string
	AirconSettings signal.AirconSettings
}

type TempretureMaxMinSettings struct {
	TooHotThreshold           float64
	TooColdThreshold          float64
	MinimumTemperatureSetting float64
	MaximumTemperatureSetting float64
}
