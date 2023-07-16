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
			fmt.Fprintf(w, "GET request handled\n")
		case http.MethodPost:
			fmt.Fprintf(w, "POST request handled\n")
		case http.MethodPatch:
			fmt.Fprintf(w, "PATCH request handled\n")
		case http.MethodDelete:
			fmt.Fprintf(w, "DELETE request handled\n")
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func main() {
	// Start a goroutine which serve a HTTP server on localhost:6060
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:8080", nil))
	// }()

	handler := setupHandler()
	quicConf := &quic.Config{}

	var err error
	server := http3.Server{
		Handler:    handler,
		Addr:       "localhost:8443",
		QuicConfig: quicConf,
	}
	fmt.Printf("Running on port 8443\n")
	err = server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		fmt.Println(err)
	}
}
