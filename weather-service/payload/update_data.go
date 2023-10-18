package payload

type UpdateDataWeather struct {
	WeatherID        int64  `json:"weather_id"`
	Country          string `json:"country"`
	City             string `json:"city"`
	Timestamp        string `json:"timestamp"`
	Temperature      int16  `json:"temperature"`
	Humidity         int8   `json:"humidity"`
	WeatherCondition string `json:"weather_condition"`
}

type UpdateDataForecast struct {
	ForecastID       int64  `json:"forecast_id"`
	Country          string `json:"country"`
	City             string `json:"city"`
	Timestamp        string `json:"timestamp"`
	TemperatureHigh  int16  `json:"temperature_high"`
	TemperatureLow   int16  `json:"temperature_low"`
	Humidity         int8   `json:"humidity"`
	WeatherCondition string `json:"weather_condition"`
}
