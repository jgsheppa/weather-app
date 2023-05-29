package weather

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

type LocationData struct {
	Name string `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Country string `json:"country"`
	State string `json:"state"`
}


const (
	geocoding = "http://api.openweathermap.org/geo/1.0/direct"
)

func GetLatAndLonByLocationName(location string) ([]LocationData, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API")
	var data []LocationData
	
	v := url.Values{}
	v.Set("appid", apiKey)
	v.Set("q", location)
	v.Set("limit", "5")
	
	baseUrl := geocoding + "?" + v.Encode()
	
	response, err := http.Get(baseUrl)
	if err != nil {
		return []LocationData{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []LocationData{}, err
	}
	json.Unmarshal(body, &data)
	
	return data, nil
}