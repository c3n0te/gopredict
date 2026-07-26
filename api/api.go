package api

type Station struct {
	StnName    string  `json:"stnname,omitempty" db:"stnname"`
	Latitude   float32 `json:"latitude,omitempty" db:"latitude"`
	Longitude  float32 `json:"longitude,omitempty" db:"longitude"`
	Altitude   float32 `json:"altitude,omitempty" db:"altitude"`
	MinHorizon float32 `json:"minhorizon,omitempty" db:"minhorizon"`
}

type TLE struct {
	SatName string `json:"satname,omitempty" db:"satname"`
	Line1   string `json:"line1,omitempty" db:"line1"`
	Line2   string `json:"line2,omitempty" db:"line2"`
}
