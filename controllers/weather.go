package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jgsheppa/weather-app/context"
	"github.com/jgsheppa/weather-app/ip"
	"github.com/jgsheppa/weather-app/models"
	"github.com/jgsheppa/weather-app/views"
	"github.com/jgsheppa/weather-app/weather"
)

func NewLocation(ls models.LocationService) *Location {
	return &Location{
		LocationDataView:   views.NewView("bootstrap", http.StatusFound, "/weather/location_data"),
		LocationResultView: views.NewView("bootstrap", http.StatusFound, "/weather/location"),
		SearchView:         views.NewView("bootstrap", http.StatusFound, "/weather/location_search"),
		ls:                 ls,
	}
}

type Location struct {
	LocationDataView   *views.View
	LocationResultView *views.View
	SearchView         *views.View
	ls                 models.LocationService
	r                  chi.Router
}

func (l *Location) Home(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-REAL-IP", "77.141.137.113")
	ipAddress := r.Header.Get("X-REAL-IP")
	ipApi := ip.ApiParams{
		IP: ipAddress,
	}
	geoLocation, err := ipApi.FindGeoLocation()

	weatherApi := weather.ApiParams{
		Lat: geoLocation.Lat,
		Lon: geoLocation.Lon,
	}
	current, err := weatherApi.GetCurrentWeatherForLocation()
	fmt.Printf("current weather: %v", current)

	var vd views.Data

	user := context.User(r.Context())
	if user == nil {
		l.SearchView.Render(w, r, vd)
		return
	}
	savedLocations, err := l.ls.GetByUserId(user.ID)
	if err != nil {
		vd.SetAlert(err)
		return
	}
	vd.Yield = savedLocations
	vd.User = user
	l.SearchView.Render(w, r, vd)
}

func (l *Location) LocationSearch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	city := r.PostFormValue("search")

	http.Redirect(w, r, "/weather/location/"+city, http.StatusFound)
}

func (l *Location) LocationData(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
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

	savedLocation, err := l.ls.FindByLonAndLat(current.Coord.Lon, current.Coord.Lat)

	vd.Yield = weather.WeatherCollection{
		Current:    current,
		AirQuality: air,
		Forecast:   forecast,
		IsSaved:    savedLocation.IsSaved,
	}

	l.LocationDataView.Render(w, r, vd)
}

func (l *Location) LocationResults(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	location := chi.URLParam(r, "name")
	api := weather.ApiParams{
		Q:     location,
		Limit: "10",
	}
	coordinates, err := api.GetCoordinatesByLocationName()
	if err != nil {
		log.Fatalf("did not receive location data: %v", err)
	}

	vd.Yield = coordinates

	l.LocationResultView.Render(w, r, vd)
}

func (l *Location) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	lon := r.URL.Query().Get("lon")
	lat := r.URL.Query().Get("lat")
	name := r.URL.Query().Get("name")

	user := context.User(r.Context())

	location := models.Location{
		Lon:     lon,
		Lat:     lat,
		UserId:  user.ID,
		Name:    name,
		IsSaved: true,
	}

	err := l.ls.Create(&location)
	if err != nil {
		vd.SetAlert(err)
		// TODO: make url prettier
		http.Redirect(w, r, "/weather/"+name+"/"+lat+"/"+lon, http.StatusFound)
		return
	}
	// TODO: make url prettier
	http.Redirect(w, r, "/weather/"+name+"/"+lat+"/"+lon, http.StatusFound)
}

func (l *Location) Delete(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	lon := r.URL.Query().Get("lon")
	lat := r.URL.Query().Get("lat")
	name := r.URL.Query().Get("name")

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		vd.SetAlert(err)
		return
	}
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		vd.SetAlert(err)
		return
	}
	location, err := l.ls.FindByLonAndLat(lonFloat, latFloat)
	if err != nil {
		vd.SetAlert(err)
		return
	}
	err = l.ls.Delete(location.ID)
	if err != nil {
		vd.SetAlert(err)
		// TODO: make url prettier
		http.Redirect(w, r, "/weather/"+name+"/"+lat+"/"+lon, http.StatusFound)
		return
	}
	// TODO: make url prettier
	http.Redirect(w, r, "/weather/"+name+"/"+lat+"/"+lon, http.StatusFound)
}

func (l *Location) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/search", l.LocationSearch)
		r.Get("/{name}/{lat}/{lon}", l.LocationData)
		r.Post("/location/delete", l.Delete)
		r.Post("/location/save", l.Create)
		r.Get("/location/{name}", l.LocationResults)
	})

	return r
}
