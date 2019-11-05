package gateway

import (
	"database/sql"
)

// User structure in database
type User struct {
	ID        string `db:"id" json:"id"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Email     string `db:"email" json:"email"`
	Birthdate string `db:"birthdate" json:"birthdate"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}

// Token structure in database
type Token struct {
	ID        string `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"userId"`
	Type      string `db:"type" json:"type"`
	IP        string `db:"ip" json:"ip"`
	UserAgent string `db:"user_agent" json:"userAgent"`
	Read      int    `db:"read" json:"read"`
	Write     int    `db:"write" json:"write"`
	Manage    int    `db:"manage" json:"manage"`
	Admin     int    `db:"admin" json:"admin"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
	ExpireAt  string `db:"expire_at" json:"expireAt"`
}

// Gateway structure in database
type Gateway struct {
	ID        string         `db:"id" json:"id"`
	HomeID    sql.NullString `db:"home_id" json:"homeId"`
	Name      sql.NullString `db:"name" json:"name"`
	Model     string         `db:"model" json:"model"`
	CreatedAt string         `db:"created_at" json:"createdAt"`
	UpdatedAt string         `db:"updated_at" json:"updatedAt"`
	CreatorID sql.NullString `db:"creator_id" json:"creatorId"`
}

// Home structure in database
type Home struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Address   string `db:"address" json:"address"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
	CreatorID string `db:"creator_id" json:"creatorId"`
}

// Room structure in database
type Room struct {
	ID        string         `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Icon      sql.NullString `db:"icon" json:"icon"`
	HomeID    string         `db:"home_id" json:"homeId"`
	CreatedAt string         `db:"created_at" json:"createdAt"`
	UpdatedAt string         `db:"updated_at" json:"updatedAt"`
	CreatorID string         `db:"creator_id" json:"creatorId"`
}

// Device structure in database
type Device struct {
	ID           string         `db:"id" json:"id"`
	GatewayID    string         `db:"gateway_id" json:"gatewayId"`
	Name         string         `db:"name" json:"name"`
	Icon         sql.NullString `db:"icon" json:"icon"`
	PhysicalID   string         `db:"physical_id" json:"physicalId"`
	PhysicalName string         `db:"physical_name" json:"physicalName"`
	Config       string         `db:"config" json:"config"`
	Plugin       string         `db:"plugin" json:"plugin"`
	RoomID       string         `db:"room_id" json:"roomId"`
	CreatedAt    string         `db:"created_at" json:"createdAt"`
	UpdatedAt    string         `db:"updated_at" json:"updatedAt"`
	CreatorID    string         `db:"creator_id" json:"creatorId"`
}

// Permission structure in database
type Permission struct {
	ID        string `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"userId"`
	Type      string `db:"type" json:"type"` //home, room, device
	TypeID    string `db:"type_id" json:"typeId"`
	Read      int    `db:"read" json:"read"`
	Write     int    `db:"write" json:"write"`
	Manage    int    `db:"manage" json:"manage"`
	Admin     int    `db:"admin" json:"admin"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}

// Automation struct in database
type Automation struct {
	ID              string   `db:"id" json:"id"`
	HomeID          string   `db:"home_id" json:"homeId"`
	Name            string   `db:"name" json:"name"`
	Trigger         []string `db:"trigger" json:"trigger"`
	TriggerKey      []string `db:"trigger_key" json:"triggerKey"`
	TriggerValue    []string `db:"trigger_value" json:"triggerValue"`
	TriggerOperator []string `db:"trigger_operator" json:"triggerOperator"`
	Action          []string `db:"action" json:"action"`
	ActionCall      []string `db:"action_call" json:"actionCall"`
	ActionValue     []string `db:"action_value" json:"actionValue"`
	Status          bool     `db:"status" json:"status"`
	CreatedAt       string   `db:"created_at" json:"createdAt"`
	UpdatedAt       string   `db:"updated_at" json:"updatedAt"`
	CreatorID       string   `db:"creator_id" json:"creatorId"`
}

// Datas struct in database
type Datas struct {
	ID        string  `db:"id" json:"id"`
	DeviceID  string  `db:"device_id" json:"deviceId"`
	Field     string  `db:"field" json:"field"`
	ValueNbr  float64 `db:"value_nbr" json:"valueNbr"`
	ValueStr  string  `db:"value_str" json:"valueStr"`
	ValueBool bool    `db:"value_bool" json:"valueBool"`
	CreatedAt string  `db:"created_at" json:"createdAt"`
}
