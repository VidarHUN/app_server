package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
)

type Message[T any] struct {
	Command string `json:"command"`
	Data    T      `json:"data"`
}

type MixedData[T any] struct {
	Id   string `json:"id"`
	Data T      `json:"data"`
}

type Id struct {
	Id string `json:"id"`
}

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

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ToJson(msg interface{}) string {
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
