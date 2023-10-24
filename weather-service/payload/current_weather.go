package payload

type Disaster struct {
	DisasterName           string `json:"disaster_name"`
	DisasterTimestampStart string `json:"disaster_timestamp_start"`
	DisasterTimestampEnd   string `json:"disaster_timestamp_end"`
	DisasterSeverity       string `json:"disaster_severity"`
	DisasterDescription    string `json:"disaster_description"`
}

type CurrentWeatherResponse struct {
	Country          string     `json:"country"`
	Location         Location   `json:"location"`
	Timestamp        string     `json:"timestamp"`
	Temperature      int16      `json:"temperature"`
	Humidity         int8       `json:"humidity"`
	WeatherCondition string     `json:"weather_condition"`
	Disasters        []Disaster `json:"disasters"`
}
