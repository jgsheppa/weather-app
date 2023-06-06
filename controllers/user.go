package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jgsheppa/weather-app/views"
)

func NewUser() *User {
	return &User{
		LoginView:    views.NewView("bootstrap", http.StatusFound, "/user/login"),
		RegisterView: views.NewView("bootstrap", http.StatusFound, "/user/register"),
	}
}

type User struct {
	LoginView    *views.View
	RegisterView *views.View
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	u.RegisterView.Render(w, r, nil)
}

func (u *User) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", u.Login)
	r.Get("/register", u.Register)

	return r
}
