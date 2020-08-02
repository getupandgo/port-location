package portdomain

import (
	"context"

	"google.golang.org/genproto/googleapis/type/latlng"

	portdomainV1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/clientapi/model"
)

type Client struct {
	conn portdomainV1.PortDomainAPIClient
}

func NewClient(conn portdomainV1.PortDomainAPIClient) Client {
	return Client{
		conn: conn,
	}
}

func (s *Client) SendPortInfo(ctx context.Context, port model.Port) error {
	_, err := s.conn.UpsertPort(ctx, &portdomainV1.UpsertPortRequest{
		Port: toGRPCPort(port),
	})

	return err
}

func (s *Client) GetPortInfoByLocode(ctx context.Context, locode string) (model.Port, error) {
	res, err := s.conn.GetPortByLocode(ctx, &portdomainV1.GetPortByLocodeRequest{
		Locode: locode,
	})
	if err != nil {
		return model.Port{}, err
	}

	return fromGRPCPort(res.Port), nil
}

func toGRPCPort(p model.Port) *portdomainV1.Port {
	return &portdomainV1.Port{
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

func fromGRPCPort(p *portdomainV1.Port) model.Port {
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
