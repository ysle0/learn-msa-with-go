package main

import (
	"log"
	"net/http"

	common "github.com/ysle0/omsv2/common"
)

var (
	httpAddr = common.EnvString("HTTP_ADDR", ":8000")
)

func main() {
	mux := http.NewServeMux()

	httpHandler := NewHttpHandler()
	httpHandler.registerRoutes(mux)

	log.Printf("HTTP server listening on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal(err)
	}
}
