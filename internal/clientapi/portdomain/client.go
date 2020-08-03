package portdomain

import (
	"context"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/common/converter"
	"port-location/internal/common/model"
)

// wrapper around grpc service to abstract transport layer logic
type Client struct {
	conn portdomainv1.PortDomainAPIClient
}

func NewClient(conn portdomainv1.PortDomainAPIClient) *Client {
	return &Client{
		conn: conn,
	}
}

func (s *Client) SendPortInfo(ctx context.Context, port model.Port) error {
	_, err := s.conn.UpsertPort(ctx, &portdomainv1.UpsertPortRequest{
		Port: converter.ToGRPCPort(port),
	})

	return err
}

func (s *Client) GetPortInfoByLocode(ctx context.Context, locode string) (model.Port, error) {
	res, err := s.conn.GetPortByLocode(ctx, &portdomainv1.GetPortByLocodeRequest{
		Locode: locode,
	})
	if err != nil {
		return model.Port{}, err
	}

	return converter.FromGRPCPort(res.Port), nil
}
