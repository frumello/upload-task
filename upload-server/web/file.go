package web

import (
	"context"
	"io/ioutil"
	"mime"
	"net/http"
	"sync"

	"upload-server/util"

	"github.com/go-chi/chi/v5"
)

const (
	fileNameParam = "fileNameParam"
)

type handler struct {
	http.Handler
	lock sync.RWMutex
}

func NewFileHandler() *handler {
	router := chi.NewRouter()

	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	h := &handler{
		Handler: router,
	}
	router.Route("/", h.files)
	return h
}

func (h *handler) files(router chi.Router) {
	router.Route("/{fileNameParam}", func(router chi.Router) {
		router.Use(fileContext)
		router.Get("/", h.downloadImage)
	})
}

func fileContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := chi.URLParam(r, fileNameParam)
		ctx := context.WithValue(r.Context(), fileNameParam, fileName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) downloadImage(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	fileName := r.Context().Value(fileNameParam).(string)
	if !util.FileExists(fileName) {
		http.NotFound(w, r)
		return
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cd := mime.FormatMediaType("attachment", map[string]string{"filename": fileName})
	w.Header().Set("Content-Disposition", cd)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(content)
}
