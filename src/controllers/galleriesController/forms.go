package galleriesController

import (
	"net/http"
)

// The contents of the gallery form which may be null
type GalleryForm struct {
	Title string
}

// The bind method checks to ensure that both email and password were provided in the form.
func (gf *GalleryForm) Bind(r *http.Request) error {
	gf.Title = r.PostFormValue("title")
	return nil
}
