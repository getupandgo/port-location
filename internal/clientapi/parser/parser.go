package parser

import (
	"io"

	"port-location/internal/clientapi/model"
)

type Service struct {
}

func (s *Service) ReadPortInfo(r io.Reader) (chan<- model.Port, error) {
	return nil, nil
}

func (s *Service) SaveUnprocessedPort(w io.Writer) error {
	return nil
}
