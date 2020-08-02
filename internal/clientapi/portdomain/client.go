package portdomain

import (
	"context"

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
		Port: model.ToGRPCPort(port),
	})

	return err
}
