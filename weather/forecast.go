package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	forecast   = "http://api.openweathermap.org/data/2.5/forecast"
)

func GetForecast() (string, error) {
	var data map[string]interface{}

	v := url.Values{}
	v.Set("appid", "f1a2ba10ec24c3ba18b106cf9fabadb9")
	v.Set("id", "524901")

	baseUrl := forecast + "?" + v.Encode()
	fmt.Printf("url: %v \n", baseUrl)

	response, err := http.Get(baseUrl)
	if err != nil {
		return "", err
	}
	fmt.Println(response)
	body, err := io.ReadAll(response.Body)
	json.Unmarshal(body, &data)
	city := data["city"]
	fmt.Printf("res: %v", city)

	return "", nil
}