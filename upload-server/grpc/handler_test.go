package grpc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"upload-server/pb"

	"github.com/stretchr/testify/require"
	gogrpc "google.golang.org/grpc"
)

func TestHandler_UploadImage(t *testing.T) {
	httpTarget := "localhost:8080"
	grpcServer := startTestImageServer(httpTarget)
	imageClient := newTestImageClient(t, grpcServer.lis.Addr().String())

	file, err := os.Open("data/cute_cat.jpg")
	require.NoError(t, err)
	fileName := filepath.Base(file.Name())
	defer func() {
		err = file.Close()
		require.NoError(t, err)
		err = os.Remove(fileName)
		require.NoError(t, err)
		grpcServer.GracefulStop()
	}()

	stream, err := imageClient.UploadImage(context.Background())
	require.NoError(t, err)

	req := &pb.UploadRequest{
		Data: &pb.UploadRequest_FileName{
			FileName: "cute_cat.jpg",
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		size += n

		req := &pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("http://%s/%s", httpTarget, fileName), res.GetUrl())
	require.FileExists(t, fileName)
}

func startTestImageServer(httpTarget string) *Server {
	server := New(":0", httpTarget)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	return server
}

func newTestImageClient(t *testing.T, addr string) pb.ImageServiceClient {
	conn, err := gogrpc.Dial(addr, gogrpc.WithBlock(), gogrpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewImageServiceClient(conn)
}
