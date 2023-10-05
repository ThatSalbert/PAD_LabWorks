package payload

type AddDataWeather struct {
	Country          string `json:"country"`
	City             string `json:"city"`
	Timestamp        string `json:"timestamp"`
	Temperature      int    `json:"temperature"`
	Humidity         int    `json:"humidity"`
	WeatherCondition string `json:"weather_condition"`
}

type AddDataForecast struct {
	Country          string `json:"country"`
	City             string `json:"city"`
	Timestamp        string `json:"timestamp"`
	TemperatureHigh  int    `json:"temperature_high"`
	TemperatureLow   int    `json:"temperature_low"`
	Humidity         int    `json:"humidity"`
	WeatherCondition string `json:"weather_condition"`
}
