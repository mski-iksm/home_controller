package temp_controller

import (
	"fmt"
	"home_controller/signal"
	"os"
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
	PreparationThreshold      float64
	MinimumTemperatureSetting float64
	MaximumTemperatureSetting float64
}

var TOO_HOT_THRESHOLD_HARD_LIMIT float64 = 30.0
var TOO_COLD_THRESHOLD_HARD_LIMIT float64 = 12.0
var PREPARATION_THRESHOLD_HARD_LIMIT float64 = 0.0
var MINIMUM_TEMPERATURE_SETTING_HARD_LIMIT float64 = 12.0
var MAXIMUM_TEMPERATURE_SETTING_HARD_LIMIT float64 = 30.0

func assertTresholdSettings(tempretureMaxMinSettings TempretureMaxMinSettings) error {
	if tempretureMaxMinSettings.TooHotThreshold > TOO_HOT_THRESHOLD_HARD_LIMIT {
		return fmt.Errorf("tooHotThresholdが閾値を超えています。閾値: %v, tooHotThreshold: %v\n", TOO_HOT_THRESHOLD_HARD_LIMIT, tempretureMaxMinSettings.TooHotThreshold)
	}
	if tempretureMaxMinSettings.TooColdThreshold < TOO_COLD_THRESHOLD_HARD_LIMIT {
		return fmt.Errorf("tooColdThresholdが閾値を超えています。閾値: %v, tooColdThreshold: %v\n", TOO_COLD_THRESHOLD_HARD_LIMIT, tempretureMaxMinSettings.TooColdThreshold)
	}
	if tempretureMaxMinSettings.PreparationThreshold < PREPARATION_THRESHOLD_HARD_LIMIT {
		return fmt.Errorf("preparationThresholdが閾値を下回っています。閾値: %v, preparationThreshold: %v\n", PREPARATION_THRESHOLD_HARD_LIMIT, tempretureMaxMinSettings.PreparationThreshold)
	}
	if tempretureMaxMinSettings.MinimumTemperatureSetting < MINIMUM_TEMPERATURE_SETTING_HARD_LIMIT {
		return fmt.Errorf("minimumTemperatureSettingが閾値を超えています。閾値: %v, minimumTemperatureSetting: %v\n", MINIMUM_TEMPERATURE_SETTING_HARD_LIMIT, tempretureMaxMinSettings.MinimumTemperatureSetting)
	}
	if tempretureMaxMinSettings.MaximumTemperatureSetting > MAXIMUM_TEMPERATURE_SETTING_HARD_LIMIT {
		return fmt.Errorf("maximumTemperatureSettingが閾値を超えています。閾値: %v, maximumTemperatureSetting: %v\n", MAXIMUM_TEMPERATURE_SETTING_HARD_LIMIT, tempretureMaxMinSettings.MaximumTemperatureSetting)
	}
	return nil
}

func ConstructTempretureMaxMinSettings(tooHotThreshold float64, tooColdThreshold float64, preparationThreshold float64, minimumTemperatureSetting float64, maximumTemperatureSetting float64) *TempretureMaxMinSettings {
	tempretureMaxMinSettings := TempretureMaxMinSettings{
		TooHotThreshold:           tooHotThreshold,
		TooColdThreshold:          tooColdThreshold,
		PreparationThreshold:      preparationThreshold,
		MinimumTemperatureSetting: minimumTemperatureSetting,
		MaximumTemperatureSetting: maximumTemperatureSetting,
	}
	hardLimitAssertionError := assertTresholdSettings(tempretureMaxMinSettings)
	if hardLimitAssertionError != nil {
		errLog.Println(hardLimitAssertionError)
		os.Exit(1)
	}
	return &tempretureMaxMinSettings
}
