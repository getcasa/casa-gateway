package devices

// Motion define xiaomi motion sensor
type Motion struct {
	status   string `json:"status"`
	NoMotion string `json:"no_motion"`
	Voltage  int    `json:"voltage"`
}
