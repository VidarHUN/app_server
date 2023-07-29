package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/utils"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var rooms []db.Room

func roomPost(hclient *http.Client) {
	user := db.User{Id: utils.GenerateRandomID(5)}

	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	rsp, err := hclient.Post("https://localhost:8443/room", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)
	}

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var room = db.Room{}
	// Unmarshal the response body into the struct.
	err = json.NewDecoder(body).Decode(&room)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Response Body:")
	fmt.Println(room)

	rooms = append(rooms, room)
}

func main() {
	var qconf quic.Config
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            nil,
			InsecureSkipVerify: true,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	roomPost(hclient)
}
