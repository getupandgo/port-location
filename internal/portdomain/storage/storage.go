package storage

import (
	"context"
	"database/sql"

	"port-location/internal/portdomain/model"
)

type Client struct {
	db *sql.DB
}

func NewClient(db *sql.DB) Client {
	return Client{
		db: db,
	}
}

func (c *Client) UpsertPort(ctx context.Context, port model.Port) error {
	_, err := c.db.ExecContext(ctx,
		`INSERT INTO ports (locode, name, city, country, alias, regions, coordinates, province, timezone, unlocs, foreign_code) 
				VALUES (?,?,?,?,?,?,?,?,?,?,?) 
				ON CONFLICT (locode) DO 
				UPDATE SET (locode, name, city, country, alias, regions, coordinates, province, timezone, unlocs, foreign_code) = 
				    (EXCLUDED.locode, EXCLUDED.city, EXCLUDED.country, EXCLUDED.alias, EXCLUDED.regions, EXCLUDED.coordinates, 
				     EXCLUDED.province, EXCLUDED.timezone, EXCLUDED.unlocs,EXCLUDED.foreign_code)`,
		port.Locode, port.Name, port.City, port.Country, port.Alias, port.Regions, port.Coordinates,
		port.Province, port.Timezone, port.Unlocs, port.ForeignCode)
	return err
}

func (c *Client) GetPort(ctx context.Context, locode string) (model.Port, error) {
	var p model.Port
	if err := c.db.QueryRowContext(ctx, "SELECT * FROM ports WHERE locode = ?", locode).Scan(&p); err != nil {
		return model.Port{}, err
	}

	return p, nil
}
