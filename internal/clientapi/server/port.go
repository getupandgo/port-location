package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"port-location/internal/clientapi/parser"
)

func (s *Server) GetPortByLocode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	locode := params["locode"]
	if locode == "" {
		writeResponse(w, http.StatusBadRequest, "locode cannot be empty")
	}

	port, err := s.portDomainClient.GetPortInfoByLocode(r.Context(), locode)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "failed to get port info")
	}

	b, err := json.Marshal(port)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "failed to marshal port info")
	}

	_, _ = w.Write(b)
}

func (s *Server) ParsePortFile(ctx context.Context, path string) error {
	portFile, err := os.Open(path)
	if err != nil {
		return err
	}

	//deadPortFile, err := os.Open(deadPortFilePath)
	//if err != nil {
	//	return err
	//}

	portCh, errCh := parser.ReadPortInfo(portFile)

	for {
		select {
		case port := <-portCh:
			if err := s.portDomainClient.SendPortInfo(ctx, port); err != nil {
				//err := parser.SaveUnprocessedPort(deadPortFile, port)
				log.Print(err.Error())
			}
		case err := <-errCh:
			return err
		}
	}
}
