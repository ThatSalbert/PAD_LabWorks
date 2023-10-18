package payload

type DisasterType struct {
	DisasterName        string `json:"disaster_name"`
	DisasterDescription string `json:"disaster_description"`
}

type DisasterTypeListResponse struct {
	DisasterTypes []DisasterType `json:"disaster_types"`
}
