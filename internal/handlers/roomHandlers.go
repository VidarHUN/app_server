package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VidarHUN/app_server/internal/config"
	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/utils"
	"github.com/gorilla/websocket"
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

func RoomPost(w http.ResponseWriter, r *http.Request, rooms *[]db.Room) {
	// Create a new struct to hold the request body.
	room := db.Room{Id: utils.GenerateRandomID(5)}
	user := db.User{}

	// Unmarshal the request body into the struct.
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}

	room.Users = append(room.Users, user)
	*rooms = append(*rooms, room)

	b, err := json.Marshal(room)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(b)
}

func CreateRoom(message map[string]interface{}, rooms *[]db.Room, conn *websocket.Conn, quicrq config.QuicrqServer) string {
	// Create a new struct to hold the request body.
	room := db.Room{Id: utils.GenerateRandomID(5)}
	user := db.User{
		Id:   message["data"].(map[string]interface{})["id"].(string),
		Conn: conn,
	}

	room.Users = append(room.Users, user)
	room.Server = quicrq
	*rooms = append(*rooms, room)

	ret := utils.Message[db.Room]{
		Command: "createRoom",
		Data:    room,
	}

	return utils.ToJson(ret)
}

func JoinRoom(message map[string]interface{}, rooms *[]db.Room, conn *websocket.Conn) string {
	var room db.Room
	for _, r := range *rooms {
		if r.Id == message["data"].(map[string]interface{})["id"] {
			user := db.User{
				Id:   message["data"].(map[string]interface{})["data"].(map[string]interface{})["id"].(string),
				Conn: conn,
			}
			go notifyUser(r, user)
			r.Users = append(r.Users, user)
			room = r
		}
	}

	ret := utils.Message[db.Room]{
		Command: "joinRoom",
		Data:    room,
	}

	return utils.ToJson(ret)
}

func DeleteRoom(message map[string]interface{}, rooms *[]db.Room) string {
	for i, r := range *rooms {
		if r.Id == message["data"].(map[string]interface{})["id"] {
			*rooms = append((*rooms)[:i], (*rooms)[i+1:]...)
			break
		}
	}
	return utils.ToJson("Room deleted")
}

func notifyUser(room db.Room, user db.User) {
	users := room.Users
	room.Users = []db.User{user}
	ret := utils.Message[db.Room]{
		Command: "newUser",
		Data:    room,
	}
	retJson := utils.ToJson(ret)
	for _, u := range users {
		err := u.Conn.WriteMessage(websocket.TextMessage, []byte(retJson))
		if err != nil {
			fmt.Println(err)
		}
	}
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
