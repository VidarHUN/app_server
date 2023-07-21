package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/redis/go-redis/v9"
)

type Room struct {
	Id    string `json:"id"`
	Users []User `json:"users"`
}

type User struct {
	Id string `json:"id"`
}

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

func handleRoomGet(w http.ResponseWriter) {
	// Dummy data
	room := &Room{Id: "DummyRoom"}
	room.Users = append(room.Users, User{Id: "DummyUser"})
	byteRoom, err := json.Marshal(room)
	if err != nil {
		w.Write(errorToBytes(err))
	}
	w.Write(byteRoom)
}

func handleRoomPost(r *http.Request) {
	fmt.Println(r.Body)
}

func handleRoomPatch(w http.ResponseWriter, id string) {
	// Dummy data
	room := &Room{Id: id}
	room.Users = append(room.Users, User{Id: "DummyUser"})
	byteRoom, err := json.Marshal(room)
	if err != nil {
		w.Write(errorToBytes(err))
	}
	w.Write(byteRoom)
}

func handleRoomDelete(w http.ResponseWriter, id string) {
	w.Write([]byte(id + "room deleted"))
}

func handleUserGet(w http.ResponseWriter) {
	// Dummy data
	user := &User{Id: "DummyUser"}
	byteUser, err := json.Marshal(user)
	if err != nil {
		w.Write(errorToBytes(err))
	}
	w.Write(byteUser)
}

func handleUserPost(r *http.Request) {
	fmt.Println(r.Body)
}

func handleUserPatch(w http.ResponseWriter, id string) {
	// Dummy data
	user := &User{Id: id}
	byteUser, err := json.Marshal(user)
	if err != nil {
		w.Write(errorToBytes(err))
	}
	w.Write(byteUser)
}

func handleUserDelete(w http.ResponseWriter, id string) {
	w.Write([]byte(id + "user deleted"))
}

func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleRoomGet(w)
		case http.MethodPost:
			handleRoomPost(r)
		case http.MethodPatch:
			id := r.URL.Query().Get("id")
			handleRoomPatch(w, id)
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			handleRoomDelete(w, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/room/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleUserGet(w)
		case http.MethodPost:
			handleUserPost(r)
		case http.MethodPatch:
			id := r.URL.Query().Get("id")
			handleUserPatch(w, id)
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			handleUserDelete(w, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func main() {
	var err error
	var client *redis.Client

	client, err = NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		fmt.Println(err)
	}
	client.Close()

	handler := setupHandler()
	quicConf := &quic.Config{}

	server := http3.Server{
		Handler:    handler,
		Addr:       "localhost:8443",
		QuicConfig: quicConf,
	}
	fmt.Printf("Running on port 8443\n")
	err = server.ListenAndServeTLS("localhost.crt", "localhost.key")
	if err != nil {
		fmt.Println(err)
	}
}
