package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/VidarHUN/app_server/internal/commands"
	"github.com/VidarHUN/app_server/internal/utils"

	"github.com/gorilla/websocket"

	"github.com/spf13/cobra"
)

var SERVER = "localhost:8080"
var PATH = "/room"
var SRC = "samples/new_video1_source.bin"
var QUICRQ = "./quicrq_app"
var USERID = utils.GenerateRandomID(5)
var in = bufio.NewReader(os.Stdin)

var ERRORS = []string{
	"Too much or less arguments",
	"Command not found",
}

func readMsg(c *websocket.Conn) {
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("ReadMessage() error:", err)
		return
	}
	log.Printf("Received: %s", message)
	commands.QuicrqProcess(message, SRC, QUICRQ, USERID)
	go readMsg(c)
}

func main() {
	input := make(chan string, 1)
	rootCmd := &cobra.Command{
		Use:   "shell",
		Short: "An interactive shell",
		// TODO: Make it scriptable from file
		Run: func(cmd *cobra.Command, args []string) {
			for {
				read_line, err := in.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					break
				}
				line := strings.TrimSuffix(read_line, "\n")
				msg := commands.Generate(line, USERID)

				if utils.Contains(ERRORS, msg) {
					log.Println(msg)
				} else {
					input <- msg
				}
			}
		},
	}
	go func() {
		err := rootCmd.Execute()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	if len(os.Args) > 1 {
		SERVER = os.Args[1]
	}

	fmt.Println("Connecting to:", SERVER, "at", PATH)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	URL := url.URL{Scheme: "ws", Host: SERVER, Path: PATH}
	c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	fmt.Println(c.LocalAddr())
	defer c.Close()
	done := make(chan struct{})

	for {
		select {
		case <-done:
			return
		case t := <-input:
			err := c.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Println("Write error:", err)
				return
			}
			go readMsg(c)
		case <-interrupt:
			log.Println("Caught interrupt signal - quitting!")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				log.Println("Write close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return
		}
	}
}
