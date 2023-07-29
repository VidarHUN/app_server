package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
)

func ErrorToBytes(err error) []byte {
	if err == nil {
		return []byte("")
	}

	errBytes, err := json.Marshal(err)
	if err != nil {
		return []byte("")
	}

	return errBytes
}

func GenerateRandomID(len int) string {
	id := make([]byte, len)
	rand.Read(id)

	// Convert the ID to a string.
	idStr := fmt.Sprintf("%x", id)

	return idStr
}
