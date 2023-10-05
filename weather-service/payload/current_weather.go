package payload

type CurrentWeatherResponse struct {
	Country          string   `json:"country"`
	Location         Location `json:"location"`
	Timestamp        string   `json:"timestamp"`
	Temperature      int      `json:"temperature"`
	Humidity         int      `json:"humidity"`
	WeatherCondition string   `json:"weather_condition"`
}

func GenerateCurrentWeatherResponse(country string, location Location, timestamp string, temperature int, humidity int, weatherCondition string) CurrentWeatherResponse {
	return CurrentWeatherResponse{
		Country:          country,
		Location:         location,
		Timestamp:        timestamp,
		Temperature:      temperature,
		Humidity:         humidity,
		WeatherCondition: weatherCondition,
	}
}
