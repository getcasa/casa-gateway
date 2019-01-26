package devices

// Gateway define xiaomi gateway
type Gateway struct {
	IP           string `json:"ip"`
	RGB          int    `json:"rgb"`
	Illumination int    `json:"illumination"`
}
