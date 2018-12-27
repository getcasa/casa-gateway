package devices

// WeatherV1 define xiaomi weather v1 sensor
type WeatherV1 struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	Pressure    string `json:"pressure"`
	Voltage     int    `json:"voltage"`
}
