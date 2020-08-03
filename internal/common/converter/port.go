package converter

import (
	"google.golang.org/genproto/googleapis/type/latlng"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/common/model"
)

func ToGRPCPort(p model.Port) *portdomainv1.Port {
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

func FromGRPCPort(p *portdomainv1.Port) model.Port {
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
