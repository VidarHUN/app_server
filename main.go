package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "GET /room request handled\n")
		case http.MethodPost:
			fmt.Fprintf(w, "POST /room request handled\n")
		case http.MethodPatch:
			fmt.Fprintf(w, "PATCH /room request handled\n")
		case http.MethodDelete:
			fmt.Fprintf(w, "DELETE /room request handled\n")
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/room/participants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "GET /room/participants request handled\n")
		case http.MethodPost:
			fmt.Fprintf(w, "POST /room/participants request handled\n")
		case http.MethodPatch:
			fmt.Fprintf(w, "PATCH /room/participants request handled\n")
		case http.MethodDelete:
			fmt.Fprintf(w, "DELETE /room/participants request handled\n")
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func main() {
	handler := setupHandler()
	quicConf := &quic.Config{}

	var err error
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
