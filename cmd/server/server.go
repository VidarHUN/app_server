package main

import (
	"flag"
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/VidarHUN/app_server/internal/commands"
	"github.com/VidarHUN/app_server/internal/config"
	"github.com/VidarHUN/app_server/internal/db"
)

var rooms []db.Room

func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection.
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Read messages from the client.
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(msg))
			retMsg := commands.Process(msg, &rooms)
			conn.WriteMessage(msgType, []byte(retMsg))
		}
	})

	return mux
}

func main() {
	var err error
	var client *redis.Client
	var path string

	flag.StringVar(&path, "path", ".", "Directory of config")
	flag.Parse()
	configuration := config.ReadConfig(path)

	redisAddress := configuration.Database.Address + ":" + string(configuration.Database.Port)
	client, err = db.NewRedisClient(redisAddress, "", 0)
	if err != nil {
		fmt.Println(err)
	}
	client.Close()

	handler := setupHandler()

	fmt.Println(http.ListenAndServe("localhost:8080", handler))

	// HTTP3
	// handler := setupHandler()
	// quicConf := &quic.Config{}

	// server := http3.Server{
	// 	Handler:    handler,
	// 	Addr:       "localhost:8443",
	// 	QuicConfig: quicConf,
	// }
	// fmt.Printf("Running on port 8443\n")
	// err = server.ListenAndServeTLS("localhost.crt", "localhost.key")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
