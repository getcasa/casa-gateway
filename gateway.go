package main

// Gateway define xiaomi gateway
type Gateway struct {
	CMD     string      `json:"cmd"`
	Model   string      `json:"model"`
	SID     string      `json:"sid"`
	ShortID interface{} `json:"short_id"`
	Token   string      `json:"token"`
	IP      string      `json:"ip"`
	Data    interface{} `json:"data"`
}
