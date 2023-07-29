package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/utils"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func roomPost(hclient *http.Client) {
	user := db.User{Id: utils.GenerateRandomID(5)}

	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	rsp, err := hclient.Post("https://localhost:8443/room", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Body:")
	fmt.Println(string(body.Bytes()))
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
