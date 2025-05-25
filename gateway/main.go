package main

import (
	"log"
	"net/http"

	common "github.com/ysle0/omsv2/common"
	pb "github.com/ysle0/omsv2/common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr         = common.EnvString("HTTP_ADDR", ":8000")
	orderServiceAddr = "localhost:3000"
)

func main() {
	conn, err := grpc.Dial(
		orderServiceAddr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
	if err != nil {
		log.Fatalf("failed to dial server: %v\n", err)
	}
	defer conn.Close()

	log.Println("Dialing orders service at ", orderServiceAddr)

	if err != nil {
		log.Fatal(err)
	}

	c := pb.NewOrderServiceClient(conn)
	mux := http.NewServeMux()

	httpHandler := NewHttpHandler(c)
	httpHandler.registerRoutes(mux)

	log.Printf("HTTP server listening on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal(err)
	}
}
