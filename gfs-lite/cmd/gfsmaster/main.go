package main

import (
	"log"
	"net"
	"os"

	gfsv1 "github.com/dsbudziwojski/gfs-lite/gen/gfs/v1"
	"github.com/dsbudziwojski/gfs-lite/internal/master"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	masterServer := master.NewServer()
	gfsv1.RegisterMasterServiceServer(server, masterServer)

	log.Println("gfsmaster listening on " + port)
	log.Fatal(server.Serve(lis))
}
