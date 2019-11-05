package gateway

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ItsJimi/casa-gateway/logger"
	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

// DB define the database object
var DB *sqlx.DB

// InitDB start the database to use it in server
func InitDB() {
	var err error
	DB, err = sqlx.Open("sqlite3", "./casa.db")
	if err != nil {
		contextLogger := logger.WithFields(logger.Fields{"code": "CGGDIDB001"})
		contextLogger.Panicf("%s", err.Error())
		return
	}

	file, err := ioutil.ReadFile("database.sql")
	if err != nil {
		contextLogger := logger.WithFields(logger.Fields{"code": "CGGDIDB002"})
		contextLogger.Panicf("%s", err.Error())
		return
	}

	_, err = DB.Exec(string(file))
	if err != nil {
		contextLogger := logger.WithFields(logger.Fields{"code": "CGGDIDB003"})
		contextLogger.Panicf("%s", err.Error())
		return
	}
}

type synced struct {
	Data struct {
		Home        Home
		Gateway     Gateway
		Users       []User
		Automations []Automation
		Devices     []Device
		Rooms       []Room
		Permissions []Permission
	} `json:"data"`
}

// SyncDB sync the DB with server's DB
func SyncDB() {
	var gateway Gateway
	err := DB.Get(&gateway, `SELECT * FROM gateways`)
	if err == nil && gateway.ID != "" {
		return
	}
	resp, err := http.Get("http://localhost:3000/v1/gateway/sync/" + utils.GetIDFile())
	if err != nil {
		contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB001"})
		contextLogger.Panicf("%s", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var dataMarshalised synced
	err = json.Unmarshal(body, &dataMarshalised)
	if err != nil {
		contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB002"})
		contextLogger.Panicf("%s", err.Error())
		return
	}

	_, err = DB.NamedExec("INSERT INTO gateways (id, home_id, name, model, created_at, updated_at, creator_id) VALUES (:id, :home_id, :name, :model, :created_at, :updated_at, :creator_id)", dataMarshalised.Data.Gateway)
	if err != nil {
		contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB003"})
		contextLogger.Panicf("%s", err.Error())
		return
	}

	_, err = DB.NamedExec("INSERT INTO homes (id, name, address, created_at, updated_at, creator_id) VALUES (:id, :name, :address, :created_at, :updated_at, :creator_id)", dataMarshalised.Data.Home)
	if err != nil {
		contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB004"})
		contextLogger.Panicf("%s", err.Error())
		return
	}

	for _, user := range dataMarshalised.Data.Users {
		_, err = DB.NamedExec("INSERT INTO users (id, firstname, lastname, email, birthdate, created_at, updated_at) VALUES (:id, :firstname, :lastname, :email, :birthdate, :created_at, :updated_at)", user)
		if err != nil {
			contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB005"})
			contextLogger.Panicf("%s", err.Error())
			return
		}
	}

	for _, automation := range dataMarshalised.Data.Automations {
		_, err := DB.Exec("INSERT INTO automations (id, home_id, name, trigger, trigger_key, trigger_value, trigger_operator, action, action_call, action_value, status, created_at, updated_at, creator_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
			automation.ID, automation.HomeID, automation.Name, pq.Array(automation.Trigger), pq.Array(automation.TriggerKey), pq.Array(automation.TriggerValue), pq.Array(automation.TriggerOperator), pq.Array(automation.Action), pq.Array(automation.ActionCall), pq.Array(automation.ActionValue), automation.Status, automation.CreatedAt, automation.UpdatedAt, automation.CreatorID)
		if err != nil {
			contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB006"})
			contextLogger.Panicf("%s", err.Error())
			return
		}
	}

	for _, device := range dataMarshalised.Data.Devices {
		_, err = DB.NamedExec("INSERT INTO devices (id, gateway_id, name, icon, physical_id, physical_name, config, plugin, room_id, created_at, updated_at, creator_id) VALUES (:id, :gateway_id, :name, :icon, :physical_id, :physical_name, :config, :plugin, :room_id, :created_at, :updated_at, :creator_id)", device)
		if err != nil {
			contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB007"})
			contextLogger.Panicf("%s", err.Error())
			return
		}
	}

	for _, room := range dataMarshalised.Data.Rooms {
		_, err = DB.NamedExec("INSERT INTO rooms (id, name, icon, home_id, created_at, updated_at, creator_id) VALUES (:id, :name, :icon, :home_id, :created_at, :updated_at, :creator_id)", room)
		if err != nil {
			contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB008"})
			contextLogger.Panicf("%s", err.Error())
			return
		}
	}

	for _, permission := range dataMarshalised.Data.Permissions {
		_, err = DB.NamedExec("INSERT INTO permissions (id, user_id, type, type_id, read, write, manage, admin, created_at, updated_at) VALUES (:id, :user_id, :type, :type_id, :read, :write, :manage, :admin, :created_at, :updated_at)", permission)
		if err != nil {
			contextLogger := logger.WithFields(logger.Fields{"code": "CGGDSDB009"})
			contextLogger.Panicf("%s", err.Error())
			return
		}
	}

	return
}
