package server

import (
	"port-location/internal/portdomain/storage"
)

type Server struct {
	storage storage.Client
}

func NewServer(storage storage.Client) *Server {
	s := &Server{
		storage: storage,
	}

	return s
}
