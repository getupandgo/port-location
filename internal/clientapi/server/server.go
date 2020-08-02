package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"port-location/internal/clientapi/portdomain"
)

type Server struct {
	Router           *mux.Router
	portDomainClient portdomain.Client
}

func NewServer(client portdomain.Client) *Server {
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
