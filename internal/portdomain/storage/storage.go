package storage

import (
	"context"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"port-location/internal/common/model"
)

type Client struct {
	db *sqlx.DB
}

func NewClient(db *sqlx.DB) Client {
	return Client{
		db: db,
	}
}

/* Database is denormalized due to lack of time.
NoSQL solution isn't chosen vecause data is strongly structured, similar, and mostly constant (city names, country names, etc)

Normalized version would seem like this:
Country (country_name, fk city)
City (city_name, province_name, fk port)
Port (locode, name, alias, lat, lon, foreign_code).

Later tables can be extended to contain airports and documents required to enter the country (visa info, medical data)
This will allow to calculate estimated route between departure and ship boarding and similar work automation processes
Also it would be great to store lat, lon as postgis geography for easier distance calculation (e.g. between ports, port/city, etc)
*/

func (c *Client) UpsertPort(ctx context.Context, port model.Port) error {
	_, err := c.db.ExecContext(ctx,
		`INSERT INTO ports (locode, name, city, country, alias, regions, lat, lon, province, timezone, unlocs, foreign_code) 
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) 
ON CONFLICT (locode) DO 
    UPDATE SET (locode, name, city, country, alias, regions, lat, lon, province, timezone, unlocs, foreign_code) = 
        (EXCLUDED.locode, EXCLUDED.name, EXCLUDED.city, EXCLUDED.country, EXCLUDED.alias, EXCLUDED.regions, 
         EXCLUDED.lat, EXCLUDED.lon, EXCLUDED.province, EXCLUDED.timezone, EXCLUDED.unlocs, EXCLUDED.foreign_code)`,
		port.Locode, port.Name, port.City, port.Country, pq.Array(port.Alias), pq.Array(port.Regions),
		floatToString(port.Coordinates.Lat), floatToString(port.Coordinates.Lon),
		port.Province, port.Timezone, pq.Array(port.Unlocs), port.ForeignCode)

	return err
}

func (c *Client) GetPort(ctx context.Context, locode string) (model.Port, error) {
	var p Port
	if err := c.db.GetContext(ctx, &p, "SELECT * FROM ports WHERE locode = $1", locode); err != nil {
		return model.Port{}, err
	}

	return toModelPort(p), nil
}

func floatToString(n float64) string {
	return strconv.FormatFloat(n, 'f', 6, 64)
}
