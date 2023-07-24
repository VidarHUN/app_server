package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/utils"
)

func RoomGet(w http.ResponseWriter) {
	// Dummy data
	room := &db.Room{Id: "DummyRoom"}
	room.Users = append(room.Users, db.User{Id: "DummyUser"})
	byteRoom, err := json.Marshal(room)
	if err != nil {
		w.Write(utils.ErrorToBytes(err))
	}
	w.Write(byteRoom)
}

func RoomPost(r *http.Request) {
	fmt.Println(r.Body)
}

func RoomPatch(w http.ResponseWriter, id string) {
	// Dummy data
	room := &db.Room{Id: id}
	room.Users = append(room.Users, db.User{Id: "DummyUser"})
	byteRoom, err := json.Marshal(room)
	if err != nil {
		w.Write(utils.ErrorToBytes(err))
	}
	w.Write(byteRoom)
}

func RoomDelete(w http.ResponseWriter, id string) {
	w.Write([]byte(id + "room deleted"))
}
