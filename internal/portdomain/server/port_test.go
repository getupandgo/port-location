package server

import (
	"context"
	"strconv"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/latlng"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/common/converter"
	"port-location/internal/common/model"
	"port-location/internal/common/test"
	"port-location/internal/portdomain/storage"
)

func TestServer_GetPortByLocode(t *testing.T) {
	tdb := test.NewTestDBClient(t)
	defer tdb.Close(t)
	grpcPorts := populateData(t, tdb.DB)
	server := &Server{storage: storage.NewClient(tdb.DB)}

	tCases := []struct {
		name         string
		locode       string
		expectedPort *portdomainv1.Port
		wantErr      bool
		expectedErr  error
	}{
		{
			name:         "success",
			locode:       "AEAJM",
			expectedPort: grpcPorts[0],
		}, {
			name:        "not found",
			locode:      "QWERTY",
			wantErr:     true,
			expectedErr: ErrNotFound,
		},
	}
	for _, tc := range tCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := server.GetPortByLocode(context.Background(),
				&portdomainv1.GetPortByLocodeRequest{Locode: tc.locode})

			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedPort, got.Port)
			}
		})
	}
}

func TestServer_UpsertPort(t *testing.T) {
	tdb := test.NewTestDBClient(t)
	defer tdb.Close(t)

	storageClient := storage.NewClient(tdb.DB)

	tCases := []struct {
		name string
		port *portdomainv1.Port
	}{
		{
			name: "success",
			port: &portdomainv1.Port{
				Locode:  "AEAJM",
				Name:    "Ajman",
				City:    "Ajman",
				Country: "United Arab Emirates",
				Coordinates: &latlng.LatLng{
					Latitude:  55.513643,
					Longitude: 25.405217,
				},
				Province:    "Ajman",
				Timezone:    "Asia/Dubai",
				Unlocs:      []string{"AEAJM"},
				ForeignCode: 52000,
			},
		}, {
			name: "success when upserting same record",
			port: &portdomainv1.Port{
				Locode:  "AEAJM",
				Name:    "Al Hidd",
				City:    "Al Hidd",
				Country: "Bahrain",
				Coordinates: &latlng.LatLng{
					Latitude:  55.513643,
					Longitude: 25.405217,
				},
				Province:    "Bahrain",
				Timezone:    "Asia/Dubai",
				Unlocs:      []string{"AEAJM"},
				ForeignCode: 52000,
			},
		},
	}
	for _, tc := range tCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			server := &Server{storage: storageClient}
			got, err := server.UpsertPort(context.Background(), &portdomainv1.UpsertPortRequest{Port: tc.port})
			require.NoError(t, err)
			require.Empty(t, got)

			res, err := storageClient.GetPort(context.Background(), tc.port.Locode)
			require.NoError(t, err)
			assert.Equal(t, tc.port, converter.ToGRPCPort(res))
		})
	}
}

func populateData(t *testing.T, db *sqlx.DB) []*portdomainv1.Port {
	ports := []model.Port{
		{
			Locode:  "AEAJM",
			Name:    "Ajman",
			City:    "Ajman",
			Country: "United Arab Emirates",
			Coordinates: model.LatLng{
				Lat: 55.513643,
				Lon: 25.405217,
			},
			Province:    "Ajman",
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAJM"},
			ForeignCode: 52000,
		}, {
			Locode:      "BHAHD",
			Name:        "Al Hidd",
			City:        "Al HIdd",
			Country:     "Bahrain",
			Province:    "Bahrain",
			Unlocs:      []string{"BHAHD"},
			ForeignCode: 52500,
		},
	}

	grpcPorts := make([]*portdomainv1.Port, 0, len(ports))

	floatToString := func(n float64) string {
		return strconv.FormatFloat(n, 'f', 6, 64)
	}

	for _, p := range ports {
		_, err := db.Exec(
			`INSERT INTO ports 
			(locode, name, city, country, alias, regions, lat, lon, province, timezone, unlocs, foreign_code) VALUES 
		    ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			p.Locode, p.Name, p.City, p.Country, pq.Array(p.Alias), pq.Array(p.Regions),
			floatToString(p.Coordinates.Lat), floatToString(p.Coordinates.Lon),
			p.Province, p.Timezone, pq.Array(p.Unlocs), p.ForeignCode)

		require.NoError(t, err)

		grpcPorts = append(grpcPorts, converter.ToGRPCPort(p))
	}

	return grpcPorts
}
