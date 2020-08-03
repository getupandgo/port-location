package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"port-location/internal/common/model"
)

type PortDomainClient interface {
	SendPortInfo(ctx context.Context, port model.Port) error
	GetPortInfoByLocode(ctx context.Context, locode string) (model.Port, error)
}

type Server struct {
	Router           *mux.Router
	portDomainClient PortDomainClient
}

func NewServer(client PortDomainClient) *Server {
	s := &Server{
		Router:           mux.NewRouter().StrictSlash(true),
		portDomainClient: client,
	}

	s.initRoutes()

	return s
}

func (s *Server) initRoutes() {
	r := s.Router

	r.HandleFunc("/port/{locode}", s.GetPortByLocode).Methods(http.MethodGet)
}

func writeResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, "%v\n", msg)
}
