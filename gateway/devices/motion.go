package devices

// Motion define xiaomi motion sensor
type Motion struct {
	Status   string `json:"status"`
	NoMotion string `json:"no_motion"`
	Voltage  int    `json:"voltage"`
}
