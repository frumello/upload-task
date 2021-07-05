package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upload-server/grpc"
	"upload-server/web"

	"golang.org/x/sync/errgroup"
)

type App struct {
	grpcServer *grpc.Server
	httpServer *web.Server
}

func New(grpcTarget, httpTarget string) (*App, error) {
	return &App{
		grpcServer: grpc.New(grpcTarget, httpTarget),
		httpServer: web.New(httpTarget),
	}, nil
}

func (a *App) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	signal.Stop(signalChan)

	g, ctx := errgroup.WithContext(ctx)

	// HTTP server
	g.Go(func() error {
		if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// gRPC server
	g.Go(func() error {
		return a.grpcServer.ListenAndServe()
	})

	println("Welcome to grpc-update")
	select {
	case <-signalChan:
		break
	case <-ctx.Done():
		break
	}

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	log.Println("received shutdown signal")

	if a.httpServer != nil {
		_ = a.httpServer.Shutdown(shutdownCtx)
	}
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	return g.Wait()
}
