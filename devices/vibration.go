package devices

// Vibration define xiaomi vibration
type Vibration struct {
	BedActivity  string `json:"bed_activity"`
	Coordination string `json:"coordination"`
	Voltage      int    `json:"voltage"`
}
