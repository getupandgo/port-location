package model

type LatLng struct {
	Lat float64
	Lon float64
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
