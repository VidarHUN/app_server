package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/VidarHUN/app_server/internal/config"
	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/handlers"
	"github.com/VidarHUN/app_server/internal/utils"
	"github.com/gorilla/websocket"
)

type message[T any] struct {
	Command string `json:"command"`
	Data    T      `json:"data"`
}

type mixedData[T any] struct {
	Id   string `json:"id"`
	Data T      `json:"data"`
}

type id struct {
	Id string `json:"id"`
}

func Generate(command string) string {
	parsed_command := strings.Split(command, " ")
	switch parsed_command[0] {
	case "createRoom":
		return createRoom()
	case "joinRoom":
		if len(parsed_command) != 2 {
			return "Too much or less arguments"
		}
		return joinRoom(parsed_command[1])
	case "deleteRoom":
		if len(parsed_command) != 2 {
			return "Too much or less arguments"
		}
		return deleteRoom(parsed_command[1])
	default:
		return "Command not found"
	}
}

func createRoom() string {
	msg := message[db.User]{
		Command: "createRoom",
		Data:    db.User{Id: utils.GenerateRandomID(5)},
	}
	return utils.ToJson(msg)
}

func joinRoom(roomId string) string {
	msg := message[mixedData[db.User]]{
		Command: "joinRoom",
		Data: mixedData[db.User]{
			Id:   roomId,
			Data: db.User{Id: utils.GenerateRandomID(5)},
		},
	}
	return utils.ToJson(msg)
}

func deleteRoom(roomId string) string {
	msg := message[id]{
		Command: "deleteRoom",
		Data:    id{Id: roomId},
	}
	return utils.ToJson(msg)
}

func Process(msg []byte, rooms *[]db.Room, conn *websocket.Conn, quicrq config.QuicrqServer) string {
	// Unmarshal the message into a map.
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println(err)
	}

	// Access the fields of the message by their names.
	command := message["command"]

	switch command {
	case "createRoom":
		return handlers.CreateRoom(message, rooms, conn, quicrq)
	case "joinRoom":
		return handlers.JoinRoom(message, rooms, conn)
	case "deleteRoom":
		return handlers.DeleteRoom(message, rooms)
	default:
		return "Unkown command: " + command.(string)
	}
}

func QuicrqProcess(msg []byte, src string) {
	// Unmarshal the message into a map.
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println(err)
	}

	// Access the fields of the message by their names.
	command := message["command"]
	var q_commands []string

	switch command {
	case "createRoom":
		q_commands = append(q_commands, quicrqPost(message, src))
	case "joinRoom":
		q_commands = append(q_commands, quicrqPost(message, src))
		q_commands = append(q_commands, quicrqGet(message, src))
	case "newUser":
		q_commands = append(q_commands, quicrqGet(message, src))
	}

	for _, cmd := range q_commands {
		go startCmd(cmd)
	}
}

func quicrqPost(message map[string]interface{}, src string) string {
	url := fmt.Sprintf("%s_%s", message["id"], message["data"].(map[string]interface{})["id"])
	return fmt.Sprintf("post:%s:%s", url, src)
}

func quicrqGet(message map[string]interface{}, src string) string {
	url := fmt.Sprintf("%s_%s", message["id"], message["data"].(map[string]interface{})["id"])
	newSrc := fmt.Sprintf("%s_%s", message["data"].(map[string]interface{})["id"], src)
	return fmt.Sprintf("get:%s:%s", url, newSrc)
}

func startCmd(command string) {
	cmd := exec.Command(command)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
