package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/redis/go-redis/v9"

	"github.com/VidarHUN/app_server/internal/db"
	"github.com/VidarHUN/app_server/internal/handlers"
)

func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.RoomGet(w)
		case http.MethodPost:
			handlers.RoomPost(r)
		case http.MethodPatch:
			id := r.URL.Query().Get("id")
			handlers.RoomPatch(w, id)
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			handlers.RoomDelete(w, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/room/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.UserGet(w)
		case http.MethodPost:
			handlers.UserPost(r)
		case http.MethodPatch:
			id := r.URL.Query().Get("id")
			handlers.UserPatch(w, id)
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			handlers.UserDelete(w, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func main() {
	var err error
	var client *redis.Client

	client, err = db.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		fmt.Println(err)
	}
	client.Close()

	handler := setupHandler()
	quicConf := &quic.Config{}

	server := http3.Server{
		Handler:    handler,
		Addr:       "localhost:8443",
		QuicConfig: quicConf,
	}
	fmt.Printf("Running on port 8443\n")
	err = server.ListenAndServeTLS("localhost.crt", "localhost.key")
	if err != nil {
		fmt.Println(err)
	}
}
