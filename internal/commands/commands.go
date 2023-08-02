package commands

import (
	"strings"
)

func Handle(command string) string {
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
	return "createRoom"
}

func joinRoom(roomId string) string {
	return "joinRoom"
}

func deleteRoom(roomId string) string {
	return "deleteRoom"
}
