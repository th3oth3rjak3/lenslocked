package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"lenslocked/views"
)

// The Users controller object.
type Users struct {
	NewView *views.View
}

// Instantiates a new Users controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}

// New is used to render the form where a user can create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// The contents of the signup form which may be null
type SignupForm struct {
	Email    string
	Password string
}

// The bind method checks to ensure that both email and password were provided in the form.
func (s *SignupForm) Bind(r *http.Request) error {
	s.Email = r.PostFormValue("email")
	s.Password = r.PostFormValue("password")
	// TODO: The empty string checking may need to be moved to a validation function.
	// TODO: It may be outisde of the single responsibility principle for the Bind method.
	if s.Email == "" || s.Password == "" {
		return errors.New("email or password were not provided")
	}
	return nil
}

// Used to process the signup request for a new user.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	newUser := &SignupForm{}
	if err := newUser.Bind(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, *newUser)
}
