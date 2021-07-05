package web

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/cors"
)

type Server struct {
	server *http.Server
}

func New(listenAddr string) *Server {
	httpServer := &http.Server{
		Handler: cors.AllowAll().Handler(NewHandler()),
		Addr:    listenAddr,
	}
	serverCtx := &Server{
		server: httpServer,
	}

	return serverCtx
}

func (s *Server) ListenAndServe() error {
	log.Printf("HTTP server serving at : %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("HTTP server shutting down")
	return s.server.Shutdown(ctx)
}
