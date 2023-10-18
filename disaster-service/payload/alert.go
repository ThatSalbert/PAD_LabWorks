package payload

type AddAlert struct {
	Country        string `json:"country"`
	City           string `json:"city"`
	DisasterName   string `json:"disaster_name"`
	TimestampStart string `json:"timestamp_start"`
	TimestampEnd   string `json:"timestamp_end"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
}

type UpdateAlert struct {
	Country        string `json:"country"`
	City           string `json:"city"`
	DisasterName   string `json:"disaster_name"`
	TimestampStart string `json:"timestamp_start"`
	TimestampEnd   string `json:"timestamp_end"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
}
