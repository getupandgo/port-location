package storage

import (
	"github.com/lib/pq"

	"port-location/internal/portdomain/model"
)

type Port struct {
	Id          string
	Locode      string
	Name        string
	City        string
	Country     string
	Alias       pq.StringArray
	Regions     pq.StringArray
	Lat         float64
	Lon         float64
	Province    string
	Timezone    string
	Unlocs      pq.StringArray
	ForeignCode int32 `db:"foreign_code"`
}

func toModelPort(p Port) model.Port {
	return model.Port{
		Locode:  p.Locode,
		Name:    p.Name,
		City:    p.City,
		Country: p.Country,
		Alias:   p.Alias,
		Regions: p.Regions,
		Coordinates: model.LatLng{
			Lat: p.Lat,
			Lon: p.Lon,
		},
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		ForeignCode: p.ForeignCode,
	}
}
