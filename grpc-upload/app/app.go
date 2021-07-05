package app

import (
	"context"
	"grpc-upload/grpc"

	gogrpc "google.golang.org/grpc"
)

type App struct {
	grpcClient *grpc.Client
}

func New(grpcPort string) (*App, error) {
	conn, err := gogrpc.Dial(grpcPort, gogrpc.WithBlock(), gogrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &App{
		grpcClient: grpc.New(conn),
	}, nil
}

func (a *App) UploadImage(filePath string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return a.grpcClient.UploadImage(ctx, filePath)
}
