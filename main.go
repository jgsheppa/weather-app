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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/", controllers.NewWeather().SearchView)

	r.Mount("/users", controllers.NewUser().Routes())
	r.Mount("/weather", controllers.NewWeather().Routes())

	fmt.Println("application running on http://localhost:3001")
	err = http.ListenAndServe(":3001", r)
	if err != nil {
		panic(err)
	}
}
