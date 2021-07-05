package grpc

import (
	"bufio"
	"context"
	"fmt"
	"grpc-upload/pb"
	"grpc-upload/util"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	gogrpc "google.golang.org/grpc"
)

type Client struct {
	lock   sync.RWMutex
	client pb.ImageServiceClient
}

func New(conn *gogrpc.ClientConn) *Client {
	return &Client{client: pb.NewImageServiceClient(conn)}
}

func (c *Client) UploadImage(ctx context.Context, filePath string) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if !util.FileExists(filePath) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stream, err := c.client.UploadImage(ctx)
	if err != nil {
		return "", err
	}

	req := NewUploadRequest(filepath.Base(file.Name()))

	err = stream.Send(req)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 512)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		req = &pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		if err = stream.Send(req); err != nil {
			return "", stream.RecvMsg(nil)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.Url, nil
}
