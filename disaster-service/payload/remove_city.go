package payload

type RemoveCity struct {
	Country string `json:"country"`
	City    string `json:"city"`
}
