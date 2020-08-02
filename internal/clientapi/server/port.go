package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

	w.Write(b)
}
