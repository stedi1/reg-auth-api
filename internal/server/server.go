package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	*http.Server
}

// создает новый сервер
func NewServer(router chi.Router) *Server {
	return &Server{
		&http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}
