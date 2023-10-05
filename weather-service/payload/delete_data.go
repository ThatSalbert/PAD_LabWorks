package payload

type DeleteData struct {
	Country   string `json:"country"`
	City      string `json:"city"`
	Timestamp string `json:"timestamp"`
}
