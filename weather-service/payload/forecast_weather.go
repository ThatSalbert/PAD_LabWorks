package payload

type ForecastWeatherDay struct {
	Timestamp        string `json:"timestamp"`
	TemperatureHigh  int    `json:"temperature_high"`
	TemperatureLow   int    `json:"temperature_low"`
	Humidity         int    `json:"humidity"`
	WeatherCondition string `json:"weather_condition"`
}

type ForecastWeatherResponse struct {
	Country             string               `json:"country"`
	Location            Location             `json:"location"`
	ForecastWeatherDays []ForecastWeatherDay `json:"forecast_weather_days"`
}

func GenerateForecastWeatherResponse(country string, location Location, forecastWeatherDays []ForecastWeatherDay) ForecastWeatherResponse {
	return ForecastWeatherResponse{
		Country:             country,
		Location:            location,
		ForecastWeatherDays: forecastWeatherDays,
	}
}
