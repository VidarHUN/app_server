package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/utils"
)

func UserGet(w http.ResponseWriter) {
	// Dummy data
	user := &db.User{Id: "DummyUser"}
	byteUser, err := json.Marshal(user)
	if err != nil {
		w.Write(utils.ErrorToBytes(err))
	}
	w.Write(byteUser)
}

func UserPost(r *http.Request) {
	fmt.Println(r.Body)
}

func UserPatch(w http.ResponseWriter, id string) {
	// Dummy data
	user := &db.User{Id: id}
	byteUser, err := json.Marshal(user)
	if err != nil {
		w.Write(utils.ErrorToBytes(err))
	}
	w.Write(byteUser)
}

func UserDelete(w http.ResponseWriter, id string) {
	w.Write([]byte(id + "user deleted"))
}
