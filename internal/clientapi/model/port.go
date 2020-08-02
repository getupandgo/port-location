package model

import (
	"google.golang.org/genproto/googleapis/type/latlng"

	portdomainV1 "port-location/api/proto/portdomain/v1"
)

type LatLng struct {
	lat float64
	lon float64
}

type Port struct {
	Locode      string   `json:"locode"`
	Name        string   `json:"name"`
	City        string   `json:"city"`
	Country     string   `json:"country"`
	Alias       []string `json:"alias"`
	Regions     []string `json:"regions"`
	Coordinates LatLng   `json:"coordinates"`
	Province    string   `json:"province"`
	Timezone    string   `json:"timezone"`
	Unlocs      []string `json:"unlocs"`
	ForeignCode int32    `json:"foreign_code"`
}

func ToGRPCPort(p Port) *portdomainV1.Port {
	return &portdomainV1.Port{
		Locode:  p.Locode,
		Name:    p.Name,
		City:    p.City,
		Country: p.Country,
		Alias:   p.Alias,
		Regions: p.Regions,
		Coordinates: &latlng.LatLng{
			Latitude:  p.Coordinates.lat,
			Longitude: 0,
		},
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		ForeignCode: p.ForeignCode,
	}
}
