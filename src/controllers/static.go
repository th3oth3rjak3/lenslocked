package controllers

import (
	"lenslocked/views"
)

type Static struct {
	Home    *views.View
	Contact *views.View
}

// Instantiates a new Users controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.html"),
		Contact: views.NewView("bootstrap", "views/static/contact.html"),
	}
}
