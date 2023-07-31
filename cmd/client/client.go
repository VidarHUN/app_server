package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/utils"

	"github.com/gorilla/websocket"
)

// var upgrader = websocket.Upgrader{} // use default options

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
	// quicConf := &quic.Config{
	// 	KeepAlivePeriod: 60,
	// }
	// roundTripper := &http3.RoundTripper{
	// 	TLSClientConfig: &tls.Config{
	// 		RootCAs:            nil,
	// 		InsecureSkipVerify: true,
	// 	},
	// 	QuicConfig: quicConf,
	// }
	// defer roundTripper.Close()
	// hclient := &http.Client{
	// 	Transport: roundTripper,
	// }

	// roomPost(hclient)
	// Create a new WebSocket client.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("connecting to ws://localhost:8080/room")

	client, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/room", nil)
	if err != nil {
		fmt.Println("dial:", err)
	}

	go func() {
		<-c
		// Run Cleanup
		err := client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			fmt.Println("write close:", err)
		}
		client.Close()
	}()

	// Create a message.
	message := map[string]string{
		"command": "createRoom",
		"userId":  utils.GenerateRandomID(5),
	}

	// Marshal the message to JSON.
	var jsonMessage []byte
	jsonMessage, err = json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		fmt.Println("write:", err)
	}

	for {
		var msg []byte
		_, msg, err = client.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Println("recv: %s", string(msg))
	}
}
