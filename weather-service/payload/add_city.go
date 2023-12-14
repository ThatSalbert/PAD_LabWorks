package payload

type AddCity struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
