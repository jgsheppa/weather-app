package controllers

import (
	"fmt"
	"net/http"

	"github.com/jgsheppa/weather-app/views"
	"github.com/jgsheppa/weather-app/weather"
)


func NewWeather() *Weather {
	return &Weather{
		NewView: views.NewView("bootstrap", http.StatusFound, "/weather/air_quality"),
	}
}

type Weather struct {
	NewView *views.View
}

func (we *Weather) AirQuality(w http.ResponseWriter, r *http.Request) {
	data, err := weather.GetAirQuality("48.2082", "16.3738")
	if err != nil {
		fmt.Errorf("did not receive air quality data: %v", err)
	}

	we.NewView.Render(w, r, data)
}