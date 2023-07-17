package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {

	// pool, err := x509.SystemCertPool()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var qconf quic.Config
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            nil,
			InsecureSkipVerify: true,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	fmt.Printf("GET %s", "https://localhost:8443/room")
	rsp, err := hclient.Get("https://localhost:8443/room")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got response for %s: %#v", "https://localhost:8443/room", rsp)

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response Body:")
	fmt.Printf("%s", body.Bytes())
}
