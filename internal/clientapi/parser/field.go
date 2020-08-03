package parser

import (
	"strconv"

	"port-location/internal/common/model"
)

type rawPort map[string]interface{}

// wrapper methods on jstream lib for casting interface{} values to concrete types
func (rp rawPort) asString(k string) string {
	var v string

	if rp[k] != nil {
		v = rp[k].(string)
	}

	return v
}

func (rp rawPort) asStringSlice(k string) []string {
	if rp[k] == nil {
		return nil
	}

	interfaceSlice := rp[k].([]interface{})
	vls := make([]string, 0, len(interfaceSlice))

	for _, v := range interfaceSlice {
		vls = append(vls, v.(string))
	}

	return vls
}

func (rp rawPort) asFloatSlice(k string) []float64 {
	if rp[k] == nil {
		return nil
	}

	interfaceSlice := rp[k].([]interface{})
	vls := make([]float64, 0, len(interfaceSlice))

	for _, v := range interfaceSlice {
		vls = append(vls, v.(float64))
	}

	return vls
}

func toModelPort(locode string, rp rawPort) (model.Port, error) {
	var (
		fc    int32
		coord model.LatLng
	)

	if foreignCode := rp.asString("code"); foreignCode != "" {
		f, err := strconv.Atoi(foreignCode)
		if err != nil {
			return model.Port{}, err
		}

		fc = int32(f)
	}

	if coordinates := rp.asFloatSlice("coordinates"); coordinates != nil {
		if len(coordinates) == 2 {
			coord = model.LatLng{
				Lat: coordinates[0],
				Lon: coordinates[1],
			}
		}
	}

	return model.Port{
		Locode:      locode,
		Name:        rp.asString("name"),
		City:        rp.asString("city"),
		Country:     rp.asString("country"),
		Alias:       rp.asStringSlice("alias"),
		Regions:     rp.asStringSlice("regions"),
		Coordinates: coord,
		Province:    rp.asString("province"),
		Timezone:    rp.asString("timezone"),
		Unlocs:      rp.asStringSlice("unlocs"),
		ForeignCode: fc,
	}, nil
}
