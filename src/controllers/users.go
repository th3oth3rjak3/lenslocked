package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"lenslocked/models"
	"lenslocked/rand"
	"lenslocked/views"
)

// The Users controller object.
type Users struct {
	NewView     *views.View
	LoginView   *views.View
	userService models.UserService
}

// Instantiates a new Users controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:     views.NewView("bootstrap", "users/new"),
		LoginView:   views.NewView("bootstrap", "users/login"),
		userService: us,
	}
}

// The contents of the signup form which may be null
type SignupForm struct {
	Name     string
	Email    string
	Password string
}

// The bind method checks to ensure that both email and password were provided in the form.
func (s *SignupForm) Bind(r *http.Request) error {
	s.Email = r.PostFormValue("email")
	s.Password = r.PostFormValue("password")
	s.Name = r.PostFormValue("name")
	if !s.IsValid() {
		return errors.New("email, name, or password were not provided")
	}
	return nil
}

// IsValid checks to see if form data for a Signup Form is empty. If any value
// is not provided, it will return false. If all values are provided, it will
// return true.
func (s *SignupForm) IsValid() bool {
	// TODO: Lookup email in the database to see if it exists?
	return s.Name != "" && s.Password != "" && s.Email != ""
}

// Used to process the signup request for a new user.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	formData := &SignupForm{}
	if err := formData.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := &models.User{
		Name:     formData.Name,
		Email:    formData.Email,
		Password: formData.Password,
	}
	if err := u.userService.Create(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// Login is used to verify the provided email address and password and log
// the user in if they have an account and the credentials are correct.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	formData := &LoginForm{}
	if err := formData.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	usr, err := u.userService.Authenticate(formData.Email, formData.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = u.signIn(w, usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// signIn is used to attach the signed-in cookie to the http response.
func (u *Users) signIn(w http.ResponseWriter, usr *models.User) error {
	// If the remember token is empty, generate one.
	if usr.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		usr.Remember = token
		err = u.userService.Update(usr)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    usr.Remember,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// Represents the form data that is required when logging in.
type LoginForm struct {
	Email    string
	Password string
}

// The bind method checks to ensure that both email and password were provided in the form.
// It also assigns the form values to the LoginForm and returns an error if any of the
// fields are empty.
func (l *LoginForm) Bind(r *http.Request) error {
	l.Email = r.PostFormValue("email")
	l.Password = r.PostFormValue("password")
	if !l.IsValid() {
		return errors.New("email or password were not provided or are invalid")
	}
	return nil
}

// IsValid checks to see if email and password are both provided and not empty
// for the Login Form.
func (l *LoginForm) IsValid() bool {
	// TODO: Lookup email in the database to see if it exists?
	return l.Email != "" && l.Password != ""
}

// CookieTest is a route handler to display cookie information for testing purposes only.
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	usr, err := u.userService.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Should be a user: %+v", usr)
}
