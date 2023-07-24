package utils

import "encoding/json"

func errorToBytes(err error) []byte {
	if err == nil {
		return []byte("")
	}

	errBytes, err := json.Marshal(err)
	if err != nil {
		return []byte("")
	}

	return errBytes
}
