package payload

type Disaster struct {
	DisasterName           string `json:"disaster_name"`
	DisasterTimestampStart string `json:"disaster_timestamp_start"`
	DisasterTimestampEnd   string `json:"disaster_timestamp_end"`
	DisasterSeverity       string `json:"disaster_severity"`
	DisasterDescription    string `json:"disaster_description"`
}

type DisasterList struct {
	Country   string     `json:"country"`
	City      string     `json:"city"`
	Disasters []Disaster `json:"disasters"`
}
