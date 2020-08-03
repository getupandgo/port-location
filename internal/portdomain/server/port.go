package server

import (
	"context"
	"database/sql"
	"errors"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/common/converter"
)

var (
	ErrNotFound      = errors.New("port not found")
	ErrInvalidLocode = errors.New("locode cannot be empty")
)

func (s *Server) UpsertPort(ctx context.Context, req *portdomainv1.UpsertPortRequest) (*portdomainv1.UpsertPortResponse, error) {
	if err := s.storage.UpsertPort(ctx, converter.FromGRPCPort(req.Port)); err != nil {
		return nil, err
	}

	return &portdomainv1.UpsertPortResponse{}, nil
}

func (s *Server) GetPortByLocode(ctx context.Context,
	req *portdomainv1.GetPortByLocodeRequest) (*portdomainv1.GetPortByLocodeResponse, error) {
	if req.Locode == "" {
		return nil, ErrInvalidLocode
	}

	p, err := s.storage.GetPort(ctx, req.Locode)
	switch err {
	case sql.ErrNoRows:
		return nil, ErrNotFound
	case nil:
	default:
		return nil, err
	}

	return &portdomainv1.GetPortByLocodeResponse{Port: converter.ToGRPCPort(p)}, nil
}
