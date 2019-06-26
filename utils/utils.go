package utils

import (
	cryptorand "crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// NewULID create an ulid
func NewULID() ulid.ULID {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), cryptorand.Reader)
	return id
}

// Check check error
func Check(e error, typ string) {
	if e != nil {
		panic(e)
	}
}
