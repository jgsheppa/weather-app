package ip

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

type GeoIP struct {
	Lon string `json:"longitude"`
	Lat string `json:"latitude"`
}

type ApiParams struct {
	IP string
}

func (a *ApiParams) FindGeoLocation() (GeoIP, error) {
	apiKey := os.Getenv("IP_GEO_LOCATION_API")
	var data GeoIP

	params := map[string]string{
		"ip":     a.IP,
		"apiKey": apiKey,
	}

	v := url.Values{}

	for key, val := range params {
		v.Set(key, val)
	}

	apiAddress := "https://api.ipgeolocation.io/ipgeo" + "?" + v.Encode()

	response, err := http.Get(apiAddress)
	if err != nil {
		return GeoIP{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return GeoIP{}, err
	}
	json.Unmarshal(body, &data)

	return data, nil
}
