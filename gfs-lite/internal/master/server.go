package master

import (
	"context"
	"log"
	"sync"

	gfsv1 "github.com/dsbudziwojski/gfs-lite/gen/gfs/v1"
	"github.com/google/uuid"
)

func uuidMaker() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return "cs-" + id.String()
}

type Server struct {
	gfsv1.UnimplementedMasterServiceServer
	mu           sync.Mutex
	chunkServers map[string]string
}

func NewServer() *Server {
	return &Server{
		chunkServers: make(map[string]string),
	}
}

func (s *Server) RegisterChunkServer(ctx context.Context, in *gfsv1.RegisterChunkServerRequest) (*gfsv1.RegisterChunkServerResponse, error) {
	ip := in.GetAddress()
	id := uuidMaker()
	s.mu.Lock()
	defer s.mu.Unlock()

	s.chunkServers[id] = ip
	return &gfsv1.RegisterChunkServerResponse{
		Accepted: true,
		ServerId: id,
	}, nil
}

func (s *Server) GetClusterStatus(ctx context.Context, in *gfsv1.GetClusterStatusRequest) (*gfsv1.GetClusterStatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := &gfsv1.GetClusterStatusResponse{
		ChunkServers: make([]*gfsv1.ChunkServerInfo, 0, len(s.chunkServers)),
	}
	for id, address := range s.chunkServers {
		out.ChunkServers = append(out.ChunkServers, &gfsv1.ChunkServerInfo{
			ServerId: id,
			Address:  address,
		})
	}
	return out, nil
}
