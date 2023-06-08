package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jgsheppa/weather-app/controllers"
	"github.com/jgsheppa/weather-app/models"
	"github.com/joho/godotenv"
)

func init() {

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s := models.NewServices("dev.db")
	err = s.AutoMigrate()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(controllers.NewUser(s.User).AuthMiddleware)

	r.Handle("/", controllers.NewWeather().SearchView)
	r.Mount("/users", controllers.NewUser(s.User).Routes())
	r.Mount("/weather", controllers.NewWeather().Routes())

	fmt.Println("application running on http://localhost:3001")
	err = http.ListenAndServe(":3002", r)
	if err != nil {
		panic(err)
	}
}
