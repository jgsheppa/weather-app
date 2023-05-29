package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	airQuality = "http://api.openweathermap.org/data/2.5/air_pollution"
)

type List struct {
	Dt   int `json:"dt"`
	Main struct {
		Aqi int `json:"aqi"`
		Description string
	}
	Components struct {
		Co float64 `json:"co"`
		No float64 `json:"no"`
		No2 float64 `json:"no2"`
		O3	float64 `json:"o3"`
		So2 float64 `json:"so2"`
		Pm25 float64 `json:"pm2_5"`
		Pm10 float64 `json:"pm10"`
		NH3 float64 `json:"nh3"`
	}
}

type AirQuality struct {
	List []List `json:"list"`
}

func GetAirQuality(lat, lon string) (List, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API")
	var data AirQuality

	v := url.Values{}
	v.Set("appid", apiKey)
	v.Set("lat", lat)
	v.Set("lon", lon)

	baseUrl := airQuality + "?" + v.Encode()

	response, err := http.Get(baseUrl)
	if err != nil {
		return List{}, err
	}
	fmt.Println(response)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return List{}, err
	}
	json.Unmarshal(body, &data)

	if len(data.List) == 0 {
		return List{}, errors.New("received empty list for air quality data")
	}

	data.List[0].Main.Description = GetAirQualityDescription(data.List[0].Main.Aqi)

	return data.List[0], nil
}

func GetAirQualityDescription(quality int) string {
	switch quality {
	case 1:
		return "Good"
	case 2:
		return "Fair"
	case 3:
		return "Moderate"
	case 4:
		return "Poor"
	case 5:
		return "Very Poor"
	default:
		return "Unable to retrieve air quality data"
	}
}