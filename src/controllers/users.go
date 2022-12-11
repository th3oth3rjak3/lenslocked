package controllers

import (
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
		NewView: views.NewView("bootstrap", "views/users/new.html"),
	}
}

// Render the new users page.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}
