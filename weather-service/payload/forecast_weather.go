package payload

type ForecastWeatherDay struct {
	Timestamp        string `json:"timestamp"`
	TemperatureHigh  int16  `json:"temperature_high"`
	TemperatureLow   int16  `json:"temperature_low"`
	Humidity         int8   `json:"humidity"`
	WeatherCondition string `json:"weather_condition"`
}

type ForecastWeatherResponse struct {
	Country             string               `json:"country"`
	Location            Location             `json:"location"`
	ForecastWeatherDays []ForecastWeatherDay `json:"forecast_weather_days"`
}
