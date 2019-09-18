package utils

import (
	cryptorand "crypto/rand"
	"io/ioutil"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
)

// NewULID create an ulid
func NewULID() ulid.ULID {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), cryptorand.Reader)
	return id
}

// GetIDFile get ID from config file
func GetIDFile() string {
	var id string

	file, err := os.OpenFile(".casa", os.O_APPEND, 0644)
	if err != nil {
		id = string(NewULID().String())
		err = ioutil.WriteFile(".casa", []byte(id), 0644)
		Check(err, "error")
	} else {
		data := make([]byte, 100)
		count, err := file.Read(data)
		Check(err, "error")
		id = string(data[:count])
	}
	return id
}

// Check check error
func Check(e error, typ string) {
	if e != nil {
		panic(e)
	}
}
