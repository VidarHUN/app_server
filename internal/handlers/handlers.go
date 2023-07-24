package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
