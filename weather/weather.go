package weather

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	BaseUrl       = "https://api.openweathermap.org"
	CurrentUrl    = "/data/2.5/weather"
	GeocodingUrl  = "/geo/1.0/direct"
	ForecastUrl   = "/data/2.5/forecast"
	AirQualityUrl = "/data/2.5/air_pollution"
)

type WeatherCollection struct {
	Current    CurrentWeather
	AirQuality AirQualityList
	Forecast
}

type ApiParams struct {
	Lat   string
	Lon   string
	Units string
	Lang  string
	Limit string
	Q     string
	Cnt   string
}

type WeatherData struct {
	Id          string `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}

type Rain struct {
	OneHour   float64 `json:"1h"`
	ThreeHour float64 `json:"3h"`
}

type WeatherMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
}

type CurrentWeather struct {
	Weather    []WeatherData `json:"weather"`
	Visibility int           `json:"visibility"`
	Main       WeatherMain   `json:"main"`
	Wind       Wind          `json:"wind"`
	City       string        `json:"name"`
	Rain       Rain          `json:"rain"`
}

func (a *ApiParams) GetCurrentWeatherForLocation() (CurrentWeather, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API")
	var data CurrentWeather

	params := map[string]string{
		"appid": apiKey,
		"lat":   a.Lat,
		"lon":   a.Lon,
		"units": a.Units,
		"lang":  a.Lang,
	}

	urlWithParams := BuildWeatherUrl(CurrentUrl, params)

	response, err := http.Get(urlWithParams)
	if err != nil {
		return CurrentWeather{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return CurrentWeather{}, err
	}
	json.Unmarshal(body, &data)

	return data, nil
}

type LocationData struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state"`
}

func (a *ApiParams) GetCoordinatesByLocationName() ([]LocationData, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API")
	var data []LocationData

	params := map[string]string{
		"appid": apiKey,
		"q":     a.Q,
		"limit": a.Limit,
	}

	urlWithParams := BuildWeatherUrl(GeocodingUrl, params)

	response, err := http.Get(urlWithParams)
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

type AirQualityList struct {
	Dt   int `json:"dt"`
	Main struct {
		Aqi         int `json:"aqi"`
		Description string
	}
	Components struct {
		Co   float64 `json:"co"`
		No   float64 `json:"no"`
		No2  float64 `json:"no2"`
		O3   float64 `json:"o3"`
		So2  float64 `json:"so2"`
		Pm25 float64 `json:"pm2_5"`
		Pm10 float64 `json:"pm10"`
		NH3  float64 `json:"nh3"`
	}
}

type AirQuality struct {
	AirQualityList []AirQualityList `json:"list"`
}

func (a *ApiParams) GetAirQuality() (AirQualityList, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API")
	var data AirQuality

	params := map[string]string{
		"appid": apiKey,
		"lat":   a.Lat,
		"lon":   a.Lon,
	}

	urlWithParams := BuildWeatherUrl(AirQualityUrl, params)

	response, err := http.Get(urlWithParams)
	if err != nil {
		return AirQualityList{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return AirQualityList{}, err
	}
	json.Unmarshal(body, &data)

	if len(data.AirQualityList) == 0 {
		return AirQualityList{}, errors.New("received empty list for air quality data")
	}

	data.AirQualityList[0].Main.Description = GetAirQualityDescription(data.AirQualityList[0].Main.Aqi)

	return data.AirQualityList[0], nil
}

type Forecast struct {
	List []ForecastList `json:"list"`
}

type ForecastList struct {
	DateTime   int64 `json:"dt"`
	Hour       int
	HourLocale string
	Main       WeatherMain   `json:"main"`
	Weather    []WeatherData `json:"weather"`
}

func (a *ApiParams) GetForecast() (Forecast, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API")
	var data Forecast

	params := map[string]string{
		"appid": apiKey,
		"lat":   a.Lat,
		"lon":   a.Lon,
		"cnt":   "12",
	}

	urlWithParams := BuildWeatherUrl(ForecastUrl, params)

	response, err := http.Get(urlWithParams)
	if err != nil {
		return Forecast{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Forecast{}, err
	}
	json.Unmarshal(body, &data)

	if len(data.List) == 0 {
		return Forecast{}, errors.New("received empty list for forecast data")
	}

	for i, point := range data.List {
		t := time.Unix(point.DateTime, 0)
		data.List[i].Hour = t.Hour()
		hour := strconv.Itoa(t.Hour())
		if len(hour) == 1 {
			data.List[i].HourLocale = "0" + hour + ".00"
			continue
		}
		data.List[i].HourLocale = hour + ".00"
	}

	return data, nil
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

func BuildWeatherUrl(path string, params map[string]string) string {
	v := url.Values{}

	for key, val := range params {
		v.Set(key, val)
	}

	return BaseUrl + path + "?" + v.Encode()
}
