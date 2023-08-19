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

func Generate(command string, userId string) string {
	parsed_command := strings.Split(command, " ")
	switch parsed_command[0] {
	case "createRoom":
		return createRoom(userId)
	case "joinRoom":
		if len(parsed_command) != 2 {
			return "Too much or less arguments"
		}
		return joinRoom(parsed_command[1], userId)
	case "deleteRoom":
		if len(parsed_command) != 2 {
			return "Too much or less arguments"
		}
		return deleteRoom(parsed_command[1])
	default:
		return "Command not found"
	}
}

func createRoom(userId string) string {
	msg := utils.Message[db.User]{
		Command: "createRoom",
		Data:    db.User{Id: userId},
	}
	return utils.ToJson(msg)
}

func joinRoom(roomId string, userId string) string {
	msg := utils.Message[utils.MixedData[db.User]]{
		Command: "joinRoom",
		Data: utils.MixedData[db.User]{
			Id:   roomId,
			Data: db.User{Id: userId},
		},
	}
	return utils.ToJson(msg)
}

func deleteRoom(roomId string) string {
	msg := utils.Message[utils.Id]{
		Command: "deleteRoom",
		Data:    utils.Id{Id: roomId},
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

func QuicrqProcess(msg []byte, src string, quicrq string, userId string) {
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
		q_commands = append(q_commands, quicrqPost(message, src, userId))
	case "joinRoom":
		q_commands = append(q_commands, quicrqPost(message, src, userId))
		q_commands = append(q_commands, quicrqGet(message, src, userId))
	case "newUser":
		q_commands = append(q_commands, quicrqGet(message, src, userId))
	}

	quicrq_server := message["data"].(map[string]interface{})["server"].(map[string]interface{})["address"]
	quicrq_port := message["data"].(map[string]interface{})["server"].(map[string]interface{})["port"]
	quicrq_port = int(quicrq_port.(float64))

	for _, cmd := range q_commands {
		newCmd := fmt.Sprintf("%s client %s d %d %s", quicrq, quicrq_server, quicrq_port, cmd)
		go startCmd(newCmd)
	}
}

func quicrqPost(message map[string]interface{}, src string, userId string) string {
	url := fmt.Sprintf("%s_%s", userId, message["data"].(map[string]interface{})["id"])
	return fmt.Sprintf("post:%s:%s", url, src)
}

func quicrqGet(message map[string]interface{}, src string, userId string) string {
	url := fmt.Sprintf("%s_%s", userId, message["data"].(map[string]interface{})["id"])
	newSrc := fmt.Sprintf("%s_%s", message["data"].(map[string]interface{})["id"], src)
	return fmt.Sprintf("get:%s:%s", url, newSrc)
}

func startCmd(command string) {
	commandList := strings.Split(command, " ")
	cmd := exec.Command(commandList[0])
	cmd.Args = commandList
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
