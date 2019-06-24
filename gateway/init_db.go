package gateway

import (
	"database/sql"
	"log"
)

// DB define the database object
var DB *sql.DB

// InitDB check and create tables
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./casa.db")
	if err != nil {
		log.Panic(err)
	}

	statement, err := DB.Prepare("CREATE TABLE IF NOT EXISTS user (id BLOB PRIMARY KEY, firstname TEXT, lastname TEXT, email TEXT, password TEXT, birthdate TEXT, created_at TEXT)")
	if err != nil {
		log.Panic(err)
	}
	statement.Exec()
	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS token (id BLOB PRIMARY KEY, user_id BLOB, type TEXT, ip TEXT, browser TEXT, os TEXT, read INTEGER, write INTEGER, manage INTEGER, created_at TEXT, expire_at TEXT)")
	if err != nil {
		log.Panic(err)
	}
	statement.Exec()
	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS home (id BLOB PRIMARY KEY, name TEXT, address TEXT, created_at TEXT, creator_id BLOB)")
	if err != nil {
		log.Panic(err)
	}
	statement.Exec()
	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS room (id BLOB PRIMARY KEY, name TEXT, home_id BLOB, created_at TEXT, creator_id BLOB)")
	if err != nil {
		log.Panic(err)
	}
	statement.Exec()
	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS device (id BLOB PRIMARY KEY, name TEXT, physical_id TEXT, room_id BLOB, created_at TEXT, creator_id BLOB)")
	if err != nil {
		log.Panic(err)
	}
	statement.Exec()
	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS permission (id BLOB PRIMARY KEY, user_id BLOB, type TEXT, type_id BLOB, read INTEGER, write INTEGER, manage INTEGER, admin INTEGER, updated_at TEXT)")
	if err != nil {
		log.Panic(err)
	}
	statement.Exec()
}
