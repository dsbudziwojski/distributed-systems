package main

import (
	"context"
	"crypto/rand"
	"flag"
	"log"
	"math/big"
	"net"
	"time"

	gfsv1 "github.com/dsbudziwojski/gfs-lite/gen/gfs/v1"
	"github.com/dsbudziwojski/gfs-lite/internal/chunkserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func nrand() string {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return string(rune(x))
}

func main() {
	addr := flag.String("addr", ":8100", "chunk server listen address")
	master := flag.String("master", ":8000", "master listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	chunkServer := chunkserver.NewServer(100)
	gfsv1.RegisterChunkServiceServer(server, chunkServer)
	log.Println("gfschunkserver listening on " + *addr)

	conn, err := grpc.NewClient(
		*master,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := gfsv1.NewMasterServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := client.RegisterChunkServer(ctx, &gfsv1.RegisterChunkServerRequest{
		Address:  *addr,
		ServerId: nrand(),
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registration accepted: %v", out.Accepted)
	log.Fatal(server.Serve(lis))
}
