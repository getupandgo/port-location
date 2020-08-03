package server

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/genproto/googleapis/type/latlng"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/portdomain/model"
)

var ErrNotFound = errors.New("port not found")

func (s *Server) UpsertPort(ctx context.Context, req *portdomainv1.UpsertPortRequest) (*portdomainv1.UpsertPortResponse, error) {
	if err := s.storage.UpsertPort(ctx, fromGRPCPort(req.Port)); err != nil {
		return nil, err
	}

	return &portdomainv1.UpsertPortResponse{}, nil
}

func (s *Server) GetPortByLocode(ctx context.Context,
	req *portdomainv1.GetPortByLocodeRequest) (*portdomainv1.GetPortByLocodeResponse, error) {
	if req.Locode == "" {
		return nil, errors.New("locode cannot be empty")
	}

	p, err := s.storage.GetPort(ctx, req.Locode)
	switch err {
	case sql.ErrNoRows:
		return nil, ErrNotFound
	case nil:
	default:
		return nil, err
	}

	return &portdomainv1.GetPortByLocodeResponse{Port: toGRPCPort(p)}, nil
}

func toGRPCPort(p model.Port) *portdomainv1.Port {
	return &portdomainv1.Port{
		Locode:  p.Locode,
		Name:    p.Name,
		City:    p.City,
		Country: p.Country,
		Alias:   p.Alias,
		Regions: p.Regions,
		Coordinates: &latlng.LatLng{
			Latitude:  p.Coordinates.Lat,
			Longitude: p.Coordinates.Lon,
		},
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		ForeignCode: p.ForeignCode,
	}
}

func fromGRPCPort(p *portdomainv1.Port) model.Port {
	return model.Port{
		Locode:  p.Locode,
		Name:    p.Name,
		City:    p.City,
		Country: p.Country,
		Alias:   p.Alias,
		Regions: p.Regions,
		Coordinates: model.LatLng{
			Lat: p.Coordinates.GetLatitude(),
			Lon: p.Coordinates.GetLongitude(),
		},
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		ForeignCode: p.ForeignCode,
	}
}
