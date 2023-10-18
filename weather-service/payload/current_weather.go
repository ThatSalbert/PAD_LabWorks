package payload

type CurrentWeatherResponse struct {
	Country          string   `json:"country"`
	Location         Location `json:"location"`
	Timestamp        string   `json:"timestamp"`
	Temperature      int16    `json:"temperature"`
	Humidity         int8     `json:"humidity"`
	WeatherCondition string   `json:"weather_condition"`
}
