package users

import (
	"net/http"
)

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
	return nil
}
