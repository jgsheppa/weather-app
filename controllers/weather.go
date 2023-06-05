package controllers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jgsheppa/weather-app/views"
	"github.com/jgsheppa/weather-app/weather"
)

func NewWeather() *Weather {
	return &Weather{
		AirQualityView: views.NewView("bootstrap", http.StatusFound, "/weather/location_data"),
		LocationView:   views.NewView("bootstrap", http.StatusFound, "/weather/location"),
		SearchView:     views.NewView("bootstrap", http.StatusFound, "/weather/location_search"),
	}
}

type Weather struct {
	AirQualityView *views.View
	LocationView   *views.View
	SearchView     *views.View
}

func (we *Weather) LocationSearch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	city := r.PostFormValue("search")

	http.Redirect(w, r, "/location/"+city, http.StatusFound)
}

func (we *Weather) LocationData(w http.ResponseWriter, r *http.Request) {
	lat := chi.URLParam(r, "lat")
	lon := chi.URLParam(r, "lon")

	api := weather.ApiParams{
		Lat: lat,
		Lon: lon,
	}

	air, err := api.GetAirQuality()
	if err != nil {
		log.Fatalf("did not receive air quality data: %v", err)
	}

	api.Units = "metric"

	current, err := api.GetCurrentWeatherForLocation()
	if err != nil {
		log.Fatalf("did not receive current weather data: %v", err)
	}

	current.Visibility = current.Visibility / 1000

	forecast, err := api.GetForecast()
	if err != nil {
		log.Fatalf("did not receive current weather data: %v", err)
	}

	data := weather.WeatherCollection{
		Current:    current,
		AirQuality: air,
		Forecast:   forecast,
	}

	we.AirQualityView.Render(w, r, data)
}

func (we *Weather) LocationResults(w http.ResponseWriter, r *http.Request) {
	location := chi.URLParam(r, "name")
	api := weather.ApiParams{
		Q:     location,
		Limit: "10",
	}
	data, err := api.GetCoordinatesByLocationName()
	if err != nil {
		log.Fatalf("did not receive location data: %v", err)
	}

	we.LocationView.Render(w, r, data)
}
