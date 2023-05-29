package main

import (
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
	r.Get("/", weatherController.AirQuality)
	http.ListenAndServe(":3000", r)
}