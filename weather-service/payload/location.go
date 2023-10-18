package payload

type Location struct {
	City      string  `json:"city"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type LocationResponse struct {
	Country   string     `json:"country"`
	Locations []Location `json:"locations"`
}
