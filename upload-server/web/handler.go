package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewHandler() http.Handler {
	router := chi.NewRouter()

	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)

	fileHandler := NewFileHandler()
	router.Route("/", fileHandler.files)

	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
}
