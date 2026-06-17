package main

import (
	"flag"
	"log"
	"net"

	gfsv1 "github.com/dsbudziwojski/gfs-lite/gen/gfs/v1"
	"github.com/dsbudziwojski/gfs-lite/internal/master"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", ":7000", "master listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	masterServer := master.NewServer()
	gfsv1.RegisterMasterServiceServer(server, masterServer)

	log.Println("gfsmaster listening on " + *addr)
	log.Fatal(server.Serve(lis))
}
