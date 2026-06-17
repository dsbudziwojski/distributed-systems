package chunkserver

import (
	"context"
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	gfsv1 "github.com/dsbudziwojski/gfs-lite/gen/gfs/v1"
)

type chunkServer struct {
	gfsv1.UnimplementedChunkServiceServer
	mu        sync.Mutex
	basePath  string
	chunkSize int64
}

func NewServer(chunkSize int64) *chunkServer {
	err := os.MkdirAll("data/chunkserver", 0755)
	if err != nil {
		log.Fatal(err)
	}
	server := &chunkServer{
		basePath:  "data/chunkserver",
		chunkSize: chunkSize,
	}
	return server
}

func (c *chunkServer) WriteChunk(ctx context.Context, in *gfsv1.WriteChunkRequest) (*gfsv1.WriteChunkResponse, error) {
	chunkHandle := in.ChunkHandle
	if !c.validHandle(chunkHandle) {
		return nil, errors.New("invalid chunk handle")
	}

	offset, byteData := in.Offset, []byte(in.Data)
	fullPath := path.Join(c.basePath, chunkHandle)
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.New("could not open file")
	}
	defer f.Close()

	if offset+int64(len(byteData)) > c.chunkSize {
		fileInfo, err := f.Stat()
		if err != nil {
			return nil, errors.New("could not stat file")
		}
		spaceLeft := c.chunkSize - fileInfo.Size()
		emptyData := make([]byte, spaceLeft)
		f.WriteAt(emptyData, fileInfo.Size())
		return nil, errors.New("chunk too big")
	}
	_, err = f.WriteAt(byteData, offset)
	if err != nil {
		return nil, errors.New("could not write to file")
	}
	out := &gfsv1.WriteChunkResponse{}
	return out, nil
}

func (c *chunkServer) ReadChunk(ctx context.Context, in *gfsv1.ReadChunkRequest) (*gfsv1.ReadChunkResponse, error) {
	chunkHandle := in.ChunkHandle
	if !c.validHandle(chunkHandle) {
		return nil, errors.New("invalid chunk handle")
	}
	offset := in.Offset
	length := in.Length

	fullPath := path.Join(c.basePath, chunkHandle)
	f, err := os.OpenFile(fullPath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, errors.New("could not open file")
	}
	defer f.Close()

	if offset > c.chunkSize || offset < 0 {
		return nil, errors.New("offset out of range")
	}
	if offset+length > c.chunkSize {
		return nil, errors.New("requested chunk too big")
	}

	buffer := make([]byte, length)
	_, err = f.ReadAt(buffer, offset)
	if err != nil {
		return nil, errors.New("could not read from file")
	}
	out := &gfsv1.ReadChunkResponse{
		Data: string(buffer),
	}
	return out, nil
}

func (c *chunkServer) validHandle(handle string) bool {
	return strings.HasPrefix(handle, "chunk-")
}
