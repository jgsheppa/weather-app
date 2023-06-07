package controllers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jgsheppa/weather-app/auth"
	"github.com/jgsheppa/weather-app/context"
	"github.com/jgsheppa/weather-app/models"
	"github.com/jgsheppa/weather-app/views"
)

func NewUser(us models.UserService) *User {
	return &User{
		LoginView:    views.NewView("bootstrap", http.StatusFound, "/user/login"),
		RegisterView: views.NewView("bootstrap", http.StatusFound, "/user/register"),
		us:           us,
	}
}

type User struct {
	LoginView    *views.View
	RegisterView *views.View
	us           models.UserService
}

func (u *User) LoginStatic(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	u.LoginView.Render(w, r, nil)
}

func (u *User) RegisterStatic(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	var form RegistrationForm
	parseURLParams(r, &form)
	u.RegisterView.Render(w, r, nil)
}

func (u *User) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", u.LoginStatic)
	r.Post("/login", u.Login)
	r.Post("/logout", u.Logout)
	r.Get("/register", u.RegisterStatic)
	r.Post("/register", u.Create)

	return r
}

type RegistrationForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to process the registration form
//
// POST /register
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form RegistrationForm
	vd.Yield = &form

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.RegisterView.Render(w, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.RegisterView.Render(w, r, vd)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/users/login", http.StatusNotFound)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Welcome to the Weather-App! You've successfully created your account.",
	}
	views.RedirectAlert(w, r, "/", http.StatusFound, alert)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)

	user := context.User(r.Context())

	token, _ := auth.RememberToken()
	user.Remember = token
	u.us.Update(user)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *User) GetUser(w http.ResponseWriter, r *http.Request) *models.User {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return user
}

func (u *User) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := auth.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)
	return nil
}

// HTTP middleware setting a value on the request context
func (u *User) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the user is logged in then pass the user for the navbar
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := u.us.ByRemember(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
