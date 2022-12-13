package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"lenslocked/models"
	"lenslocked/views"
)

// The Users controller object.
type Users struct {
	NewView     *views.View
	userService *models.UserService
}

// Instantiates a new Users controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:     views.NewView("bootstrap", "users/new"),
		userService: us,
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
	user := models.User{
		Name:  formData.Name,
		Email: formData.Email,
	}
	if err := u.userService.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, *formData)
}
