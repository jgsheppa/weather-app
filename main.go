package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jgsheppa/weather-app/controllers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	weatherController := controllers.NewWeather()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/", weatherController.SearchView)
	r.Post("/search", weatherController.LocationSearch)

	r.Route("/weather/{lat}/{lon}", func(r chi.Router) {
		r.Get("/", weatherController.LocationData)
	})
	r.Route("/location/{name}", func(r chi.Router) {
		r.Get("/", weatherController.LocationResults)
	})
	fmt.Println("application running on http://localhost:3001")
	err = http.ListenAndServe(":3001", r)
	if err != nil {
		panic(err)
	}
}
