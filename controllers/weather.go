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
		AirQualityView: views.NewView("bootstrap", http.StatusFound, "/weather/air_quality"),
		LocationView: views.NewView("bootstrap", http.StatusFound, "/weather/location"),
	}
}

type Weather struct {
	AirQualityView *views.View
	LocationView *views.View
}

func (we *Weather) AirQuality(w http.ResponseWriter, r *http.Request) {
	lat := chi.URLParam(r, "lat")
	lon := chi.URLParam(r, "lon")

	data, err := weather.GetAirQuality(lat, lon)
	if err != nil {
		log.Fatalf("did not receive air quality data: %v", err)
	}

	we.AirQualityView.Render(w, r, data)
}

func (we *Weather) Location(w http.ResponseWriter, r *http.Request) {
	location := chi.URLParam(r, "name")
	data, err := weather.GetLatAndLonByLocationName(location)
	if err != nil {
		log.Fatalf("did not receive location data: %v", err)
	}

	we.LocationView.Render(w, r, data)
}