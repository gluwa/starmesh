package util

import (
	"encoding/json"
	"net/http"
)

type GeoInfo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GetGeoLocation() (float64, float64, error) {
	resp, err := http.Get("http://ip-api.com/json")
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var info GeoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return 0, 0, err
	}
	return info.Lat, info.Lon, nil
}
