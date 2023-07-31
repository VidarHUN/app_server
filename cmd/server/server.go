package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/VidarHUN/app_server/internal/config"
	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/handlers"
)

var upgrader = websocket.Upgrader{} // use default options

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

			// Unmarshal the message into a map.
			var message map[string]string
			err = json.Unmarshal(msg, &message)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Access the fields of the message by their names.
			command := message["command"]

			switch command {
			case "createRoom":
				rsp := handlers.CreateRoom(message, &rooms)
				conn.WriteMessage(msgType, rsp)
			default:
				fmt.Println("Unknown command:", command)
				conn.WriteMessage(msgType, []byte("Unkown command: "+command))
			}
		}

		// switch r.Method {
		// case http.MethodGet:
		// 	handlers.RoomGet(w)
		// case http.MethodPost:
		// 	handlers.RoomPost(w, r, &rooms)
		// case http.MethodPatch:
		// 	id := r.URL.Query().Get("id")
		// 	handlers.RoomPatch(w, id)
		// case http.MethodDelete:
		// 	id := r.URL.Query().Get("id")
		// 	handlers.RoomDelete(w, id)
		// default:
		// 	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		// }
	})

	// mux.HandleFunc("/room/users", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		handlers.UserGet(w)
	// 	case http.MethodPost:
	// 		handlers.UserPost(r)
	// 	case http.MethodPatch:
	// 		id := r.URL.Query().Get("id")
	// 		handlers.UserPatch(w, id)
	// 	case http.MethodDelete:
	// 		id := r.URL.Query().Get("id")
	// 		handlers.UserDelete(w, id)
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

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
